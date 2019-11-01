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

package header

import (
	"fmt"
	"strings"
)

// FeatureSet represents a collection of feature tags and their values.
type FeatureSet struct {
	features map[FeatureTag][]FeatureTagValue
}

// Add introduces the provides feature tag with the provided values.
func (s FeatureSet) Add(tag FeatureTag, values ...FeatureTagValue) {
	if !s.Contains(tag) {
		s.features[tag] = values
	}
	for t, values := range s.features {
		if t.Equals(tag) {
			s.features[t] = append(s.features[t], values...)
		}
	}
}

// Contains determines if the feature set contains the provided feature tag.
func (s FeatureSet) Contains(tag FeatureTag) bool {
	for t := range s.features {
		if t.Equals(tag) {
			return true
		}
	}
	return false
}

// Values retrieves the values for the provided feature tag.
func (s FeatureSet) Values(tag FeatureTag) (v []FeatureTagValue, ok bool) {
	for t, values := range s.features {
		if t.Equals(tag) {
			ok = true
			v = append(v, values...)
		}
	}
	return
}

// String provides the textual representation of the feature set.
func (s FeatureSet) String() string {
	var r []string
	for tag, values := range s.features {
		var v []string
		for _, value := range values {
			v = append(v, value.String())
		}
		r = append(r, fmt.Sprintf("( %s , { %s } )", tag.String(), strings.Join(v, ", ")))
	}
	return fmt.Sprintf("{ %s }", r)
}
