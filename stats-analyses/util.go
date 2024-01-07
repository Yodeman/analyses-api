package statsanal

import (
	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/gonum/stat"
)

// ones generates a slice of floats filled with 1
func ones(rows, cols int) mat.Matrix {
	f := make([]float64, rows*cols)
	for i := range f {
		f[i] = float64(1)
	}

	return mat.NewDense(rows, cols, f)
}

// a helper sum function for slice of floats
func sum(elem ...float64) float64 {
	var s float64

	for _, e := range elem {
		s += e
	}

	return s
}

// mean finds the mean of each column of matrix `m`
func mean(m mat.Matrix) []float64 {
	var means []float64

	_, c := m.Dims()

	for i := 0; i < c; i++ {
		col := mat.Col(nil, i, m)
		means = append(means, stat.Mean(col, nil))
	}

	return means
}

// varianc finds the variance of each column of matrix `m`
func variance(m mat.Matrix) []float64 {
	var v []float64

	_, c := m.Dims()

	for i := 0; i < c; i++ {
		col := mat.Col(nil, i, m)
		v = append(v, stat.Variance(col, nil))
	}

	return v
}
