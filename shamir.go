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

// Peforms x^y operation and stays in integer domain.
// Using math.pow would require casting which *could* lead to incorrect values because floating points
func IPow(x, y int) int {
	if y == 0 {
		return 1
	}
	z := x
	for i := 2; i <= y; i++ {
		z *= x
	}
	return z
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

// https://en.wikipedia.org/wiki/Shamir%27s_Secret_Sharing
// Python translation

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

//No clue where this should actually go.

//Compute the Polynomial
func Polynomial(points []Point) {
	// Create equations
	//
	coeffs := Matrix{
		{},
		{},
		{},
		{},
	}
	B := Vector{}
	for i := 0; i < len(points); i++ {
		x := points[i].X
		y := points[i].Y
		//coeffs[i] = Vector{float64(pmod(x, 3)), float64(pmod(x, 2)), float64(x), 1, float64(y)}
		coeffs[i] = Vector{float64(IPow(x, 2)), float64(x), 1, float64(y)}
		B = append(B, float64(y*x))
	}
	//fmt.Println(coeffs)
	//fmt.Println(B)
	// Create V-vector

	// Create matrix
	A := AugmentedMatrix(coeffs, B)

	AugMatStr(A)

	// Solve
	//gauss_elim(A)
	dum_gauss_elim(A)
	fmt.Print("\nFinal Output was\n")
	AugMatStr(A)
	fmt.Print("\nBack-Sub was:\n")
	//X := back_substitute(A)
	X := dum_back_Sub(A)
	fmt.Printf("A:%v\nX:%v\n", A, X)
}
