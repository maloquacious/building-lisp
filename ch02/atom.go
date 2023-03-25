// Copyright (c) 2023 Michael D Henderson. All rights reserved.

package ch02

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
	symbol  []byte
}
