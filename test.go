package main

import (
	"fmt"
	"log"
	"time"

	"github.com/dmihai/stocks/pkg/api"
	"github.com/dmihai/stocks/pkg/data"
	"github.com/dmihai/stocks/pkg/store"
)

const (
	dbHost = "192.168.1.63:3306"
	dbUser = "stocks"
	dbPass = "stocks"
	dbName = "stocks"

	startDate   = "2024-07-26"
	endDate     = "2024-08-08"
	currentDate = "2024-08-09"

	serverAddr = ":3000"
)

func main() {
	db, err := store.NewConn(dbHost, dbUser, dbPass, dbName)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected!")

	store := data.NewStore()

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

	minTime, err := db.GetMinTimestampForDay(currentDate)
	if err != nil {
		log.Fatal(err)
	}

	store.UpdateIntradayMinTime(minTime)

	maxTime := minTime.Add(time.Minute * 10 * 60)

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

	server := api.NewServer(serverAddr, store)
	server.Start()
}
