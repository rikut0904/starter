package database

import "gorm.io/gorm"

// Ping verifies that the configured database can be reached.
func Ping(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}
