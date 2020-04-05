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

import "strings"

// FeatureList represents the collection of feature predicates and feature
// predicate bags that describe the quality degradation for a particular
// representation.
type FeatureList []FeatureListElement

// NewFeatureList constructs a feature list.
func NewFeatureList(features []string) (FeatureList, error) {
	var list []FeatureListElement
	for _, feature := range features {
		var le FeatureListElement
		var err error
		if strings.HasPrefix(feature, "[") && strings.HasSuffix(feature, "]") {
			if le, err = NewPredicateBagListElement(feature); err != nil {
				return list, err
			}
		} else {
			if le, err = NewPredicateListElement(feature); err != nil {
				return list, err
			}
		}
		list = append(list, le)
	}
	return FeatureList(list), nil
}

// QualityDegradation computes the overall quality degradation factor for the
// feature list based on the provided feature sets.
func (fl FeatureList) QualityDegradation(supported, unsupported FeatureSet) float32 {
	degradation := float32(1.0)
	for _, element := range fl {
		if element.Evaluate(supported, unsupported) {
			degradation = degradation * element.TrueImprovement().Float()
		} else {
			degradation = degradation * element.FalseDegradation().Float()
		}
	}
	return degradation
}

// String provides the textual representation of the feature list.
func (fl FeatureList) String() string {
	var s []string
	for _, e := range fl {
		s = append(s, e.String())
	}
	return strings.Join(s, " ")
}
