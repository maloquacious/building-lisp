// Copyright (c) 2023 Michael D Henderson. All rights reserved.

package lisp

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
)

func load_file(env Atom, path string) error {
	fmt.Printf("Reading %s...\n", path)
	input, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var expr Atom
	rest, err := read_expr(input, &expr)
	for ; err == nil; rest, err = read_expr(rest, &expr) {
		var result Atom
		if err := eval_expr(expr, env, &result); err != nil {
			fmt.Printf("error: %s in expression:\n\t%s\n", err, expr.String())
		} else {
			fmt.Printf("%s\n", result.String())
		}
	}
	if err != nil && err != Error_EndOfInput {
		fmt.Printf("error: %s in expression:\n\t%s\n", err, expr.String())
	}

	return nil
}

// read_atom reads an atom (a number or symbol) from the input.
// if it's a symbol, we assume that the caller has parsed it already
// and do no checking that it is a valid symbol.
func read_atom(input []byte, result *Atom) error {
	if val, err := strconv.Atoi(string(input)); err == nil { // it is an integer
		*result = make_int(val)
		return nil
	}
	// it is a symbol, but we must treat NIL specially.
	if label := bytes.ToUpper(input); bytes.Equal(label, []byte{'N', 'I', 'L'}) {
		// it is NIL and NIL must never be added to the symbol table.
		*result = _nil
	} else {
		*result = make_sym(label)
	}
	return nil
}

// read_list reads the next list from the input.
// it returns the remainder of the input or an error.
func read_list(input []byte, result *Atom) (remainder []byte, err error) {
	// set the result to NIL in case we read an empty list.
	*result = _nil

	var token []byte
	var tail Atom
	for token, remainder = lex(input); token != nil; token, remainder = lex(input) {
		// check for ")"
		if bytes.Equal(token, []byte{')'}) {
			// result holds the list.
			// return the remainder and no error
			return remainder, nil
		}

		// check for "."
		if bytes.Equal(token, []byte{'.'}) {
			// a dotted list must look like "(x . y)" or it is an improper list
			if nilp(tail) {
				// dot can't start a list, so this is an improper list
				return nil, Error_Syntax
			}

			// read the next expression and set the cdr of the current atom to it
			var expr Atom
			remainder, err = read_expr(remainder, &expr)
			if err != nil {
				// return the error
				return nil, err
			}
			setcdr(tail, expr)

			// read the closing paren
			token, remainder = lex(remainder)
			if !bytes.Equal(token, []byte{')'}) {
				// no closing paren, so this is an improper list
				return nil, Error_Syntax
			}

			// result holds the list.
			// return the remainder and no error
			return remainder, nil
		}

		// read the next expression
		var expr Atom
		remainder, err = read_expr(input, &expr)
		if err != nil {
			// return the error
			return nil, err
		}

		// and append it to the tail of the list
		if nilp(tail) {
			// first item in the list, so create a new list
			*result = cons(expr, _nil)
			tail = *result
		} else {
			// append to tail, then update tail
			setcdr(tail, cons(expr, _nil))
			tail = cdr(tail)
		}

		// at this point:
		//    result is the head of the list
		//    tail   is the last item in the list

		// update the input and loop back to parse the remainder of the input
		input = remainder
	}

	// eof is an error since lists must end with a close paren.
	return nil, Error_Syntax
}

// read_expr reads the next expression from the input. an expression is
// either an atom or a list of expressions. returns the expression along
// with the remainder of the input and any errors.
// returns NIL and Error_EndOfInput on end of input. the caller must
// decide how to handle it.
// todo: result is not always updated by read. does that lead to bugs later?
func read_expr(input []byte, result *Atom) (remainder []byte, err error) {
	token, rest := lex(input)
	if token == nil { // end of input
		return nil, Error_EndOfInput
	}

	switch token[0] {
	case '(':
		return read_list(rest, result)
	case ')':
		// unexpected close paren
		return nil, Error_Syntax
	case '\'':
		*result = cons(make_sym([]byte("QUOTE")), cons(_nil, _nil))
		// set car(cdr(result))
		return read_expr(rest, &result.value.pair.cdr.value.pair.car)
	}
	err = read_atom(token, result)
	return rest, err
}

// read reads the next expression from the input.
// an expression is either an atom or a list of expressions.
// returns an error for any syntax error (such as unterminated list).
// returns NIL, nil, and Error_EndOfInput on end of input.
// otherwise, returns the expression and the remainder of the input.
func read(input []byte) (expr Atom, remainder []byte, err error) {
	// stack and slice are used for building lists as we read them.
	// slice tricks cheat sheet -> https://ueokande.github.io/go-slice-tricks/
	var stack []Atom // stack of in-process lists

	for token, rest := lex(input); token != nil; token, rest = lex(rest) {
		var atom Atom

		// handle some syntax.
		//   '(' starts a new list.
		//   '.' splices in a dotted pair.
		if token[0] == '(' {
			// push a new list onto the stack
			stack = append(stack, _nil)
			continue // process the next token
		} else if bytes.Equal(token, []byte{'.'}) {
			// a dotted pair must look like "(x . y)" or it is an error.
			// when we find a dotted pair, we must splice the second atom into
			// the first atom. That first atom should already be stored as the
			// tail atom on the list on top of the stack. it's an error if there
			// is no list or if the next token in the input is not a close paren.
			if len(stack) == 0 || nilp(stack[len(stack)-1]) {
				// dot can't start a list, so this is an improper list
				return _nil, nil, Error_Syntax
			}

			// the cdr of the dotted pair is the next expression
			if atom, rest, err = read(rest); err != nil {
				return _nil, nil, err
			}

			// the dotted pair must be followed by a close paren.
			// verify by looking ahead at the next token.
			if lookAhead, _ := lex(rest); !bytes.Equal(lookAhead, []byte{')'}) {
				// no closing paren, so this is an improper list
				return _nil, nil, Error_Syntax
			}

			// get the list from the stack so that we can hack it.
			// find the last entry in that list;
			// that's the one we'll change in to a dotted pair.
			for tail := stack[len(stack)-1]; !nilp(tail); tail = cdr(tail) {
				if nilp(cdr(tail)) {
					// set the cdr of it to change it to a dotted pair.
					setcdr(tail, atom)
					break
				}
			}
			continue // process the next token
		}

		switch token[0] {
		case ')':
			// found end of a list
			if len(stack) == 0 {
				// empty stack means unexpected close paren
				return _nil, nil, Error_Syntax
			}
			// pop the list from the stack
			atom, stack = stack[len(stack)-1], stack[:len(stack)-1]

		default:
			if val, err := strconv.Atoi(string(token)); err == nil {
				// it is an integer
				atom = make_int(val)
			} else {
				// it is a symbol
				sym := bytes.ToUpper(token)
				if bytes.Equal(sym, []byte{'N', 'I', 'L'}) {
					// treat NIL specially.
					atom = _nil
				} else {
					atom = make_sym(sym)
				}
			}
		}

		// return the atom if we are not reading a list
		if len(stack) == 0 {
			return atom, rest, nil
		}

		// append the atom to the list at the top of the stack
		var list Atom
		list, stack = stack[len(stack)-1], stack[:len(stack)-1]
		if nilp(list) {
			list = cons(atom, _nil)
		} else {
			for tail := list; !nilp(tail); tail = cdr(tail) {
				if nilp(cdr(tail)) {
					setcdr(tail, cons(atom, _nil))
					break
				}
			}
		}
		stack = append(stack, list)
	}

	if len(stack) != 0 {
		// unexpected end of input
		return _nil, nil, Error_Syntax
	}

	// input contained no expressions at all
	return _nil, nil, Error_EndOfInput
}
