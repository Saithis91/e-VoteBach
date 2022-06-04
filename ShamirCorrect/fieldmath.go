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

// Finds the multiplicative inverse of a*t mod n (that is, find -t)
// https://en.wikipedia.org/wiki/Extended_Euclidean_algorithm
// Section on calculating the inverse
func Inverse(a, p int) int {
	t, nt := 0, 1
	r, nr := p, a
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
	return n * Inverse(d, p)
}

func SubField(lhs, rhs, p int) int {
	return pmod(lhs-rhs, p)
}

func MulField(lhs, rhs, p int) int {
	return pmod(lhs*rhs, p)
}

func SumField(p int, vals ...int) int {
	sum := 0
	for _, v := range vals {
		sum = pmod(sum+v, p)
	}
	return sum
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

// Peforms x^y operation and stays in integer domain.
// Using math.pow would require casting which *could* lead to incorrect values because floating points
func IPowF(x, y, p int) int {
	if y == 0 {
		return 1
	}
	z := x
	for i := 2; i <= y; i++ {
		z = MulField(z, x, p)
	}
	return z
}
