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

type AcceptLanguageTestSuite struct {
	suite.Suite
}

func TestAcceptLanguageTestSuite(t *testing.T) {
	suite.Run(t, new(AcceptLanguageTestSuite))
}

func (s AcceptLanguageTestSuite) TestAcceptLanguage_NewAcceptLanguage() {

	tests := []struct {
		name string
		in   []string
		err  error
	}{
		{"SingleRange", []string{"en-US"}, nil},
		{"MultipleRanges", []string{"en-US", "fr"}, nil},
		{"Empty", []string{}, nil},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			// action.
			ae, err := header.NewAcceptLanguage(test.in)

			// assert.
			if test.err != nil {
				s.Require().EqualError(err, test.err.Error())
				s.NotZero(ae)
			} else {
				s.Require().NoError(err)
				s.NotZero(ae)
			}
		})
	}
}

func (s AcceptLanguageTestSuite) TestAcceptLanguage_LanguageRanges() {
	english, _ := header.NewLanguageRange("en-US")
	englishWithQValue, _ := header.NewLanguageRange("en-US;q=0.4")
	frenchWithQValue, _ := header.NewLanguageRange("fr;q=0.8")
	german, _ := header.NewLanguageRange("de")

	tests := []struct {
		name string
		in   []string
		out  []header.LanguageRange
	}{
		{"MultipleRanges", []string{"en-US;q=0.4", "de", "fr;q=0.8"}, []header.LanguageRange{german, frenchWithQValue, englishWithQValue}},
		{"SingleRange", []string{"en-US"}, []header.LanguageRange{english}},
		{"Empty", []string{}, []header.LanguageRange{}},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			// action + assert.
			ae, err := header.NewAcceptLanguage(test.in)
			s.Require().NoError(err)
			s.Require().Len(ae.LanguageRanges(), len(test.out))
			for idx, lr := range ae.LanguageRanges() {
				s.Equal(test.out[idx].Tag(), lr.Tag())
				s.Equal(test.out[idx].QualityValue(), lr.QualityValue())
			}
		})
	}
}

func (s AcceptLanguageTestSuite) TestAcceptLanguage_IsEmpty() {

	tests := []struct {
		name string
		in   []string
		out  bool
	}{
		{"Empty", []string{}, true},
		{"NotEmpty", []string{"en-US", "fr"}, false},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			// action + assert.
			ae, err := header.NewAcceptLanguage(test.in)
			s.Require().NoError(err)
			s.Equal(test.out, ae.IsEmpty())
		})
	}
}

func (s AcceptLanguageTestSuite) TestAcceptLanguage_String() {

	tests := []struct {
		name string
		in   []string
		out  string
	}{
		{"Empty", []string{}, "Accept-Language: "},
		{"SingleRange", []string{"en-US"}, "Accept-Language: en-US;q=1.000"},
		{"MultipleRanges", []string{"en-US", "fr;q=0.8"}, "Accept-Language: en-US;q=1.000,fr;q=0.800"},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			// action + assert.
			ae, err := header.NewAcceptLanguage(test.in)
			s.Require().NoError(err)
			s.Equal(test.out, ae.String())
		})
	}
}
