package data

import "sort"

func (s *Store) ComputeGainers() error {
	s.gainers = make([]Gainer, len(s.symbols))

	candlesIndex := len(s.dailyDays) - 1
	for id, symbol := range s.symbols {
		currentPrice := 0.0
		if intradayPrices, ok := s.intraday[symbol]; ok {
			currentPrice = intradayPrices[s.intradayIndex].Price
		}

		dailyPrice := 0.0
		if candlePrices, ok := s.daily[symbol]; ok {
			dailyPrice = candlePrices[candlesIndex].Close
		}

		percentChange := 0.0
		if dailyPrice != 0 {
			percentChange = ((currentPrice / dailyPrice) - 1) * 100
		}

		gainer := Gainer{
			Symbol:        symbol,
			PercentChange: percentChange,
		}
		s.gainers[id] = gainer
	}

	sort.Slice(s.gainers, func(i, j int) bool {
		return s.gainers[i].PercentChange > s.gainers[j].PercentChange
	})

	return nil
}

func (s *Store) GetTopGainers(count int) []TopGainer {
	result := make([]TopGainer, count)

	for i := 0; i < count && i < len(s.gainers); i++ {
		symbol := s.gainers[i].Symbol
		result[i] = TopGainer{
			Gainer:    s.gainers[i],
			Yesterday: s.daily[symbol][len(s.dailyDays)-1],
			Current:   s.intraday[symbol][s.intradayIndex],
		}
	}

	return result
}
