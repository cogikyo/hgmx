package gorm

type DB struct{}

func (db *DB) AutoMigrate(dst ...any) error { return nil }
