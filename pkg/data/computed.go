package data

import (
	"sort"
	"time"
)

func (s *Store) ComputeGainers() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.gainers = make([]Gainer, len(s.symbols))

	candlesIndex := len(s.dailyDays) - 1
	for id, symbol := range s.symbols {
		currentPrice := 0.0
		lastIndex, ok := s.intradayLastIndex[symbol]
		if ok {
			if intradayPrice, ok := s.intraday[symbol]; ok {
				currentPrice = intradayPrice[lastIndex].Price
			}
		}

		dailyPrice := 0.0
		if candlePrices, ok := s.daily[symbol]; ok {
			dailyPrice = candlePrices[candlesIndex].Close
		}

		percentChanged := 0.0
		if dailyPrice != 0 {
			percentChanged = ((currentPrice / dailyPrice) - 1) * 100
		}

		gainer := Gainer{
			Symbol:         symbol,
			PercentChanged: percentChanged,
			intradayIndex:  lastIndex,
		}
		s.gainers[id] = gainer
	}

	sort.Slice(s.gainers, func(i, j int) bool {
		return s.gainers[i].PercentChanged > s.gainers[j].PercentChanged
	})

	return nil
}

func (s *Store) GetTopGainers(count int) []TopGainer {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]TopGainer, count)

	for i := 0; i < count && i < len(s.gainers); i++ {
		symbol := s.gainers[i].Symbol
		intradayIndex := s.gainers[i].intradayIndex

		result[i] = TopGainer{
			Gainer:      s.gainers[i],
			Yesterday:   s.daily[symbol][len(s.dailyDays)-1],
			Current:     s.intraday[symbol][intradayIndex],
			LastUpdated: s.intradayMinTime.Add(time.Minute * time.Duration(intradayIndex)),
		}
	}

	return result
}
