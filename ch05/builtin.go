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
