// Copyright (c) 2023 Michael D Henderson. All rights reserved.

package lisp

// env_create creates a new environment.
// if parent is not NIL, then parent is added to the environment.
func env_create(parent Atom) Atom {
	return cons(parent, _nil)
}

// env_get retrieves the binding for a symbol from the environment.
func env_get(env, symbol Atom) (Atom, error) {
	for bs := cdr(env); !nilp(bs); bs = cdr(bs) {
		if b := car(bs); car(b).value.symbol == symbol.value.symbol {
			return cdr(b), nil
		}
	}
	// search the parent environment (if we have one).
	if parent := car(env); !nilp(parent) {
		return env_get(parent, symbol)
	}
	// not found, so return an unbound error
	return _nil, Error_Unbound
}

// env_set creates a binding for a symbol in the environment.
func env_set(env, symbol, value Atom) {
	for bs := cdr(env); !nilp(bs); bs = cdr(bs) {
		if b := car(bs); car(b).value.symbol == symbol.value.symbol {
			b.value.pair.cdr = value
			return
		}
	}
	setcdr(env, cons(cons(symbol, value), cdr(env)))
}
