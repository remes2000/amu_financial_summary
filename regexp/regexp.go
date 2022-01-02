package regexp

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/remes2000/amu_financial_summary/common"
	"github.com/remes2000/amu_financial_summary/global"
	"gorm.io/gorm"
	"log"
	"net/http"
)

type Regexp struct {
	Id      uint   `json:"id" binding:"required" gorm:"primaryKey"`
	Content string `json:"content" binding:"required" gorm:"notNull"`
}

func (r Regexp) GetUpdateMap() map[string]interface{} {
	return map[string]interface{}{
		"Content": r.Content,
	}
}

type CreateRegexp struct {
	Content string `json:"content" binding:"required"`
}

func (c CreateRegexp) GetRegexp() Regexp {
	return Regexp{Content: c.Content}
}

func GetAllRegexp(regexp *[]Regexp) error {
	if err := global.Database.Find(regexp).Error; err != nil {
		return err
	}
	return nil
}

func GetRegexpById(regexp *Regexp, id uint) error {
	if err := global.Database.Where("id = ?", id).First(regexp).Error; err != nil {
		return err
	}
	return nil
}

func UpdateRegexp(regexpToUpdate *Regexp, updateRegexp *Regexp) error {
	if err := global.Database.Model(regexpToUpdate).Updates(updateRegexp.GetUpdateMap()).Error; err != nil {
		return err
	}
	return nil
}

func AddNewRegexp(request CreateRegexp) error {
	newRegexp := request.GetRegexp()
	if err := global.Database.Create(&newRegexp).Error; err != nil {
		return err
	}
	return nil
}

func DeleteRegexp(regexpToDelete *Regexp) error {
	if err := global.Database.Delete(regexpToDelete).Error; err != nil {
		return err
	}
	return nil
}

// ---=== REST ===---

func BindRoutes(rest *gin.Engine) {
	controllerName := "regexp"
	rest.POST(controllerName, create)
	rest.PUT(controllerName, update)
	rest.DELETE(controllerName+"/:id", delete)
	rest.GET(controllerName+"/:id", getOne)
	rest.GET(controllerName, getAll)
}

func create(context *gin.Context) {
	var createRegexpRequest CreateRegexp
	if err := context.BindJSON(&createRegexpRequest); err != nil {
		return
	}
	if err := AddNewRegexp(createRegexpRequest); err != nil {
		log.Print(err)
		context.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	context.JSON(http.StatusOK, createRegexpRequest)
}

func update(context *gin.Context) {
	var updateRegexp Regexp
	var regexpToUpdate Regexp
	if err := context.BindJSON(&updateRegexp); err != nil {
		return
	}
	if err := GetRegexpById(&regexpToUpdate, updateRegexp.Id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			context.JSON(http.StatusNotFound, gin.H{"message": "Entity with provided id not found"})
			return
		}
		log.Print(err)
		context.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	if err := UpdateRegexp(&regexpToUpdate, &updateRegexp); err != nil {
		log.Print(err)
		context.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	context.JSON(http.StatusOK, regexpToUpdate)
}

func delete(context *gin.Context) {
	var regexpToDelete Regexp
	var idUri common.IdUri

	if err := context.ShouldBindUri(&idUri); err != nil {
		context.Status(http.StatusBadRequest)
	}
	if err := GetRegexpById(&regexpToDelete, idUri.Id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			context.JSON(http.StatusNotFound, gin.H{"message": "Entity with provided id not found"})
			return
		}
		log.Print(err)
		context.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	if err := DeleteRegexp(&regexpToDelete); err != nil {
		log.Print(err)
		context.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	context.Status(http.StatusOK)
}

func getOne(context *gin.Context) {
	var requestedRegexp Regexp
	var idUri common.IdUri

	if err := context.ShouldBindUri(&idUri); err != nil {
		context.Status(http.StatusBadRequest)
	}
	if err := GetRegexpById(&requestedRegexp, idUri.Id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			context.JSON(http.StatusNotFound, gin.H{"message": "Entity with provided id not found"})
			return
		}
		log.Print(err)
		context.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	context.JSON(http.StatusOK, requestedRegexp)
}

func getAll(context *gin.Context) {
	var result []Regexp
	if err := GetAllRegexp(&result); err != nil {
		log.Print(err)
		context.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	context.JSON(http.StatusOK, result)
}
