/*
	Portfolio Management V0.1
	Perform analysis of portfolios composed of different combinations
	of securities.
*/

package portfolio

type StockScenario struct {
	StartDate  string
	EndDate    string
	Stocks     []*Stock
	PctHolding []float64
	Results    []ScenarioResults
}

type ScenarioResults struct {
	Date         string
	Shares       []float64
	StockHistIdx []int
	Value        float64
	ChangeValue  float64
}

type Stock struct {
	Ticker  string
	History []StockHistory
}

type StockHistory struct {
	Date     string
	Close    float64
	Dividend float64
}
