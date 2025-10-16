// SPDX-FileCopyrightText: 2025 Antoni Szymański
// SPDX-License-Identifier: MPL-2.0

package hoi4

import (
	"reflect"

	"github.com/antoniszymanski/hoi4-go/hoi4text"
)

type Uint uint

func (x *Uint) UnmarshalHOI4(dec *hoi4text.Decoder) error {
	t, err := dec.ReadToken()
	if err != nil {
		return err
	}
	switch t.ID() {
	case hoi4text.TokenU32:
		*x = Uint(t.U32())
	case hoi4text.TokenU64:
		*x = Uint(t.U64())
	case hoi4text.TokenI32:
		*x, err = cast[Uint](t.I32())
	case hoi4text.TokenF32:
		*x = Uint(t.F32())
	case hoi4text.TokenF64:
		*x = Uint(t.F64())
	case hoi4text.TokenI64:
		*x, err = cast[Uint](t.I64())
	default:
		err = ErrUnexpectedToken{t}
	}
	if err != nil {
		err = SemanticError{dec.Offset(), reflect.TypeOf(*x), err}
	}
	return err
}

type Uint8 uint8

func (x *Uint8) UnmarshalHOI4(dec *hoi4text.Decoder) error {
	t, err := dec.ReadToken()
	if err != nil {
		return err
	}
	switch t.ID() {
	case hoi4text.TokenU32:
		*x, err = cast[Uint8](t.U32())
	case hoi4text.TokenU64:
		*x, err = cast[Uint8](t.U64())
	case hoi4text.TokenI32:
		*x, err = cast[Uint8](t.I32())
	case hoi4text.TokenF32:
		*x = Uint8(t.F32())
	case hoi4text.TokenF64:
		*x = Uint8(t.F64())
	case hoi4text.TokenI64:
		*x, err = cast[Uint8](t.I64())
	default:
		err = ErrUnexpectedToken{t}
	}
	if err != nil {
		err = SemanticError{dec.Offset(), reflect.TypeOf(*x), err}
	}
	return err
}

type Uint16 uint16

func (x *Uint16) UnmarshalHOI4(dec *hoi4text.Decoder) error {
	t, err := dec.ReadToken()
	if err != nil {
		return err
	}
	switch t.ID() {
	case hoi4text.TokenU32:
		*x, err = cast[Uint16](t.U32())
	case hoi4text.TokenU64:
		*x, err = cast[Uint16](t.U64())
	case hoi4text.TokenI32:
		*x, err = cast[Uint16](t.I32())
	case hoi4text.TokenF32:
		*x = Uint16(t.F32())
	case hoi4text.TokenF64:
		*x = Uint16(t.F64())
	case hoi4text.TokenI64:
		*x, err = cast[Uint16](t.I64())
	default:
		err = ErrUnexpectedToken{t}
	}
	if err != nil {
		err = SemanticError{dec.Offset(), reflect.TypeOf(*x), err}
	}
	return err
}

type Uint32 uint32

func (x *Uint32) UnmarshalHOI4(dec *hoi4text.Decoder) error {
	t, err := dec.ReadToken()
	if err != nil {
		return err
	}
	switch t.ID() {
	case hoi4text.TokenU32:
		*x = Uint32(t.U32())
	case hoi4text.TokenU64:
		*x, err = cast[Uint32](t.U64())
	case hoi4text.TokenI32:
		*x, err = cast[Uint32](t.I32())
	case hoi4text.TokenF32:
		*x = Uint32(t.F32())
	case hoi4text.TokenF64:
		*x = Uint32(t.F64())
	case hoi4text.TokenI64:
		*x, err = cast[Uint32](t.I64())
	default:
		err = ErrUnexpectedToken{t}
	}
	if err != nil {
		err = SemanticError{dec.Offset(), reflect.TypeOf(*x), err}
	}
	return err
}

type Uint64 uint64

func (x *Uint64) UnmarshalHOI4(dec *hoi4text.Decoder) error {
	t, err := dec.ReadToken()
	if err != nil {
		return err
	}
	switch t.ID() {
	case hoi4text.TokenU32:
		*x = Uint64(t.U32())
	case hoi4text.TokenU64:
		*x = Uint64(t.U64())
	case hoi4text.TokenI32:
		*x, err = cast[Uint64](t.I32())
	case hoi4text.TokenF32:
		*x = Uint64(t.F32())
	case hoi4text.TokenF64:
		*x = Uint64(t.F64())
	case hoi4text.TokenI64:
		*x, err = cast[Uint64](t.I64())
	default:
		err = ErrUnexpectedToken{t}
	}
	if err != nil {
		err = SemanticError{dec.Offset(), reflect.TypeOf(*x), err}
	}
	return err
}
