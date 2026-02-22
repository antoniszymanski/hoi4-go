// SPDX-FileCopyrightText: 2025 Antoni Szymański
// SPDX-FileCopyrightText: Nick Babcock
// SPDX-License-Identifier: MPL-2.0

package hoi4date

import "github.com/antoniszymanski/checked-go"

func toInt64(s string) (int64, string, bool) {
	if s == "" {
		return 0, "", false
	}
	sign := int64(1)
	var start uint8
	switch b := s[0]; {
	case isAsciiDigit(b):
		start = b - '0'
	case b == '-':
		sign = -1
	case b == '+':
	default:
		return 0, "", false
	}
	const overflowCutoff = 20 // len(strconv.FormatUint(math.MaxUint64, 10))
	if len(s) > overflowCutoff && !checkOverflow(s) {
		return 0, "", false
	}
	u, rest := toUint64Partial(s[1:], uint64(start))
	i, ok := checked.Cast[int64](u)
	if !ok {
		return 0, "", false
	}
	i *= sign
	return i, rest, true
}

func checkOverflow(s string) bool {
	if s == "" {
		return false
	} else if s[0] == '+' || s[0] == '-' {
		s = s[1:]
	}
	acc := uint64(0)
	for _, b := range []byte(s) {
		// The input should already be validated by this point, so we just
		// return the accumulator if we find a non-digit.
		if !isAsciiDigit(b) {
			return true
		}
		var ok bool
		acc, ok = checked.Mul(acc, 10)
		if !ok {
			return false
		}
		acc, ok = checked.Add(acc, uint64(b-'0'))
		if !ok {
			return false
		}
	}
	return true
}

func toUint64Partial(s string, start uint64) (uint64, string) {
	result := start
	for len(s) > 0 {
		if !isAsciiDigit(s[0]) {
			return result, s
		}
		result *= 10
		result += uint64(s[0] - '0')
		s = s[1:]
	}
	return result, ""
}

func isAsciiDigit(b byte) bool {
	return '0' <= b && b <= '9'
}
