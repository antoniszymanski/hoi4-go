// SPDX-FileCopyrightText: 2025 Antoni Szymański
// SPDX-FileCopyrightText: Nick Babcock
// SPDX-License-Identifier: MPL-2.0

//nolint:gosec // G115
package hoi4date

import "github.com/antoniszymanski/checked-go"

func ParseBinary(data int32) (Date, bool) {
	hour := data % 24
	data /= 24

	daysSinceJan1 := data % 365
	if hour < 0 || daysSinceJan1 < 0 {
		return none()
	}

	data /= 365
	data, ok := checked.Sub(data, 5000)
	if !ok {
		return none()
	}
	year, ok := checked.Cast[int16](data)
	if !ok {
		return none()
	}
	month, day := monthDayFromJulian(daysSinceJan1)
	return some(year, month, day, uint8(hour)+1)
}

func monthDayFromJulian(daysSinceJan1 int32) (month, day uint8) {
	switch {
	case daysSinceJan1 >= 0 && daysSinceJan1 <= 30:
		return 1, uint8(daysSinceJan1 + 1)
	case daysSinceJan1 >= 31 && daysSinceJan1 <= 58:
		return 2, uint8(daysSinceJan1 - 30)
	case daysSinceJan1 >= 59 && daysSinceJan1 <= 89:
		return 3, uint8(daysSinceJan1 - 58)
	case daysSinceJan1 >= 90 && daysSinceJan1 <= 119:
		return 4, uint8(daysSinceJan1 - 89)
	case daysSinceJan1 >= 120 && daysSinceJan1 <= 150:
		return 5, uint8(daysSinceJan1 - 119)
	case daysSinceJan1 >= 151 && daysSinceJan1 <= 180:
		return 6, uint8(daysSinceJan1 - 150)
	case daysSinceJan1 >= 181 && daysSinceJan1 <= 211:
		return 7, uint8(daysSinceJan1 - 180)
	case daysSinceJan1 >= 212 && daysSinceJan1 <= 242:
		return 8, uint8(daysSinceJan1 - 211)
	case daysSinceJan1 >= 243 && daysSinceJan1 <= 272:
		return 9, uint8(daysSinceJan1 - 242)
	case daysSinceJan1 >= 273 && daysSinceJan1 <= 303:
		return 10, uint8(daysSinceJan1 - 272)
	case daysSinceJan1 >= 304 && daysSinceJan1 <= 333:
		return 11, uint8(daysSinceJan1 - 303)
	case daysSinceJan1 >= 334 && daysSinceJan1 <= 364:
		return 12, uint8(daysSinceJan1 - 333)
	default:
		panic("unreachable")
	}
}
