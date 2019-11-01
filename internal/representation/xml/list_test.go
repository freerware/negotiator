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

package xml_test

import (
	"net/http/httptest"
	"testing"

	_representation "github.com/freerware/negotiator/internal/representation"
	"github.com/freerware/negotiator/internal/representation/xml"
	"github.com/freerware/negotiator/internal/test"
	"github.com/freerware/negotiator/representation"
	"github.com/stretchr/testify/suite"
)

type XMLListTestSuite struct {
	suite.Suite
}

func TestXMLListTestSuite(t *testing.T) {
	suite.Run(t, new(XMLListTestSuite))
}

func (s *XMLListTestSuite) TestXMLList_List() {
	// arrange.
	_xml, english, ascii, gzip := "application/xml", "en-US", "ascii", "gzip"
	request := httptest.NewRequest("GET", "http://freer.ddns.net/thing", nil)
	v := _representation.NewBuilder().
		WithLocation(*request.URL).
		WithType(_xml).
		WithLanguage(english).
		WithEncoding(gzip).
		WithCharset(ascii).
		WithSourceQuality(1.0).
		Build(test.RepresentationBuilderFunc)
	variants := []representation.Representation{v}

	// action.
	list := xml.List(variants...)

	// assert.
	s.Equal("application/xml", list.ContentType())
}
