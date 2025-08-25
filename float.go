// SPDX-FileCopyrightText: 2025 Antoni Szymański
// SPDX-License-Identifier: MPL-2.0

package hoi4

import (
	"reflect"

	"github.com/antoniszymanski/hoi4-go/hoi4text"
)

type Float32 float32

func (x *Float32) UnmarshalHOI4(dec *hoi4text.Decoder) error {
	var t hoi4text.Token
	err := dec.ReadToken(&t)
	if err != nil {
		return err
	}
	switch t.ID() {
	case hoi4text.TokenU32:
		*x = Float32(t.U32())
	case hoi4text.TokenU64:
		*x = Float32(t.U64())
	case hoi4text.TokenI32:
		*x = Float32(t.I32())
	case hoi4text.TokenF32:
		*x = Float32(t.F32())
	case hoi4text.TokenF64:
		*x = Float32(t.F64())
	case hoi4text.TokenI64:
		*x = Float32(t.I64())
	default:
		err = ErrUnexpectedToken{t}
	}
	if err != nil {
		err = SemanticError{dec.Offset(), reflect.TypeOf(*x), err}
	}
	return err
}

type Float64 float64

func (x *Float64) UnmarshalHOI4(dec *hoi4text.Decoder) error {
	var t hoi4text.Token
	err := dec.ReadToken(&t)
	if err != nil {
		return err
	}
	switch t.ID() {
	case hoi4text.TokenU32:
		*x = Float64(t.U32())
	case hoi4text.TokenU64:
		*x = Float64(t.U64())
	case hoi4text.TokenI32:
		*x = Float64(t.I32())
	case hoi4text.TokenF32:
		*x = Float64(t.F32())
	case hoi4text.TokenF64:
		*x = Float64(t.F64())
	case hoi4text.TokenI64:
		*x = Float64(t.I64())
	default:
		err = ErrUnexpectedToken{t}
	}
	if err != nil {
		err = SemanticError{dec.Offset(), reflect.TypeOf(*x), err}
	}
	return err
}
