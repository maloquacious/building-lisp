// Copyright (c) 2023 Michael D Henderson. All rights reserved.

package ch02

import (
	"bytes"
)

// Bytes implements the Byter interface.
func (a Atom) Bytes() []byte {
	bb := &bytes.Buffer{}
	if _, err := a.Write(bb); err != nil {
		panic(err)
	}
	return bb.Bytes()
}
