package data

import (
	"fmt"
	"math"
	"sync"
	"time"
)

type Store struct {
	mu                sync.RWMutex
	symbolIDs         map[string]int
	symbolNames       map[int]string
	daily             [][]Candle
	dailyDays         map[string]int
	intraday          [][]Price
	intradayLastIndex []int
	intradayIndex     int
	intradayMinTime   *time.Time
	gainers           []Gainer
}

func NewStore() *Store {
	return &Store{
		symbolIDs:         make(map[string]int),
		symbolNames:       make(map[int]string),
		daily:             make([][]Candle, 0, 5000),
		intraday:          make([][]Price, 0, 5000),
		intradayLastIndex: make([]int, 0, 5000),
		gainers:           make([]Gainer, 0),
		intradayIndex:     -1,
	}
}

func (s *Store) PopulateDailyData(candles []Daily, days map[string]int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.dailyDays = days

	for _, candle := range candles {
		id := s.symbolNameToID(candle.Symbol)

		if id > len(s.daily) {
			return fmt.Errorf("invalid ID assigned to symbol %s in daily", candle.Symbol)
		} else if id == len(s.daily) {
			s.daily = append(s.daily, make([]Candle, len(days)))
		}

		s.daily[id][days[candle.Day]] = candle.Candle
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

	s.intraday = make([][]Price, len(s.daily))
	s.intradayLastIndex = make([]int, len(s.daily))

	length := 20 * 60 // 20h

	for i := range s.intraday {
		s.intraday[i] = make([]Price, length)
	}

	for _, price := range prices {
		id := s.symbolNameToID(price.Symbol)

		if id > len(s.intraday) {
			return fmt.Errorf("invalid ID assigned to symbol %s in intraday", price.Symbol)
		} else if id == len(s.intraday) {
			s.intraday = append(s.intraday, make([]Price, length))
			s.intradayLastIndex = append(s.intradayLastIndex, 0)
		}

		index := int(math.Round(price.Timestamp.Sub(*s.intradayMinTime).Minutes()))
		s.intraday[id][index] = price.Price

		if index > s.intradayIndex {
			s.intradayIndex = index
		}

		if s.intradayLastIndex[id] < index {
			s.intradayLastIndex[id] = index
		}
	}

	return nil
}

func (s *Store) symbolNameToID(name string) int {
	// checks if symbol exists
	if id, ok := s.symbolIDs[name]; ok {
		return id
	}

	// add new symbol
	newID := len(s.symbolIDs)
	s.symbolIDs[name] = newID
	s.symbolNames[newID] = name

	return newID
}

func (s *Store) symbolIDToName(id int) string {
	return s.symbolNames[id]
}
