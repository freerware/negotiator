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

type FeatureListTestSuite struct {
	suite.Suite
}

func TestFeatureListTestSuite(t *testing.T) {
	suite.Run(t, new(FeatureListTestSuite))
}

func (s FeatureListTestSuite) TestFeatureList_NewFeatureList() {

	tests := []struct {
		name string
		in   []string
		err  error
	}{
		{"Empty", []string{}, nil},
		{"WithoutPredicateBag", []string{"foo", "!bar", "baz=biz", "zip=[0-3]"}, nil},
		{"WithPredicateBag", []string{"foo", "[!bar baz=biz]", "zip=[0-3]"}, nil},
		{"InvalidPredicate", []string{""}, header.ErrInvalidPredicate},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			// action.
			_, err := header.NewFeatureList(test.in)

			// assert.
			if test.err != nil {
				s.Require().EqualError(err, test.err.Error())
			} else {
				s.Require().NoError(err)
			}
		})
	}
}

func (s FeatureListTestSuite) TestFeatureList_QualityDegredation() {

	supported, unsupported :=
		header.FeatureSet(
			map[header.FeatureTag][]header.FeatureTagValue{
				header.FeatureTag("\"foo\""): []header.FeatureTagValue{},
				header.FeatureTag("yish"):    []header.FeatureTagValue{header.FeatureTagValue("\"vam\"")},
				header.FeatureTag("baz"):     []header.FeatureTagValue{header.FeatureTagValue("1"), header.FeatureTagValue("5")},
			}),
		header.FeatureSet(
			map[header.FeatureTag][]header.FeatureTagValue{
				header.FeatureTag("\"bar\""): []header.FeatureTagValue{},
				header.FeatureTag("meep"):    []header.FeatureTagValue{header.FeatureTagValue("\"enabled\"")},
			})

	tests := []struct {
		name     string
		features []string
		out      float32
	}{
		{"AllTrue", []string{"yish=vam;+1.5", "foo"}, 1.5},
		{"AllFalse", []string{"bar;-0.5", "meep;+1.0"}, 0.5},
		{"BothTrueAndFalse", []string{"foo;+0.5", "bar;-0.5"}, 0.25},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			// action + assert.
			fl, err := header.NewFeatureList(test.features)
			s.Require().NoError(err)
			s.Equal(test.out, fl.QualityDegradation(supported, unsupported))
		})
	}
}

func (s FeatureListTestSuite) TestFeatureList_String() {

	tests := []struct {
		name     string
		features []string
		out      string
	}{
		{"Empty", []string{}, ""},
		{
			"WithoutPredicateBag",
			[]string{"foo", "!bar", "baz=biz", "zip=[0-3]"},
			"foo;+1.000-0.000 !bar;+1.000-0.000 baz=biz;+1.000-0.000 zip=[0-3];+1.000-0.000",
		},
		{
			"WithPredicateBag",
			[]string{"foo;+0.5", "[!bar baz=biz]", "zip=[0-3]"},
			"foo;+0.500-1.000 [ !bar baz=biz ];+1.000-0.000 zip=[0-3];+1.000-0.000",
		},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			// action + assert.
			fl, err := header.NewFeatureList(test.features)
			s.Require().NoError(err)
			s.Equal(test.out, fl.String())

		})
	}
}
