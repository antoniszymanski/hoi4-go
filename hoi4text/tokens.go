// SPDX-FileCopyrightText: 2025 Antoni Szymański
// SPDX-License-Identifier: MPL-2.0

package hoi4text

import (
	_ "embed"
	"encoding/binary"
)

type tokens struct {
	m map[TokenID]string
}

func (t tokens) Get(id TokenID) string {
	return t.m[id]
}

func (t tokens) Lookup(id TokenID) (string, bool) {
	s, ok := t.m[id]
	return s, ok
}

var Tokens tokens

func init() {
	Tokens.m = parseTokens(tokensData)
}

//go:embed tokens
var tokensData string

func parseTokens(data string) map[TokenID]string {
	length, data := readUint16(data)
	m := make(map[TokenID]string, length)
	var id uint16
	var token string
	for len(data) > 0 {
		id, data = readUint16(data)
		token, data = readString(data)
		m[TokenID(id)] = token
	}
	return m
}

func readUint16(data string) (uint16, string) {
	v := binary.LittleEndian.Uint16([]byte(data))
	return v, data[2:]
}

func readString(data string) (string, string) {
	length, data := data[0], data[1:]
	return data[:length], data[length:]
}
