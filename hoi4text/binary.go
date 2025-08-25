// SPDX-FileCopyrightText: 2025 Antoni Szymański
// SPDX-License-Identifier: MPL-2.0

package hoi4text

import (
	"encoding/binary"
	"io"
	"math"

	"github.com/antoniszymanski/hoi4-go/internal"
)

type BinaryReader struct {
	r      io.Reader
	buf    []byte
	offset uint64
}

func (r *BinaryReader) Offset() uint64 {
	return r.offset
}

func (r *BinaryReader) ReadToken(t *Token) error {
	t.reset()
	id, err := r.readID()
	if err != nil {
		return err
	}

	switch id {
	case TokenOpen, TokenClose, TokenEqual:
		t.putID(id)
	case TokenU32:
		v, err := r.readU32()
		if err != nil {
			return err
		}
		t.putU32(v)
	case TokenU64:
		v, err := r.readU64()
		if err != nil {
			return err
		}
		t.putU64(v)
	case TokenI32:
		v, err := r.readI32()
		if err != nil {
			return err
		}
		t.putI32(v)
	case TokenBool:
		v, err := r.readBool()
		if err != nil {
			return err
		}
		t.putBool(v)
	case TokenQuoted:
		v, err := r.readString()
		if err != nil {
			return err
		}
		t.putQuoted(v)
	case TokenUnquoted:
		v, err := r.readString()
		if err != nil {
			return err
		}
		t.putUnquoted(v)
	case TokenF32:
		v, err := r.readF32()
		if err != nil {
			return err
		}
		t.putF32(v)
	case TokenF64:
		v, err := r.readF64()
		if err != nil {
			return err
		}
		t.putF64(v)
	case TokenI64:
		v, err := r.readI64()
		if err != nil {
			return err
		}
		t.putI64(v)
	default:
		t.putID(id)
	}
	return nil
}

func (r *BinaryReader) SkipToken() (TokenID, error) {
	id, err := r.readID()
	if err != nil {
		return TokenInvalid, err
	}
	switch id {
	case TokenU32, TokenI32, TokenF32:
		return id, r.skip(4)
	case TokenU64, TokenF64, TokenI64:
		return id, r.skip(8)
	case TokenBool:
		return id, r.skip(1)
	case TokenQuoted, TokenUnquoted:
		length, err := r.readU16()
		if err != nil {
			return TokenInvalid, err
		}
		return id, r.skip(int(length))
	default:
		return id, nil
	}
}

func (r *BinaryReader) readID() (TokenID, error) {
	v, err := r.readU16()
	return TokenID(v), err
}

func (r *BinaryReader) readU16() (uint16, error) {
	b, err := r.read(2)
	if err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint16(b), nil
}

func (r *BinaryReader) readU32() (uint32, error) {
	b, err := r.read(4)
	if err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint32(b), nil
}

func (r *BinaryReader) readI32() (int32, error) {
	v, err := r.readU32()
	if err != nil {
		return 0, err
	}
	return int32(v), nil //#nosec G115
}

func (r *BinaryReader) readU64() (uint64, error) {
	b, err := r.read(8)
	if err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint64(b), nil
}

func (r *BinaryReader) readI64() (int64, error) {
	v, err := r.readU64()
	if err != nil {
		return 0, err
	}
	return int64(v), nil //#nosec G115
}

func (r *BinaryReader) readBool() (bool, error) {
	b, err := r.read(1)
	if err != nil {
		return false, err
	}
	return b[0] != 0, nil
}

func (r *BinaryReader) readString() (string, error) {
	length, err := r.readU16()
	if err != nil {
		return "", err
	}
	b, err := r.read(int(length))
	if err != nil {
		return "", err
	}
	return internal.BytesToString(b), nil
}

func (r *BinaryReader) readF32() (float32, error) {
	i, err := r.readI32()
	if err != nil {
		return 0, err
	}
	return float32(i) / 1000.0, nil
}

func (r *BinaryReader) readF64() (float64, error) {
	i, err := r.readI64()
	if err != nil {
		return 0, err
	}
	val := float64(i) / 32768.0
	return math.Floor(val*10_0000.0) / 10_0000.0, nil
}

func (r *BinaryReader) read(length int) ([]byte, error) {
	r.buf = resize(r.buf, length)
	if _, err := io.ReadFull(r.r, r.buf); err != nil {
		return nil, err
	}
	r.offset += uint64(length) //#nosec G115
	return r.buf, nil
}

func (r *BinaryReader) skip(n int) error {
	var err error
	switch rr := r.r.(type) {
	case interface{ Discard(n int) (int, error) }:
		_, err = rr.Discard(n)
	case io.Seeker:
		_, err = rr.Seek(int64(n), io.SeekCurrent)
	default:
		if rf, ok := io.Discard.(io.ReaderFrom); ok {
			_, err = rf.ReadFrom(io.LimitReader(rr, int64(n)))
		} else {
			r.buf = resize(r.buf, n)
			_, err = io.ReadFull(r.r, r.buf)
		}
	}
	if err != nil {
		return err
	}
	r.offset += uint64(n) //#nosec G115
	return nil
}
