// Copyright (c) 2023 Michael D Henderson. All rights reserved.

package ch03

import "fmt"

var (
	// Error_EndOfInput is returned at end of input.
	Error_EndOfInput = fmt.Errorf("eof")
	// Error_Syntax is returned for almost every error parsing.
	Error_Syntax = fmt.Errorf("syntax")
)
