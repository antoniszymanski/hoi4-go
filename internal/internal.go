// SPDX-FileCopyrightText: 2025 Antoni Szymański
// SPDX-License-Identifier: MPL-2.0

package internal

import "unsafe"

func BytesToString(b []byte) string {
	return unsafe.String(unsafe.SliceData(b), len(b))
}

func StringToBytes(s string) []byte {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}

//nolint:all
//go:nosplit
func NoEscape[P ~*E, E any](p P) P {
	x := uintptr(unsafe.Pointer(p))
	return P(unsafe.Pointer(x ^ 0))
}
