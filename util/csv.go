package util

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// ParseCSVToFloat parses a csv file containing numerical values in the supplied
// reader. Returns a slice of floats containing the numerical values. The data are
// stored in row major order. The numbers of rows and columns of the file is also
// returned.
//
// An error, empty slice,  zero rows and columns are returned if an error occured
// while parsing the csv string contained in the reader.
func ParseCSVToFloatSlice(r io.Reader) (rows, cols int, data []float64, err error) {
	reader := csv.NewReader(r)

	var record []string

	for {
		record, err = reader.Read()

		if err == io.EOF {
			err = nil
			return
		}

		if len(record) == 0 {
			continue
		}

		if (err != nil) && (err != io.EOF) {
			return 0, 0, []float64{}, fmt.Errorf("Error parsing file.\n%v\n", err)
		}

		if (cols != 0) && (len(record) != cols) {
			err = fmt.Errorf("All rows should have same length!!!")
			return 0, 0, []float64{}, err
		} else if cols == 0 {
			cols = len(record)
		}

		data, err = appendFloat(data, record)
		if err != nil {
			return 0, 0, []float64{}, err
		}

		rows += 1
	}
	return
}

// appendFloat converts the record/row elements to float and appends
// to the given slice. A slice is returned which might have underlying memory
// as the passed slice if the passed slice has enough capacity for new elements.
//
// An error is returned if the conversion fails.
func appendFloat(data []float64, record []string) ([]float64, error) {
	for _, elem := range record {
		f, err := strconv.ParseFloat(strings.TrimSpace(elem), 64)
		if err != nil {
			return data, err
		}
		data = append(data, f)
	}
	return data, nil
}
