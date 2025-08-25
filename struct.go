// SPDX-FileCopyrightText: 2025 Antoni Szymański
// SPDX-License-Identifier: MPL-2.0

package hoi4

import (
	"maps"
	"reflect"
	"sync"
	"unsafe"

	"github.com/antoniszymanski/hoi4-go/hoi4text"
)

type Struct[T any] struct{ V T }

func (x *Struct[T]) UnmarshalHOI4(dec *hoi4text.Decoder) error {
	return unmarshalStruct(&x.V, dec)
}

type InlineStruct[T any] struct{}

func (x *InlineStruct[T]) UnmarshalHOI4(dec *hoi4text.Decoder) error {
	return unmarshalStruct((*T)(unsafe.Pointer(x)), dec)
}

func unmarshalStruct(i any, dec *hoi4text.Decoder) error {
	dec, err := dec.Value()
	if err != nil {
		return err
	}

	val := reflect.ValueOf(i).Elem()
	fields := parseStruct(val.Type(), nil)
	for {
		var key String
		if err = key.UnmarshalHOI4(dec); err != nil {
			break
		}
		if _, err = hoi4text.SkipToken(dec); err != nil {
			break
		}
		index, ok := fields[string(key)]
		if !ok {
			if err = dec.SkipValue(); err != nil {
				break
			}
			continue
		}
		field := val.FieldByIndex(index).Addr()
		u := field.Interface().(Unmarshaler)
		if err = u.UnmarshalHOI4(dec); err != nil {
			break
		}
	}
	if err == hoi4text.ErrEndOfObject {
		err = nil
	}
	return err
}

func parseStruct(typ reflect.Type, index []int) map[string][]int {
	if fields, ok := fieldsCache.Load(typ); ok {
		return fields.(map[string][]int)
	}
	fields := make(map[string][]int)
	for i := range typ.NumField() {
		field := typ.Field(i)
		if field.Anonymous {
			maps.Copy(fields, parseStruct(field.Type, field.Index))
		} else if field.IsExported() {
			name := field.Name
			if tag := field.Tag.Get("hoi4"); tag != "" {
				name = tag
			}
			fields[name] = concat(index, field.Index)
		}
	}
	fieldsCache.Store(typ, fields)
	return fields
}

var fieldsCache sync.Map // map[reflect.Type]map[string][]int

func concat[S ~[]E, E any](a, b S) S {
	if len(a) == 0 {
		return b
	} else if len(b) == 0 {
		return a
	}
	dst := make(S, len(a)+len(b))
	n := copy(dst, a)
	copy(dst[n:], b)
	return dst
}

type Duplicated[T any] []T

func (x *Duplicated[T]) UnmarshalHOI4(dec *hoi4text.Decoder) error {
	var elem T
	u := any(&elem).(Unmarshaler)
	if err := u.UnmarshalHOI4(dec); err != nil {
		return err
	}
	*x = append(*x, elem)
	return nil
}
