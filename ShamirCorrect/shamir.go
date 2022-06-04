package main

import (
	"fmt"
	"math/rand"
)

// Secrifies the vote 'x' into three shares, using shamir sharing.
// @x = The vote (in {0, 1})
// @p = The prime number to limit Z-field (-p, p)
// @k = The polynomium degree (amount of corrupt parties we allow)
func Secrify(x, p, k int) (r1, r2, r3, r4 int) {

	// Generate random a-values
	as := make([]int, 0)
	upper := p - 1
	for i := 0; i < k; i++ {
		ai := rand.Intn(upper)
		as = append(as, ai)
	}

	// Compute shares
	r1 = Poly(1, x, p, as)
	r2 = Poly(2, x, p, as)
	r3 = Poly(3, x, p, as)
	r4 = Poly(4, x, p, as)

	// Return
	return

}

// Compute the polynomial f(x)=s+a_1x+a_2x^2+...+a_n+x^n
func Poly(x, s, p int, a []int) int {
	y := s
	for e, v := range a {
		xpow := IPow(x, e+1)
		cx := v * xpow
		y += cx
		y %= p
	}
	return y
}

// Represents a point in a coordinate system
type Point struct {
	X int // X-Value
	Y int // Y-Value
}

// Define sort by X for a point slice
type PointXSort []Point

func (a PointXSort) Len() int           { return len(a) }
func (a PointXSort) Less(i, j int) bool { return a[i].X < a[j].X }
func (a PointXSort) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

// Computes the product of all integers in integer arrray
// p = vals[1] * vals[2] * ... * vals[n]
func ProdInts(vals []int) int {
	i := 1
	for _, v := range vals {
		i *= v
	}
	return i
}

// Computes numerator
func SumNum(nums, dens []int, p, d int, points []Point) int {

	// Init
	k := len(nums)
	s := 0

	// Loop over sample points and sum them (semi-computes L(x)=delta_i(x)*y_i)
	for i := 0; i < k; i++ {
		n := nums[i] * d * points[i].Y
		n = pmod(n, p)
		s += DivMod(n, dens[i], p)
	}

	return s

}

// Computes L(x) in field p given points p
func Lagrange(x, p int, points []Point) int {

	// Init
	k := len(points)
	nums := make([]int, k)
	dens := make([]int, k)

	// Create numerators and denominators
	for i := 0; i < k; i++ {

		// Products (so n[i]=1,d[i]=1 so we multiply by 1 in first step)
		nums[i] = 1
		dens[i] = 1

		// Delta_i(x)
		for j := 0; j < k; j++ {
			if i != j {
				nums[i] *= x - points[j].X
				dens[i] *= points[i].X - points[j].X
			}
		}
	}

	// Compute denominator and numerator
	den := ProdInts(dens)
	num := SumNum(nums, dens, p, den, points)

	// Calculate (num / den mod p) ^ p -> outcome of L(x)
	return pmod(DivMod(num, den, p), p)

}

// Finds the index of the error point using the Berlekamp–Welch algorithm
// it then corrects the error by ommiting the point at the specified error index
// and computes L(0, (x,y)'...). If the correction fails, an error is returned.
func CorrectError(points []Point, prime int) (int, error) {

	// Construct system of linear equations
	A := make(IntMatrix, len(points))
	B := make(IntVector, len(points))
	for i := 0; i < len(points); i++ {
		x := points[i].X
		y := points[i].Y
		A[i] = IntVector{IPow(x, 2), x, 1.0, y}
		B[i] = y * x
	}

	// Apply gauss
	Y := GaussElim(A, B, prime)

	// Grab error value
	e := Y[3]

	// Verify e
	if e < 1 || e > len(points) {
		return -1, fmt.Errorf("invalid error index found")
	}

	// Redo Lagrange without error point
	correct := Lagrange(0, prime, append(points[0:e-1], points[e:]...))

	// Return correct vote sum and no error
	return correct, nil

}
