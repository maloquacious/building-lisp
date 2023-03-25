// Copyright (c) 2023 Michael D Henderson. All rights reserved.

package ch02

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
)

// Atom is either an Atom or a Pair
type Atom struct {
	_type AtomType
	value AtomValue
}

// AtomType is the enum for the type of value in a cell.
type AtomType int

const (
	// Nil represents the empty list.
	AtomType_Nil AtomType = iota
	// Integer is a number.
	AtomType_Integer
	// Pair is a "cons" cell holding a "car" and "cdr" pointer.
	AtomType_Pair
	// Symbol is a string of characters, converted to upper-case.
	AtomType_Symbol
)

// AtomValue is the value of an Atom.
// It can be a simple type, like an integer or symbol, or a pointer to a Pair.
type AtomValue struct {
	integer int
	pair    *Pair
	symbol  *Symbol
}

// Pair is the two elements of a cell.
// "car" is the left-hand value and "cdr" is the right-hand.
type Pair struct {
	car, cdr Atom
}

// Symbol implements data for a symbol.
type Symbol struct {
	label []byte
}

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

// _nil is the NIL symbol.
// This should be immutable, so don't change it!
var _nil = Atom{_type: AtomType_Nil}

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
		if bytes.Equal(name, car(atom).value.symbol.label) {
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

// print_expr is a helper function to write an expression
// to stdout, ignoring errors.
func print_expr(expr Atom) {
	_, _ = expr.Write(os.Stdout)
}

// Write writes the value of an Atom to the writer.
// If the atom is a pair, Write is called recursively
// to write out the entire list.
func (a Atom) Write(w io.Writer) (int, error) {
	switch a._type {
	case AtomType_Nil:
		// atom is nil, so write "NIL"
		return w.Write([]byte{'N', 'I', 'L'})
	case AtomType_Integer:
		// atom is an integer
		return w.Write([]byte(fmt.Sprintf("%d", a.value.integer)))
	case AtomType_Pair:
		// atom is a list, so write it out surrounded by ( and ).
		totalBytesWritten, err := w.Write([]byte{'('})
		if err != nil {
			return totalBytesWritten, err
		}

		// print the car of the list.
		bytesWritten, err := car(a).Write(w)
		totalBytesWritten += bytesWritten
		if err != nil {
			return totalBytesWritten, err
		}

		// write the remainder of the list
		for p := cdr(a); !nilp(p); p = cdr(p) {
			// write a space to separate expressions in the list.
			bytesWritten, err = w.Write([]byte{' '})
			totalBytesWritten += bytesWritten
			if err != nil {
				return totalBytesWritten, err
			}

			if p._type == AtomType_Pair {
				// print the car of the list
				bytesWritten, err = car(p).Write(w)
				totalBytesWritten += bytesWritten
				if err != nil {
					return totalBytesWritten, err
				}
			} else {
				// found an "improper list" (ends with a dotted pair).
				// write dot then space to separate the dotted pair.
				bytesWritten, err = w.Write([]byte{'.', ' '})
				totalBytesWritten += bytesWritten
				if err != nil {
					return totalBytesWritten, err
				}

				// print the atom
				bytesWritten, err = p.Write(w)
				totalBytesWritten += bytesWritten
				if err != nil {
					return totalBytesWritten, err
				}

				// dotted pair ends a list, so quit the loop now
				break
			}
		}

		// write the closing paren
		bytesWritten, err = w.Write([]byte{')'})
		totalBytesWritten += bytesWritten

		// and return
		return totalBytesWritten, err
	case AtomType_Symbol:
		return w.Write(a.value.symbol.label)
	}

	panic(fmt.Sprintf("assert(_type != %d)", a._type))
}

// String implements the Stringer interface.
func (a Atom) String() string {
	sb := &strings.Builder{}
	if _, err := a.Write(sb); err != nil {
		panic(err)
	}
	return sb.String()
}
