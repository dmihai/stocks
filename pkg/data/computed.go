package data

import (
	"sort"
	"time"
)

func (s *Store) GetTopGainers(count int) []TopGainer {
	s.mu.RLock()
	defer s.mu.RUnlock()

	gainers := make([]SortGainer, len(s.symbols))
	for i, symbol := range s.symbols {
		gainers[i] = SortGainer{
			symbolID:       i,
			percentChanged: symbol.percentChanged,
		}
	}

	sort.Slice(gainers, func(i, j int) bool {
		return gainers[i].percentChanged > gainers[j].percentChanged
	})

	result := make([]TopGainer, count)

	for i := 0; i < count && i < len(gainers); i++ {
		symbol := s.symbols[gainers[i].symbolID]

		result[i] = TopGainer{
			Symbol:         symbol.name,
			PercentChanged: symbol.percentChanged,
			Yesterday:      s.daily[symbol.name][len(s.dailyDays)-1],
			Current:        s.intraday[symbol.name][symbol.intradayIndex],
			LastUpdated:    s.intradayMinTime.Add(time.Minute * time.Duration(symbol.intradayIndex)),
		}
	}

	return result
}
