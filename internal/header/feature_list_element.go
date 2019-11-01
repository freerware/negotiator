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
	"regexp"
	"strconv"
	"strings"
)

var (
	// regexpPredicateBag is the regular expression responsible for parsing
	// feature predicate bags.
	regexpPredicateBag = regexp.MustCompile(regexp.QuoteMeta(`^\[\s{0,1}(.*)\s{0,1}\](;(\+(\d{0,1}\.{0,1}\d{1,3}))?(\-(\d{0,1}\.{0,1}\d{1,3}))?)?$`))
	// regexpPredicate is the regular expression responsible for parsing
	// feature predicate.
	regexpPredicate = regexp.MustCompile(regexp.QuoteMeta(`^(\!)?(\w+)((\=|\!\=)((\{(\d)?\-(\d)?\})|([0-9a-zA-Z\.]+)))?(;(\+(\d{0,1}\.{0,1}\d{1,3}))?(\-(\d{0,1}\.{0,1}\d{1,3}))?)?$`))
)

// FeatureListElement represents a single element within a feature list. An
// element can either be a feature predicate or feature predicate bag.
type FeatureListElement interface {
	evaluator
	fmt.Stringer

	TrueImprovement() TrueImprovement
	FalseDegradation() FalseDegradation
}

// featureListElement represents the core aspects of a feature list element.
type featureListElement struct {
	t *TrueImprovement
	f *FalseDegradation
}

// TrueImprovement indicates the true value for the feature list element.
func (fle featureListElement) TrueImprovement() TrueImprovement {
	if fle.t != nil {
		return *fle.t
	}
	return TrueImprovement(1.0)
}

// FalseDegradation indicates the false degradation for the feature list element.
func (fle featureListElement) FalseDegradation() FalseDegradation {
	if fle.f != nil {
		return *fle.f
	}
	if fle.t != nil {
		return FalseDegradation(1.0)
	}
	return FalseDegradation(0.0)
}

// predicateListElement represents a predict feature list element.
type predicateListElement struct {
	featureListElement

	predicate FeaturePredicate
}

// newPredicate constructs a new feature predicate.
func newPredicate(predicate string) FeaturePredicate {
	// TODO: change regex to be specifically for predicate part.
	matches := regexpPredicate.FindStringSubmatch(predicate)
	if len(matches) == 0 {
		// error.
	}

	var p FeaturePredicate
	if len(matches[1]) > 0 {
		p = absent{
			t: FeatureTag(matches[2]),
		}
	}

	switch matches[4] {
	case "=":
		if len(matches[6]) > 0 {
			low, high := "0", "0"
			if len(matches[7]) > 0 {
				low = matches[7]
			}
			if len(matches[8]) > 0 {
				high = matches[8]
			}
			p = within{
				t: FeatureTag(matches[2]),
				l: FeatureTagValue(low),
				h: FeatureTagValue(high),
			}
		} else {
			p = equals{
				t: FeatureTag(matches[2]),
				v: FeatureTagValue(matches[5]),
			}
		}
	case "!=":
		p = notEquals{
			t: FeatureTag(matches[2]),
			v: FeatureTagValue(matches[9]),
		}
	}
	return p
}

// newPredicateListElement constructs a new feature predicate list element.
func newPredicateListElement(predicate string) predicateListElement {
	// TDOD: change regex to be higher level.
	matches := regexpPredicate.FindStringSubmatch(predicate)
	p := newPredicate(predicate)

	// parse true improvement.
	var ti *TrueImprovement
	if len(matches[11]) > 0 {
		// safe to ignore error due to regexp.
		t, _ := strconv.ParseFloat(matches[12], 32)
		tt := TrueImprovement(t)
		ti = &tt
	}

	// parse false degradation.
	var fd *FalseDegradation
	if len(matches[13]) > 0 {
		// safe to ignore error due to regexp.
		f, _ := strconv.ParseFloat(matches[14], 32)
		ff := FalseDegradation(f)
		fd = &ff
	}
	return predicateListElement{
		featureListElement: featureListElement{
			t: ti,
			f: fd,
		},
		predicate: p,
	}
}

// Evaluate determines if the feature predicate list element matches based on
// the provided feature sets.
func (ple predicateListElement) Evaluate(supported, unsupported FeatureSet) bool {
	return ple.predicate.Evaluate(supported, unsupported)
}

// String provides the textual representation of the feature list element.
func (ple predicateListElement) String() string {
	return fmt.Sprintf("%s;%s%s",
		ple.predicate.String(),
		ple.TrueImprovement().String(),
		ple.FalseDegradation().String(),
	)
}

// predicateBagListElement represents a feature predicate bag list element.
type predicateBagListElement struct {
	featureListElement

	predicateBag FeaturePredicateBag
}

// newPredicateBagListElement constructs a new feature predicate bag list element.
func newPredicateBagListElement(bag string) predicateBagListElement {
	matches := regexpPredicateBag.FindStringSubmatch(bag)
	if len(matches) == 0 {
		// error.
	}

	// parse predicates.
	var predicates []FeaturePredicate
	for _, p := range strings.Split(strings.TrimSpace(matches[1]), " ") {
		predicate := newPredicate(p)
		predicates = append(predicates, predicate)
	}

	// parse true improvement.
	var ti *TrueImprovement
	if len(matches[4]) != 0 {
		// safe to ignore error due to regexp.
		t, _ := strconv.ParseFloat(matches[4], 32)
		tt := TrueImprovement(t)
		ti = &tt
	}

	// parse false degradation.
	var fd *FalseDegradation
	if len(matches[6]) != 0 {
		// safe to ignore error due to regexp.
		f, _ := strconv.ParseFloat(matches[6], 32)
		ff := FalseDegradation(f)
		fd = &ff
	}
	return predicateBagListElement{
		featureListElement: featureListElement{
			t: ti,
			f: fd,
		},
		predicateBag: FeaturePredicateBag(predicates),
	}
}

// Evaluate determines if the feature predicate bag list element matches based on
// the provided feature sets.
func (pble predicateBagListElement) Evaluate(supported, unsupported FeatureSet) bool {
	return pble.predicateBag.Evaluate(supported, unsupported)
}

// String provides the textual representation of the feature list element.
func (pble predicateBagListElement) String() string {
	return fmt.Sprintf("%s;%s%s",
		pble.predicateBag.String(),
		pble.TrueImprovement().String(),
		pble.FalseDegradation().String(),
	)
}

// TrueImprovement represents the degradation factor yielded when a feature
// list element is determined to be true.
type TrueImprovement float32

// String provides the textual representation of the true-improvement factor.
func (ti TrueImprovement) String() string {
	return fmt.Sprintf("+%.3f", ti)
}

// Float provides the floating point representation of the true-improvement factor.
func (ti TrueImprovement) Float() float32 {
	return float32(ti)
}

// FalseDegradation represents the degradation factor yielded when a feature
// list element is determined to be false.
type FalseDegradation float32

// String provides the textual representation of the false-degradation factor.
func (fd FalseDegradation) String() string {
	return fmt.Sprintf("-%.3f", fd)
}

// Float provides the floating point representation of the false-degredation factor.
func (fd FalseDegradation) Float() float32 {
	return float32(fd)
}
