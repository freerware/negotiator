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
	"sort"
	"strings"
)

var (
	EmptyFeatureSet = FeatureSet(map[FeatureTag][]FeatureTagValue{})
)

// FeatureSet represents a collection of feature tags and their values.
type FeatureSet map[FeatureTag][]FeatureTagValue

// Add introduces the provides feature tag with the provided values.
func (s FeatureSet) Add(tag FeatureTag, values ...FeatureTagValue) {
	if !s.Contains(tag) {
		s[tag] = values
		return
	}
	for t, values := range s {
		if t.Equals(tag) {
			s[t] = append(s[t], values...)
		}
	}
}

// Contains determines if the feature set contains the provided feature tag.
func (s FeatureSet) Contains(tag FeatureTag) (c bool) {
	for t := range s {
		if c = t.Equals(tag); c {
			return
		}
	}
	return
}

// Values retrieves the values for the provided feature tag.
func (s FeatureSet) Values(tag FeatureTag) (v []FeatureTagValue, ok bool) {
	for t, values := range s {
		if t.Equals(tag) {
			ok = true
			v = append(v, values...)
		}
	}
	return
}

// String provides the textual representation of the feature set.
func (s FeatureSet) String() string {
	// sort the keys for deterministic output.
	var keys []string
	for tag := range s {
		keys = append(keys, tag.String())
	}
	sort.Strings(keys)

	var r []string
	for _, key := range keys {
		values := s[FeatureTag(key)]
		var v []string
		for _, value := range values {
			v = append(v, value.String())
		}
		r = append(r, fmt.Sprintf("( %s , { %s } )", key, strings.Join(v, ", ")))
	}
	rs := strings.Join(r, " ")
	return fmt.Sprintf("{ %s }", rs)
}
