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
	"strings"
	"testing"

	"github.com/freerware/negotiator/internal/header"
	"github.com/stretchr/testify/suite"
)

type LanguageRangeTestSuite struct {
	suite.Suite
}

func TestLanguageRangeRangeTestSuite(t *testing.T) {
	suite.Run(t, new(LanguageRangeTestSuite))
}

func (s *LanguageRangeTestSuite) TestLanguageRange_NewLanguageRange() {

	tests := []struct {
		name string
		in   string
		err  error
	}{
		{"Gzip", "en-US", nil},
		//{"Invalid", "jibberish", header.ErrInvalidLanguageRange},
		//{"InvalidWithQValue", "jibberish;q=0.5", header.ErrInvalidLanguageRange},
		{"Empty", "", header.ErrEmptyLanguageRange},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			// action.
			c, err := header.NewLanguageRange(test.in)

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

func (s *LanguageRangeTestSuite) TestLanguageRange_NewLanguageRange_Compatible() {

	tests := []struct {
		name     string
		language string
		in       string
		out      bool
	}{
		{"Wildcard", "*", "en-US", true},
		{"MatchLowerCase", "en-US", "en-US", true},
		{"MatchUpperCase", "fr", "FR", true},
		{"MatchWithQValue", "fr;q=0.9", "FR", true},
		{"NoMatch", "en-US", "fr", false},
		{"Invalid", "en-US", "zippy", false},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			// action + assert.
			c, err := header.NewLanguageRange(test.language)
			s.Require().NoError(err)
			s.Require().NotZero(c)
			s.Equal(test.out, c.Compatible(test.in))
		})
	}
}

func (s *LanguageRangeTestSuite) TestLanguageRange_IsWildcard() {

	tests := []struct {
		name     string
		language string
		out      bool
	}{
		{"WithoutQValue", "*", true},
		{"WithQValue", "*;q=0.5", true},
		{"NotWildcard", "en-US", false},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			// action + assert.
			c, err := header.NewLanguageRange(test.language)
			s.Require().NoError(err)
			s.Require().NotZero(c)
			s.Equal(test.out, c.IsWildcard())
		})
	}
}

func (s *LanguageRangeTestSuite) TestLanguageRange_IsTag() {

	tests := []struct {
		name     string
		language string
		out      bool
	}{
		{"WithoutQValue", "en-US", true},
		{"WithQValue", "en-US;q=0.5", true},
		{"NotLanguageTag", "*", false},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			// action + assert.
			c, err := header.NewLanguageRange(test.language)
			s.Require().NoError(err)
			s.Require().NotZero(c)
			s.Equal(test.out, c.IsTag())
		})
	}
}

func (s *LanguageRangeTestSuite) TestLanguageRange_Tag() {

	tests := []struct {
		name     string
		language string
		out      string
	}{
		{"WithQValue", "en-US;q=0.5", "en-US"},
		{"WithoutQValue", "en-US", "en-US"},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			// action + assert.
			c, err := header.NewLanguageRange(test.language)
			s.Require().NoError(err)
			s.Require().NotZero(c)
			s.Equal(test.out, c.Tag())
		})
	}
}

func (s *LanguageRangeTestSuite) TestLanguageRange_QualityValue() {

	tests := []struct {
		name     string
		language string
		out      header.QualityValue
	}{
		{"WithoutQValue", "en-US", header.QualityValue(1.0)},
		{"WithQValue", "en-US;q=0.5", header.QualityValue(0.5)},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			// action + assert.
			c, err := header.NewLanguageRange(test.language)
			s.Require().NoError(err)
			s.Require().NotZero(c)
			s.Equal(test.out, c.QualityValue())
		})
	}
}

func (s *LanguageRangeTestSuite) TestLanguageRange_String() {

	tests := []struct {
		name     string
		language string
		out      string
	}{
		{"WithoutQValue", "en-US", "en-US;q=1.000"},
		{"WithQValue", "en-US;q=0.5", "en-US;q=0.500"},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			// action + assert.
			c, err := header.NewLanguageRange(test.language)
			s.Require().NoError(err)
			s.Require().NotZero(c)
			s.Equal(strings.ToLower(test.out), strings.ToLower(c.String()))
		})
	}
}
