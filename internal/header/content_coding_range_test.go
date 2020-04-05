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
}

func TestContentCodingTestSuite(t *testing.T) {
	suite.Run(t, new(ContentCodingTestSuite))
}

func (s *ContentCodingTestSuite) TestContentCoding_NewContentCodingRange() {

	tests := []struct {
		name string
		in   string
		err  error
	}{
		{"Gzip", "gzip", nil},
		{"XGzip", "x-gzip", nil},
		{"Compress", "compress", nil},
		{"XCompress", "x-compress", nil},
		{"Deflate", "deflate", nil},
		{"Identity", "identity", nil},
		{"Wildcard", "*", nil},
		{"GzipWithQValue", "gzip;q=0.9", nil},
		{"XGzipWithQValue", "x-gzip;q=0.9", nil},
		{"CompressWithQValue", "compress;q=0.9", nil},
		{"XCompressWithQValue", "x-compress;q=0.9", nil},
		{"DeflateWithQValue", "deflate;q=0.9", nil},
		{"Invalid", "zippy", header.ErrInvalidContentCodingRange},
		{"InvalidWithQValue", "zippy;q=0.5", header.ErrInvalidContentCodingRange},
		{"Empty", "", header.ErrEmptyContentCodingRange},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			// action.
			c, err := header.NewContentCodingRange(test.in)

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

func (s *ContentCodingTestSuite) TestContentCoding_NewContentCodingRange_Compatible() {

	tests := []struct {
		name          string
		contentCoding string
		in            string
		out           bool
	}{
		{"Wildcard", "*", "gzip", true},
		{"MatchLowerCase", "gzip", "gzip", true},
		{"MatchUpperCase", "compress", "COMPRESS", true},
		{"MatchWithQValue", "deflate;q=0.9", "DEFLATE", true},
		{"NoMatch", "deflate", "compress", false},
		{"Invalid", "deflate", "zippy", false},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			// action + assert.
			c, err := header.NewContentCodingRange(test.contentCoding)
			s.Require().NoError(err)
			s.Require().NotZero(c)
			s.Equal(test.out, c.Compatible(test.in))
		})
	}
}

func (s *ContentCodingTestSuite) TestContentCoding_IsWildcard() {

	tests := []struct {
		name          string
		contentCoding string
		out           bool
	}{
		{"WithoutQValue", "*", true},
		{"WithQValue", "*;q=0.5", true},
		{"NotWildcard", "gzip", false},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			// action + assert.
			c, err := header.NewContentCodingRange(test.contentCoding)
			s.Require().NoError(err)
			s.Require().NotZero(c)
			s.Equal(test.out, c.IsWildcard())
		})
	}
}

func (s *ContentCodingTestSuite) TestContentCoding_IsCoding() {

	tests := []struct {
		name          string
		contentCoding string
		out           bool
	}{
		{"WithoutQValue", "gzip", true},
		{"WithQValue", "gzip;q=0.5", true},
		{"NotContentCoding", "*", false},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			// action + assert.
			c, err := header.NewContentCodingRange(test.contentCoding)
			s.Require().NoError(err)
			s.Require().NotZero(c)
			s.Equal(test.out, c.IsCoding())
		})
	}
}

func (s *ContentCodingTestSuite) TestContentCoding_IsIdentity() {

	tests := []struct {
		name          string
		contentCoding string
		out           bool
	}{
		{"Identity", "identity", true},
		{"NotIdentity", "gzip;q=0.5", false},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			// action + assert.
			c, err := header.NewContentCodingRange(test.contentCoding)
			s.Require().NoError(err)
			s.Require().NotZero(c)
			s.Equal(test.out, c.IsIdentity())
		})
	}
}

func (s *ContentCodingTestSuite) TestContentCoding_QualityValue() {

	tests := []struct {
		name          string
		contentCoding string
		out           header.QualityValue
	}{
		{"WithoutQValue", "gzip", header.QualityValue(1.0)},
		{"WithQValue", "gzip;q=0.5", header.QualityValue(0.5)},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			// action + assert.
			c, err := header.NewContentCodingRange(test.contentCoding)
			s.Require().NoError(err)
			s.Require().NotZero(c)
			s.Equal(test.out, c.QualityValue())
		})
	}
}

func (s *ContentCodingTestSuite) TestContentCoding_String() {

	tests := []struct {
		name          string
		contentCoding string
		out           string
	}{
		{"WithoutQValue", "gzip", "gzip;q=1.000"},
		{"WithQValue", "gzip;q=0.5", "gzip;q=0.500"},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			// action + assert.
			c, err := header.NewContentCodingRange(test.contentCoding)
			s.Require().NoError(err)
			s.Require().NotZero(c)
			s.Equal(test.out, c.String())
		})
	}
}
