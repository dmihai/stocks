package data

import (
	"fmt"
	"math"
	"time"
)

type Store struct {
	daily           map[string][]Candle
	dailyDays       map[string]int
	intraday        map[string][]Price
	intradayIndex   int
	intradayMinTime *time.Time
}

func NewStore() *Store {
	return &Store{
		intradayIndex: -1,
	}
}

func (s *Store) PopulateDailyData(candles []Daily, days map[string]int) error {
	s.dailyDays = days
	s.daily = make(map[string][]Candle)

	for _, candle := range candles {
		if _, ok := s.daily[candle.Symbol]; !ok {
			s.daily[candle.Symbol] = make([]Candle, len(days))
		}

		s.daily[candle.Symbol][days[candle.Day]] = candle.Candle
	}

	return nil
}

func (s *Store) GetDailyData() map[string][]Candle {
	return s.daily
}

func (s *Store) UpdateIntradayMinTime(minTime *time.Time) {
	s.intradayMinTime = minTime
}

func (s *Store) UpdateIntradayData(prices []Intraday) error {
	if s.intradayMinTime == nil {
		return fmt.Errorf("intraday min time must be set before appending intraday data")
	}

	length := 20 * 60 // 20h
	s.intraday = make(map[string][]Price)

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

	return nil
}

func (s *Store) GetIntradayData() map[string][]Price {
	return s.intraday
}
