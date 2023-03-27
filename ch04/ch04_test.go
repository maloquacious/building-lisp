// Copyright (c) 2023 Michael D Henderson. All rights reserved.

package ch04

import (
	"errors"
	"testing"
)

func TestChapter04(t *testing.T) {
	env := env_create(_nil)

	for _, tc := range []struct {
		id     int
		input  string
		expect string
		err    error
	}{
		{id: 1, input: "foo", expect: "NIL", err: Error_Unbound},
		{id: 2, input: "(quote foo)", expect: "FOO"},
		{id: 3, input: "(define foo 42)", expect: "FOO"},
		{id: 4, input: "foo", expect: "42"},
		{id: 5, input: "(define foo (quote bar))", expect: "FOO"},
		{id: 6, input: "foo", expect: "BAR"},
	} {
		expr, _, err := read_expr([]byte(tc.input))
		if err != nil {
			t.Errorf("%d: read error: want nil: got %v\n", tc.id, err)
			continue
		}
		result, err := eval_expr(expr, env)
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
