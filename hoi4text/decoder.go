// SPDX-FileCopyrightText: 2025 Antoni Szymański
// SPDX-License-Identifier: MPL-2.0

package hoi4text

import "io"

type decoderState struct {
	r      Reader
	buf    []bufferedToken
	offset uint64
	depth  uint
	mode   uint8
}

func (d *decoderState) ReadToken(t *Token) error {
	if len(d.buf) > 0 {
		bt := d.buf[len(d.buf)-1]
		*t = bt.token
		d.offset = bt.offset
		d.buf = d.buf[:len(d.buf)-1]
		return nil
	}
	var err error
	switch d.mode {
	case 0:
		t.putID(TokenOpen)
		d.mode = 1
	case 1:
		err = d.r.ReadToken(t)
		if err == io.EOF {
			t.putID(TokenClose)
			err = nil
			d.mode = 2
		}
	case 2:
		t.reset()
		err = io.EOF
	}
	switch t.ID() {
	case TokenOpen:
		d.depth++
	case TokenClose:
		d.depth--
	}
	d.offset = d.r.Offset()
	return err
}

func (d *decoderState) SkipToken() (TokenID, error) {
	if len(d.buf) > 0 {
		bt := d.buf[len(d.buf)-1]
		d.offset = bt.offset
		d.buf = d.buf[:len(d.buf)-1]
		return bt.token.ID(), nil
	}
	var id TokenID
	var err error
	switch d.mode {
	case 0:
		id = TokenOpen
		d.mode = 1
	case 1:
		id, err = SkipToken(d.r)
		if err == io.EOF {
			id = TokenClose
			err = nil
			d.mode = 2
		}
	case 2:
		err = io.EOF
	}
	switch id {
	case TokenOpen:
		d.depth++
	case TokenClose:
		d.depth--
	}
	d.offset = d.r.Offset()
	return id, err
}

type Decoder struct {
	s           *decoderState
	minDepth    uint
	endOfObject bool
}

func NewDecoder(r io.Reader) (*Decoder, error) {
	tr, err := NewReader(r)
	if err != nil {
		return nil, err
	}
	return &Decoder{s: &decoderState{r: tr}}, nil
}

func (d *Decoder) Offset() uint64 {
	return d.s.offset
}

func (d *Decoder) Depth() uint {
	return d.s.depth
}

func (d *Decoder) ReadToken(t *Token) error {
	if d.minDepth == 0 {
		return d.s.ReadToken(t)
	} else if d.endOfObject {
		t.reset()
		return ErrEndOfObject
	}
	err := d.s.ReadToken(t)
	if d.Depth() < d.minDepth {
		d.endOfObject = true
		t.reset()
		err = ErrEndOfObject
	}
	return err
}

func (d *Decoder) SkipToken() (TokenID, error) {
	if d.minDepth == 0 {
		return d.s.SkipToken()
	} else if d.endOfObject {
		return TokenInvalid, ErrEndOfObject
	}
	id, err := d.s.SkipToken()
	if d.Depth() < d.minDepth {
		d.endOfObject = true
		id = TokenInvalid
		err = ErrEndOfObject
	}
	return id, err
}

func (d *Decoder) ReadAll(buf []Token) ([]Token, error) {
	buf = buf[:0]
	var t Token
	var err error
	for {
		err = d.ReadToken(&t)
		if err != nil {
			break
		}
		buf = append(buf, t)
	}
	if err != ErrEndOfObject && err != io.EOF {
		return nil, err
	}
	return buf, nil
}

func (d *Decoder) SkipAll() error {
	var t Token
	var err error
	for {
		err = d.ReadToken(&t)
		if err != nil {
			break
		}
	}
	if err != ErrEndOfObject && err != io.EOF {
		return err
	}
	return nil
}

func (d *Decoder) Value() (*Decoder, error) {
	var t Token
	if err := d.ReadToken(&t); err != nil {
		return nil, err
	}
	if t.ID() != TokenOpen {
		return nil, &SyntacticError{
			Offset: d.Offset(),
			Err:    &ErrUnexpectedToken{t},
		}
	}
	return &Decoder{s: d.s, minDepth: d.Depth()}, nil
}

func (d *Decoder) SkipValue() error {
	var t Token
	err := d.ReadToken(&t)
	if err != nil {
		return err
	}
	if t.ID() != TokenOpen {
		return nil
	}
	val := &Decoder{s: d.s, minDepth: d.Depth()}
	for {
		err = val.ReadToken(&t)
		if err != nil {
			break
		}
	}
	if err != ErrEndOfObject && err != io.EOF {
		return err
	}
	return nil
}

func (d *Decoder) Peek(noCopy bool) *Peeker {
	return &Peeker{d: d.s, noCopy: noCopy}
}

func (d *Decoder) PeekKind() (Kind, error) {
	p := d.Peek(false)
	defer p.Close()

	var t Token
	if err := p.ReadToken(&t); err != nil {
		return 0, err
	}
	if t.ID() != TokenOpen {
		return KindValue, nil
	}

	if _, err := p.SkipToken(); err != nil {
		return 0, err
	}

	if err := p.ReadToken(&t); err != nil {
		return 0, err
	}
	if t.ID() != TokenEqual {
		return KindArray, nil
	} else {
		return KindMap, nil
	}
}

type Kind uint8

const (
	KindValue Kind = iota + 1
	KindArray
	KindMap
)

func (k Kind) IsObject() bool {
	return k == KindArray || k == KindMap
}

func (k Kind) String() string {
	switch k {
	case KindValue:
		return "value"
	case KindArray:
		return "array"
	case KindMap:
		return "map"
	default:
		return "invalid"
	}
}
