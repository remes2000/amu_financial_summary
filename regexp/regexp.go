package regexp

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Regexp struct {
	Id      uint `gorm:"primary-key"`
	Content string
}

func BindRoutes(rest *gin.Engine) {
	controllerName := "regexp"
	rest.GET(controllerName, SayHello)
}

func SayHello(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{})
}
