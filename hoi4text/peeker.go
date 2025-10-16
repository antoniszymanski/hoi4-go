// SPDX-FileCopyrightText: 2025 Antoni Szymański
// SPDX-License-Identifier: MPL-2.0

package hoi4text

import "slices"

type Peeker struct {
	d   *decoderState
	buf []bufferedToken
}

type bufferedToken struct {
	token  Token
	offset uint64
}

func (p *Peeker) Close() {
	slices.Reverse(p.buf)
	p.d.buf = append(p.d.buf, p.buf...)
}

func (p *Peeker) Offset() uint64 {
	return p.d.r.Offset()
}

func (p *Peeker) ReadToken() (Token, error) {
	t, err := p.d.ReadToken()
	if err != nil {
		return t, err
	}
	p.buf = append(p.buf, bufferedToken{
		token:  t,
		offset: p.Offset(),
	})
	return t, nil
}

func (p *Peeker) SkipToken() (TokenID, error) {
	t, err := p.d.ReadToken()
	if err != nil {
		return TokenInvalid, err
	}
	p.buf = append(p.buf, bufferedToken{
		token:  t,
		offset: p.Offset(),
	})
	return t.ID(), nil
}
