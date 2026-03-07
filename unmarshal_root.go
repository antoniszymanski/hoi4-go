// SPDX-FileCopyrightText: 2026 Antoni Szymański
// SPDX-License-Identifier: MPL-2.0

package hoi4

import (
	"io"
	"reflect"

	"github.com/antoniszymanski/hoi4-go/hoi4text"
)

func unmarshalRoot(dec *hoi4text.Decoder, out reflect.Value) error {
	if out.Type() == reflect.TypeFor[any]() {
		return unmarshalAnyRoot(dec, out)
	}
	if u, ok := reflect.TypeAssert[Unmarshaler](out); ok {
		return u.UnmarshalHOI4(dec)
	}
	switch out.Kind() {
	case reflect.Interface:
		return unmarshalRootInterface(dec, out)
	case reflect.Map:
		return unmarshalRootMap(dec, out)
	case reflect.Pointer:
		return unmarshalRootPointer(dec, out)
	case reflect.Struct:
		return unmarshalRootStruct(dec, out)
	default:
		return &InvalidRootTypeError{out.Type()}
	}
}

func unmarshalRootInterface(dec *hoi4text.Decoder, out reflect.Value) error {
	if out.IsNil() {
		return ErrNilInterface
	}
	return unmarshalRoot(dec, out.Elem())
}

func unmarshalRootMap(dec *hoi4text.Decoder, out reflect.Value) error {
	return unmarshalMapContent(dec, out, io.EOF)
}

func unmarshalRootPointer(dec *hoi4text.Decoder, out reflect.Value) error {
	if out.IsNil() {
		out.Set(reflect.New(out.Type().Elem()))
	}
	return unmarshalRoot(dec, out.Elem())
}

func unmarshalRootStruct(dec *hoi4text.Decoder, out reflect.Value) error {
	return unmarshalStructContent(dec, out, io.EOF)
}
