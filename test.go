package main

import (
	"fmt"
	"log"
	"sort"
	"time"

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

	candles := store.GetDailyData()

	fmt.Printf("len(candles): %d\n", len(candles))
	fmt.Printf("candles[AAPL][0]: %+v\n", candles["AAPL"][0])
	fmt.Printf("candles[AAPL][1]: %+v\n", candles["AAPL"][1])
	fmt.Printf("candles[MSFT][0]: %+v\n", candles["MSFT"][0])

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

	intraday := store.GetIntradayData()
	fmt.Printf("len(intraday): %d\n", len(intraday))
	fmt.Printf("intraday[AAPL][300]: %+v\n", intraday["AAPL"][300])
	fmt.Printf("intraday[AAPL][550]: %+v\n", intraday["AAPL"][550])
	fmt.Printf("intraday[MSFT][0]: %+v\n", intraday["MSFT"][0])
	fmt.Printf("intraday[MSFT][580]: %+v\n", intraday["MSFT"][580])
	fmt.Printf("intraday[MSFT][599]: %+v\n", intraday["MSFT"][599])

	symbols := getSymbolMap(candles, intraday)
	fmt.Printf("len(symbols): %d\n", len(symbols))
	fmt.Printf("symbols[10]: %s\n", symbols[10])
	fmt.Printf("symbols[1000]: %s\n", symbols[1000])

	gainers := make([]data.Gainer, len(symbols))

	intradayIndex := 450
	candlesIndex := len(days) - 1
	for id, symbol := range symbols {
		currentPrice := 0.0
		if intradayPrices, ok := intraday[symbol]; ok {
			currentPrice = intradayPrices[intradayIndex].Price
		}

		dailyPrice := 0.0
		if candlePrices, ok := candles[symbol]; ok {
			dailyPrice = candlePrices[candlesIndex].Close
		}

		percentChange := 0.0
		if dailyPrice != 0 {
			percentChange = ((currentPrice / dailyPrice) - 1) * 100
		}

		gainer := data.Gainer{
			Symbol:        symbol,
			PercentChange: percentChange,
		}
		gainers[id] = gainer
	}

	sort.Slice(gainers, func(i, j int) bool {
		return gainers[i].PercentChange > gainers[j].PercentChange
	})

	i := 1
	for _, gainer := range gainers {
		fmt.Printf("%s (%f): yesterday[%+v] now[%+v]\n", gainer.Symbol, gainer.PercentChange, candles[gainer.Symbol][candlesIndex], intraday[gainer.Symbol][intradayIndex])
		i++
		if i > 20 {
			break
		}
	}
}

func getSymbolMap(daily map[string][]data.Candle, intraday map[string][]data.Price) map[int]string {
	result := make(map[int]string)
	includedSymbols := make(map[string]bool)
	i := 0

	for symbol := range daily {
		if _, ok := includedSymbols[symbol]; !ok {
			result[i] = symbol
			i += 1
			includedSymbols[symbol] = true
		}
	}

	for symbol := range intraday {
		if _, ok := includedSymbols[symbol]; !ok {
			result[i] = symbol
			i += 1
			includedSymbols[symbol] = true
		}
	}

	return result
}
