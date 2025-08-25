// SPDX-FileCopyrightText: 2025 Antoni Szymański
// SPDX-License-Identifier: MPL-2.0

package hoi4date

import "github.com/antoniszymanski/hoi4-go/internal/checked"

func FromString(in string) (Date, error) {
	None := func() (Date, error) { return Date{}, ErrInvalidDate(in) }
	Some := func(d Date) (Date, error) {
		if !d.IsValid() {
			return None()
		}
		return d, nil
	}
	data := in

	i, data, err := to_i64_t(data)
	if err != nil {
		return None()
	}
	if data == "" {
		year, ok := checked.Cast[int32](i)
		if !ok {
			return None()
		}
		return FromBinary(year)
	}

	year, ok := checked.Cast[int16](i)
	if !ok {
		return None()
	}
	if data == "" || data[0] != '.' {
		return None()
	}

	if len(data) < 1 {
		return None()
	}
	n := data[1]
	var month1 uint8
	if !is_ascii_digit(n) {
		return None()
	} else {
		month1 = n - '0'
	}

	if len(data) < 2 {
		return None()
	}
	n = data[2]
	month, offset := uint8(0), uint(0)
	switch {
	case n == '.':
		month, offset = month1, 2
	case is_ascii_digit(n):
		month, offset = month1*10+(n-'0'), 3
	default:
		return None()
	}

	if uint(len(data)) < offset {
		return None()
	}
	if data[offset] != '.' {
		return None()
	}

	if uint(len(data)) < offset+1 {
		return None()
	}
	n = data[offset+1]
	var day1 uint8
	if !is_ascii_digit(n) {
		return None()
	} else {
		day1 = n - '0'
	}

	if uint(len(data)) < offset+2 {
		return Some(Date{year, month, day1, 0})
	}
	var day uint8
	switch data[offset+2] {
	case '.':
		day, offset = day1, offset+2
	default:
		n := data[offset+2]
		if is_ascii_digit(n) {
			result := day1*10 + (n - '0')
			if uint(len(data)) != offset+3 {
				day, offset = result, offset+3
			} else {
				return Some(Date{year, month, result, 0})
			}
		} else {
			return None()
		}
	}

	if uint(len(data)) < offset {
		return None()
	}
	if data[offset] != '.' {
		return None()
	}

	if uint(len(data)) < offset+1 {
		return None()
	}
	n = data[offset+1]
	var hour1 uint8
	if !is_ascii_digit(n) || n == '0' {
		return None()
	} else {
		hour1 = n - '0'
	}

	if uint(len(data)) < offset+2 {
		return Some(Date{year, month, day, hour1})
	}
	n = data[offset+2]
	if is_ascii_digit(n) {
		result := hour1*10 + (n - '0')
		if uint(len(data)) != offset+3 {
			return None()
		} else {
			return Some(Date{year, month, day, result})
		}
	} else {
		return None()
	}
}
