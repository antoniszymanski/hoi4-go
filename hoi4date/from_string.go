// SPDX-FileCopyrightText: 2025 Antoni Szymański
// SPDX-License-Identifier: MPL-2.0

package hoi4date

import "github.com/antoniszymanski/checked-go"

func FromString(x string) (Date, bool) {
	i, x, err := to_i64_t(x)
	if err != nil {
		return none()
	}
	if x == "" {
		year, ok := checked.Cast[int32](i)
		if !ok {
			return none()
		}
		return FromBinary(year)
	}

	year, ok := checked.Cast[int16](i)
	if !ok {
		return none()
	}
	if x == "" || x[0] != '.' {
		return none()
	}

	if len(x) < 1 {
		return none()
	}
	n := x[1]
	var month1 uint8
	if !is_ascii_digit(n) {
		return none()
	} else {
		month1 = n - '0'
	}

	if len(x) < 2 {
		return none()
	}
	n = x[2]
	month, offset := uint8(0), uint(0)
	switch {
	case n == '.':
		month, offset = month1, 2
	case is_ascii_digit(n):
		month, offset = month1*10+(n-'0'), 3
	default:
		return none()
	}

	if uint(len(x)) < offset {
		return none()
	}
	if x[offset] != '.' {
		return none()
	}

	if uint(len(x)) < offset+1 {
		return none()
	}
	n = x[offset+1]
	var day1 uint8
	if !is_ascii_digit(n) {
		return none()
	} else {
		day1 = n - '0'
	}

	if uint(len(x)) < offset+2 {
		return some(year, month, day1, 0)
	}
	var day uint8
	switch x[offset+2] {
	case '.':
		day, offset = day1, offset+2
	default:
		n := x[offset+2]
		if is_ascii_digit(n) {
			result := day1*10 + (n - '0')
			if uint(len(x)) != offset+3 {
				day, offset = result, offset+3
			} else {
				return some(year, month, result, 0)
			}
		} else {
			return none()
		}
	}

	if uint(len(x)) < offset {
		return none()
	}
	if x[offset] != '.' {
		return none()
	}

	if uint(len(x)) < offset+1 {
		return none()
	}
	n = x[offset+1]
	var hour1 uint8
	if !is_ascii_digit(n) || n == '0' {
		return none()
	} else {
		hour1 = n - '0'
	}

	if uint(len(x)) < offset+2 {
		return some(year, month, day, hour1)
	}
	n = x[offset+2]
	if is_ascii_digit(n) {
		result := hour1*10 + (n - '0')
		if uint(len(x)) != offset+3 {
			return none()
		} else {
			return some(year, month, day, result)
		}
	} else {
		return none()
	}
}
