// SPDX-FileCopyrightText: 2025 Antoni Szymański
// SPDX-License-Identifier: MPL-2.0

package hoi4

import (
	"reflect"

	"github.com/antoniszymanski/hoi4-go/hoi4text"
)

type Bool bool

func (x *Bool) UnmarshalHOI4(dec *hoi4text.Decoder) error {
	var t hoi4text.Token
	err := dec.ReadToken(&t)
	if err != nil {
		return err
	}
	switch t.ID() {
	case hoi4text.TokenBool:
		*x = Bool(t.Bool())
	default:
		err = ErrUnexpectedToken{t}
	}
	if err != nil {
		err = SemanticError{dec.Offset(), reflect.TypeOf(*x), err}
	}
	return err
}
