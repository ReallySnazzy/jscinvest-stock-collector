package main

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type settings struct {
	apiKey         string
	stockDateRange string
	dbHost         string
	dbUser         string
	dbPassword     string
	dbPort         int32
}

func loadSettings() settings {
	godotenv.Load()
	settings := settings{}
	settings.apiKey = os.Getenv("API_KEY")
	settings.stockDateRange = os.Getenv("IEXCLOUD_STOCK_DATE_RANGE")
	settings.dbHost = os.Getenv("DB_HOST")
	settings.dbPassword = os.Getenv("DB_PASSWORD")
	settings.dbUser = os.Getenv("DB_USER")
	port, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		log.Fatal(err)
	}
	settings.dbPort = int32(port)
	return settings
}
