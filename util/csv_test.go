package util

import (
    "strings"
    "testing"

    "github.com/stretchr/testify/require"
)

func TestParseCSVToFloatSlice(t *testing.T) {
    // sample_text := "1,2,3,4,5,6\n"+"7,8,9,10,11,12\n"+"13,14,15,16,17,18\n"
    r, c := 10, 10
    sample_text := RandomCSV(r, c)

    reader := strings.NewReader(sample_text)

    rows, cols, data, err := ParseCSVToFloatSlice(reader)
    require.NoError(t, err)
    require.Equal(t, r, rows)
    require.Equal(t, c, cols)
    require.Equal(t, len(data), rows*cols)
}
