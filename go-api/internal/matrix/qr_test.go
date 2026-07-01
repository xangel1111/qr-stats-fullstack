package matrix

import (
	"math"
	"testing"
)

const tol = 1e-9

// multiply returns a * b for plain slice matrices.
func multiply(a, b [][]float64) [][]float64 {
	rows := len(a)
	inner := len(b)
	cols := len(b[0])
	out := make([][]float64, rows)
	for i := 0; i < rows; i++ {
		out[i] = make([]float64, cols)
		for j := 0; j < cols; j++ {
			var sum float64
			for k := 0; k < inner; k++ {
				sum += a[i][k] * b[k][j]
			}
			out[i][j] = sum
		}
	}
	return out
}

func assertClose(t *testing.T, got, want float64, msg string) {
	t.Helper()
	if math.Abs(got-want) > tol {
		t.Fatalf("%s: got %v, want %v", msg, got, want)
	}
}

// QR is not unique (sign conventions differ across implementations), so we do
// NOT assert against fixed Q/R values. Instead we verify the defining
// properties of a QR factorization.
func TestFactorize_Properties(t *testing.T) {
	cases := map[string][][]float64{
		"tall":   {{1, 2}, {3, 4}, {5, 6}},
		"square": {{4, 1}, {2, 3}},
		"single": {{2, 0}, {0, 5}, {0, 0}},
	}

	for name, a := range cases {
		t.Run(name, func(t *testing.T) {
			m, n := len(a), len(a[0])
			f := Factorize(a)

			// Dimensions: Q is m x n, R is n x n.
			if len(f.Q) != m || len(f.Q[0]) != n {
				t.Fatalf("Q dims = %dx%d, want %dx%d", len(f.Q), len(f.Q[0]), m, n)
			}
			if len(f.R) != n || len(f.R[0]) != n {
				t.Fatalf("R dims = %dx%d, want %dx%d", len(f.R), len(f.R[0]), n, n)
			}

			// A == Q * R.
			recon := multiply(f.Q, f.R)
			for i := 0; i < m; i++ {
				for j := 0; j < n; j++ {
					assertClose(t, recon[i][j], a[i][j], "reconstruction A=Q*R")
				}
			}

			// R is upper triangular.
			for i := 0; i < n; i++ {
				for j := 0; j < i; j++ {
					assertClose(t, f.R[i][j], 0, "R below-diagonal")
				}
			}

			// Q has orthonormal columns: Qᵀ*Q == I.
			for i := 0; i < n; i++ {
				for j := 0; j < n; j++ {
					var dot float64
					for k := 0; k < m; k++ {
						dot += f.Q[k][i] * f.Q[k][j]
					}
					want := 0.0
					if i == j {
						want = 1.0
					}
					assertClose(t, dot, want, "QᵀQ orthonormality")
				}
			}
		})
	}
}
