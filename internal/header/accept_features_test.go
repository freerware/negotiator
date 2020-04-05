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

type AcceptFeaturesTestSuite struct {
	suite.Suite
}

func TestAcceptFeaturesTestSuite(t *testing.T) {
	suite.Run(t, new(AcceptFeaturesTestSuite))
}

func (s AcceptFeaturesTestSuite) TestAcceptFeatures_NewAcceptFeatures() {

	tests := []struct {
		name string
		in   []string
		err  error
	}{
		{"Single", []string{"foo"}, nil},
		{"Multiple", []string{"foo", "!bar"}, nil},
		{"Empty", []string{}, nil},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			// action.
			_, err := header.NewAcceptFeatures(test.in)

			// assert.
			if test.err != nil {
				s.Require().EqualError(err, test.err.Error())
			} else {
				s.Require().NoError(err)
			}
		})
	}
}

func (s AcceptFeaturesTestSuite) TestAcceptFeatures_IsEmpty() {

	tests := []struct {
		name string
		in   []string
		out  bool
	}{
		{"Empty", []string{}, true},
		{"NotEmpty", []string{"foo", "!bar"}, false},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			// action + assert.
			af, err := header.NewAcceptFeatures(test.in)
			s.Require().NoError(err)
			s.Equal(test.out, af.IsEmpty())
		})
	}
}

func (s AcceptFeaturesTestSuite) TestAcceptFeatures_String() {

	tests := []struct {
		name string
		in   []string
		out  string
	}{
		{"Empty", []string{}, "Accept-Features: "},
		{"Single", []string{"foo"}, "Accept-Features: foo"},
		{"Multiple", []string{"foo", "!bar"}, "Accept-Features: foo,!bar"},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			// action + assert.
			af, err := header.NewAcceptFeatures(test.in)
			s.Require().NoError(err)
			s.Equal(test.out, af.String())
		})
	}
}
