package database

import (
	"fmt"
	"go-clean-arch/config"
	"go-clean-arch/pkg/hash"
	"go-clean-arch/structure/entity"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type postgreDatabase struct {
	DB *gorm.DB
}

func NewPostgreDatabase(cfg *config.Config) (Database, error) {
	db, err := gorm.Open(postgres.Open(cfg.DSN), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to PostgreSQL: %w", err)
	}
	return &postgreDatabase{DB: db}, nil
}

func SetupDatabase(cfg *config.Config) Database {
	db, err := NewPostgreDatabase(cfg)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	err = db.AutoMigrate()
	if err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	addTriggerIfNotExists(db, "before_insert_students", `
		CREATE OR REPLACE FUNCTION check_student_user_id()
		RETURNS trigger AS $$
		BEGIN
			IF EXISTS (SELECT 1 FROM teachers WHERE user_id = NEW.user_id) THEN
				RAISE EXCEPTION 'User ID already exists in teachers';
			END IF;
			RETURN NEW;
		END;
		$$ LANGUAGE plpgsql;

		CREATE TRIGGER before_insert_students
		BEFORE INSERT ON students
		FOR EACH ROW
		EXECUTE FUNCTION check_student_user_id();
	`)

	addTriggerIfNotExists(db, "before_insert_teachers", `
		CREATE OR REPLACE FUNCTION check_teacher_user_id()
		RETURNS trigger AS $$
		BEGIN
			IF EXISTS (SELECT 1 FROM students WHERE user_id = NEW.user_id) THEN
				RAISE EXCEPTION 'User ID already exists in students';
			END IF;
			RETURN NEW;
		END;
		$$ LANGUAGE plpgsql;

		CREATE TRIGGER before_insert_teachers
		BEFORE INSERT ON teachers
		FOR EACH ROW
		EXECUTE FUNCTION check_teacher_user_id();
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

	log.Println("PostgreSQL connected, migrated, and triggers added successfully!")
	return db
}

func addTriggerIfNotExists(db Database, triggerName, triggerSQL string) {
	var exists bool
	err := db.GetDB().Raw(`
		SELECT EXISTS (
			SELECT 1
			FROM pg_trigger
			WHERE tgname = ?
		)
	`, triggerName).Scan(&exists).Error

	if err != nil {
		log.Printf("Error checking trigger existence: %v", err)
		return
	}

	if !exists {
		if err := db.GetDB().Exec(triggerSQL).Error; err != nil {
			log.Printf("Failed to create trigger %s: %v", triggerName, err)
		} else {
			log.Printf("Trigger %s created successfully!", triggerName)
		}
	} else {
		log.Printf("Trigger %s already exists, skipping creation.", triggerName)
	}
}

func (p *postgreDatabase) GetDB() *gorm.DB {
	return p.DB
}
// func (p *postgreDatabase) AutoMigrate() error {
// 	// รอบแรก: Migrate ตารางที่ไม่พึ่ง FK
// 	if err := p.DB.AutoMigrate(
// 		&entity.User{},
// 		&entity.Teacher{},
// 		&entity.Student{},
// 		&entity.Faculty{},
// 		&entity.Branch{},
// 	); err != nil {
// 		return fmt.Errorf("failed to auto-migrate basic tables: %w", err)
// 	}

// 	// รอบสอง: Migrate ตารางที่มี FK
// 	if err := p.DB.AutoMigrate(
// 		&entity.News{},
// 		&entity.Event{},
// 		&entity.EventInside{},
// 		&entity.EventOutside{},
// 		&entity.Done{},
// 	); err != nil {
// 		return fmt.Errorf("failed to auto-migrate FK tables: %w", err)
// 	}

// 	return nil
// }
func (p *postgreDatabase) AutoMigrate() error {
	// Migrate 'User' table ก่อน
	if err := p.DB.AutoMigrate(&entity.User{}); err != nil {
		return fmt.Errorf("user migrate failed: %w", err)
	}

	// Migrate ตารางอื่นๆ ที่ไม่พึ่ง 'Foreign Key' จาก 'User'
	if err := p.DB.AutoMigrate(
		&entity.Teacher{},
		&entity.Student{},
		&entity.News{},
		&entity.Faculty{},
		&entity.Branch{},
	); err != nil {
		return fmt.Errorf("failed to auto-migrate basic tables: %w", err)
	}

	// Migrate ตารางที่มี FK ไปยัง 'User'
	if err := p.DB.AutoMigrate(
		&entity.Event{},
		&entity.EventInside{},
		&entity.EventOutside{},
		&entity.Done{},
	); err != nil {
		return fmt.Errorf("failed to auto-migrate FK tables: %w", err)
	}

	return nil
}
