package main

import (
	"fmt"
	"math/rand"
)

// Secrifies the vote 'x' into three shares, using shamir sharing.
// @x = The vote (in {0, 1})
// @p = The prime number to limit Z-field (-p, p)
// @k = The polynomium degree (amount of corrupt parties we allow)
func Secrify(x, p, k int) (r1, r2, r3 int) {

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

// Finds greatest common divisor based on the extended Euclidean algorithm
// https://en.wikipedia.org/wiki/Extended_Euclidean_algorithm
// Keeping this just in case the Inver(a,n) is incorrect
func GCD(a, b int) (int, int) {
	x := 0
	lx := 1
	y := 1
	ly := 0
	for b != 0 {
		q := a / b // floor division is default behaviour in Golang
		a, b = b, pmod(a, b)
		x, lx = lx-q*x, x
		y, ly = ly-q*y, y
	}
	return lx, ly
}

// Finds the multiplicative inverse of a*t mod n (that is, find -t)
// https://en.wikipedia.org/wiki/Extended_Euclidean_algorithm
// Section on calculating the inverse
func Inverse(a, n int) int {
	t := 0
	r := n
	nt := 1
	nr := a
	for nr != 0 {
		q := r / nr
		t, nt = nt, t-q*nt
		r, nr = nr, r-q*nr
	}
	if r > 1 {
		panic(fmt.Errorf("cannot invert %v given %v", a, n))
	}
	if t > 0 {
		return -t // We need the inverse, so we want the negative value here
	}
	return t
}

// Computes n/d % p
// Multiplicative inverse, which we need for staying in the field
func DivMod(n, d, p int) int {
	//inv, _ := GCD(d, p)
	//inv2 := Inverse(d, p)
	//fmt.Printf("inv = %v; inv2 = %v.\n", inv, inv2)
	return n * Inverse(d, p)
	//return n * Inverse(d, p)
}

// Computes the product of all integers in integer arrray
// p = vals[1] * vals[2] * ... * vals[n]
func ProdInts(vals []int) int {
	i := 1
	for _, v := range vals {
		i *= v
	}
	return i
}

// Positive integer mod operation (Thanks reddit)
// That is, it computes y = x mod d such that y >= 0
func pmod(x, d int) int {
	x = x % d
	if x >= 0 {
		return x
	}
	if d < 0 {
		return x - d
	}
	return x + d
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
