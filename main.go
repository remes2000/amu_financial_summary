package main

import (
	"github.com/joho/godotenv"
	"github.com/remes2000/amu_financial_summary/app"
)

func main() {
	if err := godotenv.Load(); err != nil {
		panic("Error loading .env file")
	}
	app.Initialize()
	app.Run()
}
