package main

import (
	"github.com/joho/godotenv"
	"github.com/remes2000/amu_financial_summary/app"
)

func main() {
	godotenv.Load()
	app.Initialize()
	app.Run()
}
