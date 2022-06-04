package main

import (
	"fmt"
	"math/rand"
)

// Secrifies the vote 'x' into three shares, using a polynomial of degree k
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

// Implements sorting interface for the type'PointXSort'
func (a PointXSort) Len() int           { return len(a) }
func (a PointXSort) Less(i, j int) bool { return a[i].X < a[j].X }
func (a PointXSort) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

// Finds the multiplicative inverse of a*t mod n (that is, find -t)
// https://en.wikipedia.org/wiki/Extended_Euclidean_algorithm
// Section on calculating the inverse
func Inverse(a, p int) int {
	t, nt := 0, 1
	r, nr := p, a
	for nr != 0 {
		q := r / nr // floor division is the default in Go, which Inverse uses
		r, nr = nr, r-q*nr
		t, nt = nt, t-q*nt
	}
	if r > p {
		panic(fmt.Errorf("cannot invert %v given %v", a, p))
	}
	return (1 / r) * t
}

// Computes n/d % p
func DivMod(n, d, p int) int {
	return n * Inverse(d, p)
}

// Computes the product an integer array
func Prod(vals []int) int {
	i := 1
	for _, v := range vals {
		i *= v
	}
	return i
}

// Takes the difference between a and each element in b
func DiffInts(a int, b []int) []int {
	c := make([]int, len(b))
	for i := 0; i < len(b); i++ {
		c[i] = a - b[i]
	}
	return c
}

// Positive integer mod operation (Thanks reddit)
// Code from https://www.reddit.com/r/golang/comments/bnvik4/modulo_in_golang/
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
	k := len(nums)
	s := 0
	for i := 0; i < k; i++ {
		n := nums[i] * d * points[i].Y
		n = pmod(n, p)
		v := DivMod(n, dens[i], p)
		s += v
	}
	return s
}

// Computes L(x,P) in field p given points P
func LagrangeXP(x, p int, points []Point) int {
	k := len(points)
	nums := make([]int, 3)
	dens := make([]int, 3)
	for i := 0; i < k; i++ {
		nums[i] = 1
		dens[i] = 1
		for j := 0; j < k; j++ {
			if i != j {
				nums[i] *= x - points[j].X
				dens[i] *= points[i].X - points[j].X
			}
		}
	}
	den := Prod(dens)
	num := SumNum(nums, dens, p, den, points)
	return pmod(DivMod(num, den, p), p)
}
