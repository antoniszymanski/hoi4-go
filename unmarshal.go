// SPDX-FileCopyrightText: 2026 Antoni Szymański
// SPDX-License-Identifier: MPL-2.0

package hoi4

import (
	"reflect"
	"slices"
	"strconv"
	"sync"

	"github.com/antoniszymanski/checked-go"
	"github.com/antoniszymanski/hoi4-go/hoi4date"
	"github.com/antoniszymanski/hoi4-go/hoi4text"
)

func unmarshal(dec *hoi4text.Decoder, out reflect.Value) error {
	if out.Type() == reflect.TypeFor[any]() {
		return unmarshalAny(dec, out)
	}
	if u, ok := reflect.TypeAssert[Unmarshaler](out); ok {
		return u.UnmarshalHOI4(dec)
	}
	if out, ok := reflect.TypeAssert[*hoi4date.Date](out); ok {
		return unmarshalDate(dec, out)
	}
	switch out.Kind() {
	case reflect.Bool:
		return unmarshalBool(dec, out)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return unmarshalInt(dec, out)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return unmarshalUint(dec, out)
	case reflect.Float32, reflect.Float64:
		return unmarshalFloat(dec, out)
	case reflect.Interface:
		return unmarshalInterface(dec, out)
	case reflect.Map:
		return unmarshalMap(dec, out)
	case reflect.Pointer:
		return unmarshalPointer(dec, out)
	case reflect.Slice:
		return unmarshalSlice(dec, out)
	case reflect.String:
		return unmarshalString(dec, out)
	case reflect.Struct:
		return unmarshalStruct(dec, out)
	default:
		return &InvalidTypeError{out.Type()}
	}
}

func unmarshalDate(dec *hoi4text.Decoder, out *hoi4date.Date) error {
	t, err := dec.ReadToken()
	if err != nil {
		return &ReadTokenError{dec.Offset(), err}
	}
	var x hoi4date.Date
	var ok bool
	switch t.ID() {
	case hoi4text.TokenI32:
		x, ok = hoi4date.ParseBinary(t.I32())
		if !ok {
			return &ParseDateError[int32]{t.I32()}
		}
	case hoi4text.TokenQuoted:
		x, ok = hoi4date.Parse(t.Quoted())
		if !ok {
			return &ParseDateError[string]{t.Quoted()}
		}
	default:
		return &InvalidTokenError{t, reflect.TypeFor[hoi4date.Date]()}
	}
	*out = x
	return nil
}

func unmarshalBool(dec *hoi4text.Decoder, out reflect.Value) error {
	t, err := dec.ReadToken()
	if err != nil {
		return &ReadTokenError{dec.Offset(), err}
	} else if t.ID() != hoi4text.TokenBool {
		return &InvalidTokenError{t, out.Type()}
	}
	out.SetBool(t.Bool())
	return nil
}

func unmarshalInt(dec *hoi4text.Decoder, out reflect.Value) error {
	t, err := dec.ReadToken()
	if err != nil {
		return &ReadTokenError{dec.Offset(), err}
	}
	var x int64
	var ok bool
	switch t.ID() {
	case hoi4text.TokenU32:
		x = int64(t.U32())
	case hoi4text.TokenU64:
		x, ok = checked.Cast[int64](t.U64())
		if !ok {
			return &OverflowError[uint64]{t.U64(), out.Type()}
		}
	case hoi4text.TokenI32:
		x = int64(t.I32())
	case hoi4text.TokenF32:
		x = int64(t.F32())
	case hoi4text.TokenF64:
		x = int64(t.F64())
	case hoi4text.TokenI64:
		x = t.I64()
	default:
		return &InvalidTokenError{t, out.Type()}
	}
	if out.OverflowInt(x) {
		return &OverflowError[int64]{x, out.Type()}
	}
	out.SetInt(x)
	return nil
}

func unmarshalUint(dec *hoi4text.Decoder, out reflect.Value) error {
	t, err := dec.ReadToken()
	if err != nil {
		return &ReadTokenError{dec.Offset(), err}
	}
	var x uint64
	var ok bool
	switch t.ID() {
	case hoi4text.TokenU32:
		x = uint64(t.U32())
	case hoi4text.TokenU64:
		x = t.U64()
	case hoi4text.TokenI32:
		x, ok = checked.Cast[uint64](t.I32())
		if !ok {
			return &OverflowError[int32]{t.I32(), out.Type()}
		}
	case hoi4text.TokenF32:
		x = uint64(t.F32())
	case hoi4text.TokenF64:
		x = uint64(t.F64())
	case hoi4text.TokenI64:
		x, ok = checked.Cast[uint64](t.I64())
		if !ok {
			return &OverflowError[int64]{t.I64(), out.Type()}
		}
	default:
		return &InvalidTokenError{t, out.Type()}
	}
	if out.OverflowUint(x) {
		return &OverflowError[uint64]{x, out.Type()}
	}
	out.SetUint(x)
	return nil
}

func unmarshalFloat(dec *hoi4text.Decoder, out reflect.Value) error {
	t, err := dec.ReadToken()
	if err != nil {
		return &ReadTokenError{dec.Offset(), err}
	}
	var x float64
	switch t.ID() {
	case hoi4text.TokenU32:
		x = float64(t.U32())
	case hoi4text.TokenU64:
		x = float64(t.U64())
	case hoi4text.TokenI32:
		x = float64(t.I32())
	case hoi4text.TokenF32:
		x = float64(t.F32())
	case hoi4text.TokenF64:
		x = t.F64()
	case hoi4text.TokenI64:
		x = float64(t.I64())
	default:
		return &InvalidTokenError{t, out.Type()}
	}
	if out.OverflowFloat(x) {
		return &OverflowError[float64]{x, out.Type()}
	}
	out.SetFloat(x)
	return nil
}

func unmarshalInterface(dec *hoi4text.Decoder, out reflect.Value) error {
	if out.IsNil() {
		return ErrNilInterface
	}
	return unmarshal(dec, out.Elem())
}

func unmarshalMap(dec *hoi4text.Decoder, out reflect.Value) error {
	dec, err := dec.EnterContainer()
	if err != nil {
		return &EnterContainerError{err}
	}
	if err := unmarshalMapContent(dec, out); err != hoi4text.ErrEndOfContainer {
		return err
	}
	return nil
}

func unmarshalMapContent(dec *hoi4text.Decoder, out reflect.Value) error {
	typ := out.Type()
	out.Set(reflect.MakeMap(typ))
	for {
		if err := dec.IsEndOfContainer(); err != nil {
			return err
		}
		keyPtr := reflect.New(typ.Key())
		if err := unmarshal(dec, keyPtr); err != nil {
			return err
		}
		if id, err := dec.SkipToken(); err != nil {
			return err
		} else if id != hoi4text.TokenEqual {
			return &InvalidKeyValueSeparatorError{id}
		}
		elemPtr := reflect.New(typ.Elem())
		if err := unmarshal(dec, elemPtr); err != nil {
			return err
		}
		out.SetMapIndex(keyPtr.Elem(), elemPtr.Elem())
	}
}

func unmarshalPointer(dec *hoi4text.Decoder, out reflect.Value) error {
	if out.IsNil() {
		out.Set(reflect.New(out.Type().Elem()))
	}
	return unmarshal(dec, out.Elem())
}

func unmarshalSlice(dec *hoi4text.Decoder, out reflect.Value) error {
	dec, err := dec.EnterContainer()
	if err != nil {
		return &EnterContainerError{err}
	}
	elemType := out.Type().Elem()
	for {
		if err = dec.IsEndOfContainer(); err != nil {
			break
		}
		elemPtr := reflect.New(elemType)
		if err := unmarshal(dec, elemPtr); err != nil {
			return err
		}
		out.Set(reflect.Append(out, elemPtr.Elem()))
	}
	if err != hoi4text.ErrEndOfContainer {
		return err
	}
	return nil
}

func unmarshalString(dec *hoi4text.Decoder, out reflect.Value) error {
	t, err := dec.ReadToken()
	if err != nil {
		return &ReadTokenError{dec.Offset(), err}
	}
	var x string
	switch t.ID() {
	case hoi4text.TokenQuoted:
		x = t.Quoted()
	case hoi4text.TokenUnquoted:
		x = t.Unquoted()
	default:
		if !t.ID().IsID() {
			return &InvalidTokenError{t, out.Type()}
		}
		x = hoi4text.ResolveToken(t.ID())
		if x == "" {
			return &InvalidTokenError{t, out.Type()}
		}
	}
	out.SetString(x)
	return nil
}

func unmarshalStruct(dec *hoi4text.Decoder, out reflect.Value) error {
	dec, err := dec.EnterContainer()
	if err != nil {
		return &EnterContainerError{err}
	}
	if err := unmarshalStructContent(dec, out); err != hoi4text.ErrEndOfContainer {
		return err
	}
	return nil
}

func unmarshalStructContent(dec *hoi4text.Decoder, out reflect.Value) error {
	fieldIndices, err := fieldIndices(out.Type())
	if err != nil {
		return err
	}
	for {
		if err := dec.IsEndOfContainer(); err != nil {
			return err
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
		if index := fieldIndices[key]; len(index) > 0 {
			field := fieldByIndex(out, index)
			if err := unmarshal(dec, field.Addr()); err != nil {
				return err
			}
		} else {
			if err := dec.SkipValue(); err != nil {
				return err
			}
		}
	}
}

func fieldIndices(typ reflect.Type) (m map[string][]int, err error) {
	if x, ok := cache.Load(typ); ok {
		switch x := x.(type) {
		case map[string][]int:
			return x, nil
		case error:
			return nil, x
		default:
			panic("unreachable")
		}
	}
	cache.Store(typ, map[string][]int(nil)) // prevent infinite recursion
	defer func() {
		if err == nil {
			cache.Store(typ, m)
		} else {
			cache.Store(typ, err)
		}
	}()
	m = make(map[string][]int)
	for field := range typ.Fields() {
		switch {
		case field.Anonymous:
			switch field.Type.Kind() {
			case reflect.Interface:
				return nil, ErrEmbeddedInterface
			case reflect.Pointer:
				field.Type = field.Type.Elem()
			}
			fieldIndices, err := fieldIndices(field.Type)
			if err != nil {
				return nil, err
			}
			for key, index := range fieldIndices {
				m[key] = slices.Concat(field.Index, index)
			}
		case field.IsExported():
			if tag := field.Tag.Get("hoi4"); tag != "" {
				field.Name = tag
			}
			m[field.Name] = field.Index
		}
	}
	return m, nil
}

var cache sync.Map // map[reflect.Type](map[string][]int | error)

func unmarshalObjectKey(dec *hoi4text.Decoder, out *string) error {
	t, err := dec.ReadToken()
	if err != nil {
		return &ReadTokenError{dec.Offset(), err}
	}
	var x string
	switch t.ID() {
	case hoi4text.TokenQuoted:
		x = t.Quoted()
	case hoi4text.TokenUnquoted:
		x = t.Unquoted()
	case hoi4text.TokenU32:
		x = strconv.FormatUint(uint64(t.U32()), 10)
	case hoi4text.TokenU64:
		x = strconv.FormatUint(t.U64(), 10)
	case hoi4text.TokenI32:
		x = strconv.FormatInt(int64(t.I32()), 10)
	case hoi4text.TokenI64:
		x = strconv.FormatInt(t.I64(), 10)
	default:
		if !t.ID().IsID() {
			return &InvalidObjectKeyError{t}
		}
		x = hoi4text.ResolveToken(t.ID())
		if x == "" {
			return &InvalidObjectKeyError{t}
		}
	}
	*out = x
	return nil
}

func fieldByIndex(v reflect.Value, index []int) reflect.Value {
	if len(index) == 1 {
		return v.Field(index[0])
	}
	for i, x := range index {
		if i > 0 && v.Kind() == reflect.Pointer {
			if elemType := v.Type().Elem(); elemType.Kind() == reflect.Struct {
				if v.IsNil() {
					v.Set(reflect.New(elemType))
				}
				v = v.Elem()
			}
		}
		v = v.Field(x)
	}
	return v
}
