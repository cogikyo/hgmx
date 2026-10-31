package store

import "gorm.io/gorm"

type User struct{}

func New(db *gorm.DB) error {
	return db.AutoMigrate(&User{})
}
