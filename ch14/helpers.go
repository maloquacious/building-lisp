// Copyright (c) 2023 Michael D Henderson. All rights reserved.

package lisp

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

// list_copy returns a shallow copy of a list.
// todo: define "shallow copy" and why we would create one.
func list_copy(list Atom) Atom {
	if nilp(list) {
		return _nil
	}
	a := cons(car(list), _nil)
	p := a
	for list = cdr(list); !nilp(list); list = cdr(list) {
		setcdr(p, cons(car(list), _nil))
		p = cdr(p)
	}
	return a
}

// list_get returns the k'th atom in a list.
// if the list is shorter, it returns NIL.
// assumes that CAR(NIL) is NIL and CDR(NIL) is NIL.
func list_get(list Atom, k int) Atom {
	for ; k != 0; k = k - 1 {
		list = cdr(list)
	}
	return car(list)
}

// list_reverse reverses a list in place.
func list_reverse(list *Atom) {
	tail := _nil
	for !nilp(*list) {
		p := cdr(*list)
		setcdr(*list, tail)
		tail = *list
		*list = p
	}
	*list = tail
}

// list_set updates the k'th value in a list.
// if the list is shorter, it updates the last atom in the list.
// assumes that CAR(NIL) is NIL and CDR(NIL) is NIL.
func list_set(list Atom, k int, value Atom) {
	for ; k != 0; k = k - 1 {
		list = cdr(list)
	}
	setcar(list, value)
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
