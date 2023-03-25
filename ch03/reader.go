// Copyright (c) 2023 Michael D Henderson. All rights reserved.

package ch03

import (
	"bytes"
	"strconv"
)

// read_atom reads an atom (a number or symbol) from the input.
// if it's a symbol, we assume that the caller has parsed it already
// and do no checking that it is a valid symbol.
func read_atom(input []byte) Atom {
	if val, err := strconv.Atoi(string(input)); err == nil { // it is an integer
		return make_int(val)
	}

	// it is a symbol, but we must treat NIL specially.
	sym := bytes.ToUpper(input)
	if bytes.Equal(sym, []byte{'N', 'I', 'L'}) { // it is NIL
		return _nil
	}
	return make_sym(sym)
}

// read_list reads the next list from the input.
// it returns the remainder of the input or an error.
func read_list(input []byte) (list Atom, remainder []byte, err error) {
	var token []byte
	var expr, tail Atom
	for token, remainder = lex(input); token != nil; token, remainder = lex(input) {
		// check for ")"
		if bytes.Equal(token, []byte{')'}) {
			// return the list, remainder, and no error
			return list, remainder, nil
		}

		// check for "."
		if bytes.Equal(token, []byte{'.'}) {
			// a dotted list must look like "(x . y)" or it is an improper list
			if nilp(tail) {
				// dot can't start a list, so this is an improper list
				return _nil, nil, Error_Syntax
			}

			// read the next expression and set the cdr of the current atom to it
			expr, remainder, err = read_expr(remainder)
			if err != nil {
				// return the error
				return _nil, nil, err
			}
			setcdr(tail, expr)

			// read the closing paren
			token, remainder = lex(remainder)
			if !bytes.Equal(token, []byte{')'}) {
				// no closing paren, so this is an improper list
				return _nil, nil, Error_Syntax
			}

			// return the list, remainder, and no error
			return list, remainder, nil
		}

		// read the next expression
		if expr, remainder, err = read_expr(input); err != nil {
			// return the error
			return _nil, nil, err
		}

		// and append it to the tail of the list
		if nilp(list) {
			// first item in the list, so create a new list
			list = cons(expr, _nil)
			tail = list
		} else {
			// append to tail, then update tail
			setcdr(tail, cons(expr, _nil))
			tail = cdr(tail)
		}

		// at this point:
		//    list is the head of the list
		//    tail is the last item in the list

		// update the input and loop back to parse the remainder of the input
		input = remainder
	}

	// eof is an error since lists must end with a close paren.
	return _nil, nil, Error_Syntax
}

// read_expr reads the next expression from the input. an expression is
// either an atom or a list of expressions. returns the expression along
// with the remainder of the input and any errors.
// returns NIL and Error_EndOfInput on end of input. the caller must
// decide how to handle it.
func read_expr(input []byte) (expr Atom, remainder []byte, err error) {
	token, rest := lex(input)
	if token == nil { // end of input
		return _nil, nil, Error_EndOfInput
	}

	switch token[0] {
	case '(':
		return read_list(rest)
	case ')':
		// unexpected close paren
		return _nil, nil, Error_Syntax
	}
	return read_atom(token), rest, err
}
