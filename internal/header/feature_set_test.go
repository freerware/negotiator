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

type FeatureSetTestSuite struct {
	suite.Suite
}

func TestFeatureSetTestSuite(t *testing.T) {
	suite.Run(t, new(FeatureSetTestSuite))
}

func (s FeatureSetTestSuite) TestFeatureSet_Add() {

	tests := []struct {
		name   string
		tag    header.FeatureTag
		values []header.FeatureTagValue
	}{
		{"NoValue", header.FeatureTag("yish"), []header.FeatureTagValue{}},
		{"Existing", header.FeatureTag("foo"), []header.FeatureTagValue{
			header.FeatureTagValue("bar"),
		}},
		{
			"SingleValue",
			header.FeatureTag("yish"),
			[]header.FeatureTagValue{header.FeatureTagValue("\"vam\"")},
		},
		{
			"MultipleValues",
			header.FeatureTag("yish"),
			[]header.FeatureTagValue{
				header.FeatureTagValue("\"vam\""),
				header.FeatureTagValue("\"voom\""),
			},
		},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			// action + assert.
			fs := header.FeatureSet(map[header.FeatureTag][]header.FeatureTagValue{
				header.FeatureTag("foo"): []header.FeatureTagValue{},
			})
			fs.Add(test.tag, test.values...)
		})
	}
}

func (s FeatureSetTestSuite) TestFeatureSet_Contains() {

	tests := []struct {
		name       string
		featureSet header.FeatureSet
		in         header.FeatureTag
		out        bool
	}{
		{"NoMatch", header.EmptyFeatureSet, header.FeatureTag("foo"), false},
		{
			"Match",
			header.FeatureSet(
				map[header.FeatureTag][]header.FeatureTagValue{
					header.FeatureTag("yish"): []header.FeatureTagValue{
						header.FeatureTagValue("\"vam\""),
					},
				}),
			header.FeatureTag("yish"),
			true,
		},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			// action + assert.
			s.Equal(test.out, test.featureSet.Contains(test.in))
		})
	}
}

func (s FeatureSetTestSuite) TestFeatureSet_Values() {

	fs := header.FeatureSet(
		map[header.FeatureTag][]header.FeatureTagValue{
			header.FeatureTag("\"foo\""): []header.FeatureTagValue{},
			header.FeatureTag("yish"):    []header.FeatureTagValue{header.FeatureTagValue("\"vam\"")},
			header.FeatureTag("baz"):     []header.FeatureTagValue{header.FeatureTagValue("1"), header.FeatureTagValue("5")}})

	tests := []struct {
		name       string
		featureSet header.FeatureSet
		in         header.FeatureTag
		outOK      bool
		out        []header.FeatureTagValue
	}{
		{"NoMatch", header.EmptyFeatureSet, header.FeatureTag("foo"), false, []header.FeatureTagValue{}},
		{
			"MatchEmpty",
			fs,
			header.FeatureTag("foo"),
			true,
			[]header.FeatureTagValue{},
		},
		{
			"MatchNotEmpty",
			fs,
			header.FeatureTag("baz"),
			true,
			[]header.FeatureTagValue{
				header.FeatureTagValue("1"), header.FeatureTagValue("5"),
			},
		},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			// action + assert.
			vals, ok := test.featureSet.Values(test.in)
			s.Equal(test.outOK, ok)
			s.ElementsMatch(test.out, vals)
		})
	}
}

func (s FeatureSetTestSuite) TestFeatureSet_String() {

	fs := header.FeatureSet(
		map[header.FeatureTag][]header.FeatureTagValue{
			header.FeatureTag("\"foo\""): []header.FeatureTagValue{},
			header.FeatureTag("yish"):    []header.FeatureTagValue{header.FeatureTagValue("\"vam\"")},
			header.FeatureTag("baz"):     []header.FeatureTagValue{header.FeatureTagValue("1"), header.FeatureTagValue("5")}})

	tests := []struct {
		name       string
		featureSet header.FeatureSet
		out        string
	}{
		{"Empty", header.EmptyFeatureSet, "{  }"},
		{
			"NotEmpty",
			fs,
			"{ ( \"foo\" , {  } ) ( baz , { 1, 5 } ) ( yish , { \"vam\" } ) }",
		},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			// action + assert.
			s.Equal(test.out, test.featureSet.String())
		})
	}
}
