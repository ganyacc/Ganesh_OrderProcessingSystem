package database

import "gorm.io/gorm"

type Database interface {
	GetDb() *gorm.DB
	CloseDb(db *gorm.DB) error
	AutoMigrateTables() error
}
