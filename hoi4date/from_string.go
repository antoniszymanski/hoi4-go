// SPDX-FileCopyrightText: 2025 Antoni Szymański
// SPDX-FileCopyrightText: Nick Babcock
// SPDX-License-Identifier: MPL-2.0

package hoi4date

import "github.com/antoniszymanski/checked-go"

func FromString(data string) (Date, bool) {
	i, data, ok := toInt64(data)
	if !ok {
		return none()
	}
	if data == "" {
		binary, ok := checked.Cast[int32](i)
		if !ok {
			return none()
		}
		return FromBinary(binary)
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
