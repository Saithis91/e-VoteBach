package main

import "fmt"

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
func Inverse(a, p int) int {
	t := 0
	r := p
	nt := 1
	nr := a
	for nr != 0 {
		q := r / nr
		r, nr = nr, r-q*nr
		t, nt = nt, t-q*nt
	}
	if r > p {
		panic(fmt.Errorf("cannot invert %v given %v", a, p))
	}
	return (1 / r) * t
}

// Computes n/d % p
// Multiplicative inverse, which we need for staying in the field
func DivMod(n, d, p int) int {
	//inv, _ := GCD(d, p)
	//return n * inv
	return n * Inverse(d, p)
}

func SubField(rhs, lhs, p int) int {
	return pmod(lhs-rhs, p)
}

func MulField(rhs, lhs, p int) int {
	return pmod(rhs*rhs, p)
}
