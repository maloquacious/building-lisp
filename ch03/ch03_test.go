// Copyright (c) 2023 Michael D Henderson. All rights reserved.

package ch03

import "testing"

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
		input := []byte(tc.input)
		expr, _, _ := read_expr(input)
		got := expr.String()
		if tc.expect != got {
			t.Errorf("%d: want %q: got %q\n", tc.id, tc.expect, got)
		}
	}
}
