// SPDX-FileCopyrightText: 2025 Antoni Szymański
// SPDX-License-Identifier: MPL-2.0

package hoi4

import (
	"reflect"

	"github.com/antoniszymanski/hoi4-go/hoi4text"
)

type Int int

func (x *Int) UnmarshalHOI4(dec *hoi4text.Decoder) error {
	var t hoi4text.Token
	err := dec.ReadToken(&t)
	if err != nil {
		return err
	}
	switch t.ID() {
	case hoi4text.TokenU32:
		*x, err = cast[Int](t.U32())
	case hoi4text.TokenU64:
		*x, err = cast[Int](t.U64())
	case hoi4text.TokenI32:
		*x = Int(t.I32())
	case hoi4text.TokenF32:
		*x = Int(t.F32())
	case hoi4text.TokenF64:
		*x = Int(t.F64())
	case hoi4text.TokenI64:
		*x, err = cast[Int](t.I64())
	default:
		err = ErrUnexpectedToken{t}
	}
	if err != nil {
		err = SemanticError{dec.Offset(), reflect.TypeOf(*x), err}
	}
	return err
}

type Int8 int8

func (x *Int8) UnmarshalHOI4(dec *hoi4text.Decoder) error {
	var t hoi4text.Token
	err := dec.ReadToken(&t)
	if err != nil {
		return err
	}
	switch t.ID() {
	case hoi4text.TokenU32:
		*x, err = cast[Int8](t.U32())
	case hoi4text.TokenU64:
		*x, err = cast[Int8](t.U64())
	case hoi4text.TokenI32:
		*x, err = cast[Int8](t.I32())
	case hoi4text.TokenF32:
		*x = Int8(t.F32())
	case hoi4text.TokenF64:
		*x = Int8(t.F64())
	case hoi4text.TokenI64:
		*x, err = cast[Int8](t.I64())
	default:
		err = ErrUnexpectedToken{t}
	}
	if err != nil {
		err = SemanticError{dec.Offset(), reflect.TypeOf(*x), err}
	}
	return err
}

type Int16 int16

func (x *Int16) UnmarshalHOI4(dec *hoi4text.Decoder) error {
	var t hoi4text.Token
	err := dec.ReadToken(&t)
	if err != nil {
		return err
	}
	switch t.ID() {
	case hoi4text.TokenU32:
		*x, err = cast[Int16](t.U32())
	case hoi4text.TokenU64:
		*x, err = cast[Int16](t.U64())
	case hoi4text.TokenI32:
		*x, err = cast[Int16](t.I32())
	case hoi4text.TokenF32:
		*x = Int16(t.F32())
	case hoi4text.TokenF64:
		*x = Int16(t.F64())
	case hoi4text.TokenI64:
		*x, err = cast[Int16](t.I64())
	default:
		err = ErrUnexpectedToken{t}
	}
	if err != nil {
		err = SemanticError{dec.Offset(), reflect.TypeOf(*x), err}
	}
	return err
}

type Int32 int32

func (x *Int32) UnmarshalHOI4(dec *hoi4text.Decoder) error {
	var t hoi4text.Token
	err := dec.ReadToken(&t)
	if err != nil {
		return err
	}
	switch t.ID() {
	case hoi4text.TokenU32:
		*x, err = cast[Int32](t.U32())
	case hoi4text.TokenU64:
		*x, err = cast[Int32](t.U64())
	case hoi4text.TokenI32:
		*x = Int32(t.I32())
	case hoi4text.TokenF32:
		*x = Int32(t.F32())
	case hoi4text.TokenF64:
		*x = Int32(t.F64())
	case hoi4text.TokenI64:
		*x, err = cast[Int32](t.I64())
	default:
		err = ErrUnexpectedToken{t}
	}
	if err != nil {
		err = SemanticError{dec.Offset(), reflect.TypeOf(*x), err}
	}
	return err
}

type Int64 int64

func (x *Int64) UnmarshalHOI4(dec *hoi4text.Decoder) error {
	var t hoi4text.Token
	err := dec.ReadToken(&t)
	if err != nil {
		return err
	}
	switch t.ID() {
	case hoi4text.TokenU32:
		*x = Int64(t.U32())
	case hoi4text.TokenU64:
		*x, err = cast[Int64](t.U64())
	case hoi4text.TokenI32:
		*x = Int64(t.I32())
	case hoi4text.TokenF32:
		*x = Int64(t.F32())
	case hoi4text.TokenF64:
		*x = Int64(t.F64())
	case hoi4text.TokenI64:
		*x = Int64(t.I64())
	default:
		err = ErrUnexpectedToken{t}
	}
	if err != nil {
		err = SemanticError{dec.Offset(), reflect.TypeOf(*x), err}
	}
	return err
}
