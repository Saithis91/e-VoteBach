package main

import (
	"math"
)

func gauss_elim(A [][]float64) {
	m := len(A)    // num equations
	n := len(A[0]) // num columns
	h, k := 0, 0   // first pivot point
	for h < m && k < n {
		imax := argmax_matrix(h, m, k, A)
		if A[imax][k] == 0 {
			k++
		} else {
			//fmt.Printf("A1 = %v\n", A)
			A[h], A[imax] = A[imax], A[h] // swap
			//fmt.Printf("A2 = %v\n", A)
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
		//fmt.Printf("A = %v\n", A)
	}
}

/*
func gauss_elim_vals(A [][]float64) []float64 {
	q := make([]float64, len(A))
	for i, row := range A {
		q[i] = row[len(row)-1]
	}
	return q
}
*/

// https://www.sciencedirect.com/topics/mathematics/back-substitution
func gauss_jordan_elim(A [][]float64) []float64 {
	m := len(A)
	n := len(A[0]) - 1
	B := make([]float64, len(A))
	for i := 0; i < m; i++ {
		B[i] = A[i][n]
		for j := 0; j < i-1; j++ {
			B[i] -= A[i][j] * B[j]
		}
	}
	return B
}

func argmax_matrix(min, max, k int, M [][]float64) int {
	i := min
	for j := min; j < max; j++ {
		a := math.Abs(M[i][k])
		b := math.Abs(M[j][k])
		if b >= a {
			i = j
		}
	}
	return i
}
