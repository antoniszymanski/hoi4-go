// SPDX-FileCopyrightText: 2025 Antoni Szymański
// SPDX-License-Identifier: MPL-2.0

//nolint:gosec // G115
package hoi4date

import "github.com/antoniszymanski/checked-go"

func FromBinary(x int32) (Date, bool) {
	hour := x % 24
	x /= 24

	days_since_jan1 := x % 365
	if hour < 0 || days_since_jan1 < 0 {
		return none()
	}

	x /= 365
	x, ok := checked.Sub(x, 5000)
	if !ok {
		return none()
	}
	year, ok := checked.Cast[int16](x)
	if !ok {
		return none()
	}
	month, day := month_day_from_julian(days_since_jan1)
	return some(year, month, day, uint8(hour)+1)
}

func month_day_from_julian(days_since_jan1 int32) (month, day uint8) {
	switch {
	case days_since_jan1 >= 0 && days_since_jan1 <= 30:
		return 1, uint8(days_since_jan1 + 1)
	case days_since_jan1 >= 31 && days_since_jan1 <= 58:
		return 2, uint8(days_since_jan1 - 30)
	case days_since_jan1 >= 59 && days_since_jan1 <= 89:
		return 3, uint8(days_since_jan1 - 58)
	case days_since_jan1 >= 90 && days_since_jan1 <= 119:
		return 4, uint8(days_since_jan1 - 89)
	case days_since_jan1 >= 120 && days_since_jan1 <= 150:
		return 5, uint8(days_since_jan1 - 119)
	case days_since_jan1 >= 151 && days_since_jan1 <= 180:
		return 6, uint8(days_since_jan1 - 150)
	case days_since_jan1 >= 181 && days_since_jan1 <= 211:
		return 7, uint8(days_since_jan1 - 180)
	case days_since_jan1 >= 212 && days_since_jan1 <= 242:
		return 8, uint8(days_since_jan1 - 211)
	case days_since_jan1 >= 243 && days_since_jan1 <= 272:
		return 9, uint8(days_since_jan1 - 242)
	case days_since_jan1 >= 273 && days_since_jan1 <= 303:
		return 10, uint8(days_since_jan1 - 272)
	case days_since_jan1 >= 304 && days_since_jan1 <= 333:
		return 11, uint8(days_since_jan1 - 303)
	case days_since_jan1 >= 334 && days_since_jan1 <= 364:
		return 12, uint8(days_since_jan1 - 333)
	default:
		panic("unreachable")
	}
}
