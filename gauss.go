package main

import (
	"math"
)

type Matrix = [][]float64
type Vector = []float64

func AugmentedMatrix(coeffs Matrix, B Vector) Matrix {
	n := len(coeffs)
	A := make(Matrix, n)
	for i := 0; i < n; i++ {
		A[i] = append(coeffs[i], B[i])
	}
	return A
}

func gauss_elim(A Matrix) {
	m := len(A)    // num equations
	n := len(A[0]) // num columns
	h, k := 0, 0   // first pivot point
	for h < m && k < n {
		imax := argmax_matrix(h, m, k, A) // more stable (numerically, which we may need, given integer work)
		if A[imax][k] == 0 {
			k++
		} else {
			A[h], A[imax] = A[imax], A[h]
			for i := h + 1; i < m; i++ {
				f := A[i][k] / A[h][k]
				A[i][k] = 0
				for j := k + 1; j < n; j++ {
					A[i][j] = A[i][j] - A[h][j]*f
				}
			}
			h++
			k++
		}
	}
}

func argmax_matrix(min, max, k int, M Matrix) int {
	i := min
	for j := min; j < max; j++ {
		a := math.Abs(M[i][k])
		b := math.Abs(M[j][k])
		if b > a {
			i = j
		}
	}
	return i
}

func gauss_jordan_elim(A Matrix) Vector {
	n := len(A)
	X := make(Vector, n)
	for i := n - 1; i >= 0; i-- {
		dot := 0.0
		for j := i + 1; j < n; j++ {
			dot += A[i][j] * A[j][n]
		}
		A[i][n] = (A[i][n] - dot) / A[i][i]
		X[i] = A[i][n]
	}
	return X
}
