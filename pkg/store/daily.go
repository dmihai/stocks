package store

import (
	"fmt"

	"github.com/dmihai/stocks/pkg/types"
)

func (c *Conn) GetCandlesBetweenDates(days map[string]int) (map[string][]types.Candle, error) {
	result := make(map[string][]types.Candle)

	start := "99999"
	end := "0"
	for day := range days {
		if start > day {
			start = day
		}
		if end < day {
			end = day
		}
	}

	rows, err := c.db.Query("SELECT symbol, date, open, high, low, close, volume FROM daily WHERE date >= ? AND date <= ?", start, end)
	if err != nil {
		return nil, fmt.Errorf("query failed in getCandlesBetweenDates %s - %s: %v", start, end, err)
	}
	defer rows.Close()

	for rows.Next() {
		var candle types.Candle
		var symbol string
		var date string

		if err := rows.Scan(&symbol, &date, &candle.Open, &candle.High, &candle.Low, &candle.Close, &candle.Volume); err != nil {
			return nil, fmt.Errorf("scan failed in getCandlesBetweenDates %s - %s: %v", start, end, err)
		}

		if _, ok := result[symbol]; !ok {
			result[symbol] = make([]types.Candle, len(days))
		}

		result[symbol][days[date]] = candle
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row error in getCandlesBetweenDates %s - %s: %v", start, end, err)
	}
	return result, nil
}

func (c *Conn) GetDaysBetweenDates(start, end string) (map[string]int, error) {
	result := make(map[string]int)

	rows, err := c.db.Query("SELECT DISTINCT date FROM daily WHERE date >= ? AND date <= ? ORDER BY date ASC", start, end)
	if err != nil {
		return nil, fmt.Errorf("query failed in getDaysBetweenDates %s - %s: %v", start, end, err)
	}
	defer rows.Close()

	i := 0
	for rows.Next() {
		var day string

		if err := rows.Scan(&day); err != nil {
			return nil, fmt.Errorf("scan failed in getDaysBetweenDates %s - %s: %v", start, end, err)
		}

		result[day] = i
		i += 1
	}

	return result, nil
}
