package header_test

import (
	"testing"

	"github.com/freerware/negotiator/internal/header"
	"github.com/stretchr/testify/suite"
)

type NegotiateTestSuite struct {
	suite.Suite
}

func TestNegotiateTestSuite(t *testing.T) {
	suite.Run(t, new(NegotiateTestSuite))
}

func (s NegotiateTestSuite) TestNegotiate_NewNegotiate() {

	tests := []struct {
		name string
		in   []string
		err  error
	}{
		{"SingleDirective", []string{header.NegotiateDirectiveVList.String()}, nil},
		{"MultipleDirectives", []string{header.NegotiateDirectiveTrans.String(), header.NegotiateDirectiveVList.String()}, nil},
		{"Empty", []string{}, nil},
		{"InvalidDirective", []string{""}, header.ErrEmptyNegotiateDirective},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			// action.
			_, err := header.NewNegotiate(test.in)

			// assert.
			if test.err != nil {
				s.Require().EqualError(err, test.err.Error())
			} else {
				s.Require().NoError(err)
			}
		})
	}
}

func (s NegotiateTestSuite) TestNegotiate_Directives() {

	tests := []struct {
		name       string
		directives []string
	}{
		{"SingleDirective", []string{header.NegotiateDirectiveVList.String()}},
		{"MultipleDirectives", []string{header.NegotiateDirectiveTrans.String(), header.NegotiateDirectiveVList.String()}},
		{"Empty", []string{}},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			// action + assert.
			n, err := header.NewNegotiate(test.directives)
			s.Require().NoError(err)
			for idx, d := range n.Directives() {
				s.Equal(test.directives[idx], d.String())
			}
		})
	}
}

func (s NegotiateTestSuite) TestNegotiate_Contains() {

	tests := []struct {
		name       string
		directives []string
		in         []header.NegotiateDirective
		out        bool
	}{
		{
			"MatchSingle",
			[]string{
				header.NegotiateDirectiveVList.String(),
			},
			[]header.NegotiateDirective{
				header.NegotiateDirectiveVList,
			},
			true,
		},
		{
			"MatchMultiple",
			[]string{
				header.NegotiateDirectiveVList.String(),
			},
			[]header.NegotiateDirective{
				header.NegotiateDirectiveTrans,
				header.NegotiateDirectiveVList,
			},
			true,
		},
		{
			"NoMatch",
			[]string{
				header.NegotiateDirectiveTrans.String(),
				header.NegotiateDirectiveVList.String(),
			},
			[]header.NegotiateDirective{
				header.NegotiateDirectiveGuessSmall,
			},
			false,
		},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			// action + assert.
			n, err := header.NewNegotiate(test.directives)
			s.Require().NoError(err)
			s.Equal(test.out, n.Contains(test.in...))
		})
	}
}

func (s NegotiateTestSuite) TestNegotiate_ContainsRVSA() {

	tests := []struct {
		name       string
		directives []string
		in         string
		out        bool
	}{
		{
			"MatchSingle",
			[]string{
				header.NegotiateDirective("1.0").String(),
			},
			header.NegotiateDirective("1.0").String(),
			true,
		},
		{
			"NoMatch",
			[]string{
				header.NegotiateDirective("1.0").String(),
			},
			header.NegotiateDirective("3.0").String(),
			false,
		},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			// action + assert.
			n, err := header.NewNegotiate(test.directives)
			s.Require().NoError(err)
			s.Equal(test.out, n.ContainsRVSA(test.in))
		})
	}
}

func (s NegotiateTestSuite) TestNegotiate_String() {

	tests := []struct {
		name       string
		directives []string
		out        string
	}{
		{
			"NotEmpty",
			[]string{},
			"Negotiate: ",
		},
		{
			"Empty",
			[]string{
				header.NegotiateDirective("1.0").String(),
			},
			"Negotiate: 1.0",
		},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			// action + assert.
			n, err := header.NewNegotiate(test.directives)
			s.Require().NoError(err)
			s.Equal(test.out, n.String())
		})
	}
}
