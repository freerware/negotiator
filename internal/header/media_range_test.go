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
 * See the License for the specific media governing permissions and
 * limitations under the License.
 */

package header_test

import (
	"strings"
	"testing"

	"github.com/freerware/negotiator/internal/header"
	"github.com/stretchr/testify/suite"
)

type MediaRangeTestSuite struct {
	suite.Suite
}

func TestMediaRangeTestSuite(t *testing.T) {
	suite.Run(t, new(MediaRangeTestSuite))
}

func (s *MediaRangeTestSuite) TestMediaRange_NewMediaRange() {

	tests := []struct {
		name string
		in   string
		err  error
	}{
		{"JSON", "application/json", nil},
		//{"Invalid", "jibberish", header.ErrInvalidMediaRange},
		//{"InvalidWithQValue", "jibberish;q=0.5", header.ErrInvalidMediaRange},
		{"Empty", "", header.ErrEmptyMediaRange},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			// action.
			c, err := header.NewMediaRange(test.in)

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

func (s *MediaRangeTestSuite) TestMediaRange_NewMediaRange_Compatible() {

	tests := []struct {
		name       string
		mediaRange string
		in         string
		out        bool
		err        error
	}{
		{"Wildcard", "*/*", "application/json", true, nil},
		{"SubTypeWildcard", "application/*", "application/json", true, nil},
		{"MatchLowerCase", "application/json", "application/json", true, nil},
		{"MatchUpperCase", "application/xml", "APPLICATION/XML", true, nil},
		{"MatchWithQValue", "application/xml;q=0.9", "APPLICATION/XML", true, nil},
		{"NoMatch", "application/json", "application/xml", false, nil},
		//{"Invalid", "application/json", "zoink", false, nil},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			// action + assert.
			c, err := header.NewMediaRange(test.mediaRange)
			s.Require().NoError(err)
			s.Require().NotZero(c)
			ok, err := c.Compatible(test.in)
			if test.err != nil {
				s.EqualError(err, test.err.Error())
			} else {
				s.NoError(err)
			}
			s.Equal(test.out, ok)
		})
	}
}

func (s *MediaRangeTestSuite) TestMediaRange_IsTypeWildcard() {

	tests := []struct {
		name       string
		mediaRange string
		out        bool
	}{
		{"WithoutQValue", "*/*", true},
		{"WithQValue", "*/*;q=0.5", true},
		{"NotWildcard", "application/json", false},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			// action + assert.
			c, err := header.NewMediaRange(test.mediaRange)
			s.Require().NoError(err)
			s.Require().NotZero(c)
			s.Equal(test.out, c.IsTypeWildcard())
		})
	}
}

func (s *MediaRangeTestSuite) TestMediaRange_IsSubTypeWildcard() {

	tests := []struct {
		name       string
		mediaRange string
		out        bool
	}{
		{"WithoutQValue", "*/*", true},
		{"WithQValue", "application/*;q=0.5", true},
		{"NotWildcard", "application/json", false},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			// action + assert.
			c, err := header.NewMediaRange(test.mediaRange)
			s.Require().NoError(err)
			s.Require().NotZero(c)
			s.Equal(test.out, c.IsSubTypeWildcard())
		})
	}
}

func (s *MediaRangeTestSuite) TestMediaRange_Param() {

	tests := []struct {
		name       string
		mediaRange string
		param      string
		outValue   string
		outOK      bool
	}{
		{"WithoutQValue", "application/json;foo=bar", "foo", "bar", true},
		{"WithQValue", "application/json;q=0.5;foo=bar", "foo", "bar", true},
		{"NoParams", "*", "foo", "", false},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			// action + assert.
			c, err := header.NewMediaRange(test.mediaRange)
			s.Require().NoError(err)
			s.Require().NotZero(c)
			val, ok := c.Param(test.param)
			s.Equal(test.outValue, val)
			s.Equal(test.outOK, ok)
		})
	}
}

func (s *MediaRangeTestSuite) TestMediaRange_Type() {

	tests := []struct {
		name       string
		mediaRange string
		out        string
	}{
		{"WithQValue", "application/json;q=0.5", "application"},
		{"WithoutQValue", "application/json", "application"},
		{"Uppercase", "APPLICATION/JSON", "application"},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			// action + assert.
			c, err := header.NewMediaRange(test.mediaRange)
			s.Require().NoError(err)
			s.Require().NotZero(c)
			s.Equal(test.out, c.Type())
		})
	}
}

func (s *MediaRangeTestSuite) TestMediaRange_Precedence() {

	tests := []struct {
		name       string
		mediaRange string
		out        int
	}{
		{"WithQValue", "application/json;q=0.5", 3},
		{"WithoutQValue", "application/json", 2},
		{"Uppercase", "APPLICATION/JSON", 2},
		{"Wildcard", "*/*", 0},
		{"WildcardWithParams", "*/*;foo=bar", 1},
		{"SubTypeWildcard", "application/*", 1},
		{"SubTypeWildcardWithParams", "application/*;foo=bar", 2},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			// action + assert.
			c, err := header.NewMediaRange(test.mediaRange)
			s.Require().NoError(err)
			s.Require().NotZero(c)
			s.Equal(test.out, c.Precedence())
		})
	}
}

func (s *MediaRangeTestSuite) TestMediaRange_SubType() {

	tests := []struct {
		name       string
		mediaRange string
		out        string
	}{
		{"WithQValue", "application/json;q=0.5", "json"},
		{"WithoutQValue", "application/json", "json"},
		{"Uppercase", "APPLICATION/JSON", "json"},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			// action + assert.
			c, err := header.NewMediaRange(test.mediaRange)
			s.Require().NoError(err)
			s.Require().NotZero(c)
			s.Equal(test.out, c.SubType())
		})
	}
}

func (s *MediaRangeTestSuite) TestMediaRange_QualityValue() {

	tests := []struct {
		name       string
		mediaRange string
		out        header.QualityValue
	}{
		{"WithoutQValue", "application/json", header.QualityValue(1.0)},
		{"WithQValue", "application/json;q=0.5", header.QualityValue(0.5)},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			// action + assert.
			c, err := header.NewMediaRange(test.mediaRange)
			s.Require().NoError(err)
			s.Require().NotZero(c)
			s.Equal(test.out, c.QualityValue())
		})
	}
}

func (s *MediaRangeTestSuite) TestMediaRange_String() {

	tests := []struct {
		name       string
		mediaRange string
		out        string
	}{
		{"WithoutQValueNoParams", "application/json", "application/json;q=1.000"},
		{"WithoutQValueWithParams", "application/json;foo=bar", "application/json;q=1.000;foo=bar"},
		{"WithQValue", "application/json;q=0.5", "application/json;q=0.500"},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			// action + assert.
			c, err := header.NewMediaRange(test.mediaRange)
			s.Require().NoError(err)
			s.Require().NotZero(c)
			s.Equal(strings.ToLower(test.out), strings.ToLower(c.String()))
		})
	}
}
