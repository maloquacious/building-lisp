// Copyright (c) 2023 Michael D Henderson. All rights reserved.

package lisp

// eval_expr evaluates an expression with a given environment.
// it returns the result and any errors.
func eval_expr(expr, env Atom) (Atom, error) {
	if expr._type == AtomType_Symbol {
		return env_get(env, expr)
	} else if expr._type != AtomType_Pair {
		return expr, nil
	} else if !listp(expr) {
		return _nil, Error_Syntax
	}

	op, args := car(expr), cdr(expr)
	if op._type == AtomType_Symbol {
		if op.value.symbol.EqualString("QUOTE") {
			if nilp(args) || !nilp(cdr(args)) {
				return _nil, Error_Args
			}
			return car(args), nil
		} else if op.value.symbol.EqualString("DEFINE") {
			if nilp(args) || nilp(cdr(args)) || !nilp(cdr(cdr(args))) {
				return _nil, Error_Args
			}
			sym := car(args)
			if sym._type != AtomType_Symbol {
				return _nil, Error_Type
			}
			val, err := eval_expr(car(cdr(args)), env)
			if err != nil {
				return _nil, err
			}
			env_set(env, sym, val)
			return sym, nil
		}
	}

	return _nil, Error_Syntax
}
