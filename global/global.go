package global

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var Database *gorm.DB
var Rest *gin.Engine
