// Copyright (c) 2023 Michael D Henderson. All rights reserved.

package lisp

// _nil is the NIL symbol.
// This should be immutable, so don't change it!
var _nil = Atom{_type: AtomType_Nil}

// sym_table is a global symbol table.
// it is a list of all existing symbols.
var sym_table = _nil
