// Copyright (c) 2023 Michael D Henderson. All rights reserved.

package lisp

// eval_do_exec sets expr to the next part of the function
// body, and pops the stack when we have reached end of the body.
func eval_do_exec(stack, expr, env *Atom) error {
	*env = list_get(*stack, FRAME_ENV)
	body := list_get(*stack, FRAME_BODY)
	*expr = car(body)
	if body = cdr(body); nilp(body) {
		// finished function; pop the stack
		*stack = car(*stack)
	} else {
		list_set(*stack, FRAME_BODY, body)
	}
	return nil
}

// eval_do_bind binds the function arguments into a new environment
// if they have not already been bound, then calls eval_do_exec to
// get the next expression in the body.
func eval_do_bind(stack, expr, env *Atom) error {
	body := list_get(*stack, FRAME_BODY)
	if !nilp(body) {
		return eval_do_exec(stack, expr, env)
	}
	op := list_get(*stack, FRAME_OP)
	args := list_get(*stack, FRAME_ARGS)

	*env = env_create(car(op))
	arg_names := car(cdr(op))
	body = cdr(cdr(op))
	list_set(*stack, FRAME_ENV, *env)
	list_set(*stack, FRAME_BODY, body)

	// bind the arguments
	for !nilp(arg_names) {
		if arg_names._type == AtomType_Symbol {
			_ = env_set(*env, arg_names, args)
			args = _nil
			break
		} else if nilp(args) {
			// it is an error if we have too few arguments
			return Error_Args
		}
		_ = env_set(*env, car(arg_names), car(args))
		arg_names = cdr(arg_names)
		args = cdr(args)
	}
	if !nilp(args) {
		// it is an error if we have too many arguments
		return Error_Args
	}
	list_set(*stack, FRAME_ARGS, args)

	return eval_do_exec(stack, expr, env)
}

// eval_do_apply is called once all arguments have been evaluated.
// it is responsible either generating an expression to call a builtin,
// or delegating to eval_do_bind.
func eval_do_apply(stack, expr, env, result *Atom) error {
	op := list_get(*stack, FRAME_OP)
	args := list_get(*stack, FRAME_ARGS)

	if !nilp(args) {
		list_reverse(&args)
		list_set(*stack, 4, args)
	}

	if op._type == AtomType_Symbol {
		if op.value.symbol.EqualString("APPLY") {
			// replace the current frame
			*stack = car(*stack)
			*stack = make_frame(*stack, *env, _nil)
			// update the op and args in the new frame
			op = car(args)
			list_set(*stack, FRAME_OP, op)
			if args = car(cdr(args)); !listp(args) {
				return Error_Syntax
			}
			list_set(*stack, FRAME_ARGS, args)
		}
	}

	// we must have a builtin or closure to continue
	if op._type == AtomType_Builtin {
		*stack = car(*stack)
		*expr = cons(op, args)
		return nil
	} else if op._type != AtomType_Closure {
		return Error_Type
	}

	return eval_do_bind(stack, expr, env)
}

// eval_do_return is called after an expression has been evaluated.
// is responsible for storing the result, which is either an operator,
// an argument, or an intermediate body expression, and fetching the
// next expression to evaluate.
func eval_do_return(stack, expr, env, result *Atom) error {
	var op, body, args, sym Atom

	*env = list_get(*stack, FRAME_ENV)
	op = list_get(*stack, FRAME_OP)
	body = list_get(*stack, FRAME_BODY)

	if !nilp(body) {
		// still running a procedure; ignore the intermediate result
		return eval_do_apply(stack, expr, env, result)
	}

	if nilp(op) {
		// finished evaluating operator
		op = *result
		list_set(*stack, 2, op)

		if op._type == AtomType_Macro {
			// don't evaluate macro arguments
			args = list_get(*stack, FRAME_TAIL)
			*stack = make_frame(*stack, *env, _nil)
			op._type = AtomType_Closure
			list_set(*stack, FRAME_OP, op)
			list_set(*stack, FRAME_ARGS, args)
			return eval_do_bind(stack, expr, env)
		}
	} else if op._type == AtomType_Symbol {
		// finished working on special form
		if op.value.symbol.EqualString("DEFINE") {
			sym = list_get(*stack, 4)
			_ = env_set(*env, sym, *result)
			*stack = car(*stack)
			*expr = cons(make_sym([]byte("QUOTE")), cons(sym, _nil))
			return nil
		} else if op.value.symbol.EqualString("IF") {
			args = list_get(*stack, FRAME_TAIL)
			if nilp(*result) {
				*expr = car(cdr(args))
			} else {
				*expr = car(args)
			}
			*stack = car(*stack)
			return nil
		}
		// store evaluated argument
		args = list_get(*stack, FRAME_ARGS)
		list_set(*stack, FRAME_ARGS, cons(*result, args))
	} else if op._type == AtomType_Macro {
		// finished evaluating macro
		*expr = *result
		*stack = car(*stack)
		return nil
	} else {
		// store evaluated argument
		args = list_get(*stack, FRAME_ARGS)
		list_set(*stack, FRAME_ARGS, cons(*result, args))
	}

	args = list_get(*stack, FRAME_TAIL)
	if nilp(args) {
		// no more arguments left to evaluate
		return eval_do_apply(stack, expr, env, result)
	}

	// evaluate next argument
	*expr = car(args)
	list_set(*stack, 3, cdr(args))
	return nil
}

// eval_expr evaluates an expression with a given environment and updates the result.
// much of the work is for setting up special forms; the rest is a loop to process
// then entire stack frame.
// note that the result may not be updated if we find errors.
func eval_expr(expr, env Atom, result *Atom) error {
	var stack Atom

	// do {...} while (!err);
	for {
		if expr._type == AtomType_Symbol {
			if err := env_get(env, expr, result); err != nil {
				return err
			}
		} else if expr._type != AtomType_Pair {
			*result = expr
		} else if !listp(expr) {
			return Error_Syntax
		} else {
			op, args := car(expr), cdr(expr)
			if op._type == AtomType_Symbol {
				// handle special forms
				if op.value.symbol.EqualString("QUOTE") {
					// verify number and type of args
					if nilp(args) || !nilp(cdr(args)) {
						return Error_Args
					}
					*result = car(args)
				} else if op.value.symbol.EqualString("DEFINE") {
					// verify number and type of args
					if nilp(args) || nilp(cdr(args)) {
						return Error_Args
					}
					if sym := car(args); sym._type == AtomType_Pair {
						if err := make_closure(env, cdr(sym), cdr(args), result); err != nil {
							return err
						} else if sym = car(sym); sym._type != AtomType_Symbol {
							return Error_Type
						}
						_ = env_set(env, sym, *result)
						*result = sym
					} else if sym._type == AtomType_Symbol {
						if !nilp(cdr(cdr(args))) {
							return Error_Args
						}
						stack = make_frame(stack, env, _nil)
						list_set(stack, FRAME_OP, op)
						list_set(stack, FRAME_ARGS, sym)
						expr = car(cdr(args))
						continue
					} else {
						return Error_Type
					}
				} else if op.value.symbol.EqualString("LAMBDA") {
					// verify number and type of args
					if nilp(args) || nilp(cdr(args)) {
						return Error_Args
					}
					if err := make_closure(env, car(args), cdr(args), result); err != nil {
						return err
					}
				} else if op.value.symbol.EqualString("IF") {
					// verify number and type of args
					if nilp(args) || nilp(cdr(args)) || nilp(cdr(cdr(args))) || !nilp(cdr(cdr(cdr(args)))) {
						return Error_Args
					}
					stack = make_frame(stack, env, cdr(args))
					list_set(stack, 2, op)
					expr = car(args)
					continue
				} else if op.value.symbol.EqualString("DEFMACRO") {
					// verify number and type of args
					if nilp(args) || nilp(cdr(args)) {
						return Error_Args
					} else if car(args)._type != AtomType_Pair {
						return Error_Syntax
					}
					name := car(car(args))
					if name._type != AtomType_Symbol {
						return Error_Type
					}
					var macro Atom
					if err := make_closure(env, cdr(car(args)), cdr(args), &macro); err != nil {
						return err
					}
					macro._type = AtomType_Macro
					*result = name
					_ = env_set(env, name, macro)
				} else if op.value.symbol.EqualString("APPLY") {
					// verify number and type of args
					if nilp(args) || nilp(cdr(args)) || !nilp(cdr(cdr(args))) {
						return Error_Args
					}
					stack = make_frame(stack, env, cdr(args))
					list_set(stack, FRAME_OP, op)
					expr = car(args)
					continue
				} else {
					// push a new stack frame to handle function application
					stack = make_frame(stack, env, args)
					expr = op
					continue
				}
			} else if op._type == AtomType_Builtin {
				if err := op.value.builtin.fn(args, result); err != nil {
					return err
				}
			} else {
				// push a new stack frame to handle function application
				stack = make_frame(stack, env, args)
				expr = op
				continue
			}
		}

		// terminate this loop if we've exhausted the stack
		if nilp(stack) {
			return nil
		}

		// try storing the result and fetching the next expression from the stack
		if err := eval_do_return(&stack, &expr, &env, result); err != nil {
			return nil
		}
	}
}
