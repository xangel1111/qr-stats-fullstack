package matrix

import (
	"errors"
	"testing"
)

func TestValidate(t *testing.T) {
	cases := []struct {
		name string
		in   [][]float64
		want error
	}{
		{"valid tall", [][]float64{{1, 2}, {3, 4}, {5, 6}}, nil},
		{"valid square", [][]float64{{1, 2}, {3, 4}}, nil},
		{"empty", [][]float64{}, ErrEmptyMatrix},
		{"empty row", [][]float64{{}}, ErrEmptyRow},
		{"jagged", [][]float64{{1, 2}, {3}}, ErrJagged},
		{"wide (m<n)", [][]float64{{1, 2, 3}}, ErrNotTall},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := Validate(tc.in)
			if !errors.Is(err, tc.want) {
				t.Fatalf("Validate() = %v, want %v", err, tc.want)
			}
		})
	}
}
