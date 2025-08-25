// SPDX-FileCopyrightText: 2025 Antoni Szymański
// SPDX-License-Identifier: MPL-2.0

package hoi4text

import (
	"strconv"
	"strings"
)

const ErrorPrefix = "hoi4text: "

type SyntacticError struct {
	Offset uint64
	Err    error
}

func (e SyntacticError) Error() string {
	var b []byte
	err := e.Err.Error()
	if !strings.HasPrefix(err, ErrorPrefix) {
		b = append(b, ErrorPrefix...)
	}
	b = append(b, err...)
	b = strconv.AppendUint(append(b, " after offset "...), e.Offset, 10)
	return string(b)
}

func (e SyntacticError) Unwrap() error {
	return e.Err
}

type Error string

func (e Error) Error() string {
	return ErrorPrefix + string(e)
}

const (
	ErrUnimplemented = Error("unimplemented")
	ErrUnknownHeader = Error("unknown header")
	ErrEndOfObject   = Error("end of object")
	ErrInvalidToken  = Error("invalid token")
)

type ErrUnexpectedToken struct {
	Token Token
}

func (e ErrUnexpectedToken) Error() string {
	return ErrorPrefix + "unexpected token " + e.Token.ID().String()
}
