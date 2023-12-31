package util

import (
    "strings"
    "testing"

    "github.com/stretchr/testify/require"
)

func TestParseCSVToFloatSlice(t *testing.T) {
    sample_text := "1,2,3,4,5,6\n"+"7,8,9,10,11,12\n"+"13,14,15,16,17,18\n"

    r := strings.NewReader(sample_text)

    rows, cols, data, err := ParseCSVToFloatSlice(r)
    require.NoError(t, err)
    require.Equal(t, 3, rows)
    require.Equal(t, 6, cols)
    require.Equal(t, len(data), rows*cols)
    require.Equal(t, []float64{1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18}, data)
}
