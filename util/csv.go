package util

import (
    "encoding/csv"
    "fmt"
    "io"
    "strconv"
    "strings"
)

// ParseCSVToFloat parses a csv with numerical values and returns a slice
// containing the numerical values, including the rows and cols of the file.
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
        } else if (cols == 0) {
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

// appendFloat converts the record elements to float and appends
// to the given slice
func appendFloat(data []float64, record[]string) ([]float64, error) {
    for _, elem := range record {
        f, err := strconv.ParseFloat(strings.TrimSpace(elem), 64)
        if err != nil {
            return data, err
        }
        data = append(data, f)
    }
    return data, nil
}
