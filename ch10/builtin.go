// Copyright (c) 2023 Michael D Henderson. All rights reserved.

package lisp

// Builtin is a helper for calling a native Go function.
// We define a struct around it so that we can do
// pointer comparisons for equality in other parts of this package.
type Builtin struct {
	fn Native
}

// Native is a function in Go that can evaluate expressions.
type Native func(args Atom, result *Atom) error

// builtin_add implements a function for calculating the sum of two numbers.
// note that the result may not be updated if we find errors.
func builtin_add(args Atom, result *Atom) error {
	// verify number and type of arguments
	if nilp(args) || nilp(cdr(args)) || !nilp(cdr(cdr(args))) {
		return Error_Args
	}
	a, b := car(args), car(cdr(args))
	if a._type != AtomType_Integer || b._type != AtomType_Integer {
		return Error_Type
	}

	*result = make_int(a.value.integer + b.value.integer)
	return nil
}

// builtin_car makes our native car function available to the interpreter.
// note that the result may not be updated if we find errors.
func builtin_car(args Atom, result *Atom) error {
	// verify number and type of arguments
	if nilp(args) || !nilp(cdr(args)) {
		return Error_Args
	} else if car(args)._type != AtomType_Pair {
		return Error_Type
	}

	if nilp(car(args)) {
		*result = _nil
	} else {
		*result = car(car(args))
	}
	return nil
}

// builtin_cdr makes our native cdr function available to the interpreter.
// note that the result may not be updated if we find errors.
func builtin_cdr(args Atom, result *Atom) error {
	// verify number and type of arguments
	if nilp(args) || !nilp(cdr(args)) {
		return Error_Args
	} else if car(args)._type != AtomType_Pair {
		return Error_Type
	}

	if nilp(car(args)) {
		*result = _nil
	} else {
		*result = cdr(car(args))
	}
	return nil
}

// builtin_cons makes our native cons function available to the interpreter.
// note that the result may not be updated if we find errors.
func builtin_cons(args Atom, result *Atom) error {
	// verify number and type of arguments
	if nilp(args) || nilp(cdr(args)) || !nilp(cdr(cdr(args))) {
		return Error_Args
	}

	*result = cons(car(args), car(cdr(args)))
	return nil
}

// builtin_divide implements a function for calculating the quotient of two numbers.
// note that the result may not be updated if we find errors.
// will panic on divide by zero.
func builtin_divide(args Atom, result *Atom) error {
	// verify number and type of arguments
	if nilp(args) || nilp(cdr(args)) || !nilp(cdr(cdr(args))) {
		return Error_Args
	}
	a, b := car(args), car(cdr(args))
	if a._type != AtomType_Integer || b._type != AtomType_Integer {
		return Error_Type
	}

	*result = make_int(a.value.integer / b.value.integer)
	return nil
}

// builtin_less implements a comparison operator for numbers,
// returning T if the first argument is less than the second.
// note that the result may not be updated if we find errors.
func builtin_less(args Atom, result *Atom) error {
	// verify number and type of arguments
	if nilp(args) || nilp(cdr(args)) || !nilp(cdr(cdr(args))) {
		return Error_Args
	}
	a, b := car(args), car(cdr(args))
	if a._type != AtomType_Integer || b._type != AtomType_Integer {
		return Error_Type
	}

	if a.value.integer < b.value.integer {
		// todo: should be able to assume that T is in the environment
		*result = make_sym([]byte{'T'})
	} else {
		*result = _nil
	}
	return nil
}

// builtin_multiply implements a function for calculating the product of two numbers.
// note that the result may not be updated if we find errors.
func builtin_multiply(args Atom, result *Atom) error {
	// verify number and type of arguments
	if nilp(args) || nilp(cdr(args)) || !nilp(cdr(cdr(args))) {
		return Error_Args
	}
	a, b := car(args), car(cdr(args))
	if a._type != AtomType_Integer || b._type != AtomType_Integer {
		return Error_Type
	}

	*result = make_int(a.value.integer * b.value.integer)
	return nil
}

// builtin_numeq implements a comparison operator for numbers,
// returning T if they are equal.
// note that the result may not be updated if we find errors.
func builtin_numeq(args Atom, result *Atom) error {
	// verify number and type of arguments
	if nilp(args) || nilp(cdr(args)) || !nilp(cdr(cdr(args))) {
		return Error_Args
	}
	a, b := car(args), car(cdr(args))
	if a._type != AtomType_Integer || b._type != AtomType_Integer {
		return Error_Type
	}

	if a.value.integer == b.value.integer {
		// todo: should be able to assume that T is in the environment
		*result = make_sym([]byte{'T'})
	} else {
		*result = _nil
	}
	return nil
}

// builtin_subtract implements a function for calculating the difference of two numbers.
// note that the result may not be updated if we find errors.
func builtin_subtract(args Atom, result *Atom) error {
	// verify number and type of arguments
	if nilp(args) || nilp(cdr(args)) || !nilp(cdr(cdr(args))) {
		return Error_Args
	}
	a, b := car(args), car(cdr(args))
	if a._type != AtomType_Integer || b._type != AtomType_Integer {
		return Error_Type
	}

	*result = make_int(a.value.integer - b.value.integer)
	return nil
}
