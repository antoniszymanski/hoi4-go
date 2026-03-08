// SPDX-FileCopyrightText: 2025 Antoni Szymański
// SPDX-License-Identifier: MPL-2.0

package hoi4text

import "io"

type decoderState struct {
	r     BufferedReader
	depth uint
}

func (d *decoderState) ReadToken() (Token, error) {
	t, err := d.r.ReadToken()
	d.updateDepth(t.ID())
	return t, err
}

func (d *decoderState) SkipToken() (TokenID, error) {
	id, err := d.r.SkipToken()
	d.updateDepth(id)
	return id, err
}

func (d *decoderState) updateDepth(id TokenID) {
	switch id {
	case TokenOpen:
		d.depth++
	case TokenClose:
		d.depth--
	}
}

type Decoder struct {
	s              *decoderState
	minDepth       uint
	endOfContainer bool
}

func NewDecoder(r io.Reader) (*Decoder, error) {
	tr, err := NewReader(r)
	if err != nil {
		return nil, err
	}
	br := BufferedReader{r: tr}
	return &Decoder{s: &decoderState{r: br}}, nil
}

func (d *Decoder) Offset() uint64 {
	return d.s.r.Offset()
}

func (d *Decoder) Depth() uint {
	return d.s.depth
}

func (d *Decoder) ReadToken() (Token, error) {
	if d.minDepth == 0 {
		return d.s.ReadToken()
	} else if d.endOfContainer {
		return Token{}, ErrEndOfContainer
	}
	t, err := d.s.ReadToken()
	if d.Depth() < d.minDepth {
		d.endOfContainer = true
		return Token{}, ErrEndOfContainer
	}
	return t, err
}

func (d *Decoder) SkipToken() (TokenID, error) {
	if d.minDepth == 0 {
		return d.s.SkipToken()
	} else if d.endOfContainer {
		return TokenInvalid, ErrEndOfContainer
	}
	id, err := d.s.SkipToken()
	if d.Depth() < d.minDepth {
		d.endOfContainer = true
		return TokenInvalid, ErrEndOfContainer
	}
	return id, err
}

func (d *Decoder) ReadAll(buf []Token) ([]Token, error) {
	var t Token
	var err error
	for {
		t, err = d.ReadToken()
		if err != nil {
			break
		}
		buf = append(buf, t)
	}
	if err != ErrEndOfContainer && err != io.EOF {
		return nil, err
	}
	return buf, nil
}

func (d *Decoder) SkipAll() error {
	var err error
	for {
		if _, err = d.SkipToken(); err != nil {
			break
		}
	}
	if err != ErrEndOfContainer && err != io.EOF {
		return err
	}
	return nil
}

func (d *Decoder) ReadValue(buf []Token) ([]Token, error) {
	t, err := d.ReadToken()
	if err != nil {
		return nil, err
	}
	buf = append(buf, t)
	if t.ID() != TokenOpen {
		return buf, nil
	}
	d = &Decoder{s: d.s, minDepth: d.Depth()}
	return d.ReadAll(buf)
}

func (d *Decoder) SkipValue() error {
	id, err := d.SkipToken()
	if err != nil {
		return err
	} else if id != TokenOpen {
		return nil
	}
	d = &Decoder{s: d.s, minDepth: d.Depth()}
	return d.SkipAll()
}

func (d *Decoder) EnterContainer() (*Decoder, error) {
	id, err := d.SkipToken()
	if err != nil {
		return nil, err
	} else if id != TokenOpen {
		return nil, &UnexpectedTokenError{id, BeginningOfContainer, d.Offset()}
	}
	return &Decoder{s: d.s, minDepth: d.Depth()}, nil
}

func (d *Decoder) IsEndOfContainer() error {
	if d.endOfContainer {
		return ErrEndOfContainer
	}
	t, err := d.s.r.ReadToken()
	if t.ID() == TokenClose && d.s.depth == d.minDepth {
		d.s.depth--
		d.endOfContainer = true
		return ErrEndOfContainer
	}
	d.s.r.buf = append(d.s.r.buf, bufferedToken{
		token:  t,
		err:    err,
		offset: d.s.r.Offset(),
	})
	return err
}

func (d *Decoder) Peek() Peek {
	return d.s.r.Peek()
}

func (d *Decoder) PeekKind() (Kind, error) {
	return d.s.r.PeekKind()
}
