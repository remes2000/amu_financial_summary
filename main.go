package main

import (
	"github.com/gin-gonic/gin"
	"github.com/remes2000/amu_financial_summary/regexp"
	"github.com/remes2000/amu_financial_summary/setup"
)

func main() {
	rest := gin.Default()
	setup.EstablishConnection()
	BindAllRoutes(rest)
	rest.Run()
}

func BindAllRoutes(rest *gin.Engine) {
	regexp.BindRoutes(rest)
}
