// SPDX-FileCopyrightText: 2025 Antoni Szymański
// SPDX-License-Identifier: MPL-2.0

package hoi4text

import (
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

func (id TokenID) String() string {
	switch id {
	case TokenInvalid:
		return "invalid"
	case TokenOpen:
		return "{"
	case TokenClose:
		return "}"
	case TokenEqual:
		return "="
	case TokenU32:
		return "u32"
	case TokenU64:
		return "u64"
	case TokenI32:
		return "i32"
	case TokenBool:
		return "bool"
	case TokenQuoted:
		return "quoted"
	case TokenUnquoted:
		return "unquoted"
	case TokenF32:
		return "f32"
	case TokenF64:
		return "f64"
	case TokenI64:
		return "i64"
	default:
		var b []byte
		b = append(b, '<')
		b = strconv.AppendUint(b, uint64(id), 10)
		b = append(b, '>')
		return string(b)
	}
}
