package database

import (
	"gorm.io/gorm"
)

type Database interface {
	GetDB() *gorm.DB
	AutoMigrate() error
}
