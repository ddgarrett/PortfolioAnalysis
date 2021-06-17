package portfolio

import (
	"testing"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
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

func TestRun_Part1(t *testing.T) {

	agg, err := NewStock("AGG")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	sc := NewStockScenario("1900-01-01", "3000-01-01")

	if err = sc.AddStock(agg, 1); err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	/*
		if err = sc.RunScenario(10000); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	*/

	sc = NewStockScenario("202x-01-01", "2020-12-31")
	sc.AddStock(agg, 1)
	if err = sc.RunScenario(10000); err == nil {
		t.Error("didn't catch invalid start date error")
	}

	// NOTE: invalid end date must still be < "2021-06-08"
	// the highest date in AGG stock history, or it will be changed to that date.
	sc = NewStockScenario("2020-01-01", "2020x-12-31")
	sc.AddStock(agg, 1)
	if err = sc.RunScenario(10000); err == nil {
		t.Error("didn't catch invalid end date error")
	}

	sc = NewStockScenario("2021-01-01", "2020-12-31")
	sc.AddStock(agg, 1)
	if err = sc.RunScenario(10000); err == nil {
		t.Error("didn't catch start date after end date error")
	}

	sc = NewStockScenario("2020-01-01", "2020-12-31")
	sc.AddStock(agg, 1)
	if err = sc.RunScenario(10000); err != nil {
		t.Errorf("unexpected .Run error: %v", err)
	}

	if len(sc.Results) != 254 {
		t.Errorf("invalid .Results len: %d", len(sc.Results))
	}

	if cap(sc.Results) != 366 {
		t.Errorf("invalid .Results capacity: %d", cap(sc.Results))
	}
}

func ExampleRunScenario() {
	fxaix, _ := NewStock("FXAIX")

	sc := NewStockScenario("2020-01-01", "2020-12-31")
	sc.AddStock(fxaix, 1)
	sc.RunScenario(10000)

	// fmt.Println(sc.String())
	// t.Errorf("generate error to show output")
	p := message.NewPrinter(language.English)
	p.Printf("Stock %s, StartDate: %s, EndDate: %s\n", "FXAIX", sc.StartDate, sc.EndDate)

	p.Printf("\tStartAmt: %.2f EndAmt: %.2f PctChange: %.4f%%\n",
		sc.StartAmt, sc.EndAmt, sc.PctChange*100)

	// output: Stock FXAIX, StartDate: 2020-01-01, EndDate: 2020-12-31
	//	StartAmt: 10,000.00 EndAmt: 11,839.52 PctChange: 18.3952%
}
