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
		if out.IsNil() {
			return ErrNilInterface
		}
		return unmarshalRoot(dec, out.Elem())
	case reflect.Map:
		return unmarshalRootMap(dec, out)
	case reflect.Pointer:
		if out.IsNil() {
			out.Set(reflect.New(out.Type().Elem()))
		}
		return unmarshalRoot(dec, out.Elem())
	case reflect.Struct:
		return unmarshalRootStruct(dec, out)
	default:
		return &InvalidRootTypeError{out.Type()}
	}
}

func unmarshalRootMap(dec *hoi4text.Decoder, out reflect.Value) error {
	if err := unmarshalMapContent(dec, out); err != io.EOF {
		return err
	}
	return nil
}

func unmarshalRootStruct(dec *hoi4text.Decoder, out reflect.Value) error {
	if err := unmarshalStructContent(dec, out); err != io.EOF {
		return err
	}
	return nil
}
