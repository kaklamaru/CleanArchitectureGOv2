package usecase

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"go-clean-arch/pkg/utility"
	"go-clean-arch/pkg/utility/filesystem"
	"go-clean-arch/repository"
	"go-clean-arch/structure/entity"
	"go-clean-arch/structure/request"
	"go-clean-arch/structure/response"
	"mime/multipart"
	"os"
	"strings"
)

type EventUsecase interface {
	CreateEvent(req *request.EventRequest, claims map[string]interface{}) error
	GetAllEvent() ([]response.EventResponse, error)
	GetEventByID(id uint) (*response.EventResponse, error)
	ToggleEventStatus(eventID uint, claims map[string]interface{}) (bool, error)
	DeleteEventByID(eventID uint, claims map[string]interface{}) error
	UpdateEventByID(eventID uint, claims map[string]interface{}, req request.EventRequest) error
	MyEvent(claims map[string]interface{}) ([]response.EventResponse, error)
	AllAllowedEvent() ([]response.EventResponse, error)
	AllCurrentEvent() ([]response.EventResponse, error)
	MyEventThisYear(claims map[string]interface{},year uint) ([]response.MyInside,[]response.MyOutside,error)

	JoinEvent(eventID uint, claims map[string]interface{}) error
	UnJoinEvent(eventID uint, claims map[string]interface{}) error
	UploadFile(eventID uint, claims map[string]interface{}, file *multipart.FileHeader) error
	GetFile(eventID uint, userID uint) (string, error)
	MyChecklist(eventID uint, claims map[string]interface{}) ([]response.MyChecklist, error)
	UpdateEventStatusAndComment(eventID uint, userID uint, status bool, comment string) error

	CreateEventOutside(req request.OutsideRequest,claims map[string]interface{}) error
	DeleteEventOutsideByID(eventID uint) error
	GetEventOutsideByID(eventID uint) (*response.OutsideResponse, error)
	CreateFile(eventID uint) ([]byte, string, error)
	GetFileOutside(eventID uint ,userID uint)(string,error)
	UploadFileOutside(eventID uint, claims map[string]interface{}, file *multipart.FileHeader) error 

}

type eventUsecase struct {
	userRepo    repository.UserRepository
	facultyRepo repository.FacultyBranchRepository
	eventRepo   repository.EventRepository
}

func NewEventUsecase(userRepo repository.UserRepository, facultyRepo repository.FacultyBranchRepository, eventRepo repository.EventRepository) EventUsecase {
	return &eventUsecase{
		userRepo:    userRepo,
		facultyRepo: facultyRepo,
		eventRepo:   eventRepo,
	}
}

func mapEventResponse(event entity.Event, count uint) (*response.EventResponse, error) {
	branches, err := utility.DecodeIDs(event.BranchIDs)
	if err != nil {
		return &response.EventResponse{}, err
	}

	years, err := utility.DecodeIDs(event.Years)
	if err != nil {
		return &response.EventResponse{}, err
	}
	limit := event.FreeSpace + count

	return &response.EventResponse{
		EventID:        event.EventID,
		EventName:      event.EventName,
		StartDate:      utility.FormatToThaiDate(event.StartDate),
		StartTime:      utility.FormatToThaiTime(event.StartDate),
		SchoolYear:     event.SchoolYear,
		WorkingHour:    event.WorkingHour,
		Limit:          limit,
		FreeSpace:      event.FreeSpace,
		Detail:         event.Detail,
		Location:       event.Location,
		BranchIDs:      branches,
		Status:         event.Status,
		Years:          years,
		AllowAllBranch: event.AllowAllBranch,
		AllowAllYear:   event.AllowAllYear,
		Creator: struct {
			UserID    uint   `json:"user_id"`
			TitleName string `json:"title_name"`
			FirstName string `json:"first_name"`
			LastName  string `json:"last_name"`
			Phone     string `json:"phone"`
			Code      string `json:"code"`
		}{
			UserID:    event.Teacher.UserID,
			TitleName: event.Teacher.TitleName,
			FirstName: event.Teacher.FirstName,
			LastName:  event.Teacher.LastName,
			Phone:     event.Teacher.Phone,
			Code:      event.Teacher.Code,
		},
	}, nil
}

// func mapEventOutside(outside entity.EventOutside) (*response.OutsideResponse,error){
// 	outsideRes := response.OutsideResponse{
// 		EventID:     outside.EventID,
// 		EventName:   outside.EventName,
// 		Location:    outside.Location,
// 		SchoolYear: outside.SchoolYear,
// 		StartDate:   outside.StartDate,
// 		WorkingHour: outside.WorkingHour,
// 		Intendant:   outside.Intendant,
// 		Student: response.StudentResponse{
// 			UserID:      outside.Student.UserID,
// 			TitleName:   outside.Student.TitleName,
// 			FirstName:   outside.Student.FirstName,
// 			LastName:    outside.Student.LastName,
// 			Phone:       outside.Student.Phone,
// 			Code:        outside.Student.Code,
// 			BranchID: outside.Student.BranchId,
// 			BranchName:  outside.Student.Branch.BranchName,
// 			FacultyID: outside.Student.Branch.Faculty.FacultyID,
// 			FacultyName: outside.Student.Branch.Faculty.FacultyName,
// 		},
// 	}
// 	return &outsideRes, nil
// }

func (u *eventUsecase) validateBranches(branches []uint) error {
	for _, branchID := range branches {
		exists, err := u.facultyRepo.BranchExists(branchID)
		if err != nil {
			return fmt.Errorf("error checking branch: %v", err)
		}
		if !exists {
			return fmt.Errorf("branch with ID %d does not exist", branchID)
		}
	}
	return nil
}

func (u *eventUsecase) buildPermission(branches []uint, years []uint) (*entity.Permission, error) {
	if len(branches) > 0 {
		if err := u.validateBranches(branches); err != nil {
			return nil, err
		}
	}

	createPermission := func(branches []uint, years []uint, allowAllBranch, allowAllYear bool) (*entity.Permission, error) {
		branchData, err := json.Marshal(branches)
		if err != nil {
			return nil, err
		}
		yearData, err := json.Marshal(years)
		if err != nil {
			return nil, err
		}
		return &entity.Permission{
			BranchIDs:      string(branchData),
			Years:          string(yearData),
			AllowAllBranch: allowAllBranch,
			AllowAllYear:   allowAllYear,
		}, nil
	}

	switch {
	case len(branches) > 0 && len(years) > 0:
		return createPermission(branches, years, false, false)
	case len(branches) > 0:
		return createPermission(branches, nil, false, true)
	case len(years) > 0:
		return createPermission(nil, years, true, false)
	default:
		return createPermission(nil, nil, true, true)
	}
}

func (u *eventUsecase) CreateEvent(req *request.EventRequest, claims map[string]interface{}) error {
	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return fmt.Errorf("invalid user_id in claims")
	}
	userID := uint(userIDFloat)

	permission, err := u.buildPermission(req.Branches, req.Years)
	if err != nil {
		return err
	}

	startDate, err := utility.ParseStartDate(req.StartDate)
	if err != nil {
		return err
	}

	event := &entity.Event{
		EventName:      req.EventName,
		StartDate:      startDate,
		SchoolYear:     req.SchoolYear,
		FreeSpace:      req.FreeSpace,
		WorkingHour:    req.WorkingHour,
		Detail:         req.Detail,
		Location:       req.Location,
		Creator:        userID,
		AllowAllBranch: permission.AllowAllBranch,
		AllowAllYear:   permission.AllowAllYear,
		BranchIDs:      permission.BranchIDs,
		Years:          permission.Years,
	}

	if err := u.eventRepo.CreateEvent(event); err != nil {
		return err
	}

	userIDs, err := u.userRepo.GetAllStudentID()
	if err != nil {
		return fmt.Errorf("failed to get users for event: %w", err)
	}

	for _, uid := range userIDs {
		news := entity.News{
			Title:   "กิจกรรมใหม่",
			UserID:  uid,
			Message: fmt.Sprintf("กิจกรรม'%s' '%s' '%s'", event.EventName, utility.FormatToThaiDate(event.StartDate), utility.FormatToThaiTime(event.StartDate)),
		}
		if err := u.eventRepo.NewsForUser(&news); err != nil {
			return fmt.Errorf("failed to send news to user %d: %w", uid, err)
		}
	}
	return nil
}

func (u *eventUsecase) GetAllEvent() ([]response.EventResponse, error) {
	events, err := u.eventRepo.GetAllEvent()
	if err != nil {
		return nil, err
	}
	var res []response.EventResponse
	for _, event := range events {
		count, err := u.eventRepo.CountEventInside(event.EventID)
		if err != nil {
			return nil, err
		}
		mappedEvent, err := mapEventResponse(event, count)
		if err != nil {
			return nil, err
		}
		res = append(res, *mappedEvent)
	}
	return res, nil
}

func (u *eventUsecase) GetEventByID(id uint) (*response.EventResponse, error) {
	event, err := u.eventRepo.GetEventByID(id)
	if err != nil {
		return nil, err
	}
	count, err := u.eventRepo.CountEventInside(id)
	if err != nil {
		return nil, fmt.Errorf("error calculating free space")
	}
	return mapEventResponse(*event, count)
}

func (u *eventUsecase) ToggleEventStatus(eventID uint, claims map[string]interface{}) (bool, error) {
	event, err := u.eventRepo.GetEventByID(eventID)
	if err != nil {
		return false, fmt.Errorf("event not found")
	}
	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return false, fmt.Errorf("invalid user_id in claims")
	}
	userID := uint(userIDFloat)

	if event.Creator != userID {
		return false, fmt.Errorf("you do not have permission to edit this event")
	}

	newStatus, err := u.eventRepo.ToggleEventStatus(event.EventID)
	if err != nil {
		return false, err
	}

	return newStatus, nil
}

func (u *eventUsecase) UpdateEventByID(eventID uint, claims map[string]interface{}, req request.EventRequest) error {
	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return fmt.Errorf("invalid user_id in claims")
	}
	userID := uint(userIDFloat)

	return u.eventRepo.UpdateEventWithTransaction(eventID, userID, req)
}


func (u *eventUsecase) DeleteEventByID(eventID uint, claims map[string]interface{}) error {
	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return fmt.Errorf("invalid user_id in claims")
	}
	userID := uint(userIDFloat)

	return u.eventRepo.DeleteEventWithTransaction(eventID, userID)
}

func (u *eventUsecase) MyEvent(claims map[string]interface{}) ([]response.EventResponse, error) {
	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return nil, fmt.Errorf("invalid user_id in claims")
	}
	userID := uint(userIDFloat)

	events, err := u.eventRepo.MyEvent(userID)
	if err != nil {
		return nil, err
	}
	var res []response.EventResponse
	for _, event := range events {
		count, err := u.eventRepo.CountEventInside(event.EventID)
		if err != nil {
			return nil, err
		}
		mappedEvent, err := mapEventResponse(event, count)
		if err != nil {
			return nil, err
		}
		res = append(res, *mappedEvent)
	}
	return res, nil
}

func (u *eventUsecase) AllAllowedEvent() ([]response.EventResponse, error) {
	events, err := u.eventRepo.AllAllowedEvent()
	if err != nil {
		return nil, err
	}
	var res []response.EventResponse
	for _, event := range events {
		count, err := u.eventRepo.CountEventInside(event.EventID)
		if err != nil {
			return nil, err
		}
		mappedEvent, err := mapEventResponse(event, count)
		if err != nil {
			return nil, err
		}
		res = append(res, *mappedEvent)
	}
	return res, nil
}

func (u *eventUsecase) AllCurrentEvent() ([]response.EventResponse, error) {
	events, err := u.eventRepo.AllCurrentEvent()
	if err != nil {
		return nil, err
	}
	var res []response.EventResponse
	for _, event := range events {
		count, err := u.eventRepo.CountEventInside(event.EventID)
		if err != nil {
			return nil, err
		}
		mappedEvent, err := mapEventResponse(event, count)
		if err != nil {
			return nil, err
		}
		res = append(res, *mappedEvent)
	}
	return res, nil
}

func (u *eventUsecase) MyEventThisYear(claims map[string]interface{},year uint) ([]response.MyInside,[]response.MyOutside,error){
	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return nil,nil, fmt.Errorf("invalid user_id in claims")
	}
	userID := uint(userIDFloat)
	
	inside,err:= u.eventRepo.AllEventInsideThisYear(userID,year)
	if err != nil {
		return nil,nil,err
	}
	var insideEvents []response.MyInside
	for _, event := range inside {
		mappedEvent := response.MyInside{
			EventID: event.EventId,
			EventName: event.Event.EventName,
			Location: event.Event.Location,
			StartDate: utility.FormatToThaiDate(event.Event.StartDate),
			StartTime: utility.FormatToThaiTime(event.Event.StartDate),
			WorkingHour: event.Event.WorkingHour,
			SchoolYear: event.Event.SchoolYear,
			Status:event.Status,
			Comment: event.Comment,
			File: event.File,
		}
		insideEvents = append(insideEvents, mappedEvent)
	}
	outside,err:=u.eventRepo.AllEventOutsideThisYear(userID,year)
	if err != nil {
		return nil,nil, err
	}
	var outsideEvents []response.MyOutside
	for _, event := range outside {
		mappedEvent := response.MyOutside{
			EventID: event.EventID,
			EventName: event.EventName,
			Location: event.Location,
			StartDate: utility.FormatToThaiDate(event.StartDate),
			StartTime: utility.FormatToThaiTime(event.StartDate),
			WorkingHour: event.WorkingHour,
			SchoolYear: event.SchoolYear,
			Intendant: event.Intendant,
			File: event.File,
		}
		outsideEvents = append(outsideEvents, mappedEvent)
	}
	return insideEvents,outsideEvents,nil
}

// func (u *eventUsecase) SendEvent(userID uint) error{


// }

// Inside
func checkPermission(permission *response.EventResponse, user *entity.Student) bool {
	if permission == nil || user == nil {
		return false
	}
	permissionBranch := permission.AllowAllBranch
	permissionYear := permission.AllowAllYear

	if !permissionBranch && permission.BranchIDs != nil {
		for _, branch := range permission.BranchIDs {
			if user.BranchId == branch {
				permissionBranch = true
				break
			}
		}
	}

	if !permissionYear && permission.Years != nil {
		for _, year := range permission.Years {
			if user.Year == year {
				permissionYear = true
				break
			}
		}
	}

	return permissionBranch && permissionYear
}

func (u *eventUsecase) JoinEvent(eventID uint, claims map[string]interface{}) error {
	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return fmt.Errorf("invalid user_id in claims")
	}
	userID := uint(userIDFloat)

	student, err := u.userRepo.GetStudentByID(userID)
	if err != nil || student == nil {
		return fmt.Errorf("student not found")
	}

	event, err := u.GetEventByID(eventID)
	if err != nil {
		return fmt.Errorf("event not found")
	}

	if !event.Status {
		return fmt.Errorf("event not allowed")
	}

	if event.FreeSpace == 0 {
		return fmt.Errorf("the event is full")
	}

	if !checkPermission(event, student) {
		return fmt.Errorf("user is not allowed to join this event")
	}

	eventInside := &entity.EventInside{
		EventId:   eventID,
		User:      userID,
		Status:    false,
		Certifier: event.Creator.UserID,
	}

	err = u.eventRepo.JoinEvent(eventInside)
	if err != nil {
		return fmt.Errorf("failed to join event inside: %w", err)
	}
	return nil
}

func (u *eventUsecase) UnJoinEvent(eventID uint, claims map[string]interface{}) error {
	// ดึง user_id จาก claims และแปลงเป็น uint
	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return fmt.Errorf("invalid user_id in claims")
	}
	userID := uint(userIDFloat)

	// ตรวจสอบว่าผู้ใช้เป็นนักศึกษาหรือไม่
	_, err := u.userRepo.GetStudentByID(userID)
	if err != nil {
		return fmt.Errorf("student not found")
	}

	// ตรวจสอบว่า event มีอยู่จริงหรือไม่
	_, err = u.GetEventByID(eventID)
	if err != nil {
		return fmt.Errorf("event not found")
	}

	// ดำเนินการ Unjoin Event
	if err := u.eventRepo.UnJoinEvent(eventID, userID); err != nil {
		return fmt.Errorf("failed to unjoin event: %w", err)
	}

	return nil
}

func (u *eventUsecase) UploadFile(eventID uint, claims map[string]interface{}, file *multipart.FileHeader) error {
	//  ตรวจสอบประเภทไฟล์ให้เป็น PDF
	if file.Header.Get("Content-Type") != "application/pdf" ||
		!strings.HasSuffix(file.Filename, ".pdf") {
		return fmt.Errorf("only PDF files are allowed")
	}

	//  จำกัดขนาดไฟล์ที่ 10MB
	const maxFileSize = 10 * 1024 * 1024
	if file.Size > maxFileSize {
		return fmt.Errorf("file size exceeds the 10MB limit")
	}

	//  ดึง user_id จาก claims
	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return fmt.Errorf("invalid user_id in claims")
	}
	userID := uint(userIDFloat)

	// ค้นหาไฟล์ปัจจุบันในฐานข้อมูล
	currentFilePath, err := u.eventRepo.GetFilePath(eventID, userID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("failed to fetch current file path: %w", err)
	}

	// ถ้ามีไฟล์เดิมให้ลบ
	if currentFilePath != "" {
		if removeErr := os.Remove(currentFilePath); removeErr != nil {
			return fmt.Errorf("failed to remove old file: %v", removeErr)
		}
	}

	// บันทึกไฟล์ใหม่
	path, err := filesystem.SaveFile(file, userID)
	if err != nil {
		return fmt.Errorf("failed to save file: %w", err)
	}

	// อัปเดตฐานข้อมูล
	err = u.eventRepo.UploadFile(eventID, userID, path)
	if err != nil {
		os.Remove(path) // ลบไฟล์ใหม่หากอัปเดต DB ไม่สำเร็จ
		return fmt.Errorf("failed to update database: %w", err)
	}

	return nil
}

func (u *eventUsecase) GetFile(eventID uint, userID uint) (string, error) {
	filePath, err := u.eventRepo.GetFilePath(eventID, userID)
	if err != nil {
		return "", err
	}
	return filePath, nil
}

func (u *eventUsecase) MyChecklist(eventID uint, claims map[string]interface{}) ([]response.MyChecklist, error) {
	//  ดึง user_id จาก claims
	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return nil, fmt.Errorf("invalid user_id in claims")
	}
	userID := uint(userIDFloat)

	checklist, err := u.eventRepo.MyChecklist(userID, eventID)
	if err != nil {
		return nil, err
	}
	var res []response.MyChecklist
	for _, inside := range checklist {
		mappedEvent := response.MyChecklist{
			EventID:   inside.EventId,
			UserID:    inside.User,
			TitleName: inside.Student.TitleName,
			FirstName: inside.Student.FirstName,
			LastName:  inside.Student.LastName,
			Code:      inside.Student.Code,
			Certifier: inside.Certifier,
			Status:    inside.Status,
			Comment:   inside.Comment,
			File:      inside.File,
		}
		res = append(res, mappedEvent)
	}
	return res, nil
}

func (u *eventUsecase) UpdateEventStatusAndComment(eventID uint, userID uint, status bool, comment string) error {
	return u.eventRepo.UpdateEventStatusAndComment(eventID, userID, status, comment)
}








// func (u *eventUsecase) UpdateEventByID(eventID uint, claims map[string]interface{}, req request.EventRequest) error {
// 	event, err := u.eventRepo.GetEventByID(eventID)
// 	if err != nil {
// 		return fmt.Errorf("event not found")
// 	}

// 	userIDFloat, ok := claims["user_id"].(float64)
// 	if !ok {
// 		return fmt.Errorf("invalid user_id in claims")
// 	}
// 	userID := uint(userIDFloat)
// 	if event.Creator != userID {
// 		return fmt.Errorf("you do not have permission to edit this event")
// 	}

// 	if event.FreeSpace == 0 {
// 		event.Status = false
// 	}

// 	startDate, err := utility.ParseStartDate(req.StartDate)
// 	if err != nil {
// 		return fmt.Errorf("invalid start date format")
// 	}

// 	event.EventName = req.EventName
// 	event.StartDate = startDate
// 	event.WorkingHour = req.WorkingHour
// 	event.Location = req.Location
// 	event.Detail = req.Detail

// 	userIDs, err := u.eventRepo.GroupByEvent(eventID)
// 	if err != nil {
// 		return fmt.Errorf("failed to get users for event: %w", err)
// 	}
// 	for _, uid := range userIDs {
// 		news := entity.News{
// 			Title:   "กิจกรรมมีการแก้ไขรายละเอียด",
// 			UserID:  uid,
// 			Message: fmt.Sprintf("กิจกรรม'%s' ที่คุณเข้าร่วมมีการแก้ไขรายละเอียด.", event.EventName),
// 		}
// 		if err := u.eventRepo.NewsForUser(&news); err != nil {
// 			return fmt.Errorf("failed to send news to user %d: %w", uid, err)
// 		}
// 	}

// 	return u.eventRepo.UpdateEventByID(event)
// }


// func (u *eventUsecase) DeleteEventByID(eventID uint, claims map[string]interface{}) error {
// 	event, err := u.eventRepo.GetEventByID(eventID)
// 	if err != nil {
// 		return fmt.Errorf("event not found")
// 	}

// 	userIDFloat, ok := claims["user_id"].(float64)
// 	if !ok {
// 		return fmt.Errorf("invalid user_id in claims")
// 	}
// 	userID := uint(userIDFloat)

// 	if event.Creator != userID {
// 		return fmt.Errorf("you do not have permission to delete this event")
// 	}

// 	userIDs, err := u.eventRepo.GroupByEvent(eventID)
// 	if err != nil {
// 		return fmt.Errorf("failed to get users for event: %w", err)
// 	}

// 	for _, uid := range userIDs {
// 		news := entity.News{
// 			Title:   "กิจกรรมถูกลบ",
// 			UserID:  uid,
// 			Message: fmt.Sprintf("กิจกรรม'%s' ที่คุณเข้าร่วมถูกลบแล้ว.", event.EventName),
// 		}
// 		if err := u.eventRepo.NewsForUser(&news); err != nil {
// 			return fmt.Errorf("failed to send news to user %d: %w", uid, err)
// 		}
// 	}

// 	return u.eventRepo.DeleteEventByID(event.EventID)
// }

// func (u *eventUsecase) CreateEvent(req *request.EventRequest, claims map[string]interface{}) error {
// 	userIDFloat, ok := claims["user_id"].(float64)
// 	if !ok {
// 		return fmt.Errorf("invalid user_id in claims")
// 	}
// 	userID := uint(userIDFloat)

// 	return u.eventRepo.CreateEventWithTransaction(req, userID)
// }
