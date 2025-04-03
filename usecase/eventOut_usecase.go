package usecase

import (
	"database/sql"
	"errors"
	"fmt"
	"go-clean-arch/pkg/utility"
	"go-clean-arch/pkg/utility/filesystem"
	"go-clean-arch/structure/entity"
	"go-clean-arch/structure/request"
	"go-clean-arch/structure/response"
	"mime/multipart"
	"os"
	"strings"
)

// Outside
func (u *eventUsecase) CreateEventOutside(req request.OutsideRequest,claims map[string]interface{}) error{
	startDate, err := utility.ParseStartDate(req.StartDate)
	if err != nil {
		return err
	}

	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return fmt.Errorf("invalid user_id in claims")
	}
	userID := uint(userIDFloat)

	outside := entity.EventOutside{
		User: userID,
		EventName: req.EventName,
		SchoolYear: req.SchoolYear,
		StartDate: startDate,
		Intendant: req.Intendant,
		WorkingHour: req.WorkingHour,
		Location: req.Location,
	}
	return u.eventRepo.CreateEventOutside(outside)
}

func (u *eventUsecase) GetEventOutsideByID(eventID uint) (*response.OutsideResponse, error) {
	outside, err := u.eventRepo.GetEventOutsideByID(eventID)
	if err != nil {
		return nil, err
	}
	outsideRes := response.OutsideResponse{
		EventID:     outside.EventID,
		EventName:   outside.EventName,
		Location:    outside.Location,
		SchoolYear:  outside.SchoolYear,
		StartDate:   outside.StartDate,
		WorkingHour: outside.WorkingHour,
		Intendant:   outside.Intendant,
		Student: response.StudentResponse{
			UserID:      outside.Student.UserID,
			TitleName:   outside.Student.TitleName,
			FirstName:   outside.Student.FirstName,
			LastName:    outside.Student.LastName,
			Phone:       outside.Student.Phone,
			Code:        outside.Student.Code,
			BranchID:    outside.Student.BranchId,
			BranchName:  outside.Student.Branch.BranchName,
			FacultyID:   outside.Student.Branch.Faculty.FacultyID,
			FacultyName: outside.Student.Branch.Faculty.FacultyName,
		},
	}
	return &outsideRes, nil
}

func (u *eventUsecase) DeleteEventOutsideByID(eventID uint) error{
	if err := u.eventRepo.DeleteEventOutsideByID(eventID); err != nil {
		return fmt.Errorf("failed to deleted eventoutside: %w", err)
	}
	return nil
}


func (u *eventUsecase) CreateFile(eventID uint) ([]byte, string, error) {
	data, err := u.GetEventOutsideByID(eventID)
	if err != nil {
		return nil, "", fmt.Errorf("data not found: %v", err)
	}
	pdfBytes, fileName, err := filesystem.CreatePDF(*data)
	if err != nil {
		return nil, " ", fmt.Errorf("error creating PDF: %v", err)
	}
	return pdfBytes, fileName, nil

}


func (u *eventUsecase) UploadFileOutside(eventID uint, claims map[string]interface{}, file *multipart.FileHeader) error {
	// ตรวจสอบว่า eventID มีอยู่ใน EventOutside หรือไม่
	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return fmt.Errorf("invalid user_id in claims")
	}
	userID := uint(userIDFloat)

	exists, err := u.eventRepo.EventOutsideExists(eventID, userID)
	if err != nil {
		return fmt.Errorf("failed to check event existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("event outside does not exist for this user")
	}

	// ตรวจสอบประเภทไฟล์ให้เป็น PDF
	if file.Header.Get("Content-Type") != "application/pdf" ||
		!strings.HasSuffix(file.Filename, ".pdf") {
		return fmt.Errorf("only PDF files are allowed")
	}

	// จำกัดขนาดไฟล์ที่ 10MB
	const maxFileSize = 10 * 1024 * 1024
	if file.Size > maxFileSize {
		return fmt.Errorf("file size exceeds the 10MB limit")
	}

	// ค้นหาไฟล์ปัจจุบันในฐานข้อมูล
	currentFilePath, err := u.eventRepo.GetFilePathOutside(eventID, userID)
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
	err = u.eventRepo.UploadFileOutside(eventID, userID, path)
	if err != nil {
		os.Remove(path) // ลบไฟล์ใหม่หากอัปเดต DB ไม่สำเร็จ
		return fmt.Errorf("failed to update database: %w", err)
	}

	return nil
}


func (u *eventUsecase) GetFileOutside(eventID uint ,userID uint)(string,error){
	filePath, err := u.eventRepo.GetFilePathOutside(eventID, userID)
	if err != nil {
		return "", err
	}
	return filePath, nil
}