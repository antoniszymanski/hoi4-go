// SPDX-FileCopyrightText: 2026 Antoni Szymański
// SPDX-License-Identifier: MPL-2.0

package hoi4

import (
	"io"
	"reflect"

	"github.com/antoniszymanski/hoi4-go/hoi4text"
)

func unmarshalAny(dec *hoi4text.Decoder, out reflect.Value) error {
	kind, err := dec.PeekKind()
	if err != nil {
		return &PeekKindError{err}
	}
	switch kind {
	case hoi4text.KindRoot:
		return unmarshalAnyRoot(dec, out)
	case hoi4text.KindScalar:
		return unmarshalAnyScalar(dec, out)
	case hoi4text.KindEmptyContainer:
		return unmarshalAnyEmptyContainer(dec, out)
	case hoi4text.KindArray:
		return unmarshalAnyArray(dec, out)
	case hoi4text.KindObject:
		return unmarshalAnyObject(dec, out)
	default:
		panic("unreachable")
	}
}

func unmarshalAnyRoot(dec *hoi4text.Decoder, out reflect.Value) (err error) {
	x := make(map[string][]any)
	for {
		if err = dec.IsEndOfContainer(); err != nil {
			break
		}
		var key string
		if err := unmarshalObjectKey(dec, &key); err != nil {
			return err
		}
		if id, err := dec.SkipToken(); err != nil {
			return err
		} else if id != hoi4text.TokenEqual {
			return &InvalidKeyValueSeparatorError{id}
		}
		var value any
		if err := unmarshalAny(dec, reflect.ValueOf(&value).Elem()); err != nil {
			return err
		}
		x[key] = append(x[key], value)
	}
	if err != io.EOF {
		return err
	}
	out.Set(reflect.ValueOf(x))
	return nil
}

func unmarshalAnyScalar(dec *hoi4text.Decoder, out reflect.Value) error {
	t, err := dec.ReadToken()
	if err != nil {
		return &ReadTokenError{dec.Offset(), err}
	}
	var x any
	switch t.ID() {
	case hoi4text.TokenOpen, hoi4text.TokenClose, hoi4text.TokenEqual:
		return &InvalidScalarError{t}
	case hoi4text.TokenU32:
		x = t.U32()
	case hoi4text.TokenU64:
		x = t.U64()
	case hoi4text.TokenI32:
		x = t.I32()
	case hoi4text.TokenBool:
		x = t.Bool()
	case hoi4text.TokenQuoted:
		x = t.Quoted()
	case hoi4text.TokenUnquoted:
		x = t.Unquoted()
	case hoi4text.TokenF32:
		x = t.F32()
	case hoi4text.TokenF64:
		x = t.F64()
	case hoi4text.TokenI64:
		x = t.I64()
	default:
		x = t
	}
	out.Set(reflect.ValueOf(x))
	return nil
}

func unmarshalAnyEmptyContainer(dec *hoi4text.Decoder, out reflect.Value) error {
	if id, err := dec.SkipToken(); err != nil {
		return err
	} else if id != hoi4text.TokenOpen {
		return &InvalidEmptyContainerError{id, FirstTokenOfEmptyContainer}
	}
	if id, err := dec.SkipToken(); err != nil {
		return err
	} else if id != hoi4text.TokenClose {
		return &InvalidEmptyContainerError{id, LastTokenOfEmptyContainer}
	}
	out.Set(reflect.ValueOf(struct{}{}))
	return nil
}

func unmarshalAnyArray(dec *hoi4text.Decoder, out reflect.Value) error {
	dec, err := dec.EnterContainer()
	if err != nil {
		return &EnterContainerError{err}
	}
	var x []any
	for {
		if err = dec.IsEndOfContainer(); err != nil {
			break
		}
		var elem any
		if err := unmarshalAny(dec, reflect.ValueOf(&elem).Elem()); err != nil {
			return err
		}
		x = append(x, elem)
	}
	if err != hoi4text.ErrEndOfContainer {
		return err
	}
	out.Set(reflect.ValueOf(x))
	return nil
}

func unmarshalAnyObject(dec *hoi4text.Decoder, out reflect.Value) error {
	dec, err := dec.EnterContainer()
	if err != nil {
		return &EnterContainerError{err}
	}
	x := make(map[string][]any)
	for {
		if err = dec.IsEndOfContainer(); err != nil {
			break
		}
		var key string
		if err := unmarshalObjectKey(dec, &key); err != nil {
			return err
		}
		if id, err := dec.SkipToken(); err != nil {
			return err
		} else if id != hoi4text.TokenEqual {
			return &InvalidKeyValueSeparatorError{id}
		}
		var value any
		if err := unmarshalAny(dec, reflect.ValueOf(&value).Elem()); err != nil {
			return err
		}
		x[key] = append(x[key], value)
	}
	if err != hoi4text.ErrEndOfContainer {
		return err
	}
	out.Set(reflect.ValueOf(x))
	return nil
}
