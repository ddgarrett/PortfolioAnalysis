package portfolio

import (
	"fmt"
	"strconv"
)

// NewStock returns pointer to a new Stock structure
// for a given stock ticker. Assumes data, both daily close
// and dividend files, are in the "data/" directory.
func NewStock(ticker string) (*Stock, error) {
	result := Stock{Ticker: ticker}
	err := result.readHistory()
	return &result, err
}

// ReadHistory reads CSV files to load history for a stock.
// The "Ticker" value must already be set.
//
// Current directory must contain two files:
//  - data/{Ticker}.csv  - daily history of stocks with minimum of "Date" and "Close" columns
//  - data/{Ticker}_div.csv - daily history of stock dividends with minimum of "Date" and "Dividends" columns
//
// Currently assumes that:
//  1. "Date" columns are the first (0 index) column in both CSV files
//  2. "Close" is column index 4 in the daily close CSV
//  3. "Dividends" is column index 1 in the dividends CSV
//
func (s *Stock) readHistory() error {
	var err error

	// init date and close price for stock
	csv, err := readCSVFile("data/" + s.Ticker + ".csv")
	if err != nil {
		return err
	}

	dayCount := len(csv) - 1

	s.History = make([]StockHistory, dayCount)

	dateIdx := 0
	closeIdx := 4
	dividendsIdx := 1

	for i, day := range csv {
		if i == 0 {
			if day[dateIdx] != "Date" {
				return fmt.Errorf("daily close file for %s does not start with 'Date'", s.Ticker)
			}
			if day[closeIdx] != "Close" {
				return fmt.Errorf("daily close file for %s does not have 'Close' in expected column", s.Ticker)
			}
		} else {
			s.History[i-1].Date = day[dateIdx]
			s.History[i-1].Close, err = strconv.ParseFloat(day[closeIdx], 64)
			if err != nil {
				return fmt.Errorf("invalid float in stock history file for %s, line %d, %v",
					s.Ticker, i, err)
			}
		}
	}

	// add dividends to daily history
	// Note that if the daily history file does not contain the stock dividend file date
	// the dividend will be shown on the following day (but shouldn't happen?)
	csv, err = readCSVFile("data/" + s.Ticker + "_div.csv")
	if err != nil {
		return err
	}

	for i, day := range csv {
		if i == 0 {
			if day[dateIdx] != "Date" {
				return fmt.Errorf("dividend file for %s does not start with 'Date'", s.Ticker)
			}
			if day[dividendsIdx] != "Dividends" {
				return fmt.Errorf("dividend file for %s does not have 'Dividends' in expected column", s.Ticker)
			}
		} else {
			for j, history := range s.History {
				if history.Date >= day[dateIdx] {
					s.History[j].Dividend, err = strconv.ParseFloat(day[dividendsIdx], 64)
					if err != nil {
						return fmt.Errorf("invalid float in stock dividends file for %s, line %d, %v",
							s.Ticker, i, err)
					}
					break
				}
			}
		}
	}
	return nil
}

// getHistIdx gets the index of the stock history entry
// where the date is <= a given date.
func (s *Stock) getHistIdx(date string, startIdx int) int {

	result := startIdx

	for i := startIdx + 1; i < len(s.History); i++ {
		if s.History[i].Date <= date {
			result = i
		} else {
			break
		}
	}
	return result
}
