package portfolio

import (
	"bytes"
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

	sc.StartAmt = initialAmount

	if err := sc.initResults(); err != nil {
		return err
	}

	sc.genFirstResult(initialAmount)

	date := sc.getNextResultsDate()
	for ; date <= sc.EndDate; date = sc.getNextResultsDate() {
		sr := sc.generateDaysResults(date)
		if sc.needRebalance() {
			sr.rebalanceStocks(sc)
			// sc.printScenarioResults()
		}
	}

	lastResult := sc.getLastResults()
	sc.EndAmt = lastResult.Value
	sc.PctChange = sc.EndAmt/sc.StartAmt - 1
	return nil
}

func (sc *StockScenario) String() string {
	var b bytes.Buffer
	fmt.Fprintf(&b, "%s to %s, %d stocks, %d results\n", sc.StartDate, sc.EndDate, len(sc.Stocks), len(sc.Results))
	fmt.Fprintf(&b, "Stocks: \n")
	for i, stock := range sc.Stocks {
		lastHistoryIdx := len(stock.History) - 1
		firstDate := stock.History[0].Date
		lastDate := stock.History[lastHistoryIdx].Date
		fmt.Fprintf(&b, "%s %f%% - %d history, from %s to %s \n",
			stock.Ticker, sc.PctHolding[i], len(stock.History), firstDate, lastDate)
	}
	fmt.Fprintf(&b, "\n")

	fmt.Fprintf(&b, "**** Results: \n")
	for i, result := range sc.Results {
		fmt.Fprintf(&b, "\t %d: %s", i, result.String())
	}
	fmt.Fprintf(&b, "\n")
	fmt.Fprintf(&b, "\n")
	return b.String()
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

// genFirstResult generates the first ScenarioResults entry
func (sc *StockScenario) genFirstResult(amt float64) {

	results := &ScenarioResults{Date: sc.StartDate, Value: amt}
	results.initHistIdx(sc)
	results.rebalanceStocks(sc)
	sc.Results = append(sc.Results, *results)

}

// getNextDate returns the next date for which results can be calculated.
// If no next date, returns MaxDate.
func (sc *StockScenario) getNextResultsDate() string {
	results := MaxDate

	lastResults := sc.getLastResults()
	if lastResults == nil {
		panic("lastResults nil")
	}

	for i, stock := range sc.Stocks {
		date := stock.getNextDate(lastResults.Date, lastResults.StockHistIdx[i])
		if date < results {
			results = date
		}
	}

	return results
}

// generateDaysResults generates the results for a specified day.
func (sc *StockScenario) generateDaysResults(date string) *ScenarioResults {
	results := &ScenarioResults{}
	results.initNextResults(date, sc.getLastResults(), sc.Stocks)
	sc.Results = append(sc.Results, *results)
	return results
}

// needRebalance returns true if the last days results need to be rebalanced.
// Currently rebalances on the first trading day after the 15th.
// "yyyy-mm-dd"
func (sc *StockScenario) needRebalance() bool {

	lr := sc.getLastResults()
	if lr == nil {
		// no entries to rebalance yet - shouldn't really get here
		return false
	}

	if lr.Date[8:] >= "15" {

		pr := sc.getPrevResults()
		if pr == nil {
			// no previous results - don't rebalance
			return false
		}

		if pr.Date[8:] >= "15" {
			// already rebalanced during a prior day in this month
			return false
		}

		return true
	}

	return false
}

// getLastResults returns the last entry from the Results slice
func (sc *StockScenario) getLastResults() *ScenarioResults {
	if len(sc.Results) == 0 {
		// no results yet - return nil
		return nil
	}

	results := sc.Results[len(sc.Results)-1]
	return &results
}

// getPrevResults returns the next to last results
func (sc *StockScenario) getPrevResults() *ScenarioResults {
	i := len(sc.Results)
	if i < 2 {
		return nil
	}

	return &sc.Results[i-2]
}
