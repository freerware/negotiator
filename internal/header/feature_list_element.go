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

// Regular expressions to match against the feature list grammar.
var (
	improvementDegration = regexp.MustCompile(`^(\+(\d+(\.\d{1,4})?))?(\-(\d+(\.\d{1,4})?))?$`)
)

var (
	ErrInvalidPredicateListElement = errors.New("invalid predicate list element")
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

// NewPredicateListElement constructs a new feature predicate list element.
func NewPredicateListElement(element string) (FeatureListElement, error) {
	elementParts := strings.Split(element, ";")
	predicate, err := NewFeaturePredicate(elementParts[0])
	if err != nil {
		return nil, err
	}
	if len(elementParts) == 1 {
		return predicateListElement{
			predicate: predicate,
		}, nil
	}

	if ok := improvementDegration.MatchString(elementParts[1]); !ok {
		return nil, ErrInvalidPredicateListElement
	}
	groups := improvementDegration.FindStringSubmatch(elementParts[1])

	// parse true improvement.
	var ti *TrueImprovement
	if len(groups[2]) > 0 {
		// safe to ignore error due to regexp.
		t, _ := strconv.ParseFloat(groups[2], 32)
		tt := TrueImprovement(float32(t))
		ti = &tt
	}

	// parse false degradation.
	var fd *FalseDegradation
	if len(groups[5]) > 0 {
		// safe to ignore error due to regexp.
		f, _ := strconv.ParseFloat(groups[5], 32)
		ff := FalseDegradation(float32(f))
		fd = &ff
	}
	return predicateListElement{
		featureListElement: featureListElement{
			t: ti,
			f: fd,
		},
		predicate: predicate,
	}, nil
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

// NewPredicateBagListElement constructs a new feature predicate bag list element.
func NewPredicateBagListElement(element string) (FeatureListElement, error) {
	elementParts := strings.Split(element, ";")
	predicateBag, err := NewFeaturePredicateBag(elementParts[0])
	if err != nil {
		return nil, err
	}
	if len(elementParts) == 1 {
		return predicateBagListElement{
			predicateBag: predicateBag,
		}, nil
	}

	if ok := improvementDegration.MatchString(elementParts[1]); !ok {
		return nil, ErrInvalidPredicateListElement
	}
	groups := improvementDegration.FindStringSubmatch(elementParts[1])

	// parse true improvement.
	var ti *TrueImprovement
	if len(groups[2]) > 0 {
		// safe to ignore error due to regexp.
		t, _ := strconv.ParseFloat(groups[2], 32)
		tt := TrueImprovement(float32(t))
		ti = &tt
	}

	// parse false degradation.
	var fd *FalseDegradation
	if len(groups[5]) > 0 {
		// safe to ignore error due to regexp.
		f, _ := strconv.ParseFloat(groups[5], 32)
		ff := FalseDegradation(float32(f))
		fd = &ff
	}
	return predicateBagListElement{
		featureListElement: featureListElement{
			t: ti,
			f: fd,
		},
		predicateBag: predicateBag,
	}, nil
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
