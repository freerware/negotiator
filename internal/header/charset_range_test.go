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

	sut header.CharsetRange
}

func TestCharsetTestSuite(t *testing.T) {
	suite.Run(t, new(CharsetTestSuite))
}

func (s *CharsetTestSuite) SetupTest() {
	var err error
	s.sut, err = header.NewCharsetRange("utf8;q=0.9")
	s.Require().NoError(err)
}

func (s *CharsetTestSuite) TestCharset_NewCharsetRange() {

	tests := []struct {
		in string
	}{
		{"utf8"},
		{"*"},
		{"utf8;q-0.9"},
	}

	for _, test := range tests {
		s.Run(s.T().Name(), func() {
			// action.
			c, err := header.NewCharsetRange(test.in)

			// assert.
			s.Require().NoError(err)
			s.NotZero(c)
		})
	}
}

func (s *CharsetTestSuite) TestCharset_NewCharsetRange_Errors() {

	errTests := []struct {
		in  string
		out error
	}{
		{"", header.ErrEmptyCharsetRange},
	}

	for _, test := range errTests {
		s.Run(s.T().Name(), func() {
			// action.
			c, err := header.NewCharsetRange(test.in)

			// assert.
			s.Require().EqualError(err, test.out.Error())
			s.Zero(c)
		})
	}
}

func (s *CharsetTestSuite) TestCharset_IsWildcard_WithQValue() {

	// arrange.
	var err error
	s.sut, err = header.NewCharsetRange("*;q=0.5")

	// action + assert.
	s.Require().NoError(err)
	s.Require().True(s.sut.IsWildcard())
}

func (s *CharsetTestSuite) TestCharset_IsWildcard_NoQValue() {

	// arrange.
	var err error
	s.sut, err = header.NewCharsetRange("*")

	// action + assert.
	s.Require().NoError(err)
	s.Require().True(s.sut.IsWildcard())
}

func (s *CharsetTestSuite) TestCharset_IsWildcard_False() {

	// action + assert.
	s.Require().False(s.sut.IsWildcard())
}

func (s *CharsetTestSuite) TestCharset_IsCharset_WithQValue() {

	// action + assert.
	s.Require().True(s.sut.IsCharset())
}

func (s *CharsetTestSuite) TestCharset_IsCharset_NoQValue() {

	// arrange.
	var err error
	s.sut, err = header.NewCharsetRange("utf8")

	// action + assert.
	s.Require().NoError(err)
	s.Require().True(s.sut.IsCharset())
}

func (s *CharsetTestSuite) TestCharset_QualityValue() {

	// action + assert.
	s.Require().Equal(header.QualityValue(0.9), s.sut.QualityValue())
}

func (s *CharsetTestSuite) TestCharset_String() {

	// action  + assert.
	s.Require().Equal("utf8;q=0.900", s.sut.String())
}
