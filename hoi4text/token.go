// SPDX-FileCopyrightText: 2025 Antoni Szymański
// SPDX-License-Identifier: MPL-2.0

package hoi4text

import (
	"cmp"
	"encoding"
	"encoding/binary"
	"fmt"
	"math"
	"strconv"
	"strings"
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
	switch t.id {
	case TokenBool:
		return (t.data[0] != 0) == (other.data[0] != 0)
	case TokenU32, TokenI32, TokenF32:
		return [4]byte(t.data[:4]) == [4]byte(other.data[:4])
	case TokenU64, TokenI64, TokenF64:
		return t.data == other.data
	case TokenQuoted, TokenUnquoted:
		t_len := binary.NativeEndian.Uint64(t.data[:])
		other_len := binary.NativeEndian.Uint64(other.data[:])
		return unsafe.String(t.ptr, t_len) == unsafe.String(other.ptr, other_len)
	default:
		return true
	}
}

func (t Token) Compare(other Token) int {
	if t.id != other.id {
		return t.id.Compare(other.id)
	}
	switch t.id {
	case TokenBool:
		return compareBool(t.getBool(), other.getBool())
	case TokenU32:
		return cmp.Compare(t.getU32(), other.getU32())
	case TokenI32:
		return cmp.Compare(t.getI32(), other.getI32())
	case TokenF32:
		return cmp.Compare(t.getF32(), other.getF32())
	case TokenU64:
		return cmp.Compare(t.getU64(), other.getU64())
	case TokenI64:
		return cmp.Compare(t.getI64(), other.getI64())
	case TokenF64:
		return cmp.Compare(t.getF64(), other.getF64())
	case TokenQuoted, TokenUnquoted:
		return strings.Compare(t.getString(), other.getString())
	default:
		return 0
	}
}

func compareBool(x, y bool) int {
	switch {
	case x == y:
		return 0
	case !x && y: // false < true
		return -1
	default: // true > false
		return +1
	}
}

// #region Getters

func (t Token) ID() TokenID {
	return t.id
}

func (t Token) U32() uint32 {
	if t.id != TokenU32 {
		panic("TokenID is not TokenU32")
	}
	return t.getU32()
}

func (t Token) getU32() uint32 {
	return binary.NativeEndian.Uint32(t.data[:])
}

func (t Token) U64() uint64 {
	if t.id != TokenU64 {
		panic("TokenID is not TokenU64")
	}
	return t.getU64()
}

func (t Token) getU64() uint64 {
	return binary.NativeEndian.Uint64(t.data[:])
}

func (t Token) I32() int32 {
	if t.id != TokenI32 {
		panic("TokenID is not TokenI32")
	}
	return t.getI32()
}

func (t Token) getI32() int32 {
	return int32(t.getU32()) //#nosec G115
}

func (t Token) Bool() bool {
	if t.id != TokenBool {
		panic("TokenID is not TokenBool")
	}
	return t.getBool()
}

func (t Token) getBool() bool {
	return t.data[0] != 0
}

func (t Token) Quoted() string {
	if t.id != TokenQuoted {
		panic("TokenID is not TokenQuoted")
	}
	return t.getString()
}

func (t Token) Unquoted() string {
	if t.id != TokenUnquoted {
		panic("TokenID is not TokenUnquoted")
	}
	return t.getString()
}

func (t Token) getString() string {
	length := binary.NativeEndian.Uint64(t.data[:])
	return unsafe.String(t.ptr, length)
}

func (t Token) F32() float32 {
	if t.id != TokenF32 {
		panic("TokenID is not TokenF32")
	}
	return t.getF32()
}

func (t Token) getF32() float32 {
	return math.Float32frombits(binary.NativeEndian.Uint32(t.data[:]))
}

func (t Token) F64() float64 {
	if t.id != TokenF64 {
		panic("TokenID is not TokenF64")
	}
	return t.getF64()
}

func (t Token) getF64() float64 {
	return math.Float64frombits(binary.NativeEndian.Uint64(t.data[:]))
}

func (t Token) I64() int64 {
	if t.id != TokenI64 {
		panic("TokenID is not TokenI64")
	}
	return t.getI64()
}

func (t Token) getI64() int64 {
	return int64(t.getU64()) //#nosec G115
}

// #endregion

// #region Constructors

func ID(id TokenID) Token {
	return Token{id: id}
}

func U32(i uint32) Token {
	t := Token{id: TokenU32}
	binary.NativeEndian.PutUint32(t.data[:], i)
	return t
}

func U64(i uint64) Token {
	t := Token{id: TokenU64}
	binary.NativeEndian.PutUint64(t.data[:], i)
	return t
}

func I32(i int32) Token {
	t := Token{id: TokenI32}
	binary.NativeEndian.PutUint32(t.data[:], uint32(i)) //#nosec G115
	return t
}

func Bool(b bool) Token {
	t := Token{id: TokenBool}
	if b {
		t.data[0] = 1
	}
	return t
}

func Quoted(s string) Token {
	t := Token{id: TokenQuoted, ptr: unsafe.StringData(s)}
	binary.NativeEndian.PutUint64(t.data[:], uint64(len(s)))
	return t
}

func Unquoted(s string) Token {
	t := Token{id: TokenUnquoted, ptr: unsafe.StringData(s)}
	binary.NativeEndian.PutUint64(t.data[:], uint64(len(s)))
	return t
}

func F32(f float32) Token {
	t := Token{id: TokenF32}
	binary.NativeEndian.PutUint32(t.data[:], math.Float32bits(f))
	return t
}

func F64(f float64) Token {
	t := Token{id: TokenF64}
	binary.NativeEndian.PutUint64(t.data[:], math.Float64bits(f))
	return t
}

func I64(i int64) Token {
	t := Token{id: TokenI64}
	binary.NativeEndian.PutUint64(t.data[:], uint64(i)) //#nosec G115
	return t
}

// #endregion

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
