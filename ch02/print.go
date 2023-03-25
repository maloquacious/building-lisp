// Copyright (c) 2023 Michael D Henderson. All rights reserved.

package ch02

import "os"

// print_expr is a helper function to write an expression
// to stdout, ignoring errors.
func print_expr(expr Atom) {
	_, _ = expr.Write(os.Stdout)
}
