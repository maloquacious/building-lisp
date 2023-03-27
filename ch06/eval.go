// Copyright (c) 2023 Michael D Henderson. All rights reserved.

package lisp

import "fmt"

// apply calls a native function with a list of arguments and updates the result.
// note that the result may not be updated if we find errors.
func apply(fn, args Atom, result *Atom) error {
	if fn._type != AtomType_Builtin {
		// it is an error to call this with anything other than a builtin
		return Error_Type
	}
	return fn.value.builtin.fn(args, result)
}

// eval_expr evaluates an expression with a given environment and updates the result.
// note that the result may not be updated if we find errors.
func eval_expr(expr, env Atom, result *Atom) error {
	if expr._type == AtomType_Symbol {
		return env_get(env, expr, result)
	} else if expr._type != AtomType_Pair {
		*result = expr
		return nil
	} else if !listp(expr) {
		return Error_Syntax
	}

	op, args := car(expr), cdr(expr)
	if op._type == AtomType_Symbol {
		// evaluate special forms
		if op.value.symbol.EqualString("QUOTE") {
			if nilp(args) || !nilp(cdr(args)) {
				return Error_Args
			}
			*result = car(args)
			return nil
		} else if op.value.symbol.EqualString("DEFINE") {
			if nilp(args) || nilp(cdr(args)) || !nilp(cdr(cdr(args))) {
				return Error_Args
			}
			sym := car(args)
			if sym._type != AtomType_Symbol {
				return Error_Type
			}
			var val Atom
			if err := eval_expr(car(cdr(args)), env, &val); err != nil {
				return err
			}
			env_set(env, sym, val)
			*result = sym
			return nil
		}
	}

	// evaluate and update the operator
	if err := eval_expr(op, env, &op); err != nil {
		fmt.Printf("eval op %q %v\n", op.String(), err)
		return err
	}

	// evaluate arguments by calling eval on a copy of each.
	// we have to make the copy, so we don't destroy the input.
	args = copy_list(args)
	for arg := args; !nilp(arg); arg = cdr(arg) {
		// evaluate the arg and update its value
		if err := eval_expr(car(arg), env, &arg.value.pair.car); err != nil {
			return err
		}
	}

	// return the result of applying eval on our operator and arguments
	return apply(op, args, result)
}
