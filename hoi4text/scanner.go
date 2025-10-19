// SPDX-FileCopyrightText: 2025 Antoni Szymański
// SPDX-License-Identifier: MPL-2.0

package hoi4text

import "io"

type Scanner struct {
	r     Reader
	token Token
	err   error
}

func NewScanner(r Reader) *Scanner {
	return &Scanner{r: r}
}

func (s *Scanner) Scan() bool {
	s.token, s.err = s.r.ReadToken()
	return s.err == nil
}

func (s *Scanner) Token() Token {
	return s.token
}

func (s *Scanner) Err() error {
	if s.err == io.EOF {
		return nil
	}
	return s.err
}
