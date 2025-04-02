package database

import (
	"fmt"
	"go-clean-arch/config"
	"go-clean-arch/structure/entity"
	"go-clean-arch/pkg/hash"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type mysqlDatabase struct {
	DB *gorm.DB
}

func NewMySQLDatabase(cfg *config.Config) (Database, error) {
	db, err := gorm.Open(mysql.Open(cfg.DSN), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	return &mysqlDatabase{DB: db}, nil
}

func SetupDatabase(cfg *config.Config) Database {

	db, err := NewMySQLDatabase(cfg)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	err = db.AutoMigrate()
	if err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	addTriggerIfNotExists(db, "before_insert_students", `
        CREATE TRIGGER before_insert_students
        BEFORE INSERT ON students
        FOR EACH ROW
        BEGIN
            IF EXISTS (SELECT 1 FROM teachers WHERE user_id = NEW.user_id) THEN
                SIGNAL SQLSTATE '45000'
                SET MESSAGE_TEXT = 'User ID already exists in teachers';
            END IF;
        END;
    `)

	addTriggerIfNotExists(db, "before_insert_teachers", `
        CREATE TRIGGER before_insert_teachers
        BEFORE INSERT ON teachers
        FOR EACH ROW
        BEGIN
            IF EXISTS (SELECT 1 FROM students WHERE user_id = NEW.user_id) THEN
                SIGNAL SQLSTATE '45000'
                SET MESSAGE_TEXT = 'User ID already exists in students';
            END IF;
        END;
    `)
	password, err := hash.HashPassword(cfg.Admin.Password)
	if err != nil {
		log.Fatalf("failed to hash password: %v", err)
	}

	var admin entity.User
	result := db.GetDB().Where("email = ?", cfg.Admin.Email).First(&admin)
	if result.Error == nil {
		log.Println("Admin user already exists.")
	} else if result.Error == gorm.ErrRecordNotFound {
		user := entity.User{
			Email:    cfg.Admin.Email,
			Password: password,
			Role:     "superadmin",
		}
		createResult := db.GetDB().Create(&user)
		if createResult.Error != nil {
			log.Fatalf("failed to create admin user: %v", createResult.Error)
		} else {
			log.Println("Admin user created successfully!")
		}
	} else {
		log.Fatalf("failed to check admin user existence: %v", result.Error)
	}

	log.Println("Database connected, migrated, and triggers added successfully!")
	return db
}

func addTriggerIfNotExists(db Database, triggerName, triggerSQL string) {
	var count int
	// ตรวจสอบว่า trigger มีอยู่แล้วหรือไม่
	query := fmt.Sprintf("SHOW TRIGGERS LIKE '%s'", triggerName)
	err := db.GetDB().Raw(query).Scan(&count).Error
	if err != nil {
		log.Printf("Error checking trigger existence: %v", err)
		return
	}

	if count == 0 {
		// ถ้าไม่มี trigger นี้ก็สร้างใหม่
		if err := db.GetDB().Exec(triggerSQL).Error; err != nil {
			// เพิ่มการตรวจสอบเพื่อไม่แสดงข้อความหาก trigger มีอยู่แล้ว
			if err.Error() != "Error 1359 (HY000): Trigger already exists" {
				log.Printf("Failed to create trigger %s: %v", triggerName, err)
			}
		} else {
			log.Printf("Trigger %s created successfully!", triggerName)
		}
	} else {
		log.Printf("Trigger %s already exists, skipping creation.", triggerName)
	}
}

//

func (m *mysqlDatabase) GetDB() *gorm.DB {
	return m.DB
}

func (m *mysqlDatabase) AutoMigrate() error {

	// Migrate other tables as needed
	if err := m.DB.AutoMigrate(&entity.User{}); err != nil {
		return fmt.Errorf("failed to migrate User: %w", err)
	}

	if err := m.DB.AutoMigrate(&entity.Teacher{}); err != nil {
		return fmt.Errorf("failed to migrate Teacher: %w", err)
	}
	// Migrate Faculty table first
	if err := m.DB.AutoMigrate(&entity.Faculty{}); err != nil {
		return fmt.Errorf("failed to migrate Faculty: %w", err)
	}

	// Migrate Branch table after Faculty
	if err := m.DB.AutoMigrate(&entity.Branch{}); err != nil {
		return fmt.Errorf("failed to migrate Branch: %w", err)
	}
	if err := m.DB.AutoMigrate(&entity.Student{}); err != nil {
		return fmt.Errorf("failed to migrate Student: %w", err)
	}
	if err := m.DB.AutoMigrate(&entity.Event{}); err != nil {
		return fmt.Errorf("failed to migrate Event: %w", err)
	}
	if err := m.DB.AutoMigrate(&entity.EventInside{}); err != nil {
		return fmt.Errorf("failed to migrate EventInside: %w", err)
	}
	if err := m.DB.AutoMigrate(&entity.EventOutside{}); err != nil {
		return fmt.Errorf("failed to migrate EventInside: %w", err)
	}

	if err := m.DB.AutoMigrate(&entity.Done{}); err != nil {
		return fmt.Errorf("failed to migrate EventInside: %w", err)
	}
	if err := m.DB.AutoMigrate(&entity.News{}); err != nil {
		return fmt.Errorf("failed to migrate EventInside: %w", err)
	}

	return nil
}
