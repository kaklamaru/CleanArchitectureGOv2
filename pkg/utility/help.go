package utility

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	jwtService "github.com/golang-jwt/jwt/v5"
)

func GetClaimsFromContext(ctx *fiber.Ctx) (jwtService.MapClaims, error) {
	claims, ok := ctx.Locals("claims").(jwtService.MapClaims)
	if !ok {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "Invalid claims")
	}
	return claims, nil
}

func ParseStartDate(dateStr string) (time.Time, error) {
	location, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		return time.Time{}, fiber.NewError(fiber.StatusInternalServerError, "Failed to load location")
	}
	startDate, err := time.ParseInLocation("2006-01-02 15:04:05", dateStr, location)
	if err != nil {
		return time.Time{}, fiber.NewError(fiber.StatusBadRequest, "Invalid date format, use 'YYYY-MM-DD HH:MM:SS'")
	}
	return startDate, nil
}

var thaiMonths = []string{
    "มกราคม", "กุมภาพันธ์", "มีนาคม", "เมษายน", "พฤษภาคม", "มิถุนายน",
    "กรกฎาคม", "สิงหาคม", "กันยายน", "ตุลาคม", "พฤศจิกายน", "ธันวาคม",
}

// ฟังก์ชันแปลง time.Time เป็นรูปแบบวันเดือนปีแบบไทย
func FormatToThaiDate(t time.Time) string {
    location, err := time.LoadLocation("Asia/Bangkok")
    if err != nil {
        return "Failed to load location"
    }

    t = t.In(location)

    // ดึงข้อมูลวัน เดือน ปี
    day := t.Day()
    month := thaiMonths[t.Month()-1]
    year := t.Year()+543

    // คืนค่าข้อความรูปแบบวันเดือนปี
    return fmt.Sprintf("%02d %s %d", day, month, year)
}


func FormatToThaiTime(t time.Time) string {
    location, err := time.LoadLocation("Asia/Bangkok")
    if err != nil {
        return "Failed to load location"
    }

    t = t.In(location)

    thaiTimeFormat := "15:04"

    return t.Format(thaiTimeFormat)
}

func AddHoursToTime(t time.Time, hoursToAdd uint) string {
	t = t.Add(time.Duration(hoursToAdd) * time.Hour)
	// t.Add(time.Duration(hoursToAdd) * time.Hour)
	thaiTimeFormat := "15:04"
    return t.Format(thaiTimeFormat)
}

func DecodeIDs(dataStr string) ([]uint, error) {
	var ids []uint

	// เช็คว่าเป็น string ว่างหรือไม่ (หลังจาก Trim เครื่องหมาย escape ออก)
	if strings.Trim(dataStr, "\"") == "" {
		// // ถ้าเป็น string ว่าง ให้ return slice ว่าง
		// fmt.Println("IDs: []") // กรณีที่เป็น string ว่าง
		return ids, nil
	}

	// ถ้า string ไม่ว่างให้ Trim เครื่องหมายคำพูดที่ escape ออก
	dataStr = strings.Trim(dataStr, "\"")

	// แปลง string ที่เหลือเป็น []uint
	if err := json.Unmarshal([]byte(dataStr), &ids); err != nil {
		// ถ้าเกิดข้อผิดพลาดในการ decode
		fmt.Println("Error decoding IDs:", err)
		return nil, err
	}

	return ids, nil
}
