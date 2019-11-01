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
	"fmt"
	"testing"

	"github.com/freerware/negotiator/internal/header"
	"github.com/stretchr/testify/suite"
	"golang.org/x/text/language"
)

type LanguageRangeTestSuite struct {
	suite.Suite

	sut header.LanguageRange
}

func TestLanguageRangeTestSuite(t *testing.T) {
	suite.Run(t, new(LanguageRangeTestSuite))
}

func (s *LanguageRangeTestSuite) SetupTest() {
	var err error
	s.sut, err = header.NewLanguageRange("en;q=0.9")
	s.Require().NoError(err)
}

func (s *LanguageRangeTestSuite) TestLanguageRange_NewLanguageRange() {

	tests := []struct {
		in string
	}{
		{"en-US"},
		{"*"},
	}

	for _, test := range tests {
		s.Run(s.T().Name(), func() {
			// action.
			lr, err := header.NewLanguageRange(test.in)

			// assert.
			s.Require().NoError(err)
			s.NotZero(lr)
		})
	}
}

func (s *LanguageRangeTestSuite) TestLanguageRange_NewLanguageRange_Errors() {

	invalidErr := fmt.Errorf("language: tag is not well-formed")
	errTests := []struct {
		in  string
		out error
	}{
		{"jibberish", invalidErr},
		{"jibberish;q=0.5", invalidErr},
		{"", header.ErrEmptyLanguageRange},
	}

	for _, test := range errTests {
		s.Run(s.T().Name(), func() {
			// action.
			cc, err := header.NewLanguageRange(test.in)

			// assert.
			s.Require().EqualError(err, test.out.Error())
			s.Zero(cc)
		})
	}
}

func (s *LanguageRangeTestSuite) TestLanguageRange_IsWildcard_WithQValue() {

	// arrange.
	var err error
	s.sut, err = header.NewLanguageRange("*;q=0.5")

	// action + assert.
	s.Require().NoError(err)
	s.Require().True(s.sut.IsWildcard())
}

func (s *LanguageRangeTestSuite) TestLanguageRange_IsWildcard_NoQValue() {

	// arrange.
	var err error
	s.sut, err = header.NewLanguageRange("*")

	// action + assert.
	s.Require().NoError(err)
	s.Require().True(s.sut.IsWildcard())
}

func (s *LanguageRangeTestSuite) TestLanguageRange_IsWildcard_False() {

	// action + assert.
	s.Require().False(s.sut.IsWildcard())
}

func (s *LanguageRangeTestSuite) TestLanguageRange_IsTag_WithQValue() {

	// action + assert.
	s.Require().True(s.sut.IsTag())
}

func (s *LanguageRangeTestSuite) TestLanguageRange_IsTag_NoQValue() {

	// arrange.
	var err error
	s.sut, err = header.NewLanguageRange(language.English.String())

	// action + assert.
	s.Require().NoError(err)
	s.Require().True(s.sut.IsTag())
}

func (s *LanguageRangeTestSuite) TestLanguageRange_QualityValue() {

	// action + assert.
	s.Require().Equal(header.QualityValue(0.9), s.sut.QualityValue())
}

func (s *LanguageRangeTestSuite) TestLanguageRange_String() {

	// action  + assert.
	s.Require().Equal("en;q=0.900", s.sut.String())
}
