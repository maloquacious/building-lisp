// Copyright (c) 2023 Michael D Henderson. All rights reserved.

package ch04

// this file defines all the global variables used by the implementation.
// because the implementation uses globals rather than passing state,
// we can only have one instance running per program.

// _nil is the NIL symbol.
// This should be immutable, so don't change it!
var _nil = Atom{_type: AtomType_Nil}

// sym_table is a global symbol table.
// it is a list of all existing symbols.
var sym_table = Atom{_type: AtomType_Nil}
