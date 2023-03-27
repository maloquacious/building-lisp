// Copyright (c) 2023 Michael D Henderson. All rights reserved.

package lisp

import (
	"bytes"
	"fmt"
	"io"
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
	// Builtin is a native function.
	AtomType_Builtin
	// Closure is a closure.
	AtomType_Closure
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
	builtin *Builtin
	integer int
	pair    *Pair
	symbol  *Symbol
}

// Bytes implements the Byter interface.
func (a Atom) Bytes() []byte {
	bb := &bytes.Buffer{}
	if _, err := a.Write(bb); err != nil {
		panic(err)
	}
	return bb.Bytes()
}

// String implements the Stringer interface.
func (a Atom) String() string {
	sb := &strings.Builder{}
	if _, err := a.Write(sb); err != nil {
		panic(err)
	}
	return sb.String()
}

// Write writes the value of an Atom to the writer.
// If the atom is a pair, Write is called recursively
// to write out the entire list.
func (a Atom) Write(w io.Writer) (int, error) {
	switch a._type {
	case AtomType_Nil:
		// atom is nil, so write "NIL"
		return w.Write([]byte{'N', 'I', 'L'})
	case AtomType_Builtin:
		// atom is a native function
		return w.Write([]byte(fmt.Sprintf("#<BUILTIN:%p>", a.value.builtin)))
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
