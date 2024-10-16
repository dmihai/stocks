package data

import (
	"fmt"
	"math"
	"sync"
	"time"
)

const (
	estimatedSymbolsCount = 5000
)

type Store struct {
	mu              sync.RWMutex
	symbols         []Symbol
	symbolIDs       map[string]int
	daily           map[string][]Candle
	dailyDays       map[string]int
	intraday        map[string][]Price
	intradayIndex   int
	intradayMinTime *time.Time
}

func NewStore() *Store {
	return &Store{
		symbols:       make([]Symbol, 0, estimatedSymbolsCount),
		symbolIDs:     make(map[string]int, estimatedSymbolsCount),
		daily:         make(map[string][]Candle, estimatedSymbolsCount),
		intraday:      make(map[string][]Price, estimatedSymbolsCount),
		intradayIndex: -1,
	}
}

func (s *Store) PopulateDailyData(candles []Daily, days map[string]int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.dailyDays = days

	for _, candle := range candles {
		if _, ok := s.daily[candle.Symbol]; !ok {
			s.daily[candle.Symbol] = make([]Candle, len(days))
		}

		s.daily[candle.Symbol][days[candle.Day]] = candle.Candle

		s.addSymbol(candle.Symbol)
	}

	return nil
}

func (s *Store) PopulateSymbols(symbols []Symbol) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, symbol := range symbols {
		id := s.addSymbol(symbol.Name)
		s.symbols[id].PrevDayClose = symbol.PrevDayClose
		s.symbols[id].Shares = symbol.Shares
	}

	return nil
}

func (s *Store) UpdateIntradayMinTime(minTime *time.Time) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.intradayMinTime = minTime
}

func (s *Store) UpdateIntradayData(prices []Intraday) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.intradayMinTime == nil {
		return fmt.Errorf("intraday min time must be set before updating intraday data")
	}

	length := 20 * 60 // 20h

	for _, price := range prices {
		if _, ok := s.symbolIDs[price.Symbol]; !ok {
			continue
		}

		if _, ok := s.intraday[price.Symbol]; !ok {
			s.intraday[price.Symbol] = make([]Price, length)
		}

		index := int(math.Round(price.Timestamp.Sub(*s.intradayMinTime).Minutes()))
		if index < 0 {
			continue
		}

		s.intraday[price.Symbol][index] = price.Price

		if index > s.intradayIndex {
			s.intradayIndex = index
		}

		symbolID := s.addSymbol(price.Symbol)
		if s.symbols[symbolID].intradayIndex < index {
			s.symbols[symbolID].intradayIndex = index
		}

		s.updateSymbol(symbolID)
	}

	return nil
}

func (s *Store) addSymbol(symbol string) int {
	// symbol already exists
	if id, ok := s.symbolIDs[symbol]; ok {
		return id
	}

	s.symbols = append(s.symbols, Symbol{
		Name:          symbol,
		intradayIndex: -1,
	})

	newID := len(s.symbols) - 1
	s.symbolIDs[symbol] = newID

	return newID
}

func (s *Store) updateSymbol(symbolID int) {
	symbol := s.symbols[symbolID]
	currentPrice := s.intraday[symbol.Name][symbol.intradayIndex].Price

	prevDayCandle := s.getPrevDayCandle(symbol)
	lastDayClose := prevDayCandle.Close

	percentChanged := 0.0
	if lastDayClose != 0 {
		percentChanged = ((currentPrice / lastDayClose) - 1) * 100
	}

	s.symbols[symbolID].percentChanged = percentChanged
}
