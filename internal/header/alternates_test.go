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
	"net/url"
	"testing"

	"github.com/freerware/negotiator/internal/header"
	_representation "github.com/freerware/negotiator/internal/representation"
	"github.com/freerware/negotiator/internal/test"
	"github.com/freerware/negotiator/representation"
	"github.com/stretchr/testify/suite"
)

type AlternatesTestSuite struct {
	suite.Suite
}

func TestAlternatesTestSuite(t *testing.T) {
	suite.Run(t, new(AlternatesTestSuite))
}

func (s AlternatesTestSuite) TestAlternates_NewAlternates() {

	json, html, english, ascii, gzip :=
		"application/json", "text/html", "en-US", "ascii", "gzip"
	loc, _ := url.Parse("http://www.example.com/thing")
	v1 := _representation.NewBuilder().
		WithLocation(*loc).
		WithType(json).
		WithLanguage(english).
		WithEncoding(gzip).
		WithCharset(ascii).
		WithSourceQuality(1.0).
		Build(test.RepresentationBuilderFunc)
	v2 := _representation.NewBuilder().
		WithLocation(*loc).
		WithType(html).
		WithLanguage(english).
		WithEncoding(gzip).
		WithCharset(ascii).
		WithSourceQuality(1.0).
		Build(test.RepresentationBuilderFunc)

	tests := []struct {
		name     string
		fallback representation.Representation
		reps     []representation.Representation
		err      error
	}{
		{"WithFallback", v1, []representation.Representation{v2}, nil},
		{"WithoutFallback", nil, []representation.Representation{v2}, nil},
		{"Empty", nil, []representation.Representation{}, nil},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			// action.
			_, err := header.NewAlternates(test.fallback, test.reps...)

			// assert.
			if test.err != nil {
				s.Require().EqualError(err, test.err.Error())
			} else {
				s.Require().NoError(err)
			}
		})
	}
}

func (s AlternatesTestSuite) TestAlternates_HasFallback() {

	json, english, ascii, gzip :=
		"application/json", "en-US", "ascii", "gzip"
	loc, _ := url.Parse("http://www.example.com/thing")
	v1 := _representation.NewBuilder().
		WithLocation(*loc).
		WithType(json).
		WithLanguage(english).
		WithEncoding(gzip).
		WithCharset(ascii).
		WithSourceQuality(1.0).
		Build(test.RepresentationBuilderFunc)

	tests := []struct {
		name     string
		fallback representation.Representation
		out      bool
	}{
		{"HasFallback", v1, true},
		{"NoFallback", nil, false},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			// action + assert.
			a, err := header.NewAlternates(test.fallback)
			s.Require().NoError(err)
			s.Equal(test.out, a.HasFallback())
		})
	}
}

func (s AlternatesTestSuite) TestAlternates_ValuesAsStrings() {

	json, html, english, ascii, gzip :=
		"application/json", "text/html", "en-US", "ascii", "gzip"
	loc, _ := url.Parse("http://www.example.com/thing")
	v1 := _representation.NewBuilder().
		WithLocation(*loc).
		WithType(json).
		WithLanguage(english).
		WithEncoding(gzip).
		WithCharset(ascii).
		WithSourceQuality(1.0).
		Build(test.RepresentationBuilderFunc)
	v2 := _representation.NewBuilder().
		WithLocation(*loc).
		WithType(html).
		WithLanguage(english).
		WithEncoding(gzip).
		WithCharset(ascii).
		WithSourceQuality(1.0).
		Build(test.RepresentationBuilderFunc)

	tests := []struct {
		name     string
		fallback representation.Representation
		reps     []representation.Representation
		out      []string
	}{
		{
			"WithFallback",
			v1,
			[]representation.Representation{v2},
			[]string{"{ \"http://www.example.com/thing\" 1.000 { charset ascii } { features  } { language en-US } { length 59 } { type text/html } }", "{ \"http://www.example.com/thing\" }"},
		},
		{
			"WithoutFallback",
			nil,
			[]representation.Representation{v2},
			[]string{"{ \"http://www.example.com/thing\" 1.000 { charset ascii } { features  } { language en-US } { length 59 } { type text/html } }"},
		},
		{
			"Empty",
			nil,
			[]representation.Representation{},
			[]string{},
		},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			// action + assert.
			a, err := header.NewAlternates(test.fallback, test.reps...)
			s.Require().NoError(err)
			s.Equal(test.out, a.ValuesAsStrings())
		})
	}
}

func (s AlternatesTestSuite) TestAlternates_ValuesAsString() {

	json, html, english, ascii, gzip :=
		"application/json", "text/html", "en-US", "ascii", "gzip"
	loc, _ := url.Parse("http://www.example.com/thing")
	v1 := _representation.NewBuilder().
		WithLocation(*loc).
		WithType(json).
		WithLanguage(english).
		WithEncoding(gzip).
		WithCharset(ascii).
		WithSourceQuality(1.0).
		Build(test.RepresentationBuilderFunc)
	v2 := _representation.NewBuilder().
		WithLocation(*loc).
		WithType(html).
		WithLanguage(english).
		WithEncoding(gzip).
		WithCharset(ascii).
		WithSourceQuality(1.0).
		Build(test.RepresentationBuilderFunc)

	tests := []struct {
		name     string
		fallback representation.Representation
		reps     []representation.Representation
		out      string
	}{
		{
			"WithFallback",
			v1,
			[]representation.Representation{v2},
			"{ \"http://www.example.com/thing\" 1.000 { charset ascii } { features  } { language en-US } { length 59 } { type text/html } },{ \"http://www.example.com/thing\" }",
		},
		{
			"WithoutFallback",
			nil,
			[]representation.Representation{v2},
			"{ \"http://www.example.com/thing\" 1.000 { charset ascii } { features  } { language en-US } { length 59 } { type text/html } }",
		},
		{
			"Empty",
			nil,
			[]representation.Representation{},
			"",
		},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			// action + assert.
			a, err := header.NewAlternates(test.fallback, test.reps...)
			s.Require().NoError(err)
			s.Equal(test.out, a.ValuesAsString())
		})
	}
}

func (s AlternatesTestSuite) TestAlternates_String() {

	json, html, english, ascii, gzip :=
		"application/json", "text/html", "en-US", "ascii", "gzip"
	loc, _ := url.Parse("http://www.example.com/thing")
	v1 := _representation.NewBuilder().
		WithLocation(*loc).
		WithType(json).
		WithLanguage(english).
		WithEncoding(gzip).
		WithCharset(ascii).
		WithSourceQuality(1.0).
		Build(test.RepresentationBuilderFunc)
	v2 := _representation.NewBuilder().
		WithLocation(*loc).
		WithType(html).
		WithLanguage(english).
		WithEncoding(gzip).
		WithCharset(ascii).
		WithSourceQuality(1.0).
		Build(test.RepresentationBuilderFunc)

	tests := []struct {
		name     string
		fallback representation.Representation
		reps     []representation.Representation
		out      string
	}{
		{
			"WithFallback",
			v1,
			[]representation.Representation{v2},
			"Alternates: { \"http://www.example.com/thing\" 1.000 { charset ascii } { features  } { language en-US } { length 59 } { type text/html } },{ \"http://www.example.com/thing\" }",
		},
		{
			"WithoutFallback",
			nil,
			[]representation.Representation{v2},
			"Alternates: { \"http://www.example.com/thing\" 1.000 { charset ascii } { features  } { language en-US } { length 59 } { type text/html } }",
		},
		{
			"Empty",
			nil,
			[]representation.Representation{},
			"Alternates: ",
		},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			// action + assert.
			a, err := header.NewAlternates(test.fallback, test.reps...)
			s.Require().NoError(err)
			s.Equal(test.out, a.String())
		})
	}
}
