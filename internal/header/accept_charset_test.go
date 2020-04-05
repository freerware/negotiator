package header_test

import (
	"testing"

	"github.com/freerware/negotiator/internal/header"
	"github.com/stretchr/testify/suite"
)

type AcceptCharsetTestSuite struct {
	suite.Suite
}

func TestAcceptCharsetTestSuite(t *testing.T) {
	suite.Run(t, new(AcceptCharsetTestSuite))
}

func (s AcceptCharsetTestSuite) TestAcceptCharset_NewAcceptCharset() {

	tests := []struct {
		name string
		in   []string
		err  error
	}{
		{"SingleRange", []string{"ascii"}, nil},
		{"MultipleRanges", []string{"ascii", "utf8"}, nil},
		{"Empty", []string{}, nil},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			// action.
			ae, err := header.NewAcceptCharset(test.in)

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

func (s AcceptCharsetTestSuite) TestAcceptCharset_CharsetRanges() {
	ascii, _ := header.NewCharsetRange("ascii")
	asciiWithQValue, _ := header.NewCharsetRange("ascii;q=0.4")
	utf8WithQValue, _ := header.NewCharsetRange("utf8;q=0.8")
	iso8859_1, _ := header.NewCharsetRange("iso-8859-1")

	tests := []struct {
		name string
		in   []string
		out  []header.CharsetRange
	}{
		{"MultipleRanges", []string{"ascii;q=0.4", "iso-8859-1", "utf8;q=0.8"}, []header.CharsetRange{iso8859_1, utf8WithQValue, asciiWithQValue}},
		{"SingleRange", []string{"ascii"}, []header.CharsetRange{ascii}},
		{"Empty", []string{}, []header.CharsetRange{}},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			// action + assert.
			ae, err := header.NewAcceptCharset(test.in)
			s.Require().NoError(err)
			s.Require().Len(ae.CharsetRanges(), len(test.out))
			for idx, ccr := range ae.CharsetRanges() {
				s.Equal(test.out[idx], ccr)
			}
		})
	}
}

func (s AcceptCharsetTestSuite) TestAcceptCharset_Compatible() {

	tests := []struct {
		name     string
		charsets []string
		in       string
		out      bool
		err      error
	}{
		{"Match", []string{"ascii;q=0.4", "iso-8859-1", "utf8;q=0.8"}, "utf8", true, nil},
		{"NoMatch", []string{"ascii"}, "utf8", false, nil},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			// action + assert.
			ae, err := header.NewAcceptCharset(test.charsets)
			s.Require().NoError(err)
			c, err := ae.Compatible(test.in)
			if test.err != nil {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
			}
			s.Equal(test.out, c)
		})
	}
}

func (s AcceptCharsetTestSuite) TestAcceptCharset_IsEmpty() {

	tests := []struct {
		name string
		in   []string
		out  bool
	}{
		{"Empty", []string{}, true},
		{"NotEmpty", []string{"ascii", "utf8"}, false},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			// action + assert.
			ae, err := header.NewAcceptCharset(test.in)
			s.Require().NoError(err)
			s.Equal(test.out, ae.IsEmpty())
		})
	}
}

func (s AcceptCharsetTestSuite) TestAcceptCharset_String() {

	tests := []struct {
		name string
		in   []string
		out  string
	}{
		{"Empty", []string{}, "Accept-Charset: "},
		{"SingleRange", []string{"ascii"}, "Accept-Charset: ascii;q=1.000"},
		{"MultipleRanges", []string{"ascii", "utf8;q=0.8"}, "Accept-Charset: ascii;q=1.000,utf8;q=0.800"},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			// action + assert.
			ae, err := header.NewAcceptCharset(test.in)
			s.Require().NoError(err)
			s.Equal(test.out, ae.String())
		})
	}
}
