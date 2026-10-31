package store

import "gorm.io/gorm"

var DB *gorm.DB

type Store struct{ db *gorm.DB }
