//! This module provides the Gf256 type which is used to represent
//! elements of a finite field wich 256 elements.

use std::num::Wrapping;
use std::ops::{ Add, Sub, Mul, Div };
use std::sync::{ Once, ONCE_INIT };

const POLY: u8 = 0x1D; // represents x^8 + x^4 + x^3 + x^2 + 1

/// replicates the least significant bit to every other bit
#[inline]
fn mask(bit: u8) -> u8 {
    (Wrapping(0u8) - Wrapping(bit & 1)).0
}

/// multiplies a polynomial with x and returns the residual
/// of the polynomial division with POLY as divisor
#[inline]
fn xtimes(poly: u8) -> u8 {
	(poly << 1) ^ (mask(poly >> 7) & POLY)
}

/// Tables used for multiplication and division
struct Tables {
	exp: [u8; 256],
	log: [u8; 256],
	inv: [u8; 256]
}

static INIT: Once = ONCE_INIT;
static mut TABLES: Tables = Tables {
	exp: [0; 256],
	log: [0; 256],
	inv: [0; 256]
};

fn get_tables() -> &'static Tables {
	INIT.call_once(|| {
		// mutable access is fine because of synchronization via INIT
		let tabs = unsafe { &mut TABLES };
		let mut tmp = 1;
		for power in 0..255usize {
			tabs.exp[power] = tmp;
			tabs.log[tmp as usize] = power as u8;
			tmp = xtimes(tmp);
		}
		tabs.exp[255] = 1;
		for x in 1..256usize {
			let l = tabs.log[x];
			let nl = if l==0 { 0 } else { 255 - l };
			let i = tabs.exp[nl as usize];
			tabs.inv[x] = i;
		}
	});
	// We're guaranteed to have TABLES initialized by now
	return unsafe { &TABLES };
}

/// Type for elements of a finite field with 256 elements
#[derive(Copy,Clone,PartialEq,Eq)]
pub struct Gf256 {
	pub poly: u8
}

impl Gf256 {
	/// returns the additive neutral element of the field
	#[inline]
	pub fn zero() -> Gf256 {
		Gf256 { poly: 0 }
	}
	/// returns the multiplicative neutral element of the field
	#[inline]
	pub fn one() -> Gf256 {
		Gf256 { poly: 1 }
	}
	#[inline]
	pub fn from_byte(b: u8) -> Gf256 {
		Gf256 { poly: b }
	}
	#[inline]
	pub fn to_byte(&self) -> u8 {
		self.poly
	}
	pub fn log(&self) -> Option<u8> {
		if self.poly == 0 {
			None
		} else {
			let tabs = get_tables();
			Some(tabs.log[self.poly as usize])
		}
	}
	pub fn exp(power: u8) -> Gf256 {
		let tabs = get_tables();
		Gf256 { poly: tabs.exp[power as usize] }
	}
/*
	pub fn inv(&self) -> Option<Gf256> {
		self.log().map(|l| Gf256::exp(255 - l))
	}
*/
}

impl Add<Gf256> for Gf256 {
	type Output = Gf256;
	#[inline]
	fn add(self, rhs: Gf256) -> Gf256 {
		Gf256::from_byte(self.poly ^ rhs.poly)
	}
}

impl Sub<Gf256> for Gf256 {
	type Output = Gf256;
	#[inline]
	fn sub(self, rhs: Gf256) -> Gf256 {
		Gf256::from_byte(self.poly ^ rhs.poly)
	}
}

impl Mul<Gf256> for Gf256 {
	type Output = Gf256;
	fn mul(self, rhs: Gf256) -> Gf256 {
		if let (Some(l1), Some(l2)) = (self.log(), rhs.log()) {
			let tmp = ((l1 as u16) + (l2 as u16)) % 255;
			Gf256::exp(tmp as u8)
		} else {
			Gf256 { poly: 0 }
		}
	}
}

impl Div<Gf256> for Gf256 {
	type Output = Gf256;
	fn div(self, rhs: Gf256) -> Gf256 {
		let l2 = rhs.log().expect("division by zero");
		if let Some(l1) = self.log() {
			let tmp = ((l1 as u16) + 255 - (l2 as u16)) % 255;
			Gf256::exp(tmp as u8)
		} else {
			Gf256 { poly: 0 }
		}
	}
}

fn lagrange_basis(src: &[(u8, u8)], x: Gf256, i:usize, xi: Gf256) -> Gf256 {
	let mut lix = Gf256::one();
	for (j, &(raw_xj, _)) in src.iter().enumerate() {
		if i != j {
			let xj = Gf256::from_byte(raw_xj);
			let delta = xi - xj;
			assert!(delta.poly !=0, "Duplicate shares");
			lix = lix * (x - xj) / delta;
		}
	}
	lix
}

/// evaluates an interpolated polynomial at `raw_x` where
/// the polynomial is determined using Lagrangian interpolation
/// based on the given x/y coordinates `src`.
fn lagrange_interpolate(src: &[(u8, u8)], raw_x: u8) -> u8 {
	let x = Gf256::from_byte(raw_x);
	println!("interpolating on {} using points {:?}", raw_x, src);
	let mut sum = Gf256::zero();
	for (i, &(raw_xi, raw_yi)) in src.iter().enumerate() {
		let xi = Gf256::from_byte(raw_xi);
		let yi = Gf256::from_byte(raw_yi);
		let lix = lagrange_basis(src, x, i, xi);
		println!("D_{}(0)= {}", i, lix.to_byte());
		sum = sum + lix * yi;
		println!("Sum @ {} = {}", i, sum.to_byte());
	}
	sum.to_byte()
}

fn main() {

    // Get one and one
    let x = Gf256::one();
    let y = Gf256::one();

    // Compute x + y
    let c = x + y;

    // print
    println!("1 + 1 in GF is {}", c.to_byte());

    // Log
    /*println!("Exp table = ");
    let tables = get_tables();
    for x in 0..256 {
        print!("[{} = {}],", x, tables.exp[x]);
        if x % 10 == 0 {
            println!("")
        }
    }
    println!("");
    println!("Log table = ");
    let tables = get_tables();
    for x in 0..256 {
        print!("[{} = {}],", x, tables.log[x]);
        if x % 10 == 0 {
            println!("")
        }
    }
    println!("");
    println!("Inv table = ");
    let tables = get_tables();
    for x in 0..256 {
        print!("[{} = {}],", x, tables.inv[x]);
        if x % 10 == 0 {
            println!("")
        }
    }
    println!("1^1={}", 1 ^ 1);*/

	let points = [(1, 5), (2, 8), (3, 11)];
	let p = lagrange_interpolate(&points, 0);

	println!("L(0)={}", p);

	let points = [(1, 182), (2, 113), (3, 199)];
	let p = lagrange_interpolate(&points, 0);

	println!("L(0)={}", p);

}
