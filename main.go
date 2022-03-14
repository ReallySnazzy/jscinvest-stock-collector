package main

import (
	"log"

	"github.com/reallysnazzy/jscinvest-stock-collector/database"
	"github.com/reallysnazzy/jscinvest-stock-collector/iexcloud"
)

func main() {
	settings := loadSettings()

	// Connect to database
	dbContext, err := database.CreateDatabaseContext(settings.dbHost, settings.dbPort, settings.dbUser, settings.dbPassword)
	if err != nil {
		log.Fatalf("Failed to connect to database: \n%s\n", err)
	}
	defer dbContext.Disconnect()

	// Get symbols
	symbols, err := dbContext.GetNextSymbolsToTrack()
	if err != nil {
		log.Fatalf("Failed to load symbols: \n%s\n", err)
	}

	// Update prices
	iexcloud := iexcloud.CreateIexcloudContext(settings.apiKey, settings.stockDateRange)
	for _, symbol := range symbols {
		log.Printf("Loading stock data for %s\n", symbol)
		stockEntries, err := iexcloud.GetStockHistorialPrices(symbol)
		if err == nil && stockEntries != nil {
			for _, stockEntry := range stockEntries {
				err = dbContext.AddStockPriceHistory(stockEntry)
				if err != nil {
					log.Printf("Failed to save record for %s: %s\n", symbol, err)
				}
			}
			err = dbContext.MarkSymbolCompleted(symbol)
			if err != nil {
				log.Printf("%s mark complete failed: %s\n", symbol, err)
			} else {
				log.Printf("Saved stock data for %s\n", symbol)
			}
		} else {
			if err != nil {
				log.Printf("Failed to get stock prices for %s: %s\n", symbol, err)
			} else {
				log.Printf("No stock information available for %s\n", symbol)
			}
			err = dbContext.MarkSymbolCompleted(symbol)
			if err != nil {
				log.Printf("%s mark complete failed: %s\n", symbol, err)
			} else {
				log.Printf("Marked %s as completed even though it failed to get information about it", symbol)
			}
		}
	}
}
