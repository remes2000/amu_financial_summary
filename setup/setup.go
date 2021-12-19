package setup

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var database *gorm.DB

func EstablishConnection() {
	dsn := "host=localhost user=postgres password=postgres dbname=amu_financial_summary port=5432"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database!")
	}
	database = db
}
