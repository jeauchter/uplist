package database

import (
	"github.com/jeauchter/uplist/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitDB(dbFile string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(dbFile), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Migrate the schema
	db.AutoMigrate(&models.UplistOauth{})

	return db, nil
}
