// SPDX-FileCopyrightText: 2025 Antoni Szymański
// SPDX-License-Identifier: MPL-2.0

package hoi4date

import (
	"strconv"
)

const ErrorPrefix = "hoi4date: "

type ErrInvalidDate string

func (e ErrInvalidDate) Error() string {
	return ErrorPrefix + "invalid date: " + string(e)
}

type ErrInvalidBinaryDate int32

func (e ErrInvalidBinaryDate) Error() string {
	var b []byte
	b = append(b, ErrorPrefix+"invalid binary date: "...)
	b = strconv.AppendInt(b, int64(e), 10)
	return string(b)
}
