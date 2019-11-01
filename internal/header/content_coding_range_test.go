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

type ContentCodingTestSuite struct {
	suite.Suite

	sut header.ContentCodingRange
}

func TestContentCodingTestSuite(t *testing.T) {
	suite.Run(t, new(ContentCodingTestSuite))
}

func (s *ContentCodingTestSuite) SetupTest() {
	var err error
	s.sut, err = header.NewContentCodingRange("gzip;q=0.9")
	s.Require().NoError(err)
}

func (s *ContentCodingTestSuite) TestContentCoding_NewContentCodingRange() {

	tests := []struct {
		in string
	}{
		{"gzip"},
		{"x-gzip"},
		{"compress"},
		{"x-compress"},
		{"deflate"},
		{"identity"},
		{"*"},
		{"gzip;q=0.9"},
		{"x-gzip;q=0.9"},
		{"compress;q=0.9"},
		{"x-compress;q=0.9"},
		{"deflate;q=0.9"},
	}

	for _, test := range tests {
		s.Run(s.T().Name(), func() {
			// action.
			cc, err := header.NewContentCodingRange(test.in)

			// assert.
			s.Require().NoError(err)
			s.NotZero(cc)
		})
	}
}

func (s *ContentCodingTestSuite) TestContentCoding_NewContentCodingRange_Errors() {

	errTests := []struct {
		in  string
		out error
	}{
		{"zippy", header.ErrInvalidContentCodingRange},
		{"zippy;q=0.5", header.ErrInvalidContentCodingRange},
		{"", header.ErrEmptyContentCodingRange},
	}

	for _, test := range errTests {
		s.Run(s.T().Name(), func() {
			// action.
			cc, err := header.NewContentCodingRange(test.in)

			// assert.
			s.Require().EqualError(err, test.out.Error())
			s.Zero(cc)
		})
	}
}

func (s *ContentCodingTestSuite) TestContentCoding_IsWildcard_WithQValue() {

	// arrange.
	var err error
	s.sut, err = header.NewContentCodingRange("*;q=0.5")

	// action + assert.
	s.Require().NoError(err)
	s.Require().True(s.sut.IsWildcard())
}

func (s *ContentCodingTestSuite) TestContentCoding_IsWildcard_NoQValue() {

	// arrange.
	var err error
	s.sut, err = header.NewContentCodingRange("*")

	// action + assert.
	s.Require().NoError(err)
	s.Require().True(s.sut.IsWildcard())
}

func (s *ContentCodingTestSuite) TestContentCoding_IsWildcard_False() {

	// action + assert.
	s.Require().False(s.sut.IsWildcard())
}

func (s *ContentCodingTestSuite) TestContentCoding_IsIdentity() {

	// arrange.
	var err error
	s.sut, err = header.NewContentCodingRange("identity")

	// action + assert.
	s.Require().NoError(err)
	s.Require().True(s.sut.IsIdentity())
}

func (s *ContentCodingTestSuite) TestContentCoding_IsIdentity_False() {

	// action + assert.
	s.Require().False(s.sut.IsIdentity())
}

func (s *ContentCodingTestSuite) TestContentCoding_IsCoding_WithQValue() {

	// action + assert.
	s.Require().True(s.sut.IsCoding())
}

func (s *ContentCodingTestSuite) TestContentCoding_IsCoding_NoQValue() {

	// arrange.
	var err error
	s.sut, err = header.NewContentCodingRange("gzip")

	// action + assert.
	s.Require().NoError(err)
	s.Require().True(s.sut.IsCoding())
}

func (s *ContentCodingTestSuite) TestContentCoding_QualityValue() {

	// action + assert.
	s.Require().Equal(header.QualityValue(0.9), s.sut.QualityValue())
}

func (s *ContentCodingTestSuite) TestContentCoding_String() {

	// action  + assert.
	s.Require().Equal("gzip;q=0.900", s.sut.String())
}
