// SPDX-FileCopyrightText: 2025 Antoni Szymański
// SPDX-License-Identifier: MPL-2.0

package hoi4_test

import (
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"
	"sync"
	"testing"

	"github.com/alecthomas/repr"
	"github.com/antoniszymanski/hoi4-go"
	"github.com/antoniszymanski/hoi4-go/hoi4date"
	"github.com/sergi/go-diff/diffmatchpatch"
)

func Test(t *testing.T) {
	in, err := savefile()
	if err != nil {
		t.Fatal(err)
	}
	var actual Save
	if err := hoi4.Unmarshal(in, &actual); err != nil {
		t.Fatal(err)
	}
	expected := Save{
		Player: "FRA",
		Date: hoi4date.Date{
			Year:  1936,
			Month: 1,
			Day:   1,
			Hour:  12,
		},
		PlayerCountries: map[string]PlayerCountry{
			"FRA": {
				User:          "comagoosie",
				CountryLeader: true,
				ID:            1,
			},
		},
	}
	if !reflect.DeepEqual(expected, actual) {
		dmp := diffmatchpatch.New()
		repr := func(v any) string {
			var sb strings.Builder
			repr.New(&sb).Print(v)
			return sb.String()
		}
		diffs := dmp.DiffMain(repr(expected), repr(actual), false)
		fmt.Fprint(t.Output(), dmp.DiffPrettyText(diffs)) //nolint:errcheck
		t.Fail()
	}
}

func Benchmark(b *testing.B) {
	in, err := savefile()
	if err != nil {
		b.Fatal(err)
	}
	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		var out Save
		if err := hoi4.Unmarshal(in, &out); err != nil {
			b.Fatal(err)
		}
	}
}

type Save struct {
	Player          string                   `hoi4:"player"`
	Date            hoi4date.Date            `hoi4:"date"`
	PlayerCountries map[string]PlayerCountry `hoi4:"player_countries"`
}

type PlayerCountry struct {
	User          string `hoi4:"user"`
	CountryLeader bool   `hoi4:"country_leader"`
	ID            int64  `hoi4:"id"`
}

var savefile = sync.OnceValues(func() ([]byte, error) {
	resp, err := http.Get("https://cdn-dev.pdx.tools/hoi4-saves/1.10-ironman.hoi4")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close() //nolint:errcheck
	return io.ReadAll(resp.Body)
})
