// Copyright (c) 2023 Michael D Henderson. All rights reserved.

package ch02

import (
	"bytes"
	"testing"
)

func TestPrintExpr(t *testing.T) {
	mksym := func(s string) Atom {
		return make_sym([]byte(s))
	}

	for _, tc := range []struct {
		id     int
		input  Atom
		expect string
	}{
		{1, make_int(42), "42"},
		{2, mksym("FOO"), "FOO"},
		{3, cons(mksym("X"), mksym("Y")), "(X . Y)"},
		{4, cons(make_int(1), cons(make_int(2), cons(make_int(3), _nil))), "(1 2 3)"},
	} {
		bb := &bytes.Buffer{}
		if _, err := print_expr(bb, tc.input); err != nil {
			t.Errorf("%d: error: want nil: got %v\n", tc.id, err)
		} else {
			got := string(bb.Bytes())
			if tc.expect != got {
				t.Errorf("%d: expr: want %s: got %s\n", tc.id, tc.expect, got)
			}
		}
	}
}
