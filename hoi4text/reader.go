// SPDX-FileCopyrightText: 2025 Antoni Szymański
// SPDX-License-Identifier: MPL-2.0

package hoi4text

import (
	"io"

	"github.com/antoniszymanski/hoi4-go/internal"
)

type Reader interface {
	ReadToken(*Token) error
	Offset() uint64
}

type Skipper interface {
	SkipToken() (TokenID, error)
}

const (
	HeaderTxt = "HOI4txt"
	HeaderBin = "HOI4bin"
	HeaderLen = len(HeaderTxt)
)

var _ [len(HeaderTxt)]int = [len(HeaderBin)]int{}

func NewReader(r io.Reader) (Reader, error) {
	buf, err := read(r, HeaderLen, nil)
	if err == io.ErrUnexpectedEOF {
		return nil, ErrUnknownHeader
	} else if err != nil {
		return nil, err
	}
	switch string(buf) { // does not allocate
	case HeaderBin:
		return &BinaryReader{r: r, buf: buf}, nil
	case HeaderTxt:
		return nil, ErrUnimplemented
	default:
		return nil, ErrUnknownHeader
	}
}

func SkipToken(r Reader) (TokenID, error) {
	if s, ok := r.(Skipper); ok {
		return s.SkipToken()
	}
	t := new(Token)
	if err := r.ReadToken(internal.NoEscape(t)); err != nil {
		return TokenInvalid, err
	}
	return t.ID(), nil
}

func read(r io.Reader, length int, buf []byte) ([]byte, error) {
	buf = resize(buf, length)
	if _, err := io.ReadFull(r, buf); err != nil {
		return nil, err
	}
	return buf, nil
}
