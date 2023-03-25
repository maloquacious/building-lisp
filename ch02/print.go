// Copyright (c) 2023 Michael D Henderson. All rights reserved.

package ch02

import (
	"fmt"
	"io"
)

// print_expr writes the value of an Atom to the writer.
func print_expr(w io.Writer, atom Atom) (int, error) {
	switch atom._type {
	case AtomType_Nil:
		// write "NIL"
		return w.Write([]byte{'N', 'I', 'L'})
	case AtomType_Integer:
		return w.Write([]byte(fmt.Sprintf("%d", atom.value.integer)))
	case AtomType_Pair:
		// write the open paren
		totalBytesWritten, err := w.Write([]byte{'('})
		if err != nil {
			return totalBytesWritten, err
		}

		// print the car of the atom
		bytesWritten, err := print_expr(w, car(atom))
		totalBytesWritten += bytesWritten
		if err != nil {
			return totalBytesWritten, err
		}

		// write the remainder of the list, starting with the cdr of the atom
		atom = cdr(atom)
		for !nilp(atom) {
			if atom._type == AtomType_Pair {
				// write " " to separate the atoms in the list
				bytesWritten, err = w.Write([]byte{' '})
				totalBytesWritten += bytesWritten
				if err != nil {
					return totalBytesWritten, err
				}

				// print the car of the current atom
				bytesWritten, err = print_expr(w, car(atom))
				totalBytesWritten += bytesWritten
				if err != nil {
					return totalBytesWritten, err
				}

				// advance to the next atom in the list
				atom = cdr(atom)

				continue
			}

			// found an "improper list," one that ends with a dotted pair

			// write " . " to separate the dotted pair
			bytesWritten, err = w.Write([]byte{' ', '.', ' '})
			totalBytesWritten += bytesWritten
			if err != nil {
				return totalBytesWritten, err
			}

			// print the atom
			bytesWritten, err = print_expr(w, atom)
			totalBytesWritten += bytesWritten
			if err != nil {
				return totalBytesWritten, err
			}

			// dotted pair ends a list, so quit the loop now
			break
		}

		// write the closing paren
		bytesWritten, err = w.Write([]byte{')'})
		totalBytesWritten += bytesWritten

		// and return
		return totalBytesWritten, err

	case AtomType_Symbol:
		return w.Write(atom.value.symbol)
	}

	panic(fmt.Sprintf("assert(_type != %d)", atom._type))
}
