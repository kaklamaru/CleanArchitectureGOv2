package main

import (
	"go-clean-arch/config"
	"go-clean-arch/database"
	"go-clean-arch/pkg/jwt"
	"go-clean-arch/pkg/server"
	"log"
	"time"
)
func deleteOldNews(db database.Database) {
	for {
		db.GetDB().Exec("SET TIME ZONE 'Asia/Bangkok'")

		result := db.GetDB().Exec("DELETE FROM news WHERE created_at < NOW() - INTERVAL '7 days'")
		if result.Error != nil {
			log.Printf("Error deleting old news: %v", result.Error)
		} else {
			log.Printf("Deleted %d old news data successfully.", result.RowsAffected)
		}

		time.Sleep(12 * time.Hour)
	}
}
func updatePastEventStatus(db database.Database) {
	for {
		db.GetDB().Exec("SET TIME ZONE 'Asia/Bangkok'")

		// อัปเดต status ของ event ที่ start_date < ปัจจุบัน
		result := db.GetDB().Exec("UPDATE events SET status = false WHERE start_date < NOW() AND status = true")
		
		if result.Error != nil {
			log.Printf("Error updating event statuses: %v", result.Error)
		} else {
			log.Printf("Updated %d event(s) to status=false (past start date).", result.RowsAffected)
		}

		// รอ 12 ชั่วโมงก่อนทำงานอีกครั้ง
		time.Sleep(12 * time.Hour)
	}
}

func main() {
	
	cfg := config.LoadConfig()
	db := database.SetupDatabase(cfg)

	jwt := jwt.NewJWTService(cfg)
	go deleteOldNews(db)
	go updatePastEventStatus(db)
	server,err := server.NewServer(cfg,db,jwt)
	if err != nil {
		log.Fatalf("Error creating server: %v", err)
	}

	if err := server.StartServer(); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}


}
