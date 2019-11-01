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

package representation

import (
	"sort"
)

var (
	// EmptySet is an empty representation set.
	EmptySet = Set([]RankedRepresentation{})
)

// Set represents a collection of representations.
type Set []RankedRepresentation

// Where filters the representation set using the provided predicate.
func (s Set) Where(p Predicate) Set {
	var matched []RankedRepresentation
	for _, v := range s {
		if p(v) {
			matched = append(matched, v)
		}
	}
	return Set(matched)
}

// AsSlice converts the representation set into a slice.
func (s Set) AsSlice() []RankedRepresentation {
	return s
}

// Sort sorts the representation set based on the provided less function.
func (s Set) Sort(less func(i, j int) bool) {
	sort.Slice(s, less)
}

// First task the first element of the representation set. Must check if
// the set is empty prior to invoking.
func (s Set) First() RankedRepresentation {
	return s[0]
}

// Size provides the size of the representation set.
func (s Set) Size() int {
	return len(s)
}

// Empty indicates if the representation set is empty.
func (s Set) Empty() bool {
	return s.Size() == 0
}

// Predicate represents a matching function leveraged when
// filtering the representation set.
type Predicate func(RankedRepresentation) bool
