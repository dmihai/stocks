package main

import (
	"fmt"
	"log"
	"sort"

	"github.com/dmihai/stocks/pkg/data"
	"github.com/dmihai/stocks/pkg/store"
)

func main() {
	db, err := store.NewConn("192.168.1.63:3306", "stocks", "stocks", "stocks")
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected!")

	days, err := db.GetDaysBetweenDates("2024-07-26", "2024-08-08")
	if err != nil {
		log.Fatal(err)
	}

	candles, err := db.GetCandlesBetweenDates(days)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("len(candles): %d\n", len(candles))
	fmt.Printf("candles[AAPL][0]: %+v\n", candles["AAPL"][0])
	fmt.Printf("candles[AAPL][1]: %+v\n", candles["AAPL"][1])
	fmt.Printf("candles[MSFT][0]: %+v\n", candles["MSFT"][0])

	intraday, err := db.GetIntradayCandles("2024-08-09")
	if err != nil {
		log.Fatal(err)
	}
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
			currentPrice = intradayPrices[intradayIndex].Close
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

func getSymbolMap(symbolsList ...map[string][]data.Candle) map[int]string {
	result := make(map[int]string)
	includedSymbols := make(map[string]bool)
	i := 0

	for _, symbols := range symbolsList {
		for symbol := range symbols {
			if _, ok := includedSymbols[symbol]; !ok {
				result[i] = symbol
				i += 1
				includedSymbols[symbol] = true
			}
		}
	}

	return result
}
