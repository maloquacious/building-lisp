// Copyright (c) 2023 Michael D Henderson. All rights reserved.

package lisp

import (
	"bytes"
	"errors"
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

func TestChapter03(t *testing.T) {
	// test the runof function
	for _, tc := range []struct {
		id                string
		input, delimiters string
		token, remainder  string
	}{
		{"runof: 0", "", "", "", ""},
		{"runof: 1", "42", " \t\r\n", "", "42"},
		{"runof: 2", "\n\tfoo", " \t\r\n", "\n\t", "foo"},
		{"runof: 3", "foo \n", " \t\r\n", "", "foo \n"},
	} {
		token, remainder := runof([]byte(tc.input), []byte(tc.delimiters))
		if tc.token != string(token) {
			t.Errorf("%s: token: want %q: got %q\n", tc.id, tc.token, string(token))
		} else if tc.remainder != string(remainder) {
			t.Errorf("%s: remainder: want %q: got %q\n", tc.id, tc.remainder, string(remainder))
		}
	}

	// test the runto function
	for _, tc := range []struct {
		id                string
		input, delimiters string
		token, remainder  string
	}{
		{"runto: 0", "", "", "", ""},
		{"runto: 1", "42", " \t\r\n;()", "42", ""},
		{"runto: 2", "foo", " \t\r\n;()", "foo", ""},
		{"runto: 3", "foo(42)", " \t\r\n;()", "foo", "(42)"},
		{"runto: 4", "\n\t 42;", ";", "\n\t 42", ";"},
	} {
		token, remainder := runto([]byte(tc.input), []byte(tc.delimiters))
		if tc.token != string(token) {
			t.Errorf("%s: token: want %q: got %q\n", tc.id, tc.token, string(token))
		} else if tc.remainder != string(remainder) {
			t.Errorf("%s: remainder: want %q: got %q\n", tc.id, tc.remainder, string(remainder))
		}
	}

	// test the skipws function
	for _, tc := range []struct {
		id        string
		input     string
		remainder string
	}{
		{"skipws: 0", "", ""},
		{"skipws: 1", " \n\r\t 42", "42"},
		{"skipws: 2", "f o o", "f o o"},
		{"skipws: 3", " \t\r\n", ""},
	} {
		remainder := skipws([]byte(tc.input))
		if tc.remainder != string(remainder) {
			t.Errorf("%s: remainder: want %q: got %q\n", tc.id, tc.remainder, string(remainder))
		}
	}

	// test the lexer function
	for _, tc := range []struct {
		id    int
		input string
		token []string
	}{
		{1, "", []string{}},
		{2, "42", []string{"42"}},
		{3, "(foo bar)", []string{"(", "foo", "bar", ")"}},
		{4, "(s (t . u) v . (w . nil))", []string{"(", "s", "(", "t", ".", "u", ")", "v", ".", "(", "w", ".", "nil", ")", ")"}},
		{5, "a(b)c\n", []string{"a", "(", "b", ")", "c", ""}},
	} {
		input := []byte(tc.input)
		var token []byte
		for n, want := range tc.token {
			token, input = lex(input)
			if want != string(token) {
				t.Errorf("%d:%d: token: want %q: got %q\n", tc.id, n, want, string(token))
			}
		}
		if len(input) != 0 {
			t.Errorf("%d: remainder: want %q: got %q\n", tc.id, "", string(input))

		}
	}

	// test the read_expr function
	for _, tc := range []struct {
		id     int
		input  string
		expect string
	}{
		{id: 10, input: "42", expect: "42"},
		{id: 11, input: "(foo bar)", expect: "(FOO BAR)"},
		{id: 12, input: "(s (t . u) v . (w . nil))", expect: "(S (T . U) V W)"},
		{id: 13, input: "()", expect: "NIL"},
	} {
		// reset the symbol table
		sym_table = _nil

		input := []byte(tc.input)
		var expr Atom
		_, _ = read_expr(input, &expr)
		got := expr.String()
		if tc.expect != got {
			t.Errorf("%d: want %q: got %q\n", tc.id, tc.expect, got)
		}
	}

	// reset the symbol table
	sym_table = _nil

	// test the read function
	for _, tc := range []struct {
		id     int
		input  string
		expect string
	}{
		{id: 10, input: "42", expect: "42"},
		{id: 11, input: "(foo bar)", expect: "(FOO BAR)"},
		{id: 12, input: "(s (t . u) v . (w . nil))", expect: "(S (T . U) V W)"},
		{id: 13, input: "()", expect: "NIL"},
		{id: 14, input: "(42)", expect: "(42)"},
		{id: 15, input: "(foo)", expect: "(FOO)"},
		{id: 16, input: "nil", expect: "NIL"},
		{id: 17, input: "(quote foo)", expect: "(QUOTE FOO)"},
		{id: 18, input: "(define foo 42)", expect: "(DEFINE FOO 42)"},
		{id: 19, input: "(define foo (quote bar))", expect: "(DEFINE FOO (QUOTE BAR))"},
	} {
		input := []byte(tc.input)
		expr, remainder, err := read(input)
		got := expr.String()
		if tc.expect != got {
			t.Errorf("%d: want %q: got %q\n", tc.id, tc.expect, got)
		}
		if len(remainder) != 0 {
			t.Errorf("%d: want %q: got %q\n", tc.id, "", string(remainder))
		}
		if err != nil {
			t.Errorf("%d: want nil: got %v\n", tc.id, err)
		}
	}
}

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
