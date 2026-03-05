// SPDX-FileCopyrightText: 2026 Antoni Szymański
// SPDX-License-Identifier: MPL-2.0

package hoi4

import (
	"bytes"
	"io"
	"reflect"

	"github.com/antoniszymanski/hoi4-go/hoi4text"
)

type Unmarshaler interface {
	UnmarshalHOI4(dec *hoi4text.Decoder) error
}

func Unmarshal(in []byte, out any) error {
	return UnmarshalRead(bytes.NewReader(in), out)
}

func UnmarshalRead(in io.Reader, out any) error {
	v, err := validateValue(out)
	if err != nil {
		return err
	}
	dec, err := hoi4text.NewDecoder(in)
	if err != nil {
		return &CreateDecoderError{err}
	}
	return unmarshalRoot(dec, v)
}

func UnmarshalDecode(in *hoi4text.Decoder, out any) error {
	v, err := validateValue(out)
	if err != nil {
		return err
	}
	return unmarshalRoot(in, v)
}

func validateValue(i any) (reflect.Value, error) {
	v := reflect.ValueOf(i)
	if v.Kind() != reflect.Pointer || v.IsNil() {
		return reflect.Value{}, ErrNotANonNilPointer
	}
	return v, nil
}
