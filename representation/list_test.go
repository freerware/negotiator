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

	_representation "github.com/freerware/negotiator/internal/representation"
	"github.com/freerware/negotiator/internal/test"
	"github.com/freerware/negotiator/representation"
	"github.com/stretchr/testify/suite"
)

type ListTestSuite struct {
	suite.Suite
}

func TestListTestSuite(t *testing.T) {
	suite.Run(t, new(ListTestSuite))
}

func (s *ListTestSuite) TestList_Bytes() {
	// arrange.
	json, english, ascii, gzip := "application/json", "en-US", "ascii", "gzip"
	v := _representation.NewBuilder().
		WithType(json).
		WithLanguage(english).
		WithEncoding(gzip).
		WithCharset(ascii).
		WithSourceQuality(1.0).
		Build(test.RepresentationBuilderFunc)
	variants := []representation.Representation{v}

	l := representation.List{}
	l.SetContentCharset(ascii)
	l.SetContentEncoding([]string{gzip})
	l.SetContentLanguage(english)
	l.SetContentType(json)
	l.SetRepresentations(variants...)

	// action + assert.
	b, err := l.Bytes()
	s.Require().NoError(err)
	s.Len(b, 139)
}

func (s *ListTestSuite) TestList_FromBytes() {
	// arrange.
	json, english, ascii, gzip := "application/json", "en-US", "ascii", "gzip"
	v := _representation.NewBuilder().
		WithType(json).
		WithLanguage(english).
		WithEncoding(gzip).
		WithCharset(ascii).
		WithSourceQuality(1.0).
		Build(test.RepresentationBuilderFunc)
	variants := []representation.Representation{v}

	l := representation.List{}
	l.SetContentCharset(ascii)
	l.SetContentEncoding([]string{gzip})
	l.SetContentLanguage(english)
	l.SetContentType(json)
	l.SetRepresentations(variants...)

	b, err := l.Bytes()
	s.Require().NoError(err)

	ll := representation.List{}
	ll.SetContentCharset(ascii)
	ll.SetContentEncoding([]string{gzip})
	ll.SetContentLanguage(english)
	ll.SetContentType(json)

	// action + assert.
	err = ll.FromBytes(b)
	s.Require().NoError(err)
	s.Len(ll.Representations, 1)
}
