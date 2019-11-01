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

type FeaturePredicateBagTestSuite struct {
	suite.Suite
}

func TestFeaturePredicateBagTestSuite(t *testing.T) {
	suite.Run(t, new(FeaturePredicateBagTestSuite))
}

func (s FeaturePredicateBagTestSuite) TestFeaturePredicateBag_NewFeaturePredicateBag() {
	tests := []struct {
		name string
		in   string
		err  error
	}{
		{"", "[foo !bar meep!=enabled yish=vam baz=[1-5]]", nil},
		{
			"Spaces",
			"[ foo !bar meep!=enabled yish=vam baz=[1-5] ]",
			nil,
		},
		{
			"MissingClosingBracket",
			"[foo !bar meep!=enabled yish=vam baz=[1-5]",
			header.ErrInvalidPredicateBag,
		},
		{
			"MissingOpeningBracket",
			"foo !bar meep!=enabled yish=vam baz=[1-5]]",
			header.ErrInvalidPredicateBag,
		},
		{
			"InvalidPredicate",
			"[ foo bar& meep!=enabled yish=vam baz=[1-5] ]",
			header.ErrInvalidPredicate,
		},
	}
	for _, test := range tests {
		s.Run(test.name, func() {
			_, err := header.NewFeaturePredicateBag(test.in)
			if test.err != nil {
				s.EqualError(err, test.err.Error())
			} else {
				s.NoError(err)
			}
		})
	}
}

func (s FeaturePredicateBagTestSuite) TestFeaturePredicateBag_String() {
	// arrange.
	b := "[ foo !bar meep!=enabled yish=vam baz=[1-5] ]"

	// action + assert.
	bag, err := header.NewFeaturePredicateBag(b)
	s.Require().NoError(err)
	s.Equal(b, bag.String())
}

func (s FeaturePredicateBagTestSuite) TestFeaturePredicateBag_Evaluate() {
	supported, unsupported :=
		header.FeatureSet(
			map[header.FeatureTag][]header.FeatureTagValue{
				header.FeatureTag("\"foo\""): []header.FeatureTagValue{},
				header.FeatureTag("yish"):    []header.FeatureTagValue{header.FeatureTagValue("\"vam\"")},
				header.FeatureTag("baz"):     []header.FeatureTagValue{header.FeatureTagValue("1"), header.FeatureTagValue("5")},
			},
		),
		header.FeatureSet(
			map[header.FeatureTag][]header.FeatureTagValue{
				header.FeatureTag("\"bar\""): []header.FeatureTagValue{},
				header.FeatureTag("meep"):    []header.FeatureTagValue{header.FeatureTagValue("\"enabled\"")},
			},
		)
	tests := []struct {
		name string
		bag  string
		out  bool
	}{
		{
			"Match",
			"[ foo !bar meep!=enabled yish=vam baz=[1-5] ]",
			true,
		},
		{
			"NoMatch",
			"[ !foo bar meep=enabled yish!=vam baz=[6-8] ]",
			false,
		},
	}
	for _, test := range tests {
		s.Run(test.name, func() {
			pb, err := header.NewFeaturePredicateBag(test.bag)
			s.Require().NoError(err)
			s.Equal(test.out, pb.Evaluate(supported, unsupported))
		})
	}
}

type FeaturePredicateTestSuite struct {
	suite.Suite
}

func TestFeaturePredicateTestSuite(t *testing.T) {
	suite.Run(t, new(FeaturePredicateTestSuite))
}

func (s FeaturePredicateTestSuite) TestFeaturePredicate_NewFeaturePredicate() {
	tests := []struct {
		name string
		in   string
		err  error
	}{
		{"Exists", "foo", nil},
		{"Absent", "!bar", nil},
		{"NotEquals", "meep!=enabled", nil},
		{"Equals", "yish=vam", nil},
		{"Within", "baz=[1-5]", nil},
		{"InvalidPredicate", "bar&", header.ErrInvalidPredicate},
	}
	for _, test := range tests {
		s.Run(test.name, func() {
			_, err := header.NewFeaturePredicate(test.in)
			if test.err != nil {
				s.EqualError(err, test.err.Error())
			} else {
				s.NoError(err)
			}
		})
	}
}

func (s FeaturePredicateTestSuite) TestFeaturePredicate_String() {
	tests := []struct {
		name string
		in   string
		out  string
	}{
		{"Exists", "foo", "foo"},
		{"Absent", "!bar", "!bar"},
		{"NotEquals", "meep!=enabled", "meep!=enabled"},
		{"Equals", "yish=vam", "yish=vam"},
		{"Within", "baz=[1-5]", "baz=[1-5]"},
	}
	for _, test := range tests {
		s.Run(test.name, func() {
			fp, err := header.NewFeaturePredicate(test.in)
			s.Require().NoError(err)
			s.Equal(test.out, fp.String())
		})
	}
}

func (s FeaturePredicateTestSuite) TestFeaturePredicate_Evaluate() {
	supported, unsupported :=
		header.FeatureSet(
			map[header.FeatureTag][]header.FeatureTagValue{
				header.FeatureTag("\"foo\""): []header.FeatureTagValue{},
				header.FeatureTag("yish"):    []header.FeatureTagValue{header.FeatureTagValue("\"vam\"")},
				header.FeatureTag("baz"):     []header.FeatureTagValue{header.FeatureTagValue("1"), header.FeatureTagValue("5")},
				header.FeatureTag("meep"):    []header.FeatureTagValue{},
			},
		),
		header.FeatureSet(
			map[header.FeatureTag][]header.FeatureTagValue{
				header.FeatureTag("\"bar\""): []header.FeatureTagValue{},
				header.FeatureTag("meep"):    []header.FeatureTagValue{header.FeatureTagValue("\"enabled\"")},
			},
		)
	tests := []struct {
		name      string
		predicate string
		out       bool
	}{
		{
			"Exists_Match",
			"foo",
			true,
		},
		{
			"Exists_NoMatch",
			"zap",
			false,
		},
		{
			"Absent_Match",
			"!bar",
			true,
		},
		{
			"Absent_NoMatch",
			"!zap",
			false,
		},
		{
			"Equals_Match",
			"yish=vam",
			true,
		},
		{
			"Equals_NoMatch",
			"zap=boop",
			false,
		},
		{
			"NotEquals_Match",
			"meep!=enabled",
			true,
		},
		{
			"NotEquals_NoMatch",
			"zap!=boop",
			false,
		},
		{
			"Within_Match",
			"baz=[2-5]",
			true,
		},
		{
			"Within_NoMatch",
			"baz=[6-7]",
			false,
		},
	}
	for _, test := range tests {
		s.Run(test.name, func() {
			pb, err := header.NewFeaturePredicate(test.predicate)
			s.Require().NoError(err)
			s.Equal(test.out, pb.Evaluate(supported, unsupported))
		})
	}
}
