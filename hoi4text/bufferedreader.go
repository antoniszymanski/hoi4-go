// SPDX-FileCopyrightText: 2026 Antoni Szymański
// SPDX-License-Identifier: MPL-2.0

package hoi4text

import (
	"io"
	"slices"
)

type BufferedReader struct {
	r       Reader
	offset  uint64
	buf     []bufferedToken
	peekBuf []bufferedToken
}

type bufferedToken struct {
	token  Token
	err    error
	offset uint64
}

func NewBufferedReader(r io.Reader) (*BufferedReader, error) {
	tr, err := NewReader(r)
	if err != nil {
		return nil, err
	}
	return &BufferedReader{r: tr}, nil
}

func (br *BufferedReader) Offset() uint64 {
	return br.offset
}

func (br *BufferedReader) ReadToken() (Token, error) {
	if len(br.buf) > 0 {
		last := br.buf[len(br.buf)-1]
		br.buf = br.buf[:len(br.buf)-1]
		br.offset = last.offset
		return last.token, last.err
	}
	t, err := br.r.ReadToken()
	br.offset = br.r.Offset()
	return t, err
}

func (br *BufferedReader) SkipToken() (TokenID, error) {
	if len(br.buf) > 0 {
		last := br.buf[len(br.buf)-1]
		br.buf = br.buf[:len(br.buf)-1]
		br.offset = last.offset
		return last.token.ID(), last.err
	}
	id, err := SkipToken(br.r)
	br.offset = br.r.Offset()
	return id, err
}

func (br *BufferedReader) PeekKind() (Kind, error) {
	if br.offset == 0 {
		return KindRoot, nil
	}
	p := br.Peek()
	defer p.Close()
	if id, err := p.SkipToken(); err != nil {
		return KindInvalid, err
	} else if id == TokenEqual || id == TokenClose {
		return KindInvalid, &UnexpectedTokenError{id, BeginningOfValue, br.offset}
	} else if id != TokenOpen {
		return KindScalar, nil
	}
	if id, err := p.SkipToken(); err != nil {
		return KindInvalid, err
	} else if id == TokenEqual {
		return KindInvalid, &UnexpectedTokenError{id, FirstTokenOfValue, br.offset}
	} else if id == TokenClose {
		return KindEmptyContainer, nil
	}
	if id, err := p.SkipToken(); err != nil {
		return KindInvalid, err
	} else if id != TokenEqual {
		return KindArray, nil
	} else {
		return KindObject, nil
	}
}

type Kind uint8

const (
	KindInvalid Kind = iota
	KindRoot
	KindScalar
	KindEmptyContainer
	KindArray
	KindObject
)

func (k Kind) String() string {
	switch k {
	case KindRoot:
		return "root"
	case KindScalar:
		return "scalar"
	case KindEmptyContainer:
		return "empty container"
	case KindArray:
		return "array"
	case KindObject:
		return "object"
	default:
		return "invalid"
	}
}

func (br *BufferedReader) Peek() *Peek {
	return &Peek{br}
}

type Peek struct {
	br *BufferedReader
}

func (p *Peek) Offset() uint64 {
	return p.br.offset
}

func (p *Peek) ReadToken() (Token, error) {
	t, err := p.br.ReadToken()
	p.br.peekBuf = append(p.br.peekBuf, bufferedToken{
		token:  t,
		err:    err,
		offset: p.br.Offset(),
	})
	return t, err
}

func (p *Peek) SkipToken() (TokenID, error) {
	t, err := p.ReadToken()
	return t.ID(), err
}

func (p *Peek) Close() {
	slices.Reverse(p.br.peekBuf)
	p.br.buf = append(p.br.buf, p.br.peekBuf...)
	p.br.peekBuf = p.br.peekBuf[:0]
}
