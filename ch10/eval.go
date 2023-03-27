// Copyright (c) 2023 Michael D Henderson. All rights reserved.

package lisp

import "fmt"

// apply calls a native function with a list of arguments and updates the result.
// note that the result may not be updated if we find errors.
func apply(fn, args Atom, result *Atom) error {
	// handle builtins
	if fn._type == AtomType_Builtin {
		return fn.value.builtin.fn(args, result)
	}

	// handle closure
	if fn._type == AtomType_Closure {
		// create a new environment for the closure
		env := env_create(car(fn))

		// bind the arguments
		for arg_names := car(cdr(fn)); !nilp(arg_names); arg_names = cdr(arg_names) {
			// if arg name is a symbol, apply it as rest of arguments
			if arg_names._type == AtomType_Symbol {
				_ = env_set(env, arg_names, args)
				args = _nil
				break
			}
			if nilp(args) {
				// not enough arguments passed in to bind against
				return Error_Args
			}
			// put the name and value into the environment
			_ = env_set(env, car(arg_names), car(args))
			// move on to the next argument
			args = cdr(args)
		}
		if !nilp(args) {
			// too many arguments to bind against
			return Error_Args
		}

		// evaluate every expression in the body
		for body := cdr(cdr(fn)); !nilp(body); body = cdr(body) {
			if err := eval_expr(car(body), env, result); err != nil {
				return err
			}
		}

		return nil
	}

	// any other type is an error
	return Error_Type
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
			// verify number and type of arguments
			if nilp(args) || !nilp(cdr(args)) {
				return Error_Args
			}
			*result = car(args)
			return nil
		} else if op.value.symbol.EqualString("DEFINE") {
			// verify number and type of arguments
			if nilp(args) || nilp(cdr(args)) {
				return Error_Args
			}
			var val Atom
			var err error
			sym := car(args)
			if sym._type == AtomType_Pair {
				err = make_closure(env, cdr(sym), cdr(args), &val)
				sym = car(sym)
				if sym._type != AtomType_Symbol {
					return Error_Type
				}
			} else if sym._type == AtomType_Symbol {
				if !nilp(cdr(cdr(args))) {
					return Error_Args
				}
				err = eval_expr(car(cdr(args)), env, &val)
			} else {
				return Error_Type
			}
			if err != nil {
				return err
			}
			*result = sym
			return env_set(env, sym, val)
		} else if op.value.symbol.EqualString("LAMBDA") {
			// verify number and type of arguments
			if nilp(args) || nilp(cdr(args)) {
				return Error_Args
			}
			return make_closure(env, car(args), cdr(args), result)
		} else if op.value.symbol.EqualString("IF") {
			// verify number and type of arguments
			if nilp(args) || nilp(cdr(args)) || nilp(cdr(cdr(args))) || !nilp(cdr(cdr(cdr(args)))) {
				return Error_Args
			}
			var cond Atom
			if err := eval_expr(car(args), env, &cond); err != nil {
				return err
			}
			var val Atom
			if nilp(cond) {
				val = car(cdr(cdr(args)))
			} else {
				val = car(cdr(args))
			}
			return eval_expr(val, env, result)
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
