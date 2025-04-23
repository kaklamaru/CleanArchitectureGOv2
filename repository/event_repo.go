package repository

import (
	"errors"
	"fmt"
	"go-clean-arch/pkg/utility"
	"go-clean-arch/structure/entity"
	"go-clean-arch/structure/request"
	"go-clean-arch/structure/response"
	"os"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type EventRepository interface {
	CreateEvent(event *entity.Event) error
	NewsForUser(news *entity.News) error
	GetAllEvent() ([]entity.Event, error)
	CountEventInside(eventID uint) (uint, error)
	GetEventByID(id uint) (*entity.Event, error)
	ToggleEventStatus(eventID uint) (bool, error)
	UpdateEventWithTransaction(eventID, userID uint, req request.EventRequest) error
	DeleteEventWithTransaction(eventID, userID uint) error

	GroupByEvent(eventID uint) ([]uint, error)
	JoinEvent(eventInside *entity.EventInside) error
	UnJoinEvent(eventID uint, userID uint) error
	GetFilePath(eventID uint, userID uint) (string, error)
	UploadFile(eventID uint, userID uint, filePath string) error
	MyEvent(userID uint) ([]entity.Event, error)
	AllAllowedEvent() ([]entity.Event, error)
	AllCurrentEvent() ([]entity.Event, error)
	MyChecklist(userID uint, eventID uint) ([]entity.EventInside, error)
	UpdateEventStatusAndComment(eventID uint, userID uint, status bool, comment string) error
	AllEventInsideThisYear(userID uint, year uint) ([]entity.EventInside, error)
	AllEventInside(userID uint) ([]entity.EventInside, error)

	CreateEventOutside(outside entity.EventOutside) error
	DeleteEventOutsideByID(eventID uint) error
	GetEventOutsideByID(id uint) (*entity.EventOutside, error)
	GetFilePathOutside(eventID uint, userID uint) (string, error)
	UploadFileOutside(eventID uint, userID uint, filePath string) error
	AllEventOutsideThisYear(userID uint, year uint) ([]entity.EventOutside, error)
	EventOutsideExists(eventID uint, userID uint) (bool, error)

	DashboardData(year uint) (*response.DashboardSummary, error)
	GetUpcomingEventsWithin7Days() ([]entity.Event, error)
	GetMonthlyEventStats(year uint) (eventCounts, outsideCounts []response.MonthlyEventCount, err error) 
}

// type InsideEvent struct{

// }

type eventRepository struct {
	db *gorm.DB
}

func NewEventRepository(db *gorm.DB) EventRepository {
	return &eventRepository{db: db}
}

func (r *eventRepository) UpdateEventWithTransaction(eventID, userID uint, req request.EventRequest) error {
	// เริ่ม Transaction
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// ดึงข้อมูลกิจกรรม
	var event entity.Event
	if err := tx.Where("event_id = ?", eventID).First(&event).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("event not found: %w", err)
	}

	// ตรวจสอบสิทธิ์การแก้ไข
	if event.Creator != userID {
		tx.Rollback()
		return fmt.Errorf("you do not have permission to edit this event")
	}

	// อัปเดตข้อมูล
	startDate, err := utility.ParseStartDate(req.StartDate)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("invalid start date format")
	}

	event.EventName = req.EventName
	event.StartDate = startDate
	event.WorkingHour = req.WorkingHour
	event.Location = req.Location
	event.Detail = req.Detail

	if err := tx.Save(&event).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to update event: %w", err)
	}

	// แจ้งเตือนผู้ใช้ที่เข้าร่วม
	var userIDs []uint
	if err := tx.Model(&entity.EventInside{}).Where("event_id = ?", eventID).Pluck("user", &userIDs).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to get users for event: %w", err)
	}

	for _, uid := range userIDs {
		news := entity.News{
			Title:   "กิจกรรมมีการแก้ไขรายละเอียด",
			UserID:  uid,
			Message: fmt.Sprintf("กิจกรรม '%s' ที่คุณเข้าร่วมมีการแก้ไขรายละเอียด.", event.EventName),
		}
		if err := tx.Create(&news).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to send news to user %d: %w", uid, err)
		}
	}

	return tx.Commit().Error
}

func (r *eventRepository) CreateEvent(event *entity.Event) error {
	// ตรวจสอบว่ามีกิจกรรมที่ชื่อเดียวกัน วันที่เดียวกัน และ creator เดียวกันอยู่แล้วหรือไม่
	var existingEvent entity.Event
	err := r.db.Where("event_name = ? AND start_date = ? AND creator = ?", event.EventName, event.StartDate, event.Creator).
		First(&existingEvent).Error

	if err == nil {
		// พบกิจกรรมที่ซ้ำกัน
		return fmt.Errorf("event with same name and start date already exists for this creator")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	// ถ้าไม่เจอซ้ำ สร้าง event ได้
	return r.db.Create(event).Error
}

func (r *eventRepository) NewsForUser(news *entity.News) error {
	if err := r.db.Create(news).Error; err != nil {
		return err
	}
	return nil
}

func (r *eventRepository) GetAllEvent() ([]entity.Event, error) {
	var events []entity.Event

	// คำนวณช่วงเวลา -30 วัน ถึง +30 วัน
	now := time.Now()
	start := now.AddDate(0, 0, -30)
	end := now.AddDate(0, 0, 30)

	// Query เฉพาะกิจกรรมที่ start_date อยู่ในช่วงที่กำหนด
	if err := r.db.Preload("Teacher").
		Where("start_date BETWEEN ? AND ?", start, end).
		Find(&events).Error; err != nil {
		return nil, err
	}

	return events, nil
}

func (r *eventRepository) CountEventInside(eventID uint) (uint, error) {
	var count int64
	if err := r.db.Model(&entity.EventInside{}).Where("event_id = ?", eventID).Count(&count).Error; err != nil {
		return 0, err
	}
	return uint(count), nil
}

func (r *eventRepository) GetEventByID(id uint) (*entity.Event, error) {
	var event entity.Event
	if err := r.db.Preload("Teacher").First(&event, "event_id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("event with ID %d not found", id)
		}
		return nil, fmt.Errorf("failed to retrieve event: %w", err)
	}

	return &event, nil
}

func (r *eventRepository) ToggleEventStatus(eventID uint) (bool, error) {
	// อัปเดตค่า status โดยกลับค่าด้วย NOT status
	if err := r.db.Model(&entity.Event{}).
		Where("event_id = ?", eventID).
		Update("status", gorm.Expr("NOT status")).Error; err != nil {
		return false, err
	}

	var updatedEvent entity.Event
	if err := r.db.Select("status").Where("event_id = ?", eventID).First(&updatedEvent).Error; err != nil {
		return false, err
	}

	return updatedEvent.Status, nil
}

func (r *eventRepository) DeleteEventWithTransaction(eventID, userID uint) error {
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	var event entity.Event
	if err := tx.Where("event_id = ?", eventID).First(&event).Error; err != nil {
		return fmt.Errorf("event not found: %w", err)
	}
	if event.Status {
		return fmt.Errorf("cannot delete event because status is true")
	}

	// ตรวจสอบสิทธิ์การลบ
	if event.Creator != userID {
		tx.Rollback()
		return fmt.Errorf("you do not have permission to delete this event")
	}

	// ดึงรายชื่อผู้ใช้ที่เกี่ยวข้อง
	var userIDs []uint
	if err := tx.Model(&entity.EventInside{}).Where("event_id = ?", eventID).Pluck("user", &userIDs).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to get users for event: %w", err)
	}

	// ลบกิจกรรม
	if err := tx.Where("event_id = ?", eventID).Delete(&entity.Event{}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to delete event: %w", err)
	}

	// สร้างข่าวสารแจ้งผู้ใช้
	for _, uid := range userIDs {
		news := entity.News{
			Title:   "กิจกรรมถูกลบ",
			UserID:  uid,
			Message: fmt.Sprintf("กิจกรรม '%s' ที่คุณเข้าร่วมถูกลบแล้ว.", event.EventName),
		}
		if err := tx.Create(&news).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to send news to user %d: %w", uid, err)
		}
	}

	// Commit Transaction
	return tx.Commit().Error
}

func (r *eventRepository) GroupByEvent(eventID uint) ([]uint, error) {
	var eventInsides []entity.EventInside
	err := r.db.Where("event_id = ?", eventID).Find(&eventInsides).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find users for event ID %d: %w", eventID, err)
	}
	var userIDs []uint
	for _, ei := range eventInsides {
		userIDs = append(userIDs, ei.User)
	}
	return userIDs, nil
}

func (r *eventRepository) MyEvent(userID uint) ([]entity.Event, error) {
	var events []entity.Event
	if err := r.db.Preload("Teacher").Where("creator = ?", userID).Find(&events).Error; err != nil {
		return nil, err
	}
	return events, nil
}

func (r *eventRepository) AllAllowedEvent() ([]entity.Event, error) {
	var events []entity.Event
	if err := r.db.Preload("Teacher").Where("status = true").Find(&events).Error; err != nil {
		return nil, err
	}
	return events, nil
}

func (r *eventRepository) AllCurrentEvent() ([]entity.Event, error) {
	var events []entity.Event
	today := time.Now()
	futureDate := today.AddDate(0, 1, 0)
	if err := r.db.Preload("Teacher").Where("start_date BETWEEN ? AND ?", today, futureDate).Find(&events).Error; err != nil {
		return nil, err
	}
	return events, nil
}

// inside event
func (r *eventRepository) JoinEvent(eventInside *entity.EventInside) error {
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// ตั้งค่า lock timeout เพื่อป้องกันการล็อกที่ยาวนาน
	if err := tx.Exec("SET lock_timeout TO '5s'").Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to set lock timeout: %w", err)
	}

	// ทำการล็อกแถว Event ที่ต้องการให้เป็น exclusive lock (UPDATE)
	var event entity.Event
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("event_id = ?", eventInside.EventId).
		First(&event).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to fetch event: %w", err)
	}

	// ตรวจสอบว่ามีที่ว่างสำหรับผู้เข้าร่วมกิจกรรม
	if event.FreeSpace <= 0 {
		tx.Rollback()
		return fmt.Errorf("no free space available for event")
	}

	// ลดจำนวนที่ว่างลง 1
	event.FreeSpace -= 1
	if event.FreeSpace == 0 {
		event.Status = false
	}
	if err := tx.Save(&event).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to update event free space: %w", err)
	}

	// สร้าง event_inside record
	if err := tx.Create(eventInside).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to create event inside record: %w", err)
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}

func (r *eventRepository) UnJoinEvent(eventID uint, userID uint) error {
	// เริ่มต้น Transaction
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback() // ถ้า panic ให้ rollback
		}
	}()

	// ตั้งค่า lock timeout เพื่อป้องกันการรอค้างนานเกินไป
	if err := tx.Exec("SET lock_timeout TO '5s'").Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to set lock timeout: %w", err)
	}

	// ดึงข้อมูล Event โดยล็อกแถวแบบ UPDATE
	var event entity.Event
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("event_id = ?", eventID).
		First(&event).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to fetch event: %w", err)
	}

	// ตรวจสอบว่าผู้ใช้เคยเข้าร่วมจริงหรือไม่
	var eventInside entity.EventInside
	if err := tx.Where(`event_id = ? AND "user" = ?`, eventID, userID).
		First(&eventInside).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("user has not joined this event")
	}
	if eventInside.File != "" {
		if err := os.Remove(eventInside.File); err != nil && !os.IsNotExist(err) {
			tx.Rollback()
			return fmt.Errorf("failed to delete file: %w", err)
		}
	}

	// ลบข้อมูลการเข้าร่วมของผู้ใช้
	if err := tx.Delete(&eventInside).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to remove user from event: %w", err)
	}

	// เพิ่มจำนวน FreeSpace
	event.FreeSpace += 1
	if event.FreeSpace > 0 {
		event.Status = true
	}
	if err := tx.Save(&event).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to update event free space: %w", err)
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *eventRepository) GetFilePath(eventID uint, userID uint) (string, error) {
	var filePath string

	err := r.db.Model(&entity.EventInside{}).
		Where(`event_id = ? AND "user" = ?`, eventID, userID).
		Pluck("file", &filePath).Error
	if err != nil {
		return "", fmt.Errorf("failed to retrieve file path: %w", err)
	}

	if filePath == "" {
		return "", nil
	}

	return filePath, nil
}

func (r *eventRepository) UploadFile(eventID uint, userID uint, filePath string) error {
	// 1. อัปเดต File และล้าง Comment
	if err := r.db.Model(&entity.EventInside{}).
		Where(`event_id = ? AND "user" = ?`, eventID, userID).
		Updates(map[string]interface{}{
			"file":    filePath,
			"comment": "",
		}).Error; err != nil {
		return err
	}

	// 2. ดึง EventInside เพื่อใช้ Certifier
	var inside entity.EventInside
	if err := r.db.
		Where(`event_id = ? AND "user" = ?`, eventID, userID).
		First(&inside).Error; err != nil {
		return err
	}

	// 3. ดึงข้อมูล Event เพื่อชื่อกิจกรรม
	var event entity.Event
	if err := r.db.
		Where("event_id = ?", eventID).
		First(&event).Error; err != nil {
		return err
	}

	// 4. สร้าง News โดยใช้ Certifier
	news := entity.News{
		Title:   "กิจกรรมที่ต้องตรวจสอบ",
		UserID:  inside.Certifier,
		Message: fmt.Sprintf("กิจกรรม '%s' ที่คุณต้องตรวจสอบ.", event.EventName),
	}

	if err := r.db.Create(&news).Error; err != nil {
		return fmt.Errorf("failed to create news for certifier %d: %w", inside.Certifier, err)
	}

	return nil
}

func (r *eventRepository) MyChecklist(userID uint, eventID uint) ([]entity.EventInside, error) {
	var checklist []entity.EventInside
	if err := r.db.Preload("Event").Preload("Student.Branch.Faculty").Where("event_id = ? ", eventID).Find(&checklist).Error; err != nil {
		return nil, err
	}
	return checklist, nil
}

func (r *eventRepository) UpdateEventStatusAndComment(eventID uint, userID uint, status bool, comment string) error {
	updates := map[string]interface{}{
		"status":  status,
		"comment": comment,
	}
	if err := r.db.Model(&entity.EventInside{}).
		Where(`event_id = ? AND "user" = ?`, eventID, userID).
		Updates(updates).Error; err != nil {
		return fmt.Errorf("failed to update event: %w", err)
	}

	var event entity.Event
	if err := r.db.
		Where("event_id = ?", eventID).
		First(&event).Error; err != nil {
		return err
	}

	text := "ไม่ผ่าน"
	if status {
		text = "ผ่าน"
	}

	message := fmt.Sprintf("เอกสารกิจกรรม '%s' '%s'", event.EventName, text)
	if !status && comment != "" {
		message += fmt.Sprintf(" เนื่องจาก: %s", comment)
	}

	news := entity.News{
		Title:   fmt.Sprintf("เอกสารกิจกรรม '%s'", event.EventName),
		UserID:  userID,
		Message: message,
	}

	if err := r.db.Create(&news).Error; err != nil {
		return fmt.Errorf("failed to create news: %w", err)
	}

	return nil
}

func (r *eventRepository) AllEventInsideThisYear(userID uint, year uint) ([]entity.EventInside, error) {
	var eventInsides []entity.EventInside
	// year := 2568
	err := r.db.Preload("Event").Joins("JOIN events ON events.event_id = event_insides.event_id").
		Where("event_insides.user = ?", userID).
		Where("events.school_year = ?", year).
		Find(&eventInsides).Error
	if err != nil {
		fmt.Println("Error fetching data:", err)
	}
	return eventInsides, nil
}

func (r *eventRepository) AllEventInside(userID uint) ([]entity.EventInside, error) {
	var eventInsides []entity.EventInside
	err := r.db.Preload("Event").Joins("JOIN events ON events.event_id = event_insides.event_id").
		Where("event_insides.user = ?", userID).
		Find(&eventInsides).Error
	if err != nil {
		fmt.Println("Error fetching data:", err)
	}
	return eventInsides, nil
}

// outside event
func (r *eventRepository) CreateEventOutside(outside entity.EventOutside) error {
	if err := r.db.Create(&outside).Error; err != nil {
		return err
	}
	return nil
}

func (r *eventRepository) DeleteEventOutsideByID(eventID uint) error {
	tx := r.db.Begin() // เริ่ม Transaction

	var event entity.EventOutside

	// ดึงข้อมูล event
	if err := tx.Where("event_id = ?", eventID).First(&event).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("event not found: %w", err)
	}

	if event.File != "" {
		if err := os.Remove(event.File); err != nil && !os.IsNotExist(err) {
			tx.Rollback()
			return fmt.Errorf("failed to remove file: %w", err)
		}
	}

	// ลบ Event จากฐานข้อมูล
	if err := tx.Where("event_id = ?", eventID).Delete(&entity.EventOutside{}).Error; err != nil {
		tx.Rollback() // ยกเลิก Transaction ถ้าลบข้อมูลไม่สำเร็จ
		return fmt.Errorf("failed to delete event: %w", err)
	}

	// Commit Transaction ถ้าทุกอย่างสำเร็จ
	return tx.Commit().Error
}

func (r *eventRepository) GetEventOutsideByID(id uint) (*entity.EventOutside, error) {
	var outside entity.EventOutside
	if err := r.db.Preload("Student.Branch.Faculty").First(&outside, "event_id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("event with ID %d not found", id)
		}
		return nil, fmt.Errorf("failed to retrieve event: %w", err)
	}
	return &outside, nil
}

func (r *eventRepository) AllEventOutsideThisYear(userID uint, year uint) ([]entity.EventOutside, error) {
	var eventOutside []entity.EventOutside
	if err := r.db.Where(`"user" = ? AND school_year = ?`, userID, year).Find(&eventOutside).Error; err != nil {
		return nil, err
	}
	return eventOutside, nil
}

func (r *eventRepository) GetFilePathOutside(eventID uint, userID uint) (string, error) {
	var filePath string

	err := r.db.Model(&entity.EventOutside{}).
		Where(`event_id = ? AND "user" = ?`, eventID, userID).
		Pluck("file", &filePath).Error
	if err != nil {
		return "", fmt.Errorf("failed to retrieve file path: %w", err)
	}

	if filePath == "" {
		return "", nil
	}

	return filePath, nil
}

func (r *eventRepository) UploadFileOutside(eventID uint, userID uint, filePath string) error {
	if err := r.db.Model(&entity.EventOutside{}).
		Where(`event_id = ? AND "user" = ?`, eventID, userID).
		Update("file", filePath).Error; err != nil {
		return err
	}
	return nil
}
func (r *eventRepository) EventOutsideExists(eventID uint, userID uint) (bool, error) {
	var count int64
	err := r.db.Model(&entity.EventOutside{}).
		Where(`event_id = ? AND "user" = ?`, eventID, userID).
		Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("failed to check event outside existence: %w", err)
	}
	return count > 0, nil
}

func (r *eventRepository) DashboardData(year uint) (*response.DashboardSummary, error) {
	var data response.DashboardSummary

	var countStd, countDone, countEvent, countEventAllow int64

	if err := r.db.Model(&entity.Student{}).Count(&countStd).Error; err != nil {
		return nil, err
	}
	if err := r.db.Model(&entity.Done{}).Where("year = ? AND status = ?", year, true).Count(&countDone).Error; err != nil {
		return nil, err
	}
	if err := r.db.Model(&entity.Event{}).Where("school_year = ?", year).Count(&countEvent).Error; err != nil {
		return nil, err
	}
	if err := r.db.Model(&entity.Event{}).Where("school_year = ? AND status = ?", year, true).Count(&countEventAllow).Error; err != nil {
		return nil, err
	}

	// แปลงค่าจาก int64 -> uint
	data.CountStd = uint(countStd)
	data.CountDone = uint(countDone)
	data.CountEvent = uint(countEvent)
	data.CountEventAllow = uint(countEventAllow)

	return &data, nil
}

func (r *eventRepository) GetUpcomingEventsWithin7Days() ([]entity.Event, error) {
	var events []entity.Event

	today := time.Now()
	sevenDaysLater := today.AddDate(0, 0, 6)

	if err := r.db.Where("start_date BETWEEN ? AND ?", today, sevenDaysLater).Find(&events).Error; err != nil {
		return nil, err
	}

	return events, nil
}

func (r *eventRepository) GetMonthlyEventStats(year uint) (eventCounts, outsideCounts []response.MonthlyEventCount, err error) {
	// นับจาก Event
	if err := r.db.
		Model(&entity.Event{}).
		Select("EXTRACT(MONTH FROM start_date) AS month, COUNT(*) AS count").
		Where("school_year = ?", year).
		Group("month").
		Order("month").
		Scan(&eventCounts).Error; err != nil {
		return nil, nil, err
	}

	// นับจาก EventOutside
	if err := r.db.
		Model(&entity.EventOutside{}).
		Select("EXTRACT(MONTH FROM start_date) AS month, COUNT(*) AS count").
		Where("school_year = ?", year).
		Group("month").
		Order("month").
		Scan(&outsideCounts).Error; err != nil {
		return nil, nil, err
	}

	return eventCounts, outsideCounts, nil
}