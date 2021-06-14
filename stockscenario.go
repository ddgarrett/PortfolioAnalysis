package portfolio

import (
	"errors"
	"fmt"
	"time"
)

// NewStockScenario creates a new instance of a Stock Scenario.
// The new scenario will not yet contain any stock allocations.
//
// Note that it is important that StartDate and EndDate be set before adding
// stock history since they will be adjusted accordingly if any of the
// stock history start or end outside the specified date range.
func NewStockScenario(startDate, endDate string) *StockScenario {
	return &StockScenario{StartDate: startDate, EndDate: endDate}
}

// AddStock adds a stock to the scenario
func (sc *StockScenario) AddStock(stock *Stock, pct float64) error {

	if pct > 1 {
		return errors.New("pct greater than 1")
	}

	if pct <= 0 {
		return errors.New("pct less than 0")
	}

	sc.Stocks = append(sc.Stocks, stock)
	sc.PctHolding = append(sc.PctHolding, pct)

	if sc.StartDate < stock.History[0].Date {
		sc.StartDate = stock.History[0].Date
	}

	lastHistory := len(stock.History) - 1
	if sc.EndDate > stock.History[lastHistory].Date {
		sc.EndDate = stock.History[lastHistory].Date
	}

	return nil
}

// Run runs the defined stock scenario starting with
// an initial amount of dollars.
func (sc *StockScenario) Run(initialAmount float64) error {

	if err := sc.initResults(); err != nil {
		return err
	}

	sc.genFirstResult(initialAmount)

	return nil
}

// initialize the results for a stock scenario run
//
// TODO: go through all of the stock history again to adjust
// start and end dates in case they are outside the range
// of stock history? This might happen if the date range is
// set after the stocks are added to the Stock Scenario.
//
// TODO: verify that percents add up to approx.  1
func (sc *StockScenario) initResults() error {

	timeFormat := "2006-01-02"

	start, err := time.Parse(timeFormat, sc.StartDate)
	if err != nil {
		return err
	}

	end, err := time.Parse(timeFormat, sc.EndDate)
	if err != nil {
		return err
	}

	if sc.StartDate >= sc.EndDate {
		return fmt.Errorf("StartDate '%s' not less than EndDate '%s'", sc.StartDate, sc.EndDate)
	}

	duration := end.Sub(start).Hours()/24 + 1

	sc.Results = make([]ScenarioResults, 0, int(duration))

	return nil
}

func (sc *StockScenario) genFirstResult(amt float64) {

	results := &ScenarioResults{Date: sc.StartDate, Value: amt}
	results.initHistIdx(sc)
	results.rebalanceStocks(sc)

}
