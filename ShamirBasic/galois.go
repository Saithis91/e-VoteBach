package main

import "fmt"

// https://github.com/sellibitze/secretshare/blob/master/src/gf256.rs
// Rust -> Go "compilation" - we take no credit for the implementation
// We have merely translated the Rust code to a Go equivalent

type byte = uint8
type ptr = uintptr

// This should be in every language, not just the functional ones...
type Something struct {
	// The stored value
	val byte
	// Flag marking if there's something
	some bool
}

// Yields nothing of something
var None = Something{}

// Returns Something of 'v'
func Some(v byte) Something {
	return Something{
		val:  v,
		some: true,
	}
}

// Lookup tables
type Lookups struct {
	gf_exp [256]byte
	gf_log [256]byte
	gf_inv [256]byte
}

// static lookup table
var lTables = GenLookupTable()

func wrp(b byte) byte {
	return byte(0) - (b & 1)
}

func xt(x byte) byte {
	return (x << 1) ^ (wrp(x>>7) & byte(0x1D))
}

// Generates lookups
func GenLookupTable() Lookups {
	l := Lookups{}
	tmp := uint(1)
	for i := 0; i < 255; i++ {
		p := byte(i)
		l.gf_exp[p] = byte(tmp)
		l.gf_log[tmp] = p
		tmp = uint(xt(byte(tmp))) // This hurt to write
	}
	l.gf_exp[255] = 1
	for i := 1; i < 256; i++ {
		log := l.gf_log[i]
		nl := uint8(0)
		if log != 0 {
			nl = 255 - log
		}
		l.gf_inv[i] = l.gf_exp[ptr(nl)]
	}
	return l
}

// Implements functions for doing operations within a galois field
// GF(2^8=256)
type Gf struct {
	// Polynomium base
	poly byte
}

// Get a zero element in GF
func Gf_Zero() Gf {
	return Gf{
		poly: 0,
	}
}

// Get a 'one' element in GF
func Gf_One() Gf {
	return Gf{
		poly: 1,
	}
}

// Get the galois filed of b (which is limited to 256, because of how bytes work :P )
func Gf_FromByte(b uint8) Gf {
	return Gf{
		poly: byte(b),
	}
}

// Gets the galois field element as poly
func (x Gf) ToByte() uint8 {
	return x.poly
}

// Returns something to the log at 'x'
func (x Gf) Log() Something {
	if x.poly == 0 {
		return None
	} else {
		return Some(lTables.gf_log[x.poly])
	}
}

// Returns the exponent field of Gf(2^8)
func Gf_Exp(p byte) Gf {
	return Gf{
		poly: lTables.gf_exp[p],
	}
}

func Gf_Sum(entries ...Gf) Gf {
	sum := Gf_Zero()
	for _, v := range entries {
		sum = sum.Add(v)
	}
	return sum
}

// Computes x+y in Gf(2^8)
func (x Gf) Add(y Gf) Gf {
	return Gf{
		poly: x.poly ^ y.poly,
	}
}

// Computes x-y in Gf(2^8)
// Same as addition
func (x Gf) Sub(y Gf) Gf {
	return Gf{
		poly: x.poly ^ y.poly,
	}
}

// Computes x*y in Gf(2^8)
func (x Gf) Mul(y Gf) Gf {
	v1 := x.poly
	v2 := y.poly
	if v1 != 0 && v2 != 0 {
		l1 := x.Log()
		l2 := y.Log()
		if l1.some && l2.some {
			tmp := (uint(l1.val) + uint(l2.val)) % uint(255)
			return Gf_Exp(byte(tmp))
		}
	}
	fmt.Println("Exiting Mul with 0")
	return Gf{poly: 0}
}

// Computes x/y in Gf(2^8)
func (x Gf) Div(y Gf) Gf {
	v1 := x.poly
	v2 := y.poly
	if v1 != 0 && v2 != 0 {
		l1 := x.Log()
		l2 := y.Log()
		if l1.some && l2.some {
			return Gf_Exp(byte((uint(l1.val) + uint(255) - uint(l2.val)) % uint(255)))
		}
	}
	fmt.Println("Exiting Div with 0")
	return Gf{poly: 0}
}

// Computes x^y
func (x Gf) Pow(y Gf) Gf {
	z := x
	for i := 2; i < int(y.poly); i++ {
		z = z.Mul(x)
	}
	return z
}

// Decides if x is equal to 0.
func (x Gf) Is_Zero() bool {
	return x.poly == 0
}