// SPDX-FileCopyrightText: 2025 Antoni Szymański
// SPDX-License-Identifier: MPL-2.0

package checked

import (
	"unsafe"

	"golang.org/x/exp/constraints"
)

func Add[T constraints.Integer](a, b T) (T, bool) {
	c := a + b
	if (c > a) == (b > 0) {
		return c, true
	}
	return c, false
}

func Sub[T constraints.Integer](a, b T) (T, bool) {
	c := a - b
	if (c < a) == (b > 0) {
		return c, true
	}
	return c, false
}

func Mul[T constraints.Integer](a, b T) (T, bool) {
	if a == 0 || b == 0 {
		return 0, true
	}
	c := a * b
	if (c < 0) == ((a < 0) != (b < 0)) {
		if c/b == a {
			return c, true
		}
	}
	return c, false
}

func Div[T constraints.Integer](a, b T) (T, bool) {
	q, _, ok := Quotient(a, b)
	return q, ok
}

func Quotient[T constraints.Integer](a, b T) (T, T, bool) {
	if b == 0 {
		return 0, 0, false
	}
	c := a / b
	ok := (c < 0) == ((a < 0) != (b < 0))
	return c, a % b, ok
}

func Cast[U constraints.Integer, T constraints.Integer](a T) (U, bool) {
	b := U(a)
	if T(b) == a && (a < 0) == (b < 0) {
		return b, true
	}
	return b, false
}

func MinInt[T constraints.Signed]() (x T) {
	size := unsafe.Sizeof(x) * 8
	return -1 << (size - 1)
}

func MaxInt[T constraints.Signed]() (x T) {
	size := unsafe.Sizeof(x) * 8
	return 1<<(size-1) - 1
}

func MaxUint[T constraints.Unsigned]() (x T) {
	size := unsafe.Sizeof(x) * 8
	return 1<<size - 1
}
