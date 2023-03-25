// Copyright (c) 2023 Michael D Henderson. All rights reserved.

package ch02

import "bytes"

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

// sym_table is a global symbol table.
// it is a list of all existing symbols.
var sym_table = Atom{_type: AtomType_Nil}

// make_sym returns an Atom on the stack.
// The name of the symbol is always converted to uppercase.
// If the symbol already exists in the global symbol table, that symbol is
// returned. Otherwise, a new symbol is created on the stack, added to the
// symbol table, and returned. The new symbol allocates space for the name.
func make_sym(name []byte) Atom {
	// make an upper-case copy of the name
	name = bytes.ToUpper(name)
	// search for any existing symbol with the same name
	for atom := sym_table; !nilp(atom); atom = cdr(atom) {
		if bytes.Equal(name, car(atom).value.symbol) {
			// found match, so return the existing symbol
			return atom
		}
	}
	// did not find a matching symbol, so create a new one
	atom := Atom{
		_type: AtomType_Symbol,
		value: AtomValue{
			symbol: name,
		},
	}
	// add it to the symbol_table
	sym_table = cons(atom, sym_table)
	// and return it
	return atom
}
