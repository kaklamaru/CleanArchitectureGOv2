package database

// import (
// 	"fmt"
// 	"go-clean-arch/config"
// 	"go-clean-arch/pkg/hash"
// 	"go-clean-arch/internal/entity"
// 	"log"

// 	"gorm.io/driver/postgres"
// 	"gorm.io/gorm"
// )

// // PostgreSQL database struct
// type postgresDatabase struct {
// 	DB *gorm.DB
// }

// // ฟังก์ชันสำหรับสร้าง instance ของ PostgreSQL database
// func NewPostgresDatabase(cfg *config.Config) (Database, error) {
// 	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
// 		cfg.Postgres.Host,
// 		cfg.Postgres.User,
// 		cfg.Postgres.Password,
// 		cfg.Postgres.DBName,
// 		cfg.Postgres.Port,
// 	)

// 	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &postgresDatabase{DB: db}, nil
// }

// // ฟังก์ชันสำหรับดึง instance ของฐานข้อมูล
// func (p *postgresDatabase) GetDB() *gorm.DB {
// 	return p.DB
// }

// // ฟังก์ชันสำหรับทำ AutoMigrate
// func (p *postgresDatabase) AutoMigrate() error {
// 	if err := p.DB.AutoMigrate(&entity.User{}); err != nil {
// 		return fmt.Errorf("failed to migrate User: %w", err)
// 	}

// 	return nil
// }

// // ฟังก์ชันสำหรับ setup database
// func SetupDatabase(cfg *config.Config) Database {
// 	db, err := NewPostgresDatabase(cfg)
// 	if err != nil {
// 		log.Fatalf("failed to connect to database: %v", err)
// 	}

// 	err = db.AutoMigrate()
// 	if err != nil {
// 		log.Fatalf("failed to migrate database: %v", err)
// 	}

// 	addTriggerIfNotExists(db, "before_insert_students", `
//         CREATE OR REPLACE FUNCTION before_insert_students()
//         RETURNS TRIGGER AS $$
//         BEGIN
//             IF EXISTS (SELECT 1 FROM teachers WHERE user_id = NEW.user_id) THEN
//                 RAISE EXCEPTION 'User ID already exists in teachers';
//             END IF;
//             RETURN NEW;
//         END;
//         $$ LANGUAGE plpgsql;

//         CREATE TRIGGER before_insert_students
//         BEFORE INSERT ON students
//         FOR EACH ROW EXECUTE FUNCTION before_insert_students();
//     `)

// 	// เพิ่ม Trigger สำหรับ Teachers
// 	addTriggerIfNotExists(db, "before_insert_teachers", `
//         CREATE OR REPLACE FUNCTION before_insert_teachers()
//         RETURNS TRIGGER AS $$
//         BEGIN
//             IF EXISTS (SELECT 1 FROM students WHERE user_id = NEW.user_id) THEN
//                 RAISE EXCEPTION 'User ID already exists in students';
//             END IF;
//             RETURN NEW;
//         END;
//         $$ LANGUAGE plpgsql;

//         CREATE TRIGGER before_insert_teachers
//         BEFORE INSERT ON teachers
//         FOR EACH ROW EXECUTE FUNCTION before_insert_teachers();
//     `)

// 	// สร้าง admin user หากยังไม่มี
// 	password, err := hash.HashPassword(cfg.Admin.Password)
// 	if err != nil {
// 		log.Fatalf("failed to hash password: %v", err)
// 	}

// 	var admin entity.User
// 	result := db.GetDB().Where("email = ?", cfg.Admin.Email).First(&admin)
// 	if result.Error == nil {
// 		log.Println("Admin user already exists.")
// 	} else if result.Error == gorm.ErrRecordNotFound {
// 		user := entity.User{
// 			Email:    cfg.Admin.Email,
// 			Password: password,
// 			Role:     "superadmin",
// 		}
// 		createResult := db.GetDB().Create(&user)
// 		if createResult.Error != nil {
// 			log.Fatalf("failed to create admin user: %v", createResult.Error)
// 		} else {
// 			log.Println("Admin user created successfully!")
// 		}
// 	} else {
// 		log.Fatalf("failed to check admin user existence: %v", result.Error)
// 	}

// 	log.Println("PostgreSQL database connected, migrated, and triggers added successfully!")
// 	return db
// }

// // ฟังก์ชันเพื่อเพิ่ม Trigger ถ้ามันยังไม่มี
// func addTriggerIfNotExists(db Database, triggerName, triggerSQL string) {
// 	var count int
// 	// ตรวจสอบว่า trigger มีอยู่แล้วหรือไม่
// 	query := fmt.Sprintf("SELECT count(*) FROM pg_trigger WHERE tgname = '%s'", triggerName)
// 	err := db.GetDB().Raw(query).Scan(&count).Error
// 	if err != nil {
// 		log.Printf("Error checking trigger existence: %v", err)
// 		return
// 	}

// 	if count == 0 {
// 		// ถ้าไม่มี trigger นี้ก็สร้างใหม่
// 		if err := db.GetDB().Exec(triggerSQL).Error; err != nil {
// 			log.Printf("Failed to create trigger %s: %v", triggerName, err)
// 		} else {
// 			log.Printf("Trigger %s created successfully!", triggerName)
// 		}
// 	} else {
// 		log.Printf("Trigger %s already exists, skipping creation.", triggerName)
// 	}
// }
