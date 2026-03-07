// SPDX-FileCopyrightText: 2025 Antoni Szymański
// SPDX-License-Identifier: MPL-2.0

package hoi4text

import (
	_ "embed"
	"strings"

	"github.com/antoniszymanski/hoi4-go/hoi4text/tokenmap"
)

func LookupToken(id TokenID) string {
	return tokens[uint16(id)]
}

func init() {
	var err error
	tokens, err = tokenmap.Decode(strings.NewReader(tokensData))
	if err != nil {
		panic(err)
	}
}

var tokens map[uint16]string

//go:embed tokens
var tokensData string
