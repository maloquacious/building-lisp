// Copyright (c) 2023 Michael D Henderson. All rights reserved.

package ch04

// car returns the first item from a list.
// It will panic if p is not a Pair
func car(p Atom) Atom {
	return p.value.pair.car
}

// cdr returns the remainder of a list.
// It will panic if p is not a Pair
func cdr(p Atom) Atom {
	return p.value.pair.cdr
}

// listp returns true if the expression is a proper list or is NIL.
func listp(expr Atom) bool {
	for ; !nilp(expr); expr = cdr(expr) {
		if expr._type != AtomType_Pair {
			return false
		}
	}
	return true
}

// nilp is a predicate function. It returns true if the atom is NIL.
func nilp(atom Atom) bool {
	return atom._type == AtomType_Nil
}

// setcar is a helper function to set the car of a pair.
// panics if p is not a pair.
func setcar(p, a Atom) {
	p.value.pair.car = a
}

// setcdr is a helper function to set the cdr of a pair.
// panics if p is not a pair.
func setcdr(p, a Atom) {
	p.value.pair.cdr = a
}
