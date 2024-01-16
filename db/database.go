package db

import (
	"spotify-api/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectionDB() {
	dsn := "your_username:your_password@tcp(localhost:3306)/spotify?parseTime=true"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to the database")
	}

	db.AutoMigrate(&models.Artist{}, &models.Track{})

	DB = db
}
