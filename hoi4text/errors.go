// SPDX-FileCopyrightText: 2026 Antoni Szymański
// SPDX-License-Identifier: MPL-2.0

package hoi4text

import (
	"errors"
	"strconv"
)

var (
	ErrUnknownHeader  = errors.New("unknown header")
	ErrUnimplemented  = errors.New("unimplemented")
	ErrInvalidToken   = errors.New("invalid token")
	ErrEndOfContainer = errors.New("end of container")
)

type UnexpectedTokenError struct {
	TokenID TokenID
	Where   Where
	Offset  uint64
}

func (e *UnexpectedTokenError) Error() string {
	var dst []byte
	dst = append(dst, "unexpected token "...)
	dst = append(dst, e.TokenID.String()...)
	dst = append(dst, " at "...)
	dst = append(dst, e.Where...)
	dst = append(dst, " at offset "...)
	dst = strconv.AppendUint(dst, e.Offset, 10)
	return string(dst)
}

type Where string

const (
	BeginningOfContainer Where = "the beginning of a container"
	BeginningOfValue     Where = "the beginning of a value"
	FirstTokenOfValue    Where = "the first token of a value"
)
