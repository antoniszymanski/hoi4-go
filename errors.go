// SPDX-FileCopyrightText: 2026 Antoni Szymański
// SPDX-License-Identifier: MPL-2.0

package hoi4

import (
	"errors"
	"fmt"
	"reflect"

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

type UnmarshalDateError[T int32 | string] struct {
	Input T
}

func (e *UnmarshalDateError[T]) Error() string {
	return fmt.Sprintf("failed to unmarshal date: %q is not valid", e.Input)
}

type InvalidTokenError struct {
	Token hoi4text.Token
	Type  reflect.Type
}

func (e *InvalidTokenError) Error() string {
	return fmt.Sprintf("cannot unmarshal %v token into Go value of type %v", e.Token.ID(), e.Type)
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
