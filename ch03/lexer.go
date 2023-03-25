// Copyright (c) 2023 Michael D Henderson. All rights reserved.

package ch03

import (
	"bytes"
)

var (
	// whitespace separates tokens and is generally ignored
	whitespace = []byte{' ', '\t', '\r', '\n'}
	// reserved characters are part of the syntax of our lisp
	reserved = []byte{'(', ')'}
	// delimiters are characters that are not allowed in a symbol.
	// at the minimum, this must include all whitespace and
	// reserved characters.
	delimiters = []byte{' ', '\t', '\r', '\n', '(', ')'}
)

// lex extracts the next token from the input after skipping
// comments and leading whitespace. on end of input, it
// returns nil for both the token and the remainder.
func lex(input []byte) (token []byte, remainder []byte) {
	// skip whitespace
	input = skipws(input)

	// check for end of input
	if len(input) == 0 {
		return nil, nil
	}

	// check for reserved tokens
	if bytes.IndexByte(reserved, input[0]) >= 0 { // input is reserved
		token, remainder = input[:1], input[1:]
		return token, remainder
	}

	// if we get here, the token is a symbol.
	// collect and return all the characters up to the first delimiter.
	token, remainder = runto(input, delimiters)
	return token, remainder
}

// runof splits the input in two. the first part is the prefix from input that
// includes delimiters. the second is the remainder of the input.
func runof(input, delim []byte) ([]byte, []byte) {
	for n, ch := range input {
		if bytes.IndexByte(delim, ch) == -1 {
			// ch is NOT a delimiter, so split here
			return input[:n], input[n:]
		}
	}
	// the entire input consists delimiters
	return input, nil
}

// runto splits the input in two. the first part is the prefix from input that
// does not include any delimiter. the second is the remainder of the input,
func runto(input, delim []byte) ([]byte, []byte) {
	for n, ch := range input {
		if bytes.IndexByte(delim, ch) != -1 {
			// ch IS a delimiter, so split here
			return input[:n], input[n:]
		}
	}
	// there are no delimiters in the input
	return input, nil
}

// skipws skips whitespace characters.
func skipws(input []byte) []byte {
	_, input = runof(input, whitespace)
	return input
}
