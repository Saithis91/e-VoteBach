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

// Secrifies the vote 'x' into three shares using shamir sharing.
// Returns r1, r2, r3 values in Gf(2^8)
// @x: The secret to secrify
// @k: The threshold to accept leaks
func SecrifyGf(x, k int) (r1, r2, r3 uint8) {

	secrifyFail := 0

	for {

		// secret to Gf
		s := Gf_FromByte(uint8(x))

		// Generate random a-values
		as := make([]uint8, 0)
		for i := 0; i < k; i++ {
			ai := uint8(rand.Intn(256))
			as = append(as, ai)
		}

		// Compute shares
		r1 = PolyGf(1, s, as)
		r2 = PolyGf(2, s, as)
		r3 = PolyGf(3, s, as)

		// Compute L(0) to verify
		if Lagrange0Gf(r1, r2, r3) == uint8(x) {
			return
		} else {
			secrifyFail++
			fmt.Printf("Failed to secrify X on attempt %v.\n", secrifyFail)
		}

	}

}

// Compute the polynomial f(x)=s+a_1x+a_2x^2+...+a_n+x^n
func Poly(x, s, p int, a []int) int {
	y := s                // f(0) = s
	for e, v := range a { // for e=exponent(+1), v = a_i
		xpow := IMod(IPow(x, e+1), p) // x**(e+1) % p
		cx := IMod(v*xpow, p)         // a_i * xpow % p
		y = IMod(y+cx, p)             // y + a_i * xpow % p
	}
	// "overuse" of mod p is based on comments in last paragraph of
	// https://medium.com/partisia-blockchain/mpc-techniques-series-part-3-secret-sharing-shamir-style-f2a952fa7828
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

// Compute the polynomial f(x)=s+a_1x+a_2x^2+...+a_n+x^n
// Guaranteed to be inside the field Gf(2^8)
func PolyGf(x uint8, s Gf, a []uint8) uint8 {
	xg := Gf_FromByte(x)
	y := s
	for e, v := range a {
		xpow := xg.Pow(Gf_FromByte(uint8(e + 1)))
		cx := Gf_FromByte(v).Mul(xpow)
		y = y.Add(cx)
	}
	return y.ToByte()
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

type GfPoint struct {
	X Gf // X-value in Gf(2^8)
	Y Gf // Y-value in Gf(2^8)
}

// Define sort by X for a point slice
type GfPointXSort []GfPoint

func (a GfPointXSort) Len() int           { return len(a) }
func (a GfPointXSort) Less(i, j int) bool { return a[i].X.ToByte() < a[j].X.ToByte() }
func (a GfPointXSort) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

// Computes L(x) for the set of given points.
func Lagrange(x int, points ...Point) int {
	l := 0
	k := len(points)
	for j := 0; j < k; j++ {
		l += points[j].Y * LagrangeBasis(x, j, k, points)
	}
	return l
}

// Computes L(x) for the set of given points
func LagrangeGf(x uint8, points ...GfPoint) Gf {
	//l := Gf_Zero()
	k := len(points)
	x_gf := Gf_FromByte(x)
	ts := make([]Gf, 0)
	for j := 0; j < k; j++ {
		t := points[j].Y.Mul(LagrangeBasisGf(j, k, x_gf, points))
		ts = append(ts, t)
		//l = l.Add(t)
		fmt.Printf("t=%v\n", t)
	}
	l := Gf_Sum(ts...)
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

// Computes delta(x) for the set of given points at x w.r.t i
func LagrangeBasisGf(i, k int, x Gf, points []GfPoint) Gf {
	l := Gf_One()
	x_m := points[i].X
	for j := 0; j < k; j++ {
		if i != j {
			x_j := points[j].X
			delta := x_m.Sub(x_j)
			if delta.Is_Zero() {
				panic(fmt.Errorf("division by zero! "))
			}
			top := x.Sub(x_j)
			l = l.Mul(top.Div(delta))
		}
	}
	return l
}

// Computes L(x) from r1, r2, r3 values
func LagrangeXGf(x uint8, r1, r2, r3 uint8) uint8 {
	l0 := LagrangeGf(x, GfPoint{Y: Gf_FromByte(r1), X: Gf_FromByte(1)}, GfPoint{Y: Gf_FromByte(r2), X: Gf_FromByte(2)}, GfPoint{Y: Gf_FromByte(r3), X: Gf_FromByte(3)})
	return l0.ToByte()
}

// Computes L(0) from r1, r2, r3 values
func Lagrange0Gf(r1, r2, r3 uint8) uint8 {
	return LagrangeXGf(0, r1, r2, r3)
}

// Computes L(0) from the set of r1, r2, r3 values
func Lagrange0(r1, r2, r3 int) int {
	return Lagrange(0, Point{X: 1, Y: r1}, Point{X: 2, Y: r2}, Point{X: 3, Y: r3})
}
