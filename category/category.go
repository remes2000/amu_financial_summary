package category

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Category struct {
	Id   uint   `json:"id" binding:"required" gorm:"primaryKey"`
	Name string `json:"name" binding:"required" gorm:"notNull"`
}

func (c Category) GetUpdateMap() map[string]interface{} {
	return map[string]interface{}{
		"Name": c.Name,
	}
}

type CreateCategory struct {
	Name string `json:"name" binding:"required"`
}

func (c CreateCategory) GetCategory() Category {
	return Category{Name: c.Name}
}

// ---=== REST ===---

func BindRoutes(rest *gin.Engine) {
	controllerName := "category"
	rest.GET(controllerName, getAll)
}

func getAll(context *gin.Context) {
	//var result []Regexp
	//if err := GetAllRegexp(&result); err != nil {
	//	log.Print(err)
	//	context.AbortWithStatus(http.StatusInternalServerError)
	//	return
	//}
	context.Status(http.StatusOK)
}
