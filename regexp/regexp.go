package regexp

import (
	"errors"
	"github.com/gin-gonic/gin"
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
		"content": r.Content,
	}
}

type CreateRegexp struct {
	Content string `json:"content" binding:"required"`
}

func (c CreateRegexp) GetRegexp() Regexp {
	return Regexp{Content: c.Content}
}

type UpdateRegexp struct {
	Content string
}

func BindRoutes(rest *gin.Engine) {
	controllerName := "regexp"
	rest.POST(controllerName, Create)
	rest.PUT(controllerName, Update)
	rest.DELETE(controllerName+"/:id", Delete)
}

func Create(context *gin.Context) {
	var createRegexpRequest CreateRegexp
	if err := context.BindJSON(&createRegexpRequest); err != nil {
		return
	}
	newRegexp := createRegexpRequest.GetRegexp()
	if err := global.Database.Create(&newRegexp).Error; err != nil {
		log.Print(err)
		context.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	context.JSON(http.StatusOK, newRegexp)
}

func Update(context *gin.Context) {
	var regexpToUpdate Regexp
	if err := context.BindJSON(&regexpToUpdate); err != nil {
		return
	}
	updateMap := regexpToUpdate.GetUpdateMap()
	if err := global.Database.Where("id = ?", regexpToUpdate.Id).First(&regexpToUpdate).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			context.JSON(http.StatusNotFound, gin.H{"message": "Entity with provided id not found"})
			return
		}
		log.Print(err)
		context.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	global.Database.Model(&regexpToUpdate).Updates(updateMap)
	context.JSON(http.StatusOK, regexpToUpdate)
}

func Delete(context *gin.Context) {
	// todo handle if not exist throw 404
	if err := global.Database.Where("id = ?", context.Param("id")).Delete(&Regexp{}).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			context.JSON(http.StatusNotFound, gin.H{"message": "Entity with provided id not found"})
			return
		}
		log.Print(err)
		context.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	context.Status(http.StatusOK)
}
