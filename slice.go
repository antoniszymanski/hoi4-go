// SPDX-FileCopyrightText: 2025 Antoni Szymański
// SPDX-License-Identifier: MPL-2.0

package hoi4

import "github.com/antoniszymanski/hoi4-go/hoi4text"

type Slice[T any] []T

func (x *Slice[T]) UnmarshalHOI4(dec *hoi4text.Decoder) error {
	dec, err := dec.Value()
	if err != nil {
		return err
	}

	for {
		var elem T
		u := any(&elem).(Unmarshaler)
		if err = u.UnmarshalHOI4(dec); err != nil {
			break
		}
		*x = append(*x, elem)
	}
	if err == hoi4text.ErrEndOfObject {
		err = nil
	}
	return err
}
