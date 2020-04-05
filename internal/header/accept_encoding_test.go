package header_test

import (
	"testing"

	"github.com/freerware/negotiator/internal/header"
	"github.com/stretchr/testify/suite"
)

type AcceptEncodingTestSuite struct {
	suite.Suite
}

func TestAcceptEncodingTestSuite(t *testing.T) {
	suite.Run(t, new(AcceptEncodingTestSuite))
}

func (s AcceptEncodingTestSuite) TestAcceptEncoding_NewAcceptEncoding() {

	tests := []struct {
		name string
		in   []string
		err  error
	}{
		{"SingleRange", []string{"gzip"}, nil},
		{"MultipleRanges", []string{"gzip", "compress"}, nil},
		{"InvalidContentCodingRange", []string{"zippy;q=0.5"}, header.ErrInvalidContentCodingRange},
		{"Empty", []string{}, nil},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			// action.
			ae, err := header.NewAcceptEncoding(test.in)

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

func (s AcceptEncodingTestSuite) TestAcceptEncoding_CodingRanges() {
	gzip, _ := header.NewContentCodingRange("gzip")
	gzipWithQValue, _ := header.NewContentCodingRange("gzip;q=0.4")
	compressWithQValue, _ := header.NewContentCodingRange("compress;q=0.8")
	deflate, _ := header.NewContentCodingRange("deflate")

	tests := []struct {
		name string
		in   []string
		out  []header.ContentCodingRange
	}{
		{"MultipleRanges", []string{"gzip;q=0.4", "deflate", "compress;q=0.8"}, []header.ContentCodingRange{deflate, compressWithQValue, gzipWithQValue}},
		{"SingleRange", []string{"gzip"}, []header.ContentCodingRange{gzip}},
		{"Empty", []string{}, []header.ContentCodingRange{}},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			// action + assert.
			ae, err := header.NewAcceptEncoding(test.in)
			s.Require().NoError(err)
			s.Require().Len(ae.CodingRanges(), len(test.out))
			for idx, ccr := range ae.CodingRanges() {
				s.Equal(test.out[idx], ccr)
			}
		})
	}
}

func (s AcceptEncodingTestSuite) TestAcceptEncoding_IsEmpty() {

	tests := []struct {
		name string
		in   []string
		out  bool
	}{
		{"Empty", []string{}, true},
		{"NotEmpty", []string{"gzip", "compress"}, false},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			// action + assert.
			ae, err := header.NewAcceptEncoding(test.in)
			s.Require().NoError(err)
			s.Equal(test.out, ae.IsEmpty())
		})
	}
}

func (s AcceptEncodingTestSuite) TestAcceptEncoding_String() {

	tests := []struct {
		name string
		in   []string
		out  string
	}{
		{"Empty", []string{}, "Accept-Encoding: "},
		{"SingleRange", []string{"gzip"}, "Accept-Encoding: gzip;q=1.000"},
		{"MultipleRanges", []string{"gzip", "compress;q=0.8"}, "Accept-Encoding: gzip;q=1.000,compress;q=0.800"},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			// action + assert.
			ae, err := header.NewAcceptEncoding(test.in)
			s.Require().NoError(err)
			s.Equal(test.out, ae.String())
		})
	}
}
