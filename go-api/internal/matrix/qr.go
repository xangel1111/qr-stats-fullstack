package matrix

import "gonum.org/v1/gonum/mat"

// Factorization holds the reduced (thin) QR factors of a matrix.
type Factorization struct {
	Q [][]float64 // m x n, orthonormal columns
	R [][]float64 // n x n, upper triangular
}

// Factorize computes the reduced QR factorization of the given m x n matrix
// (with m >= n) such that A = Q * R.
//
// We return the reduced form on purpose: Q is m x n and R is the n x n upper
// triangular block. Keeping R square is what makes the downstream "is diagonal"
// check semantically meaningful.
//
// Callers MUST call Validate first; Factorize assumes a well-formed input.
func Factorize(data [][]float64) Factorization {
	rows := len(data)
	cols := len(data[0])

	// Flatten into gonum's row-major dense representation.
	flat := make([]float64, 0, rows*cols)
	for _, row := range data {
		flat = append(flat, row...)
	}
	a := mat.NewDense(rows, cols, flat)

	var qr mat.QR
	qr.Factorize(a)

	// gonum returns the full factors: Q is m x m and R is m x n.
	var qFull, rFull mat.Dense
	qr.QTo(&qFull)
	qr.RTo(&rFull)

	// Slice down to the reduced/thin factors.
	q := qFull.Slice(0, rows, 0, cols) // m x n
	r := rFull.Slice(0, cols, 0, cols) // n x n

	return Factorization{
		Q: toSlice(q),
		R: toSlice(r),
	}
}

// toSlice converts any gonum matrix (including a sliced view) into a plain
// [][]float64 for JSON transport.
func toSlice(m mat.Matrix) [][]float64 {
	rows, cols := m.Dims()
	out := make([][]float64, rows)
	for i := 0; i < rows; i++ {
		out[i] = make([]float64, cols)
		for j := 0; j < cols; j++ {
			out[i][j] = m.At(i, j)
		}
	}
	return out
}
