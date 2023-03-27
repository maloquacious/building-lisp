// Copyright (c) 2023 Michael D Henderson. All rights reserved.

package lisp

import (
	"errors"
	"testing"
)

func TestChapter05(t *testing.T) {
	env := env_create_default()

	for _, tc := range []struct {
		id     int
		input  string
		expect string
		err    error
	}{
		{id: 1, input: "(define foo 1)", expect: "FOO"},
		{id: 2, input: "(define bar 2)", expect: "BAR"},
		{id: 3, input: "(cons foo bar)", expect: "(1 . 2)"},
		{id: 4, input: "(define baz (quote (a b c)))", expect: "BAZ"},
		{id: 5, input: "(car baz)", expect: "A"},
		{id: 6, input: "(cdr baz)", expect: "(B C)"},
	} {
		var expr Atom
		_, err := read_expr([]byte(tc.input), &expr)
		if err != nil {
			t.Errorf("%d: read error: want nil: got %v\n", tc.id, err)
			continue
		}

		var result Atom
		err = eval_expr(expr, env, &result)

		if tc.err == nil && err == nil {
			// yay
		} else if tc.err == nil && err != nil {
			t.Errorf("%d: error: want nil: got %v\n", tc.id, err)
		} else if !errors.Is(err, tc.err) {
			t.Errorf("%d: error: want %v: got %v\n", tc.id, tc.err, err)
		}
		if got := result.String(); tc.expect != got {
			t.Errorf("%d: eval: want %q: got %q\n", tc.id, tc.expect, got)
		}
	}
}
