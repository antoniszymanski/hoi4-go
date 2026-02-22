// SPDX-FileCopyrightText: 2025 Antoni Szymański
// SPDX-FileCopyrightText: Nick Babcock
// SPDX-License-Identifier: MPL-2.0

package hoi4date

import "github.com/antoniszymanski/checked-go"

func Parse(data string) (Date, bool) {
	i, data, ok := toInt64(data)
	if !ok {
		return none()
	}
	if data == "" {
		binary, ok := checked.Cast[int32](i)
		if !ok {
			return none()
		}
		return ParseBinary(binary)
	}

	year, ok := checked.Cast[int16](i)
	if !ok || data[0] != '.' || len(data) <= 1 {
		return none()
	}

	n := data[1]
	if !isAsciiDigit(n) {
		return none()
	}
	month1 := n - '0'
	if len(data) <= 2 {
		return none()
	}

	n = data[2]
	month, offset := uint8(0), uint(0)
	switch {
	case n == '.':
		month, offset = month1, 2
	case isAsciiDigit(n):
		month, offset = month1*10+(n-'0'), 3
	default:
		return none()
	}
	if uint(len(data)) <= offset || data[offset] != '.' || uint(len(data)) <= offset+1 {
		return none()
	}

	n = data[offset+1]
	if !isAsciiDigit(n) {
		return none()
	}
	day1 := n - '0'
	if uint(len(data)) <= offset+2 {
		return some(year, month, day1, 0)
	}

	var day uint8
	switch data[offset+2] {
	case '.':
		day, offset = day1, offset+2
	default:
		n := data[offset+2]
		if !isAsciiDigit(n) {
			return none()
		}
		day = day1*10 + (n - '0')
		if uint(len(data)) == offset+3 {
			return some(year, month, day, 0)
		}
		offset += 3
	}
	if uint(len(data)) <= offset || data[offset] != '.' || uint(len(data)) <= offset+1 {
		return none()
	}

	n = data[offset+1]
	if !isAsciiDigit(n) || n == '0' {
		return none()
	}
	hour1 := n - '0'
	if uint(len(data)) <= offset+2 {
		return some(year, month, day, hour1)
	}

	n = data[offset+2]
	if !isAsciiDigit(n) || uint(len(data)) != offset+3 {
		return none()
	}
	hour := hour1*10 + (n - '0')
	return some(year, month, day, hour)
}

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
