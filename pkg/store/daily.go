package store

import (
	"fmt"

	"github.com/dmihai/stocks/pkg/data"
)

func (c *Conn) GetCandlesBetweenDates(start, end string) ([]data.Daily, error) {
	result := make([]data.Daily, 0)

	rows, err := c.db.Query("SELECT symbol, date, open, high, low, close, volume FROM daily WHERE date >= ? AND date <= ?", start, end)
	if err != nil {
		return nil, fmt.Errorf("query failed in getCandlesBetweenDates %s - %s: %v", start, end, err)
	}
	defer rows.Close()

	for rows.Next() {
		var daily data.Daily

		if err := rows.Scan(&daily.Symbol, &daily.Day, &daily.Open, &daily.High, &daily.Low, &daily.Close, &daily.Volume); err != nil {
			return nil, fmt.Errorf("scan failed in getCandlesBetweenDates %s - %s: %v", start, end, err)
		}

		result = append(result, daily)
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
