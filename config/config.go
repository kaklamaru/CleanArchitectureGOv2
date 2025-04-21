// package config

// import (
// 	"fmt"
// 	"log"
// 	"os"
// 	"strconv"

// 	"github.com/joho/godotenv"
// )

// // Config struct สำหรับการเก็บค่าคอนฟิก
// type Config struct {
// 	DSN        string
// 	JWTSecret  string
// 	ServerPort int
// 	Admin      struct {
// 		Email    string
// 		Password string
// 	}
// }

// // LoadConfig โหลดค่าคอนฟิกจากไฟล์ .env
// func LoadConfig() *Config {
// 	// โหลดค่าคอนฟิกจากไฟล์ .env
// 	if err := godotenv.Load(); err != nil {
// 		log.Println("No .env file found, using environment variables")
// 	}

// 	// สร้าง DSN สำหรับการเชื่อมต่อฐานข้อมูล MySQL
// 	// dsn := fmt.Sprintf(
// 	// 	"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
// 	// 	getEnv("DB_USER", ""),          // ดึงค่า DB_USER จากไฟล์ .env
// 	// 	getEnv("DB_PASSWORD", ""),      // ดึงค่า DB_PASSWORD จากไฟล์ .env
// 	// 	getEnv("DB_HOST", "localhost"), // ค่าพื้นฐานถ้าไม่ได้ระบุ
// 	// 	getEnv("DB_PORT", "3306"),      // ค่าพื้นฐานถ้าไม่ได้ระบุ
// 	// 	getEnv("DB_NAME", ""),          // ดึงค่า DB_NAME จากไฟล์ .env
// 	// )
// 	// postgre
// 	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
//     os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))

// 	// ดึงค่า JWT_SECRET จากไฟล์ .env
// 	jwtSecret := getEnv("JWT_SECRET", "")

// 	// ดึงค่า SERVER_PORT จากไฟล์ .env และแปลงเป็น int
// 	serverPort, err := strconv.Atoi(getEnv("SERVER_PORT", "8080"))
// 	if err != nil {
// 		log.Printf("Invalid SERVER_PORT value, using default: %d", serverPort)
// 	}

// 	// ตรวจสอบว่าค่าที่จำเป็นถูกตั้งค่าแล้ว
// 	if dsn == "" || jwtSecret == "" {
// 		log.Fatalf("Required environment variables are missing")
// 	}

// 	// คืนค่า Config struct ที่มีค่าคอนฟิกทั้งหมด
// 	return &Config{
// 		DSN:        dsn,
// 		JWTSecret:  jwtSecret,
// 		ServerPort: serverPort,
// 		Admin: struct {
// 			Email    string
// 			Password string
// 		}{
// 			Email:    getEnv("USER", ""),
// 			Password: getEnv("PASSWORD", ""),
// 		},
// 	}
// }

// // getEnv
// func getEnv(key, defaultVal string) string {
// 	val := os.Getenv(key)
// 	if val == "" {
// 		return defaultVal
// 	}
// 	return val
// }
//

// package config

// import (
// 	"fmt"
// 	"log"
// 	"net/url"
// 	"os"
// 	"strconv"
// )

// // Config struct สำหรับการเก็บค่าคอนฟิก
// type Config struct {
// 	DSN        string
// 	JWTSecret  string
// 	ServerPort int
// 	Admin      struct {
// 		Email    string
// 		Password string
// 	}
// }

// // LoadConfig โหลดค่าคอนฟิกจาก environment variables
// func LoadConfig() *Config {
// 	// รับค่า DATABASE_URL
// 	databaseURL := os.Getenv("DATABASE_URL")
// 	if databaseURL == "" {
// 		log.Fatal("DATABASE_URL is required but not set")
// 	}

// 	// ใช้ net/url ในการแยกค่าจาก URL
// 	parsedURL, err := url.Parse(databaseURL)
// 	if err != nil {
// 		log.Fatalf("Invalid DATABASE_URL: %v", err)
// 	}

// 	user := parsedURL.User.Username()
// 	password, _ := parsedURL.User.Password()
// 	host := parsedURL.Hostname()
// 	port := parsedURL.Port()
// 	dbname := parsedURL.Path[1:] // ละเครื่องหมาย '/' หน้า database name
// 	sslmode := parsedURL.Query().Get("sslmode")
// 	if sslmode == "" {
// 		sslmode = "disable"
// 	}

// 	// สร้าง DSN สำหรับเชื่อมต่อ PostgreSQL
// 	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
// 		host, port, user, password, dbname, sslmode)

// 	// โหลดค่า JWT_SECRET
// 	jwtSecret := os.Getenv("JWT_SECRET")
// 	if jwtSecret == "" {
// 		log.Fatal("JWT_SECRET is required but not set")
// 	}

// 	// โหลด SERVER_PORT และแปลงเป็น int
// 	serverPort, err := strconv.Atoi(getEnv("SERVER_PORT", "8080"))
// 	if err != nil {
// 		log.Printf("Invalid SERVER_PORT value, using default: 8080")
// 		serverPort = 8080
// 	}

// 	// คืนค่า config struct
// 	return &Config{
// 		DSN:        dsn,
// 		JWTSecret:  jwtSecret,
// 		ServerPort: serverPort,
// 		Admin: struct {
// 			Email    string
// 			Password string
// 		}{
// 			Email:    getEnv("USER", ""),
// 			Password: getEnv("PASSWORD", ""),
// 		},
// 	}
// }

// // getEnv ดึงค่า environment variable หรือค่า default
// func getEnv(key, defaultVal string) string {
// 	val := os.Getenv(key)
// 	if val == "" {
// 		return defaultVal
// 	}
// 	return val
// }

package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
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

// LoadConfig โหลดค่าคอนฟิกจาก environment variables
func LoadConfig() *Config {
	// 
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}
	// 
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL is required but not set")
	}

	// ดึงค่า JWT_SECRET จาก environment
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET is required but not set")
	}

	// ดึง SERVER_PORT และแปลงเป็น int
	serverPort, err := strconv.Atoi(getEnv("SERVER_PORT", "8080"))
	if err != nil {
		log.Printf("Invalid SERVER_PORT value, using default 8080")
		serverPort = 8080
	}

	// คืนค่า Config struct
	return &Config{
		DSN:        databaseURL, // ใช้ URL ตรง ๆ ได้เลย
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

// getEnv ดึงค่า environment variable หรือคืนค่า default ถ้าไม่เจอ
func getEnv(key, defaultVal string) string {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}
	return val
}
