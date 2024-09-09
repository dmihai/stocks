package store

import (
	"database/sql"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/dmihai/stocks/pkg/data"
)

func (c *Conn) GetIntradayCandles(day string) (map[string][]data.Candle, error) {
	result := make(map[string][]data.Candle)

	table := "intraday_" + strings.Replace(day, "-", "", -1)

	minTime, err := c.getMinTimestampForTable(table)
	if err != nil {
		return nil, err
	}

	length := 10 * 60
	maxTime := minTime.Add(time.Minute * 10 * 60)

	rows, err := c.db.Query("SELECT symbol, timestamp, open, high, low, close, volume FROM "+table+" WHERE timestamp < ?", maxTime)
	if err != nil {
		return nil, fmt.Errorf("query failed in getIntradayCandles %s: %v", day, err)
	}
	defer rows.Close()

	for rows.Next() {
		var candle data.Candle
		var symbol string
		var timestamp time.Time

		if err := rows.Scan(&symbol, &timestamp, &candle.Open, &candle.High, &candle.Low, &candle.Close, &candle.Volume); err != nil {
			return nil, fmt.Errorf("scan failed in getIntradayCandles %s: %v", day, err)
		}

		if _, ok := result[symbol]; !ok {
			result[symbol] = make([]data.Candle, length)
		}

		index := int(math.Round(timestamp.Sub(*minTime).Minutes()))
		result[symbol][index] = candle
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row error in getIntradayCandles %s: %v", day, err)
	}
	return result, nil
}

func (c *Conn) getMinTimestampForTable(table string) (*time.Time, error) {
	var minTime time.Time

	row := c.db.QueryRow("SELECT MIN(timestamp) FROM " + table)
	if err := row.Scan(&minTime); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("getMinTimestampForTable %s: no result", table)
		}
		return nil, fmt.Errorf("getMinTimestampForTable %s: %v", table, err)
	}

	return &minTime, nil
}
