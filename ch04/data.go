// Copyright (c) 2023 Michael D Henderson. All rights reserved.

package ch04

import "bytes"

// functions in this file create cells on the stack or on the heap.

// cons returns a new Pair created on the heap.
func cons(car, cdr Atom) Atom {
	return Atom{
		_type: AtomType_Pair,
		value: AtomValue{
			pair: &Pair{
				car: car,
				cdr: cdr,
			},
		},
	}
}

// make_int returns an Atom on the stack.
func make_int(x int) Atom {
	return Atom{
		_type: AtomType_Integer,
		value: AtomValue{
			integer: x,
		},
	}
}

// make_sym returns an Atom on the stack.
// The name of the symbol is always converted to uppercase.
// If the symbol already exists in the global symbol table, that symbol is
// returned. Otherwise, a new symbol is created on the stack, added to the
// symbol table, and returned. The new symbol allocates space for the name.
func make_sym(name []byte) Atom {
	// make an upper-case copy of the name
	name = bytes.ToUpper(name)
	// search for any existing symbol with the same name
	for p := sym_table; !nilp(p); p = cdr(p) {
		if atom := car(p); bytes.Equal(name, atom.value.symbol.label) {
			// found match, so return the existing symbol
			return atom
		}
	}
	// did not find a matching symbol, so create a new one
	atom := Atom{
		_type: AtomType_Symbol,
		value: AtomValue{
			symbol: &Symbol{
				label: name,
			},
		},
	}
	// add it to the symbol_table
	sym_table = cons(atom, sym_table)
	// and return it
	return atom
}
