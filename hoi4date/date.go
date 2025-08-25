// SPDX-FileCopyrightText: 2025 Antoni Szymański
// SPDX-License-Identifier: MPL-2.0

package hoi4date

import "fmt"

type Date struct {
	Year  int16
	Month uint8
	Day   uint8
	Hour  uint8
}

func (d Date) IsValid() bool {
	return d.Year != 0 &&
		d.Month != 0 && d.Month < 13 &&
		d.Day != 0 && d.Day < daysPerMonth[d.Month] &&
		d.Hour != 0 && d.Hour < 25
}

var daysPerMonth = [...]uint8{0, 31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}

type DateFormat uint8

const (
	// Y.M.D[.H] where month, day, and hour don't have zero padding
	DotShort DateFormat = iota + 1

	// Y.M.D[.H] where month, day, and hour are zero padded to two digits
	DotWide

	// ISO-8601 format
	ISO8601
)

func (d Date) String() string {
	return d.Format(DotShort)
}

func (d Date) Format(format DateFormat) string {
	return string(d.AppendFormat(nil, format))
}

func (d Date) AppendFormat(b []byte, format DateFormat) []byte {
	switch format {
	default: // [DotShort]
		return fmt.Appendf(b, "%d.%d.%d.%d", d.Year, d.Month, d.Day, d.Hour)
	case DotWide:
		return fmt.Appendf(b, "%d.%02d.%02d.%02d", d.Year, d.Month, d.Day, d.Hour)
	case ISO8601:
		return fmt.Appendf(b, "%04d-%02d-%02dT%02d", d.Year, d.Month, d.Day, d.Hour-1)
	}
}
