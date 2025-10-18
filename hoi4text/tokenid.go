// SPDX-FileCopyrightText: 2025 Antoni Szymański
// SPDX-License-Identifier: MPL-2.0

package hoi4text

import (
	"encoding"
	"fmt"
	"strconv"
)

type TokenID uint16

const (
	TokenInvalid  TokenID = 0
	TokenOpen     TokenID = 0x0003
	TokenClose    TokenID = 0x0004
	TokenEqual    TokenID = 0x0001
	TokenU32      TokenID = 0x0014
	TokenU64      TokenID = 0x029c
	TokenI32      TokenID = 0x000c
	TokenBool     TokenID = 0x000e
	TokenQuoted   TokenID = 0x000f
	TokenUnquoted TokenID = 0x0017
	TokenF32      TokenID = 0x000d
	TokenF64      TokenID = 0x0167
	TokenI64      TokenID = 0x0317
)

// Identifies if the given ID does not match of the predefined [TokenID]
// constants, and thus can be considered an ID token.
func (id TokenID) IsID() bool {
	switch id {
	case TokenInvalid, TokenOpen, TokenClose, TokenEqual,
		TokenU32, TokenU64, TokenI32, TokenBool, TokenQuoted,
		TokenUnquoted, TokenF32, TokenF64, TokenI64:
		return false
	default:
		return true
	}
}

var (
	_ fmt.Stringer           = TokenID(0)
	_ encoding.TextMarshaler = TokenID(0)
	_ encoding.TextAppender  = TokenID(0)
)

func (id TokenID) String() string {
	b, _ := id.MarshalText()
	return string(b)
}

func (id TokenID) MarshalText() ([]byte, error) {
	return id.AppendText(nil)
}

func (id TokenID) AppendText(b []byte) ([]byte, error) {
	switch id {
	case TokenInvalid:
		b = append(b, "invalid"...)
	case TokenOpen:
		b = append(b, '{')
	case TokenClose:
		b = append(b, '}')
	case TokenEqual:
		b = append(b, '=')
	case TokenU32:
		b = append(b, "u32"...)
	case TokenU64:
		b = append(b, "u64"...)
	case TokenI32:
		b = append(b, "i32"...)
	case TokenBool:
		b = append(b, "bool"...)
	case TokenQuoted:
		b = append(b, "quoted"...)
	case TokenUnquoted:
		b = append(b, "unquoted"...)
	case TokenF32:
		b = append(b, "f32"...)
	case TokenF64:
		b = append(b, "f64"...)
	case TokenI64:
		b = append(b, "i64"...)
	default:
		b = append(b, '<')
		b = strconv.AppendUint(b, uint64(id), 10)
		b = append(b, '>')
	}
	return b, nil
}
