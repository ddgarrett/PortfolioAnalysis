package portfolio

// Calculate the StockHistIdx for each stock
// such that the StockHistory.Date is <= sr.Date.
// sr.Date must be set before calling this method.
func (sr *ScenarioResults) initHistIdx(sc *StockScenario) {
	sr.StockHistIdx = make([]int, len(sc.Stocks))

	for i, stock := range sc.Stocks {
		sr.StockHistIdx[i] = stock.getHistIdx(sr.Date, 0)
	}
}

// Buy/Sell stocks to rebalance the stock portfolio to the scenario defined percents.
func (sr *ScenarioResults) rebalanceStocks(sc *StockScenario) {

}
