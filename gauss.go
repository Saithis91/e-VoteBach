package main

import (
	"fmt"
	"math"
)

type Matrix = [][]float64
type Vector = []float64

type IntMatrix = [][]int
type IntVector = []int

func AugmentedMatrix(coeffs Matrix, B Vector) Matrix {
	n := len(coeffs)
	A := make(Matrix, n)
	for i := 0; i < n; i++ {
		A[i] = append(coeffs[i], B[i])
	}
	return A
}

func AugMatStr(A Matrix) {
	for _, row := range A {
		var rowstr string
		for i, col := range row {
			if i == 0 {
				rowstr = fmt.Sprintf("%v", col)
			} else if i == len(row)-1 {
				rowstr = fmt.Sprintf("%s = %v", rowstr, col)
			} else {
				rowstr = fmt.Sprintf("%s %v", rowstr, col)
			}
		}
		fmt.Println(rowstr)
	}
}

func AugmentedIntMatrix(coeffs IntMatrix, B IntVector) IntMatrix {
	n := len(coeffs)
	A := make(IntMatrix, n)
	for i := 0; i < n; i++ {
		A[i] = append(coeffs[i], B[i])
	}
	return A
}

func gauss_elim(A Matrix) {
	m := len(A)    // num equations
	n := len(A[0]) // num columns (num variables + B or Y column vector)
	if n-1 > m {
		panic(fmt.Errorf("cannot solve a system with %v unknowns and %v equations - more unknowns than equations", n, m))
	}
	h, k := 0, 0 // first pivot point
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
		fmt.Println(" ")
		AugMatStr(A)
	}
}

func dum_gauss_elim(A Matrix) {
	n := len(A)
	for i := 0; i < n; i++ {
		if A[i][i] == 0 {
			continue
		}
		for j := i + 1; j < n; j++ {
			ratio := A[j][i] / A[i][i]

			for k := 0; k < n+1; k++ {
				A[j][k] = A[j][k] - ratio*A[i][k]
			}
		}
		fmt.Println(" ")
		AugMatStr(A)
	}
}

func dum_back_Sub(A Matrix) Vector {
	n := len(A)
	X := make(Vector, n)
	X[n-1] = A[n-1][n] / A[n-1][n-1]

	for i := n - 2; i > -1; i-- {
		X[i] = A[i][n]

		for j := i + 1; j < n; j++ {
			X[i] = X[i] - A[i][j]*X[j]
		}
		X[i] = X[i] / A[i][i]
	}
	return X
}

func gauss_elim_field(A IntMatrix, p int) {
	m := len(A)    // num equations
	n := len(A[0]) // num columns
	h, k := 0, 0   // first pivot point
	for h < m && k < n {
		imax := argmax_imatrix(h, m, k, A) // more stable (numerically, which we may need, given integer work)
		if A[imax][k] == 0 {
			k++
		} else {
			A[h], A[imax] = A[imax], A[h]
			for i := h + 1; i < m; i++ {
				f := DivMod(A[i][k], A[h][k], p)
				A[i][k] = 0
				for j := k + 1; j < n; j++ {
					prod := MulField(A[h][j], f, p)
					A[i][j] = SubField(A[i][j], prod, p)
				}
			}
			h++
			k++
		}
	}
}

// Finds in the index of the largest (absolute) value in a matrix column 'k'
// within the row range min <= x <= max
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

// Finds in the index of the largest (absolute) value in a matrix column 'k'
// within the row range min <= x <= max
func argmax_imatrix(min, max, k int, M IntMatrix) int {
	i := min
	for j := min; j < max; j++ {
		a := math.Abs(float64(M[i][k])) // tmp solution
		b := math.Abs(float64(M[j][k]))
		if b > a {
			i = j
		}
	}
	return i
}

func back_substitute(A Matrix) Vector {
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

func back_substitute_field(A IntMatrix, p int) IntVector {
	n := len(A)
	X := make(IntVector, n)
	for i := n - 1; i >= 0; i-- {
		dot := 0
		for j := i + 1; j < n; j++ {
			dot += MulField(A[i][j], A[j][n], p)
		}
		num := SubField(A[i][n], dot, p)
		A[i][n] = DivMod(num, A[i][i], p)
		X[i] = A[i][n]
	}
	return X
}
