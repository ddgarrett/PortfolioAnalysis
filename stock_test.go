package portfolio

import (
	"testing"
)

func TestNewStock(t *testing.T) {
	result, err := NewStock("AGG")

	if err != nil {
		t.Errorf("error creating new stock for AGG: %v", result)
	}

	if result.Ticker != "AGG" {
		t.Error("NewStock ticker != 'AGG'")
	}

	if len(result.History) != 4454 {
		t.Error("NewStock did not read 4,454 dates")
	}

	for _, history := range result.History {
		if history.Dividend != 0 {
			if history.Date != "2003-11-03" {
				t.Errorf("first AGG stock dividend on %s instead of 2003-11-03", history.Date)
			}
			break
		}
	}

}
