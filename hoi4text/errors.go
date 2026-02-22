// SPDX-FileCopyrightText: 2026 Antoni Szymański
// SPDX-License-Identifier: MPL-2.0

package hoi4text

import (
	"errors"
	"fmt"
)

var (
	ErrUnknownHeader = errors.New("unknown header")
	ErrUnimplemented = errors.New("unimplemented")
	ErrInvalidToken  = errors.New("invalid token")
	ErrEndOfObject   = errors.New("end of object")
)

type NotAContainerError struct {
	Offset  uint64
	TokenID TokenID
}

func (e *NotAContainerError) Error() string {
	return fmt.Sprintf("expected start of container at offset %d, got %v", e.Offset, e.TokenID)
}
