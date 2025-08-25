// SPDX-FileCopyrightText: 2025 Antoni Szymański
// SPDX-License-Identifier: MPL-2.0

package hoi4

import "github.com/antoniszymanski/hoi4-go/hoi4text"

type Map[K ~string, V any] map[K]V

func (x *Map[K, V]) UnmarshalHOI4(dec *hoi4text.Decoder) error {
	dec, err := dec.Value()
	if err != nil {
		return err
	}

	*x = make(Map[K, V])
	for {
		var key String
		if err = key.UnmarshalHOI4(dec); err != nil {
			break
		}
		if _, err = hoi4text.SkipToken(dec); err != nil {
			break
		}
		var elem V
		if err = any(&elem).(Unmarshaler).UnmarshalHOI4(dec); err != nil {
			break
		}
		(*x)[K(key)] = elem
	}
	if err == hoi4text.ErrEndOfObject {
		err = nil
	}
	return err
}

type MultiMap[K ~string, V any] map[K][]V

func (x *MultiMap[K, V]) UnmarshalHOI4(dec *hoi4text.Decoder) error {
	dec, err := dec.Value()
	if err != nil {
		return err
	}

	*x = make(MultiMap[K, V])
	for {
		var key String
		if err = key.UnmarshalHOI4(dec); err != nil {
			break
		}
		if _, err = hoi4text.SkipToken(dec); err != nil {
			break
		}
		var elem V
		if err = any(&elem).(Unmarshaler).UnmarshalHOI4(dec); err != nil {
			break
		}
		(*x)[K(key)] = append((*x)[K(key)], elem)
	}
	if err == hoi4text.ErrEndOfObject {
		err = nil
	}
	return err
}
