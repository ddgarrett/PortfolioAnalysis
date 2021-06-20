package portfolio

import (
	"fmt"
	"testing"
)

func TestAddStock(t *testing.T) {
	sc := NewStockScenario("1900-01-01", "3000-01-01")

	agg, err := NewStock("AGG")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if err = sc.AddStock(agg, 2); err == nil {
		t.Error("missed error pct > 1")
	}

	if err = sc.AddStock(agg, 0); err == nil {
		t.Error("missed error pct <= 0")
	}

	if err = sc.AddStock(agg, 1); err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if sc.StartDate != "2003-09-29" || sc.EndDate != "2021-06-08" {
		t.Errorf("start/end dates not set to stock start/end dates: %s %s", sc.StartDate, sc.EndDate)
	}

	if len(sc.Stocks) != 1 || sc.Stocks[0] != agg {
		t.Error("stock agg not added correctly")
	}
	if len(sc.PctHolding) != 1 || sc.PctHolding[0] != 1 {
		t.Error("pct holding not set correctly")
	}

}

func TestCalcResults(t *testing.T) {

	agg, err := NewStock("AGG")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	sc := NewStockScenario("1900-01-01", "3000-01-01")

	if err = sc.AddStock(agg, 1); err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if err = sc.CalcResults(10000); err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	sc = NewStockScenario("202x-01-01", "2020-12-31")
	sc.AddStock(agg, 1)
	if err = sc.CalcResults(10000); err == nil {
		t.Error("didn't catch invalid start date error")
	}

	// NOTE: invalid end date must still be < "2021-06-08"
	// the highest date in AGG stock history, or it will be changed to that date.
	sc = NewStockScenario("2020-01-01", "2020x-12-31")
	sc.AddStock(agg, 1)
	if err = sc.CalcResults(10000); err == nil {
		t.Error("didn't catch invalid end date error")
	}

	sc = NewStockScenario("2021-01-01", "2020-12-31")
	sc.AddStock(agg, 1)
	if err = sc.CalcResults(10000); err == nil {
		t.Error("didn't catch start date after end date error")
	}

	sc = NewStockScenario("2020-01-01", "2020-12-31")
	sc.AddStock(agg, 1)
	if err = sc.CalcResults(10000); err != nil {
		t.Errorf("unexpected .Run error: %v", err)
	}

	if len(sc.Results) != 254 {
		t.Errorf("invalid .Results len: %d", len(sc.Results))
	}

	if cap(sc.Results) != 366 {
		t.Errorf("invalid .Results capacity: %d", cap(sc.Results))
	}
}

func TestString(t *testing.T) {
	fxaix, _ := NewStock("FXAIX")
	sc := NewStockScenario("2020-01-01", "2021-01-01")
	sc.AddStock(fxaix, 1)
	sc.CalcResults(10000)

	fmt.Println(len(sc.String()))

	// output: 16902
	if len(sc.String()) != 16902 {
		t.Errorf("invalid .Results capacity: %d", cap(sc.Results))
	}
}

func TestCalcResults_Part02(t *testing.T) {

	stockTickers := []string{"FXAIX", "FXNAX", "VDADX"}
	years := []string{"2016", "2017", "2018", "2019", "2020"}

	expectedResult := [][]string{
		{"11.97", "21.81", "-4.40", "31.47", "18.40"},
		{"2.51", "3.49", "0.03", "8.48", "7.80"},
		{"11.79", "12.22", "-2.03", "29.68", "15.46"},
	}

	stocks := []*Stock{}

	// read stock history files
	for _, stockTicker := range stockTickers {
		stock, err := NewStock(stockTicker)
		if err != nil {
			t.Errorf("unexpected error reading stock: %v", err)
		}
		stocks = append(stocks, stock)

	}

	for i, year := range years {
		for j, stock := range stocks {
			sc := NewStockScenario(year+"-01-01", year+"-12-31")
			sc.AddStock(stock, 1)
			sc.CalcResults(10000)

			pctChg := fmt.Sprintf("%.2f", sc.PctChange*100.0)
			expectPctChg := expectedResult[j][i]

			if pctChg != expectPctChg {
				// show 4 decimal digits on actual pct change
				pctChg = fmt.Sprintf("%.4f", sc.PctChange*100.0)
				t.Errorf("stock %s, year %s expected %s%% change, got %s%% change",
					stock.Ticker, year, expectPctChg, pctChg)
			}
		}
	}

}
