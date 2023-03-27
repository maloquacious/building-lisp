// Copyright (c) 2023 Michael D Henderson. All rights reserved.

package lisp

import "bytes"

// Symbol implements data for a symbol.
type Symbol struct {
	label []byte
}

func (s *Symbol) EqualString(str string) bool {
	return bytes.Equal(s.label, []byte(str))
}
