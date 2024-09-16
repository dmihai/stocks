package store

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/dmihai/stocks/pkg/data"
)

func (c *Conn) GetIntradayCandles(day string, maxTime time.Time) ([]data.Intraday, error) {
	result := make([]data.Intraday, 0)

	table := getTableForDay(day)

	rows, err := c.db.Query("SELECT symbol, timestamp, close, volume FROM "+table+" WHERE timestamp < ? ORDER BY symbol, timestamp", maxTime)
	if err != nil {
		return nil, fmt.Errorf("query failed in GetIntradayCandles %s max %v: %v", day, maxTime, err)
	}
	defer rows.Close()

	for rows.Next() {
		var price data.Intraday

		if err := rows.Scan(&price.Symbol, &price.Timestamp, &price.Price.Price, &price.Volume); err != nil {
			return nil, fmt.Errorf("scan failed in GetIntradayCandles %s max %v: %v", day, maxTime, err)
		}

		result = append(result, price)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row error in GetIntradayCandles %s max %v: %v", day, maxTime, err)
	}

	return result, nil
}

func (c *Conn) GetMinTimestampForDay(day string) (*time.Time, error) {
	var minTime time.Time

	table := getTableForDay(day)

	row := c.db.QueryRow("SELECT MIN(timestamp) FROM " + table)
	if err := row.Scan(&minTime); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("GetMinTimestampForDay %s: no result", day)
		}
		return nil, fmt.Errorf("GetMinTimestampForDay %s: %v", day, err)
	}

	return &minTime, nil
}

func getTableForDay(day string) string {
	return "intraday_" + strings.Replace(day, "-", "", -1)
}
