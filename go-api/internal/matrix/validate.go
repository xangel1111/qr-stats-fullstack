// Package matrix contains the pure numerical domain logic: input validation
// and QR factorization. It has no knowledge of HTTP, JSON transport or the
// database, which keeps it trivially unit-testable.
package matrix

import "errors"

// Validation errors. They are returned as-is to the service layer, which maps
// them to a 422 VALIDATION_ERROR response.
var (
	ErrEmptyMatrix = errors.New("matrix must not be empty")
	ErrEmptyRow    = errors.New("matrix rows must not be empty")
	ErrJagged      = errors.New("matrix rows must all have the same length")
	ErrNotTall     = errors.New("matrix must have rows >= cols (m >= n) for QR factorization")
)

// Validate checks that data is a well-formed rectangular matrix suitable for
// reduced QR factorization.
//
// gonum's QR requires the number of rows to be greater than or equal to the
// number of columns, so wide matrices (m < n) are rejected here. This is a
// documented limitation of the endpoint (see README / API contract).
func Validate(data [][]float64) error {
	rows := len(data)
	if rows == 0 {
		return ErrEmptyMatrix
	}

	cols := len(data[0])
	if cols == 0 {
		return ErrEmptyRow
	}

	for _, row := range data {
		if len(row) != cols {
			return ErrJagged
		}
	}

	if rows < cols {
		return ErrNotTall
	}

	return nil
}
