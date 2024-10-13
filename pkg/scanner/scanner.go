package scanner

import (
	"log"
	"time"

	"github.com/dmihai/stocks/pkg/data"
	"github.com/dmihai/stocks/pkg/providers/fmp"
	"github.com/dmihai/stocks/pkg/store"
)

type Scanner struct {
	db    *store.Conn
	store *data.Store
	fmp   *fmp.Client
}

func NewScanner(db *store.Conn, store *data.Store, fmp *fmp.Client) *Scanner {
	return &Scanner{
		db:    db,
		store: store,
		fmp:   fmp,
	}
}

func (s *Scanner) Start(startDate, endDate, currentDate string) error {
	symbols, err := s.fmp.GetAvailableSymbolsByExchange(fmp.ExchangeNASDAQ)
	if err != nil {
		return err
	}
	log.Printf("symbols count: %d\n", len(symbols))
	log.Printf("first symbol: %+v\n", symbols[0])

	err = s.populateDailyData(startDate, endDate)
	if err != nil {
		return err
	}

	minTime, err := s.db.GetMinTimestampForDay(currentDate)
	if err != nil {
		return err
	}

	s.store.UpdateIntradayMinTime(minTime)

	for i := 1; i < 20*60; i++ {
		maxTime := minTime.Add(time.Minute * time.Duration(i))

		err = s.updateIntradayData(currentDate, maxTime)
		if err != nil {
			log.Fatal(err)
		}

		err = s.store.ComputeGainers()
		if err != nil {
			log.Fatal(err)
		}

		time.Sleep(time.Millisecond * 500)
	}

	return nil
}

func (s *Scanner) populateDailyData(startDate, endDate string) error {
	defer timeTrack(time.Now(), "populateDailyData")

	days, err := s.db.GetDaysBetweenDates(startDate, endDate)
	if err != nil {
		return err
	}

	daily, err := s.db.GetCandlesBetweenDates(startDate, endDate)
	if err != nil {
		return err
	}

	err = s.store.PopulateDailyData(daily, days)
	if err != nil {
		return err
	}

	return nil
}

func (s *Scanner) updateIntradayData(currentDate string, maxTime time.Time) error {
	defer timeTrack(time.Now(), "updateIntradayData")

	intradayPrices, err := s.db.GetIntradayCandles(currentDate, maxTime)
	if err != nil {
		log.Fatal(err)
	}

	err = s.store.UpdateIntradayData(intradayPrices)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}
