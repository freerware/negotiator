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
	"strconv"
	"strings"
)

// FeaturePredicateBag represents a collection of feature predicates.
type FeaturePredicateBag []FeaturePredicate

// Evaluate determines if the predicate bag matches with the provided feature sets.
func (fpb FeaturePredicateBag) Evaluate(supported, unsupported FeatureSet) bool {
	for _, fp := range fpb {
		if fp.Evaluate(supported, unsupported) {
			return true
		}
	}
	return false
}

// String provides the textual representation of the feature predicate bag.
func (fpb FeaturePredicateBag) String() string {
	var s []string
	for _, p := range fpb {
		s = append(s, p.String())
	}
	return fmt.Sprintf("[ %s ]", strings.Join(s, " "))
}

// FeaturePredicate represents a predicate used to express support for a
// particular feature.
type FeaturePredicate interface {
	evaluator
	fmt.Stringer
}

// exists represents a feature predicate that ensures a feature is present.
type exists struct {
	t FeatureTag
}

// Evaluate determines if the predicate matches with the provided feature sets.
func (e exists) Evaluate(supported, unsupported FeatureSet) bool {
	return supported.Contains(e.t)
}

// String provides the textual representation of the predicate.
func (e exists) String() string {
	return e.t.String()
}

// absent represents a feature predicate that ensures a feature is not supported.
type absent struct {
	t FeatureTag
}

// Evaluate determines if the predicate matches with the provided feature sets.
func (a absent) Evaluate(supported, unsupported FeatureSet) bool {
	values, ok := unsupported.Values(a.t)
	return !supported.Contains(a.t) && ok && len(values) == 0
}

// Strings provides the textual representation of the predicate.
func (a absent) String() string {
	return fmt.Sprintf("!%s", a.t.String())
}

// equals represents a feature predicate that ensures a feature is present
// with a particular value.
type equals struct {
	t FeatureTag
	v FeatureTagValue
}

// Evaluate determines if the predicate matches with the provided feature sets.
func (e equals) Evaluate(supported, unsupported FeatureSet) bool {
	// ensure the feature is supported and has a matching value.
	if sValues, ok := supported.Values(e.t); ok {
		for _, val := range sValues {
			if e.v.Equals(val) {
				return true
			}
		}
	}
	return false
}

// String provides the textual representation of the predicate.
func (e equals) String() string {
	return fmt.Sprintf("%s=%s", e.t.String(), e.v.String())
}

// equals represents a feature predicate that ensures a feature is present
// without a particular value.
type notEquals struct {
	t FeatureTag
	v FeatureTagValue
}

// Evaluate determines if the predicate matches with the provided feature sets.
func (ne notEquals) Evaluate(supported, unsupported FeatureSet) bool {
	// ensure the feature is supported.
	sValues, s := supported.Values(ne.t)
	if !s {
		return false
	}

	// ensure the unsupported value is not supported.
	for _, val := range sValues {
		if ne.v.Equals(val) {
			return false
		}
	}
	return true
}

// String provides the textual representation of the predicate.
func (ne notEquals) String() string {
	return fmt.Sprintf("%s!=%s", ne.t.String(), ne.v.String())
}

// within represents a feature predicate that ensures a feature is present
// with a values within a specified range.
type within struct {
	t FeatureTag
	l FeatureTagValue
	h FeatureTagValue
}

// Evaluate determines if the predicate matches with the provided feature sets.
func (w within) Evaluate(supported, unsupported FeatureSet) bool {
	// ensure the feature tag is supported.
	sValues, s := supported.Values(w.t)
	if !s {
		return false
	}

	// ensure that at least one supported value falls within the range.
	for _, val := range sValues {
		if val.IsNumeric() && w.l.IsNumeric() && w.h.IsNumeric() {
			num, _ := strconv.ParseFloat(val.String(), 64)
			low, _ := strconv.ParseFloat(w.l.String(), 64)
			high, _ := strconv.ParseFloat(w.h.String(), 64)
			if num > low && num < high {
				return true
			}
		}
	}
	return false
}

// String provides the textual representation of the predicate.
func (w within) String() string {
	return fmt.Sprintf("%s=[%s-%s]", w.t.String(), w.l.String(), w.h.String())
}
