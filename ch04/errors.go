// Copyright (c) 2023 Michael D Henderson. All rights reserved.

package lisp

import "fmt"

var (
	// Error_Args is returned when a list expression was shorter or longer than anticipated.
	Error_Args = fmt.Errorf("args")
	// Error_EndOfInput is returned at end of input.
	Error_EndOfInput = fmt.Errorf("eof")
	// Error_Syntax is returned for almost every error parsing.
	Error_Syntax = fmt.Errorf("syntax")
	// Error_Type is returned when an object in an expression isn't the expected type.
	Error_Type = fmt.Errorf("type")
	// Error_Unbound is returned when we attempt to evaluate an unbound symbol.
	Error_Unbound = fmt.Errorf("unbound")
)
