// SPDX-FileCopyrightText: 2025 Antoni Szymański
// SPDX-License-Identifier: MPL-2.0

package hoi4_test

import (
	"archive/zip"
	"bytes"
	"io"
	"net/http"
	"reflect"
	"testing"

	"github.com/antoniszymanski/hoi4-go"
	"github.com/k0kubun/pp/v3"
)

type Save struct {
	hoi4.InlineStruct[Save]
	Player          hoi4.String                     `hoi4:"player"`
	Date            hoi4.Date                       `hoi4:"date"`
	PlayerCountries hoi4.Map[string, PlayerCountry] `hoi4:"player_countries"`
}

type PlayerCountry struct {
	hoi4.InlineStruct[PlayerCountry]
	User          hoi4.String `hoi4:"user"`
	CountryLeader hoi4.Bool   `hoi4:"country_leader"`
	ID            hoi4.Int64  `hoi4:"id"`
}

func Benchmark(b *testing.B) {
	resp, err := http.Get("https://hoi4saves-test-cases.s3.us-west-002.backblazeb2.com/1.10-ironman.zip")
	if err != nil {
		b.Fatal(err)
	}
	defer resp.Body.Close() //nolint:errcheck

	in, err := io.ReadAll(resp.Body)
	if err != nil {
		b.Fatal(err)
	}

	r := bytes.NewReader(in)
	zr, err := zip.NewReader(r, r.Size())
	if err != nil {
		b.Fatal(err)
	}

	f, err := zr.Open("1.10-ironman.hoi4")
	if err != nil {
		b.Fatal(err)
	}
	defer f.Close() //nolint:errcheck

	in, err = io.ReadAll(f)
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
	b.StopTimer()

	var actual Save
	if err := hoi4.Unmarshal(in, &actual); err != nil {
		b.Fatal(err)
	}
	expected := Save{
		InlineStruct: hoi4.InlineStruct[Save]{},
		Player:       "FRA",
		Date: hoi4.Date{
			Year:  1936,
			Month: 1,
			Day:   1,
			Hour:  12,
		},
		PlayerCountries: hoi4.Map[string, PlayerCountry]{
			"FRA": PlayerCountry{
				InlineStruct:  hoi4.InlineStruct[PlayerCountry]{},
				User:          "comagoosie",
				CountryLeader: true,
				ID:            1,
			},
		},
	}
	if !reflect.DeepEqual(&actual, &expected) {
		b.Fatal(pp.Sprintf("Not equal:\nactual: %v\nexpected: %v", &actual, &expected))
	}
}
