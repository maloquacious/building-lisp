// Copyright (c) 2023 Michael D Henderson. All rights reserved.

package lisp

// Pair is the two elements of a cell.
// "car" is the left-hand value and "cdr" is the right-hand.
type Pair struct {
	car, cdr Atom
}
