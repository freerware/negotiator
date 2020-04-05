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

type CharsetTestSuite struct {
	suite.Suite
}

func TestCharsetTestSuite(t *testing.T) {
	suite.Run(t, new(CharsetTestSuite))
}

func (s *CharsetTestSuite) TestCharset_NewCharsetRange() {

	tests := []struct {
		name string
		in   string
		err  error
	}{
		{"WithoutQValue", "utf8", nil},
		{"Wildcard", "*", nil},
		{"WithQValue", "utf8;q=0.9", nil},
		{"Empty", "", header.ErrEmptyCharsetRange},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			// action.
			c, err := header.NewCharsetRange(test.in)

			// assert.
			if test.err != nil {
				s.Require().EqualError(err, test.err.Error())
				s.Zero(c)
			} else {
				s.Require().NoError(err)
				s.NotZero(c)
			}
		})
	}
}

func (s *CharsetTestSuite) TestCharset_NewCharsetRange_Compatible() {

	tests := []struct {
		name    string
		charset string
		in      string
		out     bool
	}{
		{"Wildcard", "*", "utf8", true},
		{"MatchLowerCase", "utf8", "utf8", true},
		{"MatchUpperCase", "utf8", "UTF8", true},
		{"MatchWithQValue", "utf8;q=0.9", "UTF8", true},
		{"NoMatch", "utf8", "ascii", false},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			// action + assert.
			c, err := header.NewCharsetRange(test.charset)
			s.Require().NoError(err)
			s.Require().NotZero(c)
			s.Equal(test.out, c.Compatible(test.in))
		})
	}
}

func (s *CharsetTestSuite) TestCharset_IsWildcard() {

	tests := []struct {
		name    string
		charset string
		out     bool
	}{
		{"WithoutQValue", "*", true},
		{"WithQValue", "*;q=0.5", true},
		{"NotWildcard", "utf8", false},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			// action + assert.
			c, err := header.NewCharsetRange(test.charset)
			s.Require().NoError(err)
			s.Require().NotZero(c)
			s.Equal(test.out, c.IsWildcard())
		})
	}
}

func (s *CharsetTestSuite) TestCharset_IsCharset_NoQValue() {

	tests := []struct {
		name    string
		charset string
		out     bool
	}{
		{"WithoutQValue", "utf8", true},
		{"WithQValue", "utf;q=0.5", true},
		{"NotCharset", "*", false},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			// action + assert.
			c, err := header.NewCharsetRange(test.charset)
			s.Require().NoError(err)
			s.Require().NotZero(c)
			s.Equal(test.out, c.IsCharset())
		})
	}
}

func (s *CharsetTestSuite) TestCharset_QualityValue() {

	tests := []struct {
		name    string
		charset string
		out     header.QualityValue
	}{
		{"WithoutQValue", "utf8", header.QualityValue(1.0)},
		{"WithQValue", "utf;q=0.5", header.QualityValue(0.5)},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			// action + assert.
			c, err := header.NewCharsetRange(test.charset)
			s.Require().NoError(err)
			s.Require().NotZero(c)
			s.Equal(test.out, c.QualityValue())
		})
	}
}

func (s *CharsetTestSuite) TestCharset_String() {

	tests := []struct {
		name    string
		charset string
		out     string
	}{
		{"WithoutQValue", "utf8", "utf8;q=1.000"},
		{"WithQValue", "utf;q=0.5", "utf;q=0.500"},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			// action + assert.
			c, err := header.NewCharsetRange(test.charset)
			s.Require().NoError(err)
			s.Require().NotZero(c)
			s.Equal(test.out, c.String())
		})
	}
}
