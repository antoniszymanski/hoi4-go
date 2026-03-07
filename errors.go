// SPDX-FileCopyrightText: 2026 Antoni Szymański
// SPDX-License-Identifier: MPL-2.0

package hoi4

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"unsafe"

	"github.com/antoniszymanski/hoi4-go/hoi4text"
)

var (
	ErrNotANonNilPointer = errors.New("passed value is not a non-nil pointer")
	ErrNilInterface      = errors.New("cannot unmarshal into nil interface")
	ErrEmbeddedInterface = errors.New("cannot unmarshal into embedded interface")
)

type CreateDecoderError struct {
	Err error
}

func (e *CreateDecoderError) Error() string {
	return fmt.Sprintf("failed to create a new decoder: %v", e.Err)
}

func (e *CreateDecoderError) Unwrap() error {
	return e.Err
}

type InvalidTypeError struct {
	Type reflect.Type
}

func (e *InvalidTypeError) Error() string {
	return fmt.Sprintf("cannot unmarshal into Go value of type %v", e.Type)
}

type ReadTokenError struct {
	Offset uint64
	Err    error
}

func (e *ReadTokenError) Error() string {
	return fmt.Sprintf("failed to read token at offset %d: %v", e.Offset, e.Err)
}

func (e *ReadTokenError) Unwrap() error {
	return e.Err
}

type ParseDateError[T int32 | string] struct {
	Input T
}

func (e *ParseDateError[T]) Error() string {
	var dst []byte
	dst = append(dst, "failed to parse "...)
	switch unsafe.Sizeof(e.Input) {
	case unsafe.Sizeof(int32(0)):
		dst = strconv.AppendInt(dst, int64(*(*int32)(unsafe.Pointer(&e.Input))), 10)
	case unsafe.Sizeof(""):
		dst = appendQuote(dst, *(*string)(unsafe.Pointer(&e.Input)))
	}
	dst = append(dst, " as a date"...)
	return string(dst)
}

func appendQuote(dst []byte, s string) []byte {
	if strconv.CanBackquote(s) {
		dst = append(dst, '`')
		dst = append(dst, s...)
		dst = append(dst, '`')
		return dst
	} else {
		return strconv.AppendQuoteToGraphic(dst, s)
	}
}

type InvalidTokenError struct {
	Token hoi4text.Token
	Type  reflect.Type
}

func (e *InvalidTokenError) Error() string {
	return fmt.Sprintf("cannot unmarshal token %v into Go value of type %v", e.Token.ID(), e.Type)
}

type OverflowError[T int64 | uint64 | float64 | int32] struct {
	Value T
	Type  reflect.Type
}

func (e *OverflowError[T]) Error() string {
	return fmt.Sprintf("%v(%v) cannot be represented by Go value of type %v", reflect.TypeFor[T](), e.Value, e.Type)
}

type EnterContainerError struct {
	Err error
}

func (e *EnterContainerError) Error() string {
	return fmt.Sprintf("failed to enter the container: %v", e.Err)
}

func (e *EnterContainerError) Unwrap() error {
	return e.Err
}

type InvalidKeyValueSeparatorError struct {
	TokenID hoi4text.TokenID
}

func (e *InvalidKeyValueSeparatorError) Error() string {
	return fmt.Sprintf("token %v is not a valid key-value separator", e.TokenID)
}

type InvalidObjectKeyError struct {
	Token hoi4text.Token
}

func (e *InvalidObjectKeyError) Error() string {
	return fmt.Sprintf("token %v is not a object key", e.Token.ID())
}

type PeekKindError struct {
	Err error
}

func (e *PeekKindError) Error() string {
	return fmt.Sprintf("failed to peek kind: %v", e.Err)
}

type InvalidRootTypeError struct {
	Type reflect.Type
}

func (e *InvalidRootTypeError) Error() string {
	return fmt.Sprintf("cannot unmarshal root into Go value of type %v", e.Type)
}

type InvalidScalarError struct {
	Token hoi4text.Token
}

func (e *InvalidScalarError) Error() string {
	return fmt.Sprintf("token %v is not a scalar", e.Token.ID())
}
