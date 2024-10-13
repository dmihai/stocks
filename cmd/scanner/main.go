package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"

	"github.com/dmihai/stocks/pkg/api"
	"github.com/dmihai/stocks/pkg/auth"
	"github.com/dmihai/stocks/pkg/data"
	"github.com/dmihai/stocks/pkg/scanner"
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
	scan := scanner.NewScanner(db, store)

	go func() {
		err := scan.Start(startDate, endDate, currentDate)
		if err != nil {
			log.Fatal(err)
		}
	}()

	authnz := auth.NewAuth()

	server := api.NewServer(os.Getenv("SERVER_ADDR"), authnz, store)
	server.Start()
}
