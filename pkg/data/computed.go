package data

import (
	"sort"
	"time"
)

func (s *Store) GetTopGainers(count int) []TopGainer {
	s.mu.RLock()
	defer s.mu.RUnlock()

	gainers := s.getGainers()

	result := make([]TopGainer, count)

	for i := 0; i < count && i < len(gainers); i++ {
		symbol := s.symbols[gainers[i].symbolID]

		result[i] = TopGainer{
			Symbol:         symbol.Name,
			PercentChanged: symbol.percentChanged,
			Yesterday:      s.getPrevDayCandle(symbol),
			Current:        s.getCurrentPrice(symbol),
			LastUpdated:    s.intradayMinTime.Add(time.Minute * time.Duration(symbol.intradayIndex)),
		}
	}

	return result
}

func (s *Store) getGainers() []Gainer {
	gainers := make([]Gainer, len(s.symbols))
	for i, symbol := range s.symbols {
		gainers[i] = Gainer{
			symbolID:       i,
			percentChanged: symbol.percentChanged,
		}
	}

	sort.Slice(gainers, func(i, j int) bool {
		return gainers[i].percentChanged > gainers[j].percentChanged
	})

	return gainers
}

func (s *Store) getPrevDayCandle(symbol Symbol) Candle {
	if len(s.dailyDays) > 0 {
		return s.daily[symbol.Name][len(s.dailyDays)-1]
	}

	return Candle{
		Close: symbol.PrevDayClose,
	}
}

func (s *Store) getCurrentPrice(symbol Symbol) Price {
	if symbol.intradayIndex >= 0 {
		return s.intraday[symbol.Name][symbol.intradayIndex]
	}

	return Price{}
}
