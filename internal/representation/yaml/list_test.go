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

package yaml_test

import (
	"net/http/httptest"
	"testing"

	_representation "github.com/freerware/negotiator/internal/representation"
	"github.com/freerware/negotiator/internal/representation/yaml"
	"github.com/freerware/negotiator/internal/test"
	"github.com/freerware/negotiator/representation"
	"github.com/stretchr/testify/suite"
)

type YAMLListTestSuite struct {
	suite.Suite
}

func TestYAMLListTestSuite(t *testing.T) {
	suite.Run(t, new(YAMLListTestSuite))
}

func (s *YAMLListTestSuite) TestYAMLList_List() {
	// arrange.
	_yaml, english, ascii, gzip := "application/yaml", "en-US", "ascii", "gzip"
	request := httptest.NewRequest("GET", "http://freer.ddns.net/thing", nil)
	v := _representation.NewBuilder().
		WithLocation(*request.URL).
		WithType(_yaml).
		WithLanguage(english).
		WithEncoding(gzip).
		WithCharset(ascii).
		WithSourceQuality(1.0).
		Build(test.RepresentationBuilderFunc)
	variants := []representation.Representation{v}

	// action.
	list := yaml.List(variants...)

	// assert.
	s.Equal("application/yaml", list.ContentType())
}
