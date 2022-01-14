package app

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/remes2000/amu_financial_summary/account_transaction"
	"github.com/remes2000/amu_financial_summary/backup"
	"github.com/remes2000/amu_financial_summary/category"
	"github.com/remes2000/amu_financial_summary/global"
	"github.com/remes2000/amu_financial_summary/regexp"
	"github.com/remes2000/amu_financial_summary/report"
	"github.com/remes2000/amu_financial_summary/validators"
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
	db.AutoMigrate(&account_transaction.AccountTransaction{})
}

func initRestApi() *gin.Engine {
	gin.SetMode(os.Getenv("GIN_MODE"))
	rest := gin.Default()
	registerValidators()
	bindAllRoutes(rest)
	return rest
}

func registerValidators() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("validdate", validators.ValidDate)
		v.RegisterValidation("currency", validators.Currency)
	}
}

func bindAllRoutes(rest *gin.Engine) {
	category.BindRoutes(rest)
	account_transaction.BindRoutes(rest)
	report.BindRoutes(rest)
	backup.BindRoutes(rest)
}
