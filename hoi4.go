// SPDX-FileCopyrightText: 2025 Antoni Szymański
// SPDX-License-Identifier: MPL-2.0

package hoi4

import (
	"bytes"
	"io"

	"github.com/antoniszymanski/hoi4-go/hoi4text"
)

type Unmarshaler interface {
	UnmarshalHOI4(dec *hoi4text.Decoder) error
}

func Unmarshal(in []byte, out Unmarshaler) error {
	return UnmarshalRead(bytes.NewReader(in), out)
}

func UnmarshalRead(in io.Reader, out Unmarshaler) error {
	dec, err := hoi4text.NewDecoder(in)
	if err != nil {
		return err
	}
	return UnmarshalDecode(dec, out)
}

func UnmarshalDecode(in *hoi4text.Decoder, out Unmarshaler) error {
	return out.UnmarshalHOI4(in)
}
