package scanner

import (
	"log"
	"time"

	"github.com/dmihai/stocks/pkg/data"
	"github.com/dmihai/stocks/pkg/stocks"
	"github.com/dmihai/stocks/pkg/store"
)

type Scanner struct {
	db     *store.Conn
	store  *data.Store
	stocks stocks.Client
}

func NewScanner(db *store.Conn, store *data.Store, stocks stocks.Client) *Scanner {
	return &Scanner{
		db:     db,
		store:  store,
		stocks: stocks,
	}
}

func (s *Scanner) Start(startDate, endDate, currentDate string) error {
	exchanges := []string{
		stocks.ExchangeNASDAQ,
		stocks.ExchangeNYSE,
		stocks.ExchangeAMEX,
	}

	for _, exchange := range exchanges {
		err := s.populateSymbols(exchange)
		if err != nil {
			return err
		}
	}

	// err = s.populateDailyData(startDate, endDate)
	// if err != nil {
	// 	return err
	// }

	minTime := time.Now()

	s.store.UpdateIntradayMinTime(&minTime)

	for i := 1; i < 20*60; i++ {
		err := s.updateIntradayData()
		if err != nil {
			log.Fatal(err)
		}

		time.Sleep(time.Minute)
	}

	return nil
}

func (s *Scanner) populateSymbols(exchange string) error {
	defer timeTrack(time.Now(), "populateSymbols")

	symbols, err := s.stocks.GetAvailableSymbolsByExchange(exchange)
	if err != nil {
		return err
	}

	err = s.store.PopulateSymbols(symbols)
	if err != nil {
		return err
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

func (s *Scanner) updateIntradayData() error {
	defer timeTrack(time.Now(), "updateIntradayData")

	intradayPrices, err := s.stocks.GetAllRealtimePrices()
	if err != nil {
		return err
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
