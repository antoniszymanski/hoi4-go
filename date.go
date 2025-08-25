// SPDX-FileCopyrightText: 2025 Antoni Szymański
// SPDX-License-Identifier: MPL-2.0

package hoi4

import (
	"reflect"

	"github.com/antoniszymanski/hoi4-go/hoi4date"
	"github.com/antoniszymanski/hoi4-go/hoi4text"
)

type Date hoi4date.Date

func (x *Date) UnmarshalHOI4(dec *hoi4text.Decoder) error {
	var t hoi4text.Token
	err := dec.ReadToken(&t)
	if err != nil {
		return err
	}
	wrapDate := func(d hoi4date.Date, err error) (Date, error) { return Date(d), err }
	switch t.ID() {
	case hoi4text.TokenU32:
		var i int32
		i, err = cast[int32](t.U32())
		if err == nil {
			*x, err = wrapDate(hoi4date.FromBinary(i))
		}
	case hoi4text.TokenU64:
		var i int32
		i, err = cast[int32](t.U64())
		if err == nil {
			*x, err = wrapDate(hoi4date.FromBinary(i))
		}
	case hoi4text.TokenI32:
		*x, err = wrapDate(hoi4date.FromBinary(t.I32()))
	case hoi4text.TokenF32:
		*x, err = wrapDate(hoi4date.FromBinary(int32(t.F32())))
	case hoi4text.TokenF64:
		*x, err = wrapDate(hoi4date.FromBinary(int32(t.F64())))
	case hoi4text.TokenI64:
		var i int32
		i, err = cast[int32](t.I64())
		if err == nil {
			*x, err = wrapDate(hoi4date.FromBinary(i))
		}
	case hoi4text.TokenQuoted:
		*x, err = wrapDate(hoi4date.FromString(t.Quoted()))
	default:
		err = ErrUnexpectedToken{t}
	}
	if err != nil {
		err = SemanticError{dec.Offset(), reflect.TypeOf(*x), err}
	}
	return err
}
