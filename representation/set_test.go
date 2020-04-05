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

package representation_test

import (
	"testing"

	"github.com/freerware/negotiator/representation"
	"github.com/stretchr/testify/suite"
)

type RepresentationSetTestSuite struct {
	suite.Suite

	sut representation.Set
}

func TestRepresentationSet(t *testing.T) {
	suite.Run(t, new(RepresentationSetTestSuite))
}

func (s *RepresentationSetTestSuite) SetupTest() {
	var set representation.Set
	set = append(set, representation.RankedRepresentation{
		SourceQualityValue:    representation.SourceQualityAcceptable,
		MediaTypeQualityValue: 0.8,
		CharsetQualityValue:   0.8,
		LanguageQualityValue:  0.8,
		EncodingQualityValue:  0.8,
		FeatureQualityValue:   0.8,
	})
	set = append(set, representation.RankedRepresentation{
		SourceQualityValue:    representation.SourceQualityNearlyPerfect,
		MediaTypeQualityValue: 0.9,
		CharsetQualityValue:   0.9,
		LanguageQualityValue:  0.9,
		EncodingQualityValue:  0.9,
		FeatureQualityValue:   0.9,
	})
	s.sut = set
}

func (s *RepresentationSetTestSuite) TestSet_Where_WithMatches() {

	// action.
	matches := s.sut.Where(func(r representation.RankedRepresentation) bool {
		return r.SourceQualityValue < 0.9
	})

	// assert.
	s.Require().Len(matches, 1)
	s.Equal(s.sut[0], matches[0])
}

func (s *RepresentationSetTestSuite) TestSet_Where_WithoutMatches() {

	// action.
	matches := s.sut.Where(func(r representation.RankedRepresentation) bool {
		return r.SourceQualityValue < 0.8
	})

	// assert.
	s.Require().Len(matches, 0)
}

func (s *RepresentationSetTestSuite) TestSet_AsSlice() {

	// action.
	slice := s.sut.AsSlice()

	// assert.
	s.Require().Len(slice, len(s.sut))
	s.ElementsMatch(s.sut, slice)
}

func (s *RepresentationSetTestSuite) TestSet_Sort() {

	// arrange.
	first, second := s.sut[0], s.sut[1]

	// action.
	s.sut.Sort(func(i, j int) bool {
		return s.sut[i].SourceQualityValue > s.sut[j].SourceQualityValue
	})

	// assert.
	s.Equal(first, s.sut[1])
	s.Equal(second, s.sut[0])
}

func (s *RepresentationSetTestSuite) TestSet_First() {

	// action + assert.
	s.Equal(s.sut[0], s.sut.First())
}

func (s *RepresentationSetTestSuite) TestSet_First_NoElements() {

	// arrange.
	s.sut = representation.EmptySet

	// action + assert.
	s.Panics(func() {
		s.sut.First()
	})
}

func (s *RepresentationSetTestSuite) TestSet_Size() {

	// action + assert.
	s.Equal(len(s.sut), s.sut.Size())
}

func (s *RepresentationSetTestSuite) TestSet_Empty_IsEmpty() {

	// action + assert.
	s.False(s.sut.Empty())
}

func (s *RepresentationSetTestSuite) TestSet_Empty_IsNotEmpty() {

	// arrange.
	s.sut = representation.EmptySet

	// action + assert.
	s.True(s.sut.Empty())
}

func (s *RepresentationSetTestSuite) TearDownTest() {
	s.sut = nil
}
