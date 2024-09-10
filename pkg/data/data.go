package data

type Store struct {
	daily         map[string][]Candle
	dailyDays     map[string]int
	intraday      map[string][]Candle
	intradayIndex int
}

func NewStore() *Store {
	return &Store{
		intradayIndex: -1,
	}
}

func (s *Store) PopulateDailyData(daily []Daily, days map[string]int) error {
	s.dailyDays = days
	s.daily = make(map[string][]Candle)

	for _, day := range daily {
		if _, ok := s.daily[day.Symbol]; !ok {
			s.daily[day.Symbol] = make([]Candle, len(days))
		}

		s.daily[day.Symbol][days[day.Day]] = day.Candle
	}

	return nil
}

func (s *Store) GetDailyData() map[string][]Candle {
	return s.daily
}
