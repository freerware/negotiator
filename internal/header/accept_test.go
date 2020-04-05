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

type AcceptTestSuite struct {
	suite.Suite
}

func TestAcceptTestSuite(t *testing.T) {
	suite.Run(t, new(AcceptTestSuite))
}

func (s AcceptTestSuite) TestAccept_NewAccept() {

	tests := []struct {
		name string
		in   []string
		err  error
	}{
		{"SingleRange", []string{"application/json"}, nil},
		{"WildcardSubType", []string{"application/*"}, nil},
		{"Wildcard", []string{"*/*"}, nil},
		{"MultipleRanges", []string{"application/json", "application/yaml;q=0.9"}, nil},
		{"Empty", []string{}, nil},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			// action.
			ae, err := header.NewAccept(test.in)

			// assert.
			if test.err != nil {
				s.Require().EqualError(err, test.err.Error())
			} else {
				s.Require().NoError(err)
				s.NotZero(ae)
			}
		})
	}
}

func (s AcceptTestSuite) TestAccept_MediaRanges() {
	json, _ := header.NewMediaRange("application/json")
	jsonWithQValue, _ := header.NewMediaRange("application/json;q=0.4")
	yamlWithQValue, _ := header.NewMediaRange("text/yaml;q=0.8")
	xml, _ := header.NewMediaRange("application/xml")

	tests := []struct {
		name string
		in   []string
		out  []header.MediaRange
	}{
		{
			"MultipleRanges",
			[]string{
				"application/json;q=0.4",
				"application/xml",
				"text/yaml;q=0.8",
			},
			[]header.MediaRange{
				xml,
				yamlWithQValue,
				jsonWithQValue,
			},
		},
		{"SingleRange", []string{"application/json"}, []header.MediaRange{json}},
		{"Empty", []string{}, []header.MediaRange{}},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			// action + assert.
			ae, err := header.NewAccept(test.in)
			s.Require().NoError(err)
			s.Require().Len(ae.MediaRanges(), len(test.out))
			for idx, ccr := range ae.MediaRanges() {
				s.Equal(test.out[idx], ccr)
			}
		})
	}
}

func (s AcceptTestSuite) TestAccept_IsEmpty() {

	tests := []struct {
		name string
		in   []string
		out  bool
	}{
		{"Empty", []string{}, true},
		{"NotEmpty", []string{"application/json", "application/xml"}, false},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			// action + assert.
			ae, err := header.NewAccept(test.in)
			s.Require().NoError(err)
			s.Equal(test.out, ae.IsEmpty())
		})
	}
}

func (s AcceptTestSuite) TestAccept_String() {

	tests := []struct {
		name string
		in   []string
		out  string
	}{
		{"Empty", []string{}, "Accept: "},
		{
			"SingleRange",
			[]string{"application/json"},
			"Accept: application/json;q=1.000",
		},
		{
			"MultipleRanges",
			[]string{
				"application/json",
				"application/xml;q=0.8",
			},
			"Accept: application/json;q=1.000,application/xml;q=0.800",
		},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			// action + assert.
			ae, err := header.NewAccept(test.in)
			s.Require().NoError(err)
			s.Equal(test.out, ae.String())
		})
	}
}

func (s *AcceptTestSuite) TestAccept_Compatible() {

	tests := []struct {
		name       string
		mediaRange []string
		in         string
		out        bool
		err        error
	}{
		{"Wildcard", []string{"*/*"}, "application/json", true, nil},
		{"SubTypeWildcard", []string{"application/*"}, "application/json", true, nil},
		{"MatchLowerCase", []string{"application/json"}, "application/json", true, nil},
		{"MatchUpperCase", []string{"application/xml"}, "APPLICATION/XML", true, nil},
		{"MatchWithQValue", []string{"application/xml;q=0.9"}, "APPLICATION/XML", true, nil},
		{"NoMatch", []string{"application/json"}, "application/xml", false, nil},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			// action + assert.
			c, err := header.NewAccept(test.mediaRange)
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
