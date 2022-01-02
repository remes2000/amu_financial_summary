package setup

import (
	"github.com/remes2000/amu_financial_summary/category"
	"github.com/remes2000/amu_financial_summary/global"
	"github.com/remes2000/amu_financial_summary/regexp"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func EstablishConnection() {
	dsn := "host=localhost user=postgres password=postgres dbname=amu_financial_summary port=5432"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database!")
	}
	global.Database = db
	Migrate()
}

func Migrate() {
	global.Database.AutoMigrate(&regexp.Regexp{})
	global.Database.AutoMigrate(&category.Category{})
}
