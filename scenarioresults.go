package portfolio

import (
	"bytes"
	"fmt"
	"math"
)

func (sr *ScenarioResults) String() string {
	var b bytes.Buffer
	fmt.Fprintf(&b, "%s  %v  %v  %.2f  %.2f\n", sr.Date, sr.Shares, sr.StockHistIdx, sr.Value, sr.ChangeValue)
	return b.String()
}

// Calculate the StockHistIdx for each stock
// such that the StockHistory.Date is <= sr.Date.
// sr.Date must be set before calling this method.
func (sr *ScenarioResults) initHistIdx(sc *StockScenario) {
	sr.StockHistIdx = make([]int, len(sc.Stocks))

	sr.Shares = make([]float64, len(sc.Stocks))

	for i, stock := range sc.Stocks {
		sr.StockHistIdx[i] = stock.getHistIdx(sr.Date, 0)
	}
}

// initNextResults initializes a new ScenarioResults struct.
func (sr *ScenarioResults) initNextResults(date string, prevSR *ScenarioResults, stocks []*Stock) {
	sr.Date = date

	sr.Shares = make([]float64, len(prevSR.Shares))
	copy(sr.Shares, prevSR.Shares)

	sr.StockHistIdx = make([]int, len(sr.Shares))

	for i, stock := range stocks {
		lastIdx := prevSR.StockHistIdx[i]
		closeIdx := stock.getCloseDateIdx(date, lastIdx)
		sr.StockHistIdx[i] = closeIdx

		shares := sr.Shares[i]
		close := stock.History[closeIdx].Close
		dividend := stock.History[closeIdx].Dividend
		dividend += stock.History[closeIdx].Distribution

		if dividend != 0 {
			// round to dividend amount to nearest cent
			// use roundToEven to eliminate bias for .5 cents
			dividendTotal := math.RoundToEven(shares * dividend * 100)
			dividendTotal = dividendTotal / 100

			// calc new shares
			newShares := dividendTotal / close

			// round new shares to 3 decimal points, again using roundToEven
			newShares = math.RoundToEven(newShares * 1000)
			newShares = newShares / 1000

			// add new shares to holdings
			shares += newShares
			sr.Shares[i] = shares
		}

		sr.Value += (close * shares)
	}

	sr.ChangeValue = sr.Value - prevSR.Value
	sr.PctChange = sr.ChangeValue / prevSR.Value
}

// Buy/Sell stocks to rebalance the stock portfolio to the scenario defined percents.
func (sr *ScenarioResults) rebalanceStocks(sc *StockScenario) {

	for i, stock := range sc.Stocks {
		histIdx := sr.StockHistIdx[i]
		close := stock.History[histIdx].Close
		pct := sc.PctHolding[i]

		stkValue := sr.Value * pct
		shares := stkValue / close

		// round number of shares to 3 decimal places
		shares = math.RoundToEven(shares * 1000)
		shares = shares / 1000

		sr.Shares[i] = shares
	}
}
