package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config struct สำหรับการเก็บค่าคอนฟิก
type Config struct {
	DSN        string
	JWTSecret  string
	ServerPort int
	Admin      struct {
		Email    string
		Password string
	}
}

// LoadConfig โหลดค่าคอนฟิกจากไฟล์ .env
func LoadConfig() *Config {
	// โหลดค่าคอนฟิกจากไฟล์ .env
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// สร้าง DSN สำหรับการเชื่อมต่อฐานข้อมูล MySQL
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		getEnv("DB_USER", ""),          // ดึงค่า DB_USER จากไฟล์ .env
		getEnv("DB_PASSWORD", ""),      // ดึงค่า DB_PASSWORD จากไฟล์ .env
		getEnv("DB_HOST", "localhost"), // ค่าพื้นฐานถ้าไม่ได้ระบุ
		getEnv("DB_PORT", "3306"),      // ค่าพื้นฐานถ้าไม่ได้ระบุ
		getEnv("DB_NAME", ""),          // ดึงค่า DB_NAME จากไฟล์ .env
	)

	// ดึงค่า JWT_SECRET จากไฟล์ .env
	jwtSecret := getEnv("JWT_SECRET", "")

	// ดึงค่า SERVER_PORT จากไฟล์ .env และแปลงเป็น int
	serverPort, err := strconv.Atoi(getEnv("SERVER_PORT", "8080"))
	if err != nil {
		log.Printf("Invalid SERVER_PORT value, using default: %d", serverPort)
	}

	// ตรวจสอบว่าค่าที่จำเป็นถูกตั้งค่าแล้ว
	if dsn == "" || jwtSecret == "" {
		log.Fatalf("Required environment variables are missing")
	}

	// คืนค่า Config struct ที่มีค่าคอนฟิกทั้งหมด
	return &Config{
		DSN:        dsn,
		JWTSecret:  jwtSecret,
		ServerPort: serverPort,
		Admin: struct {
			Email    string
			Password string
		}{
			Email:    getEnv("USER", ""),
			Password: getEnv("PASSWORD", ""),
		},
	}
}

// getEnv
func getEnv(key, defaultVal string) string {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}
	return val
}
