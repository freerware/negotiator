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

type FeatureListElementTestSuite struct {
	suite.Suite
}

func TestFeatureListElementTestSuite(t *testing.T) {
	suite.Run(t, new(FeatureListElementTestSuite))
}

func (s FeatureListElementTestSuite) TestFeatureListElement_NewPredicateListElement() {

	tests := []struct {
		name string
		in   string
		err  error
	}{
		{"WithoutImprovementOrDegradation", "foo", nil},
		{"WithImprovement", "foo;+0.4", nil},
		{"WithDegradation", "foo;-0.4", nil},
		{"WithImprovementAndDegradation", "foo;+0.4-0.4", nil},
		{"InvalidPredicate", "", header.ErrInvalidPredicate},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			// action.
			_, err := header.NewPredicateListElement(test.in)

			// assert.
			if test.err != nil {
				s.Require().EqualError(err, test.err.Error())
			} else {
				s.Require().NoError(err)
			}
		})
	}
}

func (s FeatureListElementTestSuite) TestPredicateListElement_TrueImprovement() {

	tests := []struct {
		name string
		in   string
		out  header.TrueImprovement
	}{
		{"WithoutImprovementOrDegradation", "foo", header.TrueImprovement(1.0)},
		{"WithImprovement", "foo;+0.4", header.TrueImprovement(0.4)},
		{"WithDegradation", "foo;-0.4", header.TrueImprovement(1.0)},
		{"WithImprovementAndDegradation", "foo;+0.4-0.4", header.TrueImprovement(0.4)},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			// action.
			le, err := header.NewPredicateListElement(test.in)
			s.Require().NoError(err)
			s.Equal(test.out, le.TrueImprovement())
		})
	}
}

func (s FeatureListElementTestSuite) TestPredicateListElement_FalseDegradation() {

	tests := []struct {
		name string
		in   string
		out  header.FalseDegradation
	}{
		{"WithoutImprovementOrDegradation", "foo", header.FalseDegradation(0.0)},
		{"WithImprovement", "foo;+0.4", header.FalseDegradation(1.0)},
		{"WithDegradation", "foo;-0.4", header.FalseDegradation(0.4)},
		{"WithImprovementAndDegradation", "foo;+0.4-0.4", header.FalseDegradation(0.4)},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			// action + assert.
			le, err := header.NewPredicateListElement(test.in)
			s.Require().NoError(err)
			s.Equal(test.out, le.FalseDegradation())
		})
	}
}

func (s FeatureListElementTestSuite) TestPredicateListElement_String() {

	tests := []struct {
		name      string
		predicate string
		out       string
	}{
		{
			"WithoutImprovementOrDegredation",
			"foo",
			"foo;+1.000-0.000",
		},
		{
			"WithImprovement",
			"foo;+0.5",
			"foo;+0.500-1.000",
		},
		{
			"WithDegradation",
			"foo;-0.5",
			"foo;+1.000-0.500",
		},
		{
			"WithImprovementAndDegradation",
			"foo;+0.5-0.5",
			"foo;+0.500-0.500",
		},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			// action + assert.
			fl, err := header.NewPredicateListElement(test.predicate)
			s.Require().NoError(err)
			s.Equal(test.out, fl.String())
		})
	}
}

func (s FeatureListElementTestSuite) TestFeatureListElement_NewPredicateBagListElement() {

	tests := []struct {
		name string
		in   string
		err  error
	}{
		{"WithoutImprovementOrDegradation", "[ !bar baz=biz ]", nil},
		{"WithImprovement", "[ !bar baz=biz ];+0.4", nil},
		{"WithDegradation", "[ !bar baz=biz ];-0.4", nil},
		{"WithImprovementAndDegradation", "[ !bar baz=biz ];+0.4-0.4", nil},
		{"InvalidPredicateBag", "", header.ErrInvalidPredicate},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			// action.
			_, err := header.NewPredicateBagListElement(test.in)

			// assert.
			if test.err != nil {
				s.Require().EqualError(err, test.err.Error())
			} else {
				s.Require().NoError(err)
			}
		})
	}
}

func (s FeatureListElementTestSuite) TestPredicateBagListElement_TrueImprovement() {

	tests := []struct {
		name string
		in   string
		out  header.TrueImprovement
	}{
		{"WithoutImprovementOrDegradation", "[ !bar baz=biz ]", header.TrueImprovement(1.0)},
		{"WithImprovement", "[ !bar baz=biz ];+0.4", header.TrueImprovement(0.4)},
		{"WithDegradation", "[ !bar baz=biz ];-0.4", header.TrueImprovement(1.0)},
		{"WithImprovementAndDegradation", "[ !bar baz=biz ];+0.4-0.4", header.TrueImprovement(0.4)},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			// action.
			le, err := header.NewPredicateBagListElement(test.in)
			s.Require().NoError(err)
			s.Equal(test.out, le.TrueImprovement())
		})
	}
}

func (s FeatureListElementTestSuite) TestPredicateBagListElement_FalseDegradation() {

	tests := []struct {
		name string
		in   string
		out  header.FalseDegradation
	}{
		{"WithoutImprovementOrDegradation", "[ !bar baz=biz ]", header.FalseDegradation(0.0)},
		{"WithImprovement", "[ !bar baz=biz ];+0.4", header.FalseDegradation(1.0)},
		{"WithDegradation", "[ !bar baz=biz ];-0.4", header.FalseDegradation(0.4)},
		{"WithImprovementAndDegradation", "[ !bar baz=biz ];+0.4-0.4", header.FalseDegradation(0.4)},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			// action + assert.
			le, err := header.NewPredicateBagListElement(test.in)
			s.Require().NoError(err)
			s.Equal(test.out, le.FalseDegradation())
		})
	}
}

func (s FeatureListElementTestSuite) TestPredicateBagListElement_String() {

	tests := []struct {
		name         string
		predicateBag string
		out          string
	}{
		{
			"WithoutImprovementOrDegredation",
			"[ !bar baz=biz ]",
			"[ !bar baz=biz ];+1.000-0.000",
		},
		{
			"WithImprovement",
			"[ !bar baz=biz ];+0.5",
			"[ !bar baz=biz ];+0.500-1.000",
		},
		{
			"WithDegradation",
			"[ !bar baz=biz ];-0.5",
			"[ !bar baz=biz ];+1.000-0.500",
		},
		{
			"WithImprovementAndDegradation",
			"[ !bar baz=biz ];+0.5-0.5",
			"[ !bar baz=biz ];+0.500-0.500",
		},
	}

	for _, test := range tests {
		s.Run(test.name, func() {
			// action + assert.
			fl, err := header.NewPredicateBagListElement(test.predicateBag)
			s.Require().NoError(err)
			s.Equal(test.out, fl.String())
		})
	}
}
