package data

import (
	"fmt"
	"math"
	"time"
)

type Store struct {
	symbols         map[int]string
	daily           map[string][]Candle
	dailyDays       map[string]int
	intraday        map[string][]Price
	intradayIndex   int
	intradayMinTime *time.Time
	gainers         []Gainer
}

func NewStore() *Store {
	return &Store{
		symbols:       make(map[int]string),
		daily:         make(map[string][]Candle),
		intraday:      make(map[string][]Price),
		intradayIndex: -1,
	}
}

func (s *Store) PopulateDailyData(candles []Daily, days map[string]int) error {
	s.dailyDays = days

	for _, candle := range candles {
		if _, ok := s.daily[candle.Symbol]; !ok {
			s.daily[candle.Symbol] = make([]Candle, len(days))
		}

		s.daily[candle.Symbol][days[candle.Day]] = candle.Candle
	}

	s.updateSymbolMap()

	return nil
}

func (s *Store) UpdateIntradayMinTime(minTime *time.Time) {
	s.intradayMinTime = minTime
}

func (s *Store) UpdateIntradayData(prices []Intraday) error {
	if s.intradayMinTime == nil {
		return fmt.Errorf("intraday min time must be set before updating intraday data")
	}

	length := 20 * 60 // 20h

	for _, price := range prices {
		if _, ok := s.intraday[price.Symbol]; !ok {
			s.intraday[price.Symbol] = make([]Price, length)
		}

		index := int(math.Round(price.Timestamp.Sub(*s.intradayMinTime).Minutes()))
		s.intraday[price.Symbol][index] = price.Price

		if index > s.intradayIndex {
			s.intradayIndex = index
		}
	}

	s.updateSymbolMap()

	return nil
}

func (s *Store) updateSymbolMap() {
	includedSymbols := make(map[string]bool)
	for _, symbol := range s.symbols {
		includedSymbols[symbol] = true
	}

	for symbol := range s.daily {
		if _, ok := includedSymbols[symbol]; !ok {
			s.symbols[len(s.symbols)] = symbol
			includedSymbols[symbol] = true
		}
	}

	for symbol := range s.intraday {
		if _, ok := includedSymbols[symbol]; !ok {
			s.symbols[len(s.symbols)] = symbol
			includedSymbols[symbol] = true
		}
	}
}
