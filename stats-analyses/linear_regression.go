package statsanal

import (
	"fmt"
	"math"

	"gonum.org/v1/gonum/mat"
)

// LinearRegression computes the statistical multivariable linear regression
// on the given matrix `m`, using the explanation found in:
//
//	https://developer.ibm.com/articles/linear-regression-from-scratch/
//
// The last column of the matrix is used as the target `Y` and the rest of the
// columns is taken as the predictor `X`. Successful computation returns the
// coefficients of regression and the t-test statistics for each columns, with
// the first element being the intercept|bias.
//
// Returns a non-nil error if an error occured during computation.
func LinearRegression(m *mat.Dense) (coeffs, tstat string, err error) {
	// Calculate the regression coefficients.
	var x, inv, xTransDotx, invDotxTrans, res mat.Dense

	r, c := m.Dims()
	x.Stack(ones(1, r), m.Slice(0, r, 0, c-1).T())
	X := x.T()
	Y := m.Slice(0, r, c-1, c)

	xTransDotx.Mul(X.T(), X)
	err = inv.Inverse(&xTransDotx)
	if err != nil {
		err = fmt.Errorf("Error calculating regression coefficients.\n%v\n", err)
		return
	}
	invDotxTrans.Mul(&inv, X.T())
	res.Mul(&invDotxTrans, Y)

	fr := mat.Formatted(&res, mat.FormatPython()) //, mat.Prefix("    "), mat.Squeeze())
	coeffs = fmt.Sprintf("%.5f", fr)

	// Calculate t-statistics
	var yHat, residual, residualSquare, coeffVariance mat.Dense
	var coeffVarianceRoot, coeffVarianceDiag, tst mat.Dense

	yHat.Mul(X, &res)
	residual.Sub(Y, &yHat)
	residualSquare.Apply(
		func(i, j int, elem float64) float64 {
			return math.Pow(elem, 2)
		},
		&residual,
	)

	sigmaHat := mat.Sum(&residualSquare) / float64(r-c-2)

	coeffVariance.Apply(
		func(i, j int, elem float64) float64 {
			return sigmaHat * elem
		},
		&inv,
	)

	diag := coeffVariance.DiagView()
	dr, _ := diag.Dims()
	coeffVarianceDiag.Mul(diag, ones(dr, 1))
	coeffVarianceRoot.Apply(
		func(i, j int, elem float64) float64 {
			return math.Sqrt(elem)
		},
		&coeffVarianceDiag,
	)

	tst.DivElem(&res, &coeffVarianceRoot)

	ft := mat.Formatted(&tst, mat.FormatPython()) //, mat.Prefix("    "), mat.Squeeze())
	tstat = fmt.Sprintf("%.5f", ft)

	return
}
