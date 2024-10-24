package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"

	"github.com/dmihai/stocks/pkg/api"
	"github.com/dmihai/stocks/pkg/auth"
	"github.com/dmihai/stocks/pkg/data"
	"github.com/dmihai/stocks/pkg/scanner"
	"github.com/dmihai/stocks/pkg/stocks"
	"github.com/dmihai/stocks/pkg/store"
)

const (
	startDate   = "2024-07-26"
	endDate     = "2024-08-08"
	currentDate = "2024-08-09"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db, err := store.NewConn(
		os.Getenv("MYSQL_HOST"),
		os.Getenv("MYSQL_USER"),
		os.Getenv("MYSQL_PASS"),
		os.Getenv("MYSQL_DB"),
	)
	if err != nil {
		log.Fatal(err)
	}

	store := data.NewStore()

	fmp := stocks.NewFMPClient(os.Getenv("FMP_API_URL"), os.Getenv("FMP_API_KEY"))
	polygon := stocks.NewPolygonClient(os.Getenv("POLYGON_API_URL"), os.Getenv("POLYGON_API_KEY"))
	stocksClient := stocks.NewHybridClient(fmp, polygon)

	scan := scanner.NewScanner(db, store, stocksClient)

	go startScanner(scan)

	authnz := auth.NewAuth()

	server := api.NewServer(os.Getenv("SERVER_ADDR"), authnz, store, stocksClient)
	server.Start()
}

func startScanner(scan *scanner.Scanner) {
	err := scan.Start(startDate, endDate, currentDate)
	if err != nil {
		log.Fatal(err)
	}
}
