// SPDX-FileCopyrightText: 2025 Antoni Szymański
// SPDX-License-Identifier: MPL-2.0

package hoi4

import (
	"errors"
	"reflect"
	"strconv"
	"strings"

	"github.com/antoniszymanski/checked-go"
	"github.com/antoniszymanski/hoi4-go/hoi4date"
	"github.com/antoniszymanski/hoi4-go/hoi4text"
)

const ErrorPrefix = "hoi4: "

type SemanticError struct {
	Offset uint64
	GoType reflect.Type
	Err    error
}

func (e SemanticError) Error() string {
	var b []byte
	b = append(b, ErrorPrefix+"unable to unmarshal into Go "...)
	b = append(b, e.GoType.String()...)
	b = append(b, ": "...)
	err := e.Err.Error()
	switch reflect.TypeOf(e.Err).PkgPath() {
	case "github.com/antoniszymanski/hoi4-go":
		err = strings.TrimPrefix(err, ErrorPrefix)
	case "github.com/antoniszymanski/hoi4-go/hoi4text":
		err = strings.TrimPrefix(err, hoi4text.ErrorPrefix)
	case "github.com/antoniszymanski/hoi4-go/hoi4date":
		err = strings.TrimPrefix(err, hoi4date.ErrorPrefix)
	}
	b = append(b, err...)
	b = strconv.AppendUint(append(b, " at offset "...), e.Offset, 10)
	return string(b)
}

func (e SemanticError) Unwrap() error {
	return e.Err
}

type ErrUnexpectedToken struct {
	Token hoi4text.Token
}

func (e ErrUnexpectedToken) Error() string {
	return ErrorPrefix + `unexpected token "` + e.Token.ID().String() + `"`
}

var ErrOutOfRange = errors.New(ErrorPrefix + "out of range")

func cast[U, T checked.Integer](a T) (U, error) {
	b, ok := checked.Cast[U](a)
	if !ok {
		return 0, ErrOutOfRange
	}
	return b, nil
}
