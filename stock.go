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
// Current directory must contain three files:
//  - data/{Ticker}.csv  - daily history of stocks
// 		with minimum of "Date" and "Close" columns
//  - data/{Ticker}_div.csv - history of stock dividends
//		with minimum of "Date" and "Dividends" columns
//	- data/{Ticker}_distr.csv - history of capital gains distriubtions
//      with minimum of "Date" and "Distribution" columns
//
// Currently assumes that:
//  1. "Date" columns are the first (0 index) column in all three CSV files
//  2. "Close" is column index 4 in the daily close CSV
//  3. "Dividends" is column index 1 in the dividends CSV
//  4. "Distributions" is column index 1 in the distributions CSV
//
func (s *Stock) readHistory() error {
	if err := s.readCloseData(); err != nil {
		return err
	}

	if err := s.readDivData(); err != nil {
		return err
	}

	return s.readDistrData()
}

// readCloseData reads the stock market date
// and close amount.
func (s *Stock) readCloseData() error {
	// init date and close price for stock
	csv, err := readCSVFile("data/" + s.Ticker + ".csv")
	if err != nil {
		return err
	}

	dayCount := len(csv) - 1

	s.History = make([]StockHistory, dayCount)

	dateIdx := 0
	closeIdx := 4

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

	return nil
}

// readDivData reads dividend amounts and adds them
// to the existing close data.
func (s *Stock) readDivData() error {
	dateIdx := 0
	dividendsIdx := 1

	csv, err := readCSVFile("data/" + s.Ticker + "_div.csv")
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
					dividend, err := strconv.ParseFloat(day[dividendsIdx], 64)
					if err != nil {
						return fmt.Errorf("invalid float in stock dividends file for %s, line %d, %v",
							s.Ticker, i, err)
					}
					s.History[j].Dividend += dividend
					break
				}
			}
		}
	}

	return nil
}

// readDistrData reads capital gains distribution data
// and adds them to the existing close and dividends data.
// Note that if the daily history file does not contain the stock distribution file date
// the dividend will be shown on the following day (but shouldn't happen?)
func (s *Stock) readDistrData() error {
	dateIdx := 0
	distributionIdx := 1

	csv, err := readCSVFile("data/" + s.Ticker + "_distr.csv")
	if err != nil {
		return err
	}

	for i, day := range csv {
		if i == 0 {
			if day[dateIdx] != "Date" {
				return fmt.Errorf("distribution file for %s does not start with 'Date'", s.Ticker)
			}
			if day[distributionIdx] != "Distributions" {
				return fmt.Errorf("distributions file for %s does not have 'Distributions' in expected column", s.Ticker)
			}
		} else {
			for j, history := range s.History {
				if history.Date >= day[dateIdx] {
					distribution, err := strconv.ParseFloat(day[distributionIdx], 64)
					if err != nil {
						return fmt.Errorf("invalid float in stock distributions file for %s, line %d, %v",
							s.Ticker, i, err)
					}
					s.History[j].Distribution += distribution
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

// getNextDate returns the next stock history date
// which is greater than lastDate.
// Starts searching stock history at beginIdx.
// If no next history date available, returns the constant MaxDate
func (s *Stock) getNextDate(lastDate string, beginIdx int) string {

	for i := beginIdx; i < len(s.History); i++ {
		if s.History[i].Date > lastDate {
			return s.History[i].Date
		}
	}
	return MaxDate
}

// getCloseDateIdx returns the index of the stock close date
// which is on or before the specified date.
// Starts search in stock history at beginIdx.
func (s *Stock) getCloseDateIdx(closeDate string, beginIdx int) int {

	result := beginIdx
	for i := beginIdx + 1; i < len(s.History); i++ {
		if s.History[i].Date <= closeDate {
			result = i
		} else {
			return result
		}
	}

	if result != beginIdx {
		return result
	}

	return -1
}
