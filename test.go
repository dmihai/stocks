package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"

	"github.com/dmihai/stocks/pkg/api"
	"github.com/dmihai/stocks/pkg/auth"
	"github.com/dmihai/stocks/pkg/data"
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

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected!")

	store := data.NewStore()

	startDaily := time.Now()

	days, err := db.GetDaysBetweenDates(startDate, endDate)
	if err != nil {
		log.Fatal(err)
	}

	daily, err := db.GetCandlesBetweenDates(startDate, endDate)
	if err != nil {
		log.Fatal(err)
	}

	err = store.PopulateDailyData(daily, days)
	if err != nil {
		log.Fatal(err)
	}

	elapsedDaily := time.Since(startDaily)
	log.Printf("Daily load took %s", elapsedDaily)

	minTime, err := db.GetMinTimestampForDay(currentDate)
	if err != nil {
		log.Fatal(err)
	}

	store.UpdateIntradayMinTime(minTime)

	go func() {
		for i := 1; i < 20*60; i++ {
			startIntraday := time.Now()

			maxTime := minTime.Add(time.Minute * time.Duration(i))

			intradayPrices, err := db.GetIntradayCandles(currentDate, maxTime)
			if err != nil {
				log.Fatal(err)
			}

			err = store.UpdateIntradayData(intradayPrices)
			if err != nil {
				log.Fatal(err)
			}

			err = store.ComputeGainers()
			if err != nil {
				log.Fatal(err)
			}

			elapsedIntraday := time.Since(startIntraday)
			log.Printf("Intraday load to %s took %s", maxTime, elapsedIntraday)

			time.Sleep(time.Millisecond * 500)
		}
	}()

	authnz := auth.NewAuth()

	server := api.NewServer(os.Getenv("SERVER_ADDR"), authnz, store)
	server.Start()
}
