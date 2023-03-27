// Copyright (c) 2023 Michael D Henderson. All rights reserved.

package lisp

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

	return Error_Syntax
}
