// SPDX-FileCopyrightText: 2025 Antoni Szymański
// SPDX-License-Identifier: MPL-2.0

package hoi4text

import (
	"encoding/binary"
	"math"
	"strconv"

	"github.com/antoniszymanski/hoi4-go/internal"
)

type Token struct {
	id   TokenID
	data []byte
}

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
		b = appendUint(b, id)
		b = append(b, '>')
		return internal.BytesToString(b)
	}
}

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

func (t *Token) reset() {
	t.id = TokenInvalid
	t.data = t.data[:0]
}

func (t *Token) Copy(other *Token) {
	other.id = t.id
	if t.data == nil {
		other.data = nil
	} else {
		other.data = resize(other.data, len(t.data))
		copy(other.data, t.data)
	}
}

func (t *Token) ID() TokenID {
	return t.id
}

func (t *Token) putID(v TokenID) {
	t.id = v
	t.data = t.data[:0]
}

func (t *Token) U32() uint32 {
	if t.id != TokenU32 {
		panic("TokenID is not TokenU32")
	}
	return binary.LittleEndian.Uint32(t.data)
}

func (t *Token) putU32(v uint32) {
	t.id = TokenU32
	t.data = binary.LittleEndian.AppendUint32(t.data[:0], v)
}

func (t *Token) U64() uint64 {
	if t.id != TokenU64 {
		panic("TokenID is not TokenU64")
	}
	return binary.LittleEndian.Uint64(t.data)
}

func (t *Token) putU64(v uint64) {
	t.id = TokenU64
	t.data = binary.LittleEndian.AppendUint64(t.data[:0], v)
}

func (t *Token) I32() int32 {
	if t.id != TokenI32 {
		panic("TokenID is not TokenI32")
	}
	return int32(binary.LittleEndian.Uint32(t.data)) //#nosec G115
}

func (t *Token) putI32(v int32) {
	t.id = TokenI32
	t.data = binary.LittleEndian.AppendUint32(t.data[:0], uint32(v)) //#nosec G115
}

func (t *Token) Bool() bool {
	if t.id != TokenBool {
		panic("TokenID is not TokenBool")
	}
	return t.data[0] != 0
}

func (t *Token) putBool(v bool) {
	t.id = TokenBool
	if v {
		t.data = append(t.data[:0], 1)
	} else {
		t.data = append(t.data[:0], 0)
	}
}

func (t *Token) Quoted() string {
	if t.id != TokenQuoted {
		panic("TokenID is not TokenQuoted")
	}
	return internal.BytesToString(t.data)
}

func (t *Token) putQuoted(v string) {
	t.id = TokenQuoted
	t.data = append(t.data[:0], v...)
}

func (t *Token) Unquoted() string {
	if t.id != TokenUnquoted {
		panic("TokenID is not TokenUnquoted")
	}
	return internal.BytesToString(t.data)
}

func (t *Token) putUnquoted(v string) {
	t.id = TokenUnquoted
	t.data = append(t.data[:0], v...)
}

func (t *Token) F32() float32 {
	if t.id != TokenF32 {
		panic("TokenID is not TokenF32")
	}
	return math.Float32frombits(binary.LittleEndian.Uint32(t.data))
}

func (t *Token) putF32(v float32) {
	t.id = TokenF32
	t.data = binary.LittleEndian.AppendUint32(t.data[:0], math.Float32bits(v))
}

func (t *Token) F64() float64 {
	if t.id != TokenF64 {
		panic("TokenID is not TokenF64")
	}
	return math.Float64frombits(binary.LittleEndian.Uint64(t.data))
}

func (t *Token) putF64(v float64) {
	t.id = TokenF64
	t.data = binary.LittleEndian.AppendUint64(t.data[:0], math.Float64bits(v))
}

func (t *Token) I64() int64 {
	if t.id != TokenI64 {
		panic("TokenID is not TokenI64")
	}
	return int64(binary.LittleEndian.Uint64(t.data)) //#nosec G115
}

func (t *Token) putI64(v int64) {
	t.id = TokenI64
	t.data = binary.LittleEndian.AppendUint64(t.data[:0], uint64(v)) //#nosec G115
}

func (t Token) String() string {
	b, err := t.TextMarshaler()
	if err != nil {
		return "<invalid token>"
	}
	return internal.BytesToString(b)
}

func (t *Token) TextMarshaler() ([]byte, error) {
	return t.TextAppender(nil)
}

func (t *Token) TextAppender(b []byte) ([]byte, error) {
	var err error
	switch t.id {
	case TokenInvalid:
		err = ErrInvalidToken
	case TokenOpen:
		b = append(b, '{')
	case TokenClose:
		b = append(b, '}')
	case TokenEqual:
		b = append(b, '=')
	case TokenU32:
		b = appendUint(b, t.U32())
	case TokenU64:
		b = appendUint(b, t.U64())
	case TokenI32:
		b = appendInt(b, t.I32())
	case TokenBool:
		b = strconv.AppendBool(b, t.Bool())
	case TokenQuoted:
		b = strconv.AppendQuote(b, t.Quoted())
	case TokenUnquoted:
		b = append(b, t.Unquoted()...)
	case TokenF32:
		b = appendFloat32(b, t.F32())
	case TokenF64:
		b = appendFloat64(b, t.F64())
	case TokenI64:
		b = appendInt(b, t.I64())
	default:
		b = append(b, '<')
		b = appendUint(b, t.id)
		b = append(b, '>')
	}
	return b, err
}
