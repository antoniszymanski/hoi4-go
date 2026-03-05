// SPDX-FileCopyrightText: 2026 Antoni Szymański
// SPDX-License-Identifier: MPL-2.0

package tokenmap

import (
	"bytes"
	"io"
	"os"
	"testing"
)

func BenchmarkEncode(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		if _, err := Encode(io.Discard, decoded); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkDecode(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		if _, err := Decode(bytes.NewReader(encoded)); err != nil {
			b.Fatal(err)
		}
	}
}

func init() {
	var err error
	encoded, err = os.ReadFile("../tokens")
	if err != nil {
		panic(err)
	}
	decoded, err = Decode(bytes.NewReader(encoded))
	if err != nil {
		panic(err)
	}
}

var (
	encoded []byte
	decoded map[uint16]string
)
