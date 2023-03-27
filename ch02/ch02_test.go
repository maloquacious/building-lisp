// Copyright (c) 2023 Michael D Henderson. All rights reserved.

package lisp

import (
	"bytes"
	"testing"
)

func TestChapter02(t *testing.T) {
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
		// test Writer interface
		w := &bytes.Buffer{}
		if _, err := tc.input.Write(w); err != nil {
			t.Errorf("%d: write: error: want nil: got %v\n", tc.id, err)
		} else if got := string(w.Bytes()); tc.expect != got {
			t.Errorf("%d: write: want %s: got %s\n", tc.id, tc.expect, got)
		}

		// test Byter interface
		if got := string(tc.input.Bytes()); tc.expect != got {
			t.Errorf("%d: byter: want %s: got %s\n", tc.id, tc.expect, got)
		}

		// test Stringer interface
		if got := tc.input.String(); tc.expect != got {
			t.Errorf("%d: stringer: want %s: got %s\n", tc.id, tc.expect, got)
		}
	}
}
