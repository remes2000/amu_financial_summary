package app

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/remes2000/amu_financial_summary/category"
	"github.com/remes2000/amu_financial_summary/global"
	"github.com/remes2000/amu_financial_summary/regexp"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
)

func Initialize() {
	global.Database = establishDatabaseConnection()
	global.Rest = initRestApi()
}

func Run() {
	global.Rest.Run()
}

func establishDatabaseConnection() *gorm.DB {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s",
		os.Getenv("DB_HOSTNAME"),
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"))
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database!")
	}
	migrateDatabase(db)
	return db
}

func migrateDatabase(db *gorm.DB) {
	db.AutoMigrate(&category.Category{})
	db.AutoMigrate(&regexp.Regexp{})
}

func initRestApi() *gin.Engine {
	gin.SetMode(os.Getenv("GIN_MODE"))
	rest := gin.Default()
	bindAllRoutes(rest)
	return rest
}

func bindAllRoutes(rest *gin.Engine) {
	regexp.BindRoutes(rest)
	category.BindRoutes(rest)
}
