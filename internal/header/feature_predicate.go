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
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Regular expressions to match against the feature predicate grammar.
var (
	existsPredicate    = regexp.MustCompile(`^(\w+)$`)
	absentPredicate    = regexp.MustCompile(`^(!(\w+))$`)
	equalsPredicate    = regexp.MustCompile(`^(\w+)\s?=\s?(\w+)$`)
	notEqualsPredicate = regexp.MustCompile(`^(\w+)\s?!=\s?(\w+)$`)
	rangePredicate     = regexp.MustCompile(`^(\w+)\s?=\s?\[\s?(\d+)?\s?-\s?(\d+)?\s?\]$`)
)

// Errors that can be encountered when interacting with feature predicates.
var (

	// ErrInvalidPredicate represents an error that occurs when the
	// feature predicate is invalid.
	ErrInvalidPredicate = errors.New("invalid feature predicate")

	// ErrInvalidPredicateBag represents an error that occurs when the
	// feature predicate bag is invalid.
	ErrInvalidPredicateBag = errors.New("invalid feature predicate bag")
)

// FeaturePredicateBag represents a collection of feature predicates.
type FeaturePredicateBag []FeaturePredicate

// NewFeaturePredicateBag constructs a feature predicate bag.
func NewFeaturePredicateBag(predicateBag string) (FeaturePredicateBag, error) {
	openCount, closeCount :=
		strings.Count(predicateBag, "["),
		strings.Count(predicateBag, "]")

	if openCount != closeCount {
		return nil, ErrInvalidPredicateBag
	}
	predicateBag =
		strings.TrimSpace(
			strings.TrimPrefix(strings.TrimSuffix(predicateBag, "]"), "["))
	predicates := strings.Split(predicateBag, " ")

	var fpb FeaturePredicateBag
	for _, predicate := range predicates {
		fp, err := NewFeaturePredicate(predicate)
		if err != nil {
			return nil, err
		}
		fpb = append(fpb, fp)
	}

	return fpb, nil
}

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

// newPredicate constructs a new feature predicate.
func NewFeaturePredicate(predicate string) (FeaturePredicate, error) {

	parseExists := func(p string) (fp FeaturePredicate, ok bool) {
		if ok = existsPredicate.MatchString(p); !ok {
			return
		}
		groups := existsPredicate.FindStringSubmatch(p)
		fp = exists{t: FeatureTag(groups[1])}
		return
	}

	parseAbsent := func(p string) (fp FeaturePredicate, ok bool) {
		if ok = absentPredicate.MatchString(p); !ok {
			return
		}
		groups := absentPredicate.FindStringSubmatch(p)
		fp = absent{t: FeatureTag(groups[2])}
		return
	}

	parseEquals := func(p string) (fp FeaturePredicate, ok bool) {
		if ok = equalsPredicate.MatchString(p); !ok {
			return
		}
		groups := equalsPredicate.FindStringSubmatch(p)
		fp = equals{t: FeatureTag(groups[1]), v: FeatureTagValue(groups[2])}
		return
	}

	parseNotEquals := func(p string) (fp FeaturePredicate, ok bool) {
		if ok = notEqualsPredicate.MatchString(p); !ok {
			return
		}
		groups := notEqualsPredicate.FindStringSubmatch(p)
		fp = notEquals{t: FeatureTag(groups[1]), v: FeatureTagValue(groups[2])}
		return
	}

	parseRange := func(p string) (fp FeaturePredicate, ok bool) {
		if ok = rangePredicate.MatchString(p); !ok {
			return
		}
		groups := rangePredicate.FindStringSubmatch(p)
		t, l, h := groups[1], groups[2], groups[3]
		fp = within{
			t: FeatureTag(t),
			l: FeatureTagValue(l),
			h: FeatureTagValue(h),
		}
		return
	}

	parsers := []func(p string) (FeaturePredicate, bool){
		parseExists,
		parseAbsent,
		parseEquals,
		parseNotEquals,
		parseRange,
	}

	// parse.
	for _, parser := range parsers {
		if fp, ok := parser(predicate); ok {
			return fp, nil
		}
	}
	return nil, ErrInvalidPredicate
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

	low, _ := strconv.Atoi(w.l.String())
	high, _ := strconv.Atoi(w.h.String())

	// ensure that the highest numeric value for the feature in the set falls
	// within the range inclusively.
	var highest int
	for _, val := range sValues {
		if val.IsNumeric() && w.l.IsNumeric() && w.h.IsNumeric() {
			num, _ := strconv.Atoi(val.String())
			if num > highest {
				highest = num
			}
		}
	}
	return highest >= low && highest <= high
}

// String provides the textual representation of the predicate.
func (w within) String() string {
	return fmt.Sprintf("%s=[%s-%s]", w.t.String(), w.l.String(), w.h.String())
}
