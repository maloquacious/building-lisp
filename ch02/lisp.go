// Copyright (c) 2023 Michael D Henderson. All rights reserved.

package ch02

// _nil is the NIL symbol.
// This should be immutable, so don't change it!
var _nil = Atom{_type: AtomType_Nil}

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

// nilp is a predicate function. It returns true if the atom is NIL.
func nilp(atom Atom) bool {
	return atom._type == AtomType_Nil
}
