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

var (
	// headerAcceptFeatures is the header key for the Accept-Features header.
	headerAcceptFeatures = "Accept-Features"

	// DefaultAcceptFeatures is an empty Accept-Features header.
	DefaultAcceptFeatures = EmptyAcceptFeatures

	// EmptyAcceptFeatures is an empty Accept-Features header.
	EmptyAcceptFeatures = AcceptFeatures([]FeatureExpression{})
)

// AcceptFeatures represents the Accept-Features header.
//
// The Accept-Features request header can be used by a user agent to
// give information about the presence or absence of certain features in
// the feature set of the current request.  Servers can use this
// information when running a remote variant selection algorithm.
type AcceptFeatures []FeatureExpression

// NewAcceptFeatures constructs an Accept-Features header with the provided
// feature expressions.
func NewAcceptFeatures(acceptFeatures []string) (AcceptFeatures, error) {
	if len(acceptFeatures) == 0 {
		return EmptyAcceptFeatures, nil
	}
	var expressions []FeatureExpression
	for _, e := range acceptFeatures {
		expressions = append(expressions, FeatureExpression(e))
	}
	return AcceptFeatures(expressions), nil
}

// Empty indicates if the Accept-Header is empty.
func (f AcceptFeatures) IsEmpty() bool {
	return len(f) == len(EmptyAcceptFeatures)
}

// AsFeatureSets utilizes the feature expressions within the Accept-Feature
// header to construct a partial view of the user agent's supported and
// unsupported feature sets.
func (f AcceptFeatures) AsFeatureSets() (supported, unsupported FeatureSet) {
	supported, unsupported = make(FeatureSet), make(FeatureSet)
	for _, e := range f {
		switch e.Type() {
		case FeatureExpressionTypeExists:
			supported.Add(e.Tag())
		case FeatureExpressionTypeNotExists:
			unsupported.Add(e.Tag())
		case FeatureExpressionTypeEquals,
			FeatureExpressionTypeExclusiveEquals:
			supported.Add(e.Tag(), e.Value())
		case FeatureExpressionTypeNotEquals:
			supported.Add(e.Tag())
			unsupported.Add(e.Tag(), e.Value())
		}
	}
	return
}

// String provides a textual representation of the Accept-Features header.
func (f AcceptFeatures) String() string {
	var expressions []string
	for _, e := range f {
		expressions = append(expressions, e.String())
	}
	return fmt.Sprintf("%s: %s", headerAcceptFeatures, strings.Join(expressions, ","))
}
