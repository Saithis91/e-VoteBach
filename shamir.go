package main

import (
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

// Compute the polynomial f(x)=s+a_1x+a_2x^2+...+a_n+x^n
// Was used when testing
func Poly2(x, s int, a []int) int {
	y := s
	for e, v := range a {
		xpow := IPow(x, e+1)
		cx := v * xpow
		y = y + cx
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

// Apparently Go has a 'botched' modulo operator implementation
// Which can yield negative numbers - which does not adhere to the strict
// mathemtatical modulo operation we require.
// Code from https://www.reddit.com/r/golang/comments/bnvik4/modulo_in_golang/
func IMod(a int, b int) int {
	return (a + b) % b
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

// Computes L(x) for the set of given points.
func Lagrange(x int, points ...Point) int {
	l := 0
	k := len(points)
	for j := 0; j < k; j++ {
		l += points[j].Y * LagrangeBasis(x, j, k, points)
	}
	return l
}

// Computes eel(x) for the set of given points at x w.r.t i
func LagrangeBasis(x, j, k int, points []Point) int {
	l := 1
	x_j := points[j].X
	for m := 0; m < k; m++ {
		if m != j {
			l *= (x - points[m].X) / (x_j - points[m].X)
		}
	}
	return l // TODO: Multiplicative inverse of (x_j - points[m].X)
}

// Computes L(0) from the set of r1, r2, r3 values
func Lagrange0(r1, r2, r3 int) int {
	return Lagrange(0, Point{X: 1, Y: r1}, Point{X: 2, Y: r2}, Point{X: 3, Y: r3})
}

// https://en.wikipedia.org/wiki/Shamir%27s_Secret_Sharing
// Python translation

// Finds greatest common divisor based on  Euclidean algorithm
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
	//fmt.Printf("LX,LY=%v,%v\n", lx, ly)
	return lx, ly
}

// Computes n/d % p in the field of integers
func DivMod(n, d, p int) int {
	inv, _ := GCD(d, p)
	return n * inv
}

// Sums integers
func SumInts(vals []int) int {
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

// Sums all values in vals
func SumInt(vals []int) int {
	i := 1
	for _, v := range vals {
		i *= v
	}
	return i
}

// Positive integer mod operation (Thanks reddit)
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
		//fmt.Printf("v = %v\n", v)
		s += v
	}
	//fmt.Printf("s = %v\n", s)
	return s
}

// Computes L(x) in field p given points p
func LagrangeXP(x, p int, points []Point) int {
	k := len(points)
	nums := make([]int, k)
	dens := make([]int, k)
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
	den := SumInts(dens)
	num := SumNum(nums, dens, p, den, points)
	return pmod(DivMod(num, den, p), p)
}
