// SPDX-FileCopyrightText: 2026 Antoni Szymański
// SPDX-License-Identifier: MPL-2.0

package hoi4

import "github.com/antoniszymanski/hoi4-go/hoi4text"

type Value []hoi4text.Token

func (v *Value) UnmarshalHOI4(dec *hoi4text.Decoder) error {
	x, err := dec.ReadValue(nil)
	if err != nil {
		return err
	}
	*v = x
	return nil
}
