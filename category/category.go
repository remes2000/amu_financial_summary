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
	Id      uint            `json:"id" binding:"required" gorm:"primaryKey"`
	Name    string          `json:"name" binding:"required" gorm:"notNull"`
	Regexps []regexp.Regexp `json:"regexps" binding:"required,dive" gorm:"constraint:OnDelete:CASCADE;"`
}

type CreateCategory struct {
	Name    string                `json:"name" binding:"required"`
	Regexps []regexp.CreateRegexp `json:"regexps" binding:"required,dive"`
}

func (c CreateCategory) GetCategory() Category {
	var regexps []regexp.Regexp
	for _, regexp := range c.Regexps {
		regexps = append(regexps, regexp.GetRegexp())
	}
	return Category{Name: c.Name, Regexps: regexps}
}

func GetAllCategories(categories *[]Category) error {
	if err := global.Database.Preload("Regexps").Find(categories).Error; err != nil {
		return err
	}
	return nil
}

func GetCategoryById(category *Category, id uint) error {
	if err := global.Database.Preload("Regexps").Where("id = ?", id).First(category).Error; err != nil {
		return err
	}
	return nil
}

func AddNewCategory(category *Category) error {
	if err := global.Database.Create(category).Error; err != nil {
		return err
	}
	return nil
}

func UpdateCategory(categoryToUpdate *Category, updateCategory *Category) error {
	return global.Database.Transaction(func(tx *gorm.DB) error {
		if err := global.Database.Model(categoryToUpdate).Association("Regexps").Replace(updateCategory.Regexps); err != nil {
			return err
		}
		if err := global.Database.Session(&gorm.Session{FullSaveAssociations: true}).Save(updateCategory).Error; err != nil {
			return err
		}
		// Delete orphans
		if err := global.Database.Where("category_id is null").Delete(&regexp.Regexp{}).Error; err != nil {
			return err
		}
		return nil
	})
}

// ---=== REST ===---

func BindRoutes(rest *gin.Engine) {
	controllerName := "category"
	rest.GET(controllerName, getAll)
	rest.GET(controllerName+"/:id", getOne)
	rest.POST(controllerName, create)
	rest.PUT(controllerName, update)
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
	categoryToCreate := createCategoryRequest.GetCategory()
	if err := AddNewCategory(&categoryToCreate); err != nil {
		log.Print(err)
		context.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	context.JSON(http.StatusOK, categoryToCreate)
}

func update(context *gin.Context) {
	var categoryToUpdate Category
	var updateCategory Category
	if err := context.BindJSON(&updateCategory); err != nil {
		return
	}
	if err := GetCategoryById(&categoryToUpdate, updateCategory.Id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			context.JSON(http.StatusNotFound, gin.H{"message": "Entity with provided id not found"})
			return
		}
		log.Print(err)
		context.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	if err := UpdateCategory(&categoryToUpdate, &updateCategory); err != nil {
		log.Print(err)
		context.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	context.JSON(http.StatusOK, updateCategory)
}
