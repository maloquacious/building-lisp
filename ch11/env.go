// Copyright (c) 2023 Michael D Henderson. All rights reserved.

package lisp

// env_create creates a new environment.
// if parent is not NIL, then parent is added to the environment.
func env_create(parent Atom) Atom {
	return cons(parent, _nil)
}

// env_create_default creates a new environment with some native
// functions added to the symbol table.
func env_create_default() Atom {
	// create a new environment
	env := env_create(_nil)
	// add the default list of native functions to the environment
	env_set(env, make_sym([]byte("CAR")), make_builtin(builtin_car))
	env_set(env, make_sym([]byte("CDR")), make_builtin(builtin_cdr))
	env_set(env, make_sym([]byte("CONS")), make_builtin(builtin_cons))
	env_set(env, make_sym([]byte{'+'}), make_builtin(builtin_add))
	env_set(env, make_sym([]byte{'-'}), make_builtin(builtin_subtract))
	env_set(env, make_sym([]byte{'*'}), make_builtin(builtin_multiply))
	env_set(env, make_sym([]byte{'/'}), make_builtin(builtin_divide))
	env_set(env, make_sym([]byte{'T'}), make_sym([]byte{'T'}))
	env_set(env, make_sym([]byte{'='}), make_builtin(builtin_numeq))
	env_set(env, make_sym([]byte{'<'}), make_builtin(builtin_less))

	// return the new environment
	return env
}

// env_get retrieves the binding for a symbol from the environment.
// does not update result unless it finds a symbol in the environment.
func env_get(env, symbol Atom, result *Atom) error {
	for bs := cdr(env); !nilp(bs); bs = cdr(bs) {
		if b := car(bs); car(b).value.symbol == symbol.value.symbol {
			*result = cdr(b)
			return nil
		}
	}
	// search the parent environment (if we have one).
	if parent := car(env); !nilp(parent) {
		return env_get(parent, symbol, result)
	}
	// not found, so return an unbound error
	return Error_Unbound
}

// env_set creates a binding for a symbol in the environment.
func env_set(env, symbol, value Atom) error {
	for bs := cdr(env); !nilp(bs); bs = cdr(bs) {
		if b := car(bs); car(b).value.symbol == symbol.value.symbol {
			b.value.pair.cdr = value
			return nil
		}
	}
	setcdr(env, cons(cons(symbol, value), cdr(env)))
	return nil
}
