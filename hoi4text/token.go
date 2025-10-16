// SPDX-FileCopyrightText: 2025 Antoni Szymański
// SPDX-License-Identifier: MPL-2.0

package hoi4text

import (
	"encoding"
	"encoding/binary"
	"fmt"
	"math"
	"strconv"
	"unsafe"
)

type Token struct {
	_    nonComparable
	id   TokenID
	data [8]byte
	ptr  *byte
}

type nonComparable [0]func()

func (t Token) Equal(other Token) bool {
	if t.id != other.id {
		return false
	}
	if t.id != TokenQuoted && t.id != TokenUnquoted {
		return true
	}
	lengthA := binary.NativeEndian.Uint64(t.data[:])
	lengthB := binary.NativeEndian.Uint64(other.data[:])
	return unsafe.String(t.ptr, lengthA) == unsafe.String(other.ptr, lengthB)
}

func (t *Token) reset() {
	t.putID(TokenInvalid)
}

func (t *Token) ID() TokenID {
	return t.id
}

func (t *Token) putID(v TokenID) {
	t.id = v
	t.ptr = nil
}

func (t *Token) U32() uint32 {
	if t.id != TokenU32 {
		panic("TokenID is not TokenU32")
	}
	return binary.NativeEndian.Uint32(t.data[:])
}

func (t *Token) putU32(v uint32) {
	t.putID(TokenU32)
	binary.NativeEndian.PutUint32(t.data[:], v)
}

func (t *Token) U64() uint64 {
	if t.id != TokenU64 {
		panic("TokenID is not TokenU64")
	}
	return binary.NativeEndian.Uint64(t.data[:])
}

func (t *Token) putU64(v uint64) {
	t.putID(TokenU64)
	binary.NativeEndian.PutUint64(t.data[:], v)
}

func (t *Token) I32() int32 {
	if t.id != TokenI32 {
		panic("TokenID is not TokenI32")
	}
	return int32(binary.NativeEndian.Uint32(t.data[:])) //#nosec G115
}

func (t *Token) putI32(v int32) {
	t.putID(TokenI32)
	binary.NativeEndian.PutUint32(t.data[:], uint32(v)) //#nosec G115
}

func (t *Token) Bool() bool {
	if t.id != TokenBool {
		panic("TokenID is not TokenBool")
	}
	return t.data[0] != 0
}

func (t *Token) putBool(v bool) {
	t.putID(TokenBool)
	if v {
		t.data[0] = 1
	} else {
		t.data[0] = 0
	}
}

func (t *Token) Quoted() string {
	if t.id != TokenQuoted {
		panic("TokenID is not TokenQuoted")
	}
	length := binary.NativeEndian.Uint64(t.data[:])
	return unsafe.String(t.ptr, length)
}

func (t *Token) putQuoted(v string) {
	t.putID(TokenQuoted)
	binary.NativeEndian.PutUint64(t.data[:], uint64(len(v)))
	t.ptr = unsafe.StringData(v)
}

func (t *Token) Unquoted() string {
	if t.id != TokenUnquoted {
		panic("TokenID is not TokenUnquoted")
	}
	length := binary.NativeEndian.Uint64(t.data[:])
	return unsafe.String(t.ptr, length)
}

func (t *Token) putUnquoted(v string) {
	t.putID(TokenUnquoted)
	binary.NativeEndian.PutUint64(t.data[:], uint64(len(v)))
	t.ptr = unsafe.StringData(v)
}

func (t *Token) F32() float32 {
	if t.id != TokenF32 {
		panic("TokenID is not TokenF32")
	}
	return math.Float32frombits(binary.NativeEndian.Uint32(t.data[:]))
}

func (t *Token) putF32(v float32) {
	t.putID(TokenF32)
	binary.NativeEndian.PutUint32(t.data[:], math.Float32bits(v))
}

func (t *Token) F64() float64 {
	if t.id != TokenF64 {
		panic("TokenID is not TokenF64")
	}
	return math.Float64frombits(binary.NativeEndian.Uint64(t.data[:]))
}

func (t *Token) putF64(v float64) {
	t.putID(TokenF64)
	binary.NativeEndian.PutUint64(t.data[:], math.Float64bits(v))
}

func (t *Token) I64() int64 {
	if t.id != TokenI64 {
		panic("TokenID is not TokenI64")
	}
	return int64(binary.NativeEndian.Uint64(t.data[:])) //#nosec G115
}

func (t *Token) putI64(v int64) {
	t.putID(TokenI64)
	binary.NativeEndian.PutUint64(t.data[:], uint64(v)) //#nosec G115
}

var (
	_ fmt.Stringer           = Token{}
	_ encoding.TextMarshaler = Token{}
	_ encoding.TextAppender  = Token{}
)

func (t Token) String() string {
	b, err := t.MarshalText()
	if err != nil {
		return "<invalid token>"
	}
	return string(b)
}

func (t Token) MarshalText() ([]byte, error) {
	return t.AppendText(nil)
}

func (t Token) AppendText(b []byte) ([]byte, error) {
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
		b = strconv.AppendUint(b, uint64(t.U32()), 10)
	case TokenU64:
		b = strconv.AppendUint(b, uint64(t.U64()), 10)
	case TokenI32:
		b = strconv.AppendInt(b, int64(t.I32()), 10)
	case TokenBool:
		b = strconv.AppendBool(b, t.Bool())
	case TokenQuoted:
		b = strconv.AppendQuote(b, t.Quoted())
	case TokenUnquoted:
		b = append(b, t.Unquoted()...)
	case TokenF32:
		b = strconv.AppendFloat(b, float64(t.F32()), 'g', -1, 32)
	case TokenF64:
		b = strconv.AppendFloat(b, t.F64(), 'g', -1, 64)
	case TokenI64:
		b = strconv.AppendInt(b, t.I64(), 10)
	default:
		b = append(b, '<')
		b = strconv.AppendUint(b, uint64(t.id), 10)
		b = append(b, '>')
	}
	return b, err
}
