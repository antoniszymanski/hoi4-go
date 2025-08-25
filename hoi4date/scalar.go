// SPDX-FileCopyrightText: 2025 Antoni Szymański
// SPDX-License-Identifier: MPL-2.0

//nolint:gosec // G115
package hoi4date

import (
	"errors"
	"math"

	"github.com/antoniszymanski/hoi4-go/internal/checked"
)

var (
	errAllDigits = errors.New("did not contain all digits")
	errOverflow  = errors.New("caused an overflow")
)

func to_i64_t(d string) (int64, string, error) {
	if d == "" {
		return 0, "", errAllDigits
	}
	c, data := d[0], d[1:]
	sign := int64(1)

	var start uint8
	switch {
	case is_ascii_digit(c):
		start = c - '0'
	case c == '-':
		sign = -1
	case c == '+':
	default:
		return 0, "", errAllDigits
	}

	val, rest := to_u64_partial(data, uint64(start))
	OVERFLOW_CUTOFF := digits_in(math.MaxUint64)
	if uint(len(d)) >= OVERFLOW_CUTOFF+1 {
		_, err := check_overflow(d)
		if err != nil {
			return 0, "", errOverflow
		}
	}

	{
		val, ok := checked.Cast[int64](val)
		if !ok {
			return 0, "", errOverflow
		}
		val *= sign
		return val, rest, nil
	}
}

func is_ascii_digit(c byte) bool {
	return '0' <= c && c <= '9'
}

func to_u64_partial(d string, start uint64) (uint64, string) {
	result := start

	for len(d) > 0 {
		c, rest := d[0], d[1:]
		if !is_ascii_digit(c) {
			return result, d
		}

		result *= 10
		result += uint64(c - '0')
		d = rest
	}

	return result, ""
}

func digits_in(n uint64) uint {
	return log10(n) + 1
}

func log10(n uint64) uint {
	switch {
	case n == 0: // special case
		return 0
	case n < 1e1:
		return 1
	case n < 1e2:
		return 2
	case n < 1e3:
		return 3
	case n < 1e4:
		return 4
	case n < 1e5:
		return 5
	case n < 1e6:
		return 6
	case n < 1e7:
		return 7
	case n < 1e8:
		return 8
	case n < 1e9:
		return 9
	case n < 1e10:
		return 10
	case n < 1e11:
		return 11
	case n < 1e12:
		return 12
	case n < 1e13:
		return 13
	case n < 1e14:
		return 14
	case n < 1e15:
		return 15
	case n < 1e16:
		return 16
	case n < 1e17:
		return 17
	case n < 1e18:
		return 18
	case n < 1e19:
		return 19
	default:
		return 20
	}
}

func check_overflow(d string) (uint64, error) {
	if d == "" {
		return 0, errAllDigits
	}

	if d[0] == '+' || d[0] == '-' {
		d = d[1:]
	}

	return check_overflow_init(d, 0)
}

func check_overflow_init(d string, start uint64) (uint64, error) {
	acc := start
	for _, x := range []byte(d) {
		// The input should already be validated by this point, so we just
		// return the accumulator if we find a non-digit.
		if !is_ascii_digit(x) {
			return acc, nil
		}

		var ok bool
		acc, ok = checked.Mul(acc, 10)
		if !ok {
			return 0, errOverflow
		}
		acc, ok = checked.Add(acc, uint64(x-'0'))
		if !ok {
			return 0, errOverflow
		}
	}

	return acc, nil
}
