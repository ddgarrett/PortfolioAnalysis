/*
	Portfolio Management V0.1
	Perform analysis of portfolios composed of different combinations
	of securities.
*/

package portfolio

// StockScenario defines a scenario for a set of securities and timeframe.
// Each stock is assigned a given percent of the portfolio. The stock is
// rebalanced at specific times. Currently rebalance is the 15th of the month
// but this may change in the future.
type StockScenario struct {
	StartDate  string
	EndDate    string
	Stocks     []*Stock
	PctHolding []float64
	Results    []ScenarioResults
}

// Daily results of the portfolio value.
type ScenarioResults struct {
	Date         string
	Shares       []float64
	StockHistIdx []int
	Value        float64
	ChangeValue  float64
}

// Stock information, ticker and history.
type Stock struct {
	Ticker  string
	History []StockHistory
}

// Stock history
type StockHistory struct {
	Date     string
	Close    float64
	Dividend float64
}

const MaxDate = "4000-01-01"
