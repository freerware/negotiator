/* Copyright 2020 Freerware
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package header_test

import (
	"testing"

	"github.com/freerware/negotiator/internal/header"
	"github.com/stretchr/testify/suite"
)

type NegotiateDirectiveTestSuite struct {
	suite.Suite
}

func TestNegotiateDirectiveTestSuite(t *testing.T) {
	suite.Run(t, new(NegotiateDirectiveTestSuite))
}

func (s NegotiateDirectiveTestSuite) TestNegotiateDirective_NewNegotiateDirective() {

	tests := []struct {
		name string
		in   string
		err  error
	}{
		{"ValidDirective", header.NegotiateDirectiveGuessSmall.String(), nil},
		{"InvalidDirective", "", header.ErrEmptyNegotiateDirective},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			// action.
			_, err := header.NewNegotiateDirective(test.in)

			// assert.
			if test.err != nil {
				s.Require().EqualError(err, test.err.Error())
			} else {
				s.Require().NoError(err)
			}
		})
	}
}

func (s NegotiateDirectiveTestSuite) TestNegotiateDirective_IsWildcard() {

	tests := []struct {
		name string
		in   string
		out  bool
	}{
		{"NotWildcard", header.NegotiateDirectiveGuessSmall.String(), false},
		{"Wildcard", "*", true},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			// action + assert.
			d, err := header.NewNegotiateDirective(test.in)
			s.Require().NoError(err)
			s.Equal(test.out, d.IsWildcard())
		})
	}
}

func (s NegotiateDirectiveTestSuite) TestNegotiateDirective_IsRVSAVersion() {

	tests := []struct {
		name string
		in   string
		out  bool
	}{
		{"NotRVSAVersion", header.NegotiateDirectiveGuessSmall.String(), false},
		{"RVSAVersion", "1.0", true},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			// action + assert.
			d, err := header.NewNegotiateDirective(test.in)
			s.Require().NoError(err)
			s.Equal(test.out, d.IsRVSAVersion())
		})
	}
}

func (s NegotiateDirectiveTestSuite) TestNegotiateDirective_IsExtension() {

	tests := []struct {
		name string
		in   string
		out  bool
	}{
		{"NotExtension", header.NegotiateDirectiveGuessSmall.String(), false},
		{"Extension", "X", true},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			// action + assert.
			d, err := header.NewNegotiateDirective(test.in)
			s.Require().NoError(err)
			s.Equal(test.out, d.IsExtension())
		})
	}
}

func (s NegotiateDirectiveTestSuite) TestNegotiateDirective_String() {

	tests := []struct {
		name string
		in   string
		out  string
	}{
		{"Wildcard", "*", "*"},
		{"RVSAVersion", "1.0", "1.0"},
		{
			"NotExtension",
			header.NegotiateDirectiveTrans.String(),
			header.NegotiateDirectiveTrans.String(),
		},
		{"Extension", "X", "X"},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			// action + assert.
			d, err := header.NewNegotiateDirective(test.in)
			s.Require().NoError(err)
			s.Equal(test.out, d.String())
		})
	}
}
