package main

type Matrix = [][]float64
type Vector = []float64
type IntMatrix [][]int
type IntVector = []int

// Applies the gaussian elimination algorithm on the system of linear equations
// Based on the C++ code https://github.com/makramkd/gaussian-elimination
func GaussElim(A IntMatrix, B IntVector, prime int) IntVector {

	// Define pivots
	piv := make(IntVector, len(A))
	cpiv := make(IntVector, len(A[0]))

	for i := 0; i < len(piv); i++ {
		piv[i] = i
		cpiv[i] = i
	}

	n := len(A)

	for i := 0; i < n-1; i++ {
		magnitude := 0
		row_index := -1
		col_index := -1
		for j := i; j < n-1; j++ {
			for k := i; k < n-1; k++ {
				if A[piv[j]][piv[k]] > magnitude {
					magnitude = A[piv[j]][piv[k]]
					row_index = j
					col_index = k
				}
			}
		}

		// Swap rows
		if row_index != -1 {
			piv[i], piv[row_index] = piv[row_index], piv[i]
		}

		// Swap columns
		if col_index != -1 {
			cpiv[i], cpiv[col_index] = cpiv[col_index], cpiv[i]
		}

		for j := i + 1; j < n; j++ {
			ratio := DivMod(A[piv[j]][cpiv[i]], A[piv[i]][cpiv[i]], prime)
			for k := i; k < n; k++ {
				A[piv[j]][cpiv[k]] = SubField(A[piv[j]][cpiv[k]], MulField(ratio, A[piv[i]][cpiv[k]], prime), prime)
			}
			B[piv[j]] = SubField(B[piv[j]], MulField(ratio, B[piv[i]], prime), prime)
		}

	}

	return BacksubField(A, B, piv, cpiv, prime)
}

// Applies backsubstitution on the the row-echeclon formed matrix and solves the system
func BacksubField(A IntMatrix, B IntVector, rowPivots, columnPivots IntVector, prime int) IntVector {

	// Number of rows (eqs)
	m := len(A)

	// Prepare temp solution vector
	s := make(IntVector, len(B))

	// Standard back sub
	for i := m - 1; i >= 0; i-- {
		s[rowPivots[i]] = B[rowPivots[i]]
		for j := i + 1; j < m; j++ {
			s[rowPivots[i]] = SubField(s[rowPivots[i]], MulField(A[rowPivots[i]][columnPivots[j]], s[rowPivots[j]], prime), prime)
		}
		s[rowPivots[i]] = DivMod(s[rowPivots[i]], A[rowPivots[i]][columnPivots[i]], prime)
	}

	// Prepare solution vector
	solution := make(IntVector, len(B))

	// Correct pivot points
	for i := 0; i < len(rowPivots); i++ {
		solution[columnPivots[i]] = pmod(s[rowPivots[i]], prime)
	}

	// Return solution
	return solution

}
