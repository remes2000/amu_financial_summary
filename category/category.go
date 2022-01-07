package category

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/remes2000/amu_financial_summary/common"
	"github.com/remes2000/amu_financial_summary/global"
	"github.com/remes2000/amu_financial_summary/regexp"
	"gorm.io/gorm"
	"log"
	"net/http"
)

type Category struct {
	gorm.Model
	Id      uint            `json:"id" binding:"required" gorm:"primaryKey"`
	Name    string          `json:"name" binding:"required" gorm:"notNull"`
	Regexps []regexp.Regexp `json:"regexps"`
}

func (c Category) GetUpdateMap() map[string]interface{} {
	return map[string]interface{}{
		"Name": c.Name,
	}
}

type CreateCategory struct {
	Name    string                `json:"name" binding:"required"`
	Regexps []regexp.CreateRegexp `json:"regexps"`
}

func (c CreateCategory) GetCategory() Category {
	return Category{Name: c.Name}
}

func GetAllCategories(categories *[]Category) error {
	if err := global.Database.Find(categories).Error; err != nil {
		return err
	}
	return nil
}

func GetCategoryById(category *Category, id uint) error {
	if err := global.Database.Where("id = ?", id).First(category).Error; err != nil {
		return err
	}
	return nil
}

// ---=== REST ===---

func BindRoutes(rest *gin.Engine) {
	controllerName := "category"
	rest.GET(controllerName, getAll)
	rest.GET(controllerName+"/:id", getOne)
	rest.POST(controllerName, create)
}

func getAll(context *gin.Context) {
	var result []Category
	if err := GetAllCategories(&result); err != nil {
		log.Print(err)
		context.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	context.JSON(http.StatusOK, result)
}

func getOne(context *gin.Context) {
	var requestedCategory Category
	var idUri common.IdUri

	if err := context.ShouldBindUri(&idUri); err != nil {
		context.Status(http.StatusBadRequest)
		return
	}
	if err := GetCategoryById(&requestedCategory, idUri.Id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			context.JSON(http.StatusNotFound, gin.H{"message": fmt.Sprintf("Cannot find entity with id %d", idUri.Id)})
			return
		}
		log.Print(err)
		context.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	context.JSON(http.StatusOK, requestedCategory)
}

func create(context *gin.Context) {
	var createCategoryRequest CreateCategory
	if err := context.BindJSON(&createCategoryRequest); err != nil {
		return
	}
	context.Status(404)
	//if err := AddNewRegexp(createRegexpRequest); err != nil {
	//	log.Print(err)
	//	context.AbortWithStatus(http.StatusInternalServerError)
	//	return
	//}
	//context.JSON(http.StatusOK, createRegexpRequest)
}
