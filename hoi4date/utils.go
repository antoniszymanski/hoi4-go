// SPDX-FileCopyrightText: 2026 Antoni Szymański
// SPDX-License-Identifier: MPL-2.0

package hoi4date

func some(year int16, month, day, hour uint8) (Date, bool) {
	d := Date{year, month, day, hour}
	if !d.IsValid() {
		return none()
	}
	return d, true
}

func none() (Date, bool) {
	return Date{}, false
}
