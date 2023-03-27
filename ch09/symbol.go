// Copyright (c) 2023 Michael D Henderson. All rights reserved.

package lisp

import "bytes"

// Symbol implements data for a symbol.
// We define a struct around it so that we can do
// pointer comparisons for equality in other parts of this package.
type Symbol struct {
	label []byte
}

func (s *Symbol) EqualString(str string) bool {
	return bytes.Equal(s.label, []byte(str))
}
