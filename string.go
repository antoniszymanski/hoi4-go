// SPDX-FileCopyrightText: 2025 Antoni Szymański
// SPDX-License-Identifier: MPL-2.0

package hoi4

import (
	"reflect"

	"github.com/antoniszymanski/hoi4-go/hoi4text"
)

type String string

func (x *String) UnmarshalHOI4(dec *hoi4text.Decoder) error {
	var t hoi4text.Token
	err := dec.ReadToken(&t)
	if err != nil {
		return err
	}
	switch t.ID() {
	case hoi4text.TokenQuoted:
		*x = String(t.Quoted())
	case hoi4text.TokenUnquoted:
		*x = String(t.Unquoted())
	default:
		if t.ID().IsID() {
			s, ok := hoi4text.Tokens.Lookup(t.ID())
			if ok {
				*x = String(s)
			} else {
				*x = String(t.String())
			}
		} else {
			err = ErrUnexpectedToken{t}
		}
	}
	if err != nil {
		err = SemanticError{dec.Offset(), reflect.TypeOf(*x), err}
	}
	return err
}
