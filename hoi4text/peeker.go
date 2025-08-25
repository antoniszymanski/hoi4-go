// SPDX-FileCopyrightText: 2025 Antoni Szymański
// SPDX-License-Identifier: MPL-2.0

package hoi4text

import "slices"

type Peeker struct {
	d      *decoderState
	noCopy bool
	buf    []bufferedToken
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

func (p *Peeker) ReadToken(t *Token) error {
	if err := p.d.ReadToken(t); err != nil {
		return err
	}
	if !p.noCopy {
		var copy Token
		t.Copy(&copy)
		p.buf = append(p.buf, bufferedToken{
			token:  copy,
			offset: p.Offset(),
		})
	} else {
		p.buf = append(p.buf, bufferedToken{
			token:  *t,
			offset: p.Offset(),
		})
	}
	return nil
}

func (p *Peeker) SkipToken() (TokenID, error) {
	var t Token
	if err := p.d.ReadToken(&t); err != nil {
		return TokenInvalid, err
	}
	p.buf = append(p.buf, bufferedToken{
		token:  t,
		offset: p.Offset(),
	})
	return t.ID(), nil
}
