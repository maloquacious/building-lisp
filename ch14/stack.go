// Copyright (c) 2023 Michael D Henderson. All rights reserved.

package lisp

// create constants for reaching into the stack frame
const (
	FRAME_PARENT = 0
	FRAME_ENV    = 1
	FRAME_OP     = 2
	FRAME_TAIL   = 3
	FRAME_ARGS   = 4
	FRAME_BODY   = 5
)

// make_frame returns a frame.
// the standard layout of a frame makes it easy to use
// list_get to fetch values and list_set to update them.
func make_frame(parent, env, tail Atom) Atom {
	op, args, body := _nil, _nil, _nil
	return cons(parent, // depth == 0
		cons(env, // depth == 1
			cons(op, // depth == 2
				cons(tail, // depth == 3
					cons(args, // depth == 4
						cons(body, // depth == 5
							_nil))))))
}
