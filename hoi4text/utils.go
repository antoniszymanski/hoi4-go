// SPDX-FileCopyrightText: 2025 Antoni Szymański
// SPDX-License-Identifier: MPL-2.0

package hoi4text

import (
	"strconv"

	"golang.org/x/exp/constraints"
)

func appendInt[T constraints.Signed](dst []byte, i T) []byte {
	return strconv.AppendInt(dst, int64(i), 10)
}

func appendUint[T constraints.Unsigned](dst []byte, i T) []byte {
	return strconv.AppendUint(dst, uint64(i), 10)
}

func appendFloat32(dst []byte, f float32) []byte {
	return strconv.AppendFloat(dst, float64(f), 'g', -1, 32)
}

func appendFloat64(dst []byte, f float64) []byte {
	return strconv.AppendFloat(dst, f, 'g', -1, 64)
}

func resize[S ~[]E, E any](s S, length int) S {
	if length <= cap(s) {
		return s[:length]
	}
	return append(s[:cap(s)], make(S, length-cap(s))...)
}
