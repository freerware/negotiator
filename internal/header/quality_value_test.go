package header_test

import (
	"testing"

	"github.com/freerware/negotiator/internal/header"
	"github.com/stretchr/testify/suite"
)

type QualityValueTestSuite struct {
	suite.Suite
}

func TestQualityValueTestSuite(t *testing.T) {
	suite.Run(t, new(QualityValueTestSuite))
}

func (s QualityValueTestSuite) TestQualityValue_NewQualityValue() {

	tests := []struct {
		name string
		in   float32
		err  error
	}{
		{"WithinRange", float32(0.5), nil},
		{"Negative", float32(-1.0), header.ErrInvalidQualityValue},
		{"TooBig", float32(2.0), header.ErrInvalidQualityValue},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			// action.
			qv, err := header.NewQualityValue(test.in)

			// assert.
			if test.err != nil {
				s.Require().EqualError(err, test.err.Error())
			} else {
				s.Require().NoError(err)
				s.Equal(test.in, qv.Float())
			}
		})
	}
}

func (s QualityValueTestSuite) TestQualityValue_Equals() {

	tests := []struct {
		name string
		qv   float32
		in   header.QualityValue
		out  bool
	}{
		{"Match", 0.5, header.QualityValue(0.5), true},
		{"NoMatch", 1.0, header.QualityValueMinimum, false},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			// action + assert.
			qv, err := header.NewQualityValue(test.qv)
			s.Require().NoError(err)
			s.Equal(test.out, qv.Equals(test.in))
		})
	}
}

func (s QualityValueTestSuite) TestQualityValue_LessThan() {

	tests := []struct {
		name string
		qv   float32
		in   header.QualityValue
		out  bool
	}{
		{"Less", 0.4, header.QualityValue(0.5), true},
		{"NotLess", 1.0, header.QualityValueMinimum, false},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			// action + assert.
			qv, err := header.NewQualityValue(test.qv)
			s.Require().NoError(err)
			s.Equal(test.out, qv.LessThan(test.in))
		})
	}
}

func (s QualityValueTestSuite) TestQualityValue_GreaterThan() {

	tests := []struct {
		name string
		qv   float32
		in   header.QualityValue
		out  bool
	}{
		{"Less", 0.4, header.QualityValue(0.5), false},
		{"NotLess", 1.0, header.QualityValueMinimum, true},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			// action + assert.
			qv, err := header.NewQualityValue(test.qv)
			s.Require().NoError(err)
			s.Equal(test.out, qv.GreaterThan(test.in))
		})
	}
}

func (s QualityValueTestSuite) TestQualityValue_Multiply() {

	tests := []struct {
		name string
		qv   float32
		in   header.QualityValue
		out  header.QualityValue
	}{
		{"", 0.5, header.QualityValue(0.5), header.QualityValue(0.25)},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			// action + assert.
			qv, err := header.NewQualityValue(test.qv)
			s.Require().NoError(err)
			s.Equal(test.out, qv.Multiply(test.in))
		})
	}
}

func (s QualityValueTestSuite) TestQualityValue_Round() {

	tests := []struct {
		name string
		qv   float32
		in   int
		out  header.QualityValue
	}{
		{"Tenths", 0.123, 1, header.QualityValue(0.1)},
		{"Hundredths", 0.987, 2, header.QualityValue(0.99)},
		{"Thousandths", 0.555, 3, header.QualityValue(0.555)},
		{"Tens", 0.555, 0, header.QualityValueMaximum},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			// action + assert.
			qv, err := header.NewQualityValue(test.qv)
			s.Require().NoError(err)
			s.Equal(test.out, qv.Round(test.in))
		})
	}
}

func (s QualityValueTestSuite) TestQualityValue_String() {

	tests := []struct {
		name string
		qv   float32
		out  string
	}{
		{
			"Tenths",
			0.5,
			"0.500",
		},
		{
			"Hundredths",
			0.55,
			"0.550",
		},
		{
			"Thousandths",
			0.555,
			"0.555",
		},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			// action + assert.
			qv, err := header.NewQualityValue(test.qv)
			s.Require().NoError(err)
			s.Equal(test.out, qv.String())
		})
	}
}
