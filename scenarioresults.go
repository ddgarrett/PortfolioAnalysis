package portfolio

import "log"

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

	if len(sr.Shares) != len(prevSR.Shares) {
		log.Fatal("misunderstood copy")
	}

	sr.StockHistIdx = make([]int, len(sr.Shares))

	for i, stock := range stocks {
		lastIdx := prevSR.StockHistIdx[i]
		closeIdx := stock.getCloseDateIdx(date, lastIdx)
		sr.StockHistIdx[i] = closeIdx

		shares := sr.Shares[i]
		close := stock.History[closeIdx].Close
		dividend := stock.History[closeIdx].Dividend

		sr.Value += ((close + dividend) * shares)
	}

	sr.ChangeValue = sr.Value - prevSR.Value
}

// Buy/Sell stocks to rebalance the stock portfolio to the scenario defined percents.
func (sr *ScenarioResults) rebalanceStocks(sc *StockScenario) {

	for i, stock := range sc.Stocks {
		histIdx := sr.StockHistIdx[i]
		close := stock.History[histIdx].Close
		pct := sc.PctHolding[i]

		stkValue := sr.Value * pct
		sr.Shares[i] = stkValue / close
	}
}
