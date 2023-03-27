// Copyright (c) 2023 Michael D Henderson. All rights reserved.

package lisp

import "os"

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

// print_expr is a helper function to write an expression
// to stdout, ignoring errors.
func print_expr(expr Atom) {
	_, _ = expr.Write(os.Stdout)
}
