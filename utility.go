package portfolio

import (
	"encoding/csv"
	"os"
)

// Read a CSV array from a file
// where the first row contains the column names
// and subsequent columns contain column values in string format
func readCSVFile(file string) ([][]string, error) {

	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}

	records, err := csv.NewReader(f).ReadAll()
	f.Close()

	return records, err
}
