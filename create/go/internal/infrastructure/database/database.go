package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// New opens the application database. Repositories should receive *gorm.DB through dependency injection.
func New(dsn string) (*gorm.DB, error) {
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}
