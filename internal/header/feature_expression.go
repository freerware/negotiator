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

// FeatureExpressionType indicates the operation the feature expression is
// performing.pr
type FeatureExpressionType int

// String provides the textual representation of the feature expression type.
func (t FeatureExpressionType) String() string {
	return []string{
		"EXISTS",
		"NOT_EXISTS",
		"EQUALS",
		"EXCLUSIVE_EQUALS",
		"NOT_EQUALS",
		"WILDCARD",
	}[t]
}

const (

	// FeatureExpressionTypeExists indicates the feature expression is testing
	// for existance of a particular feature.
	FeatureExpressionTypeExists = iota

	// FeatureExpressionTypeNotExists indicates the feature expression is
	// testing for the absence of a particular feature.
	FeatureExpressionTypeNotExists

	// FeatureExpressionTypeEquals indicates the feature expression is
	// testing for a feature with a particular value.
	FeatureExpressionTypeEquals

	// FeatureExpressionTypeExclusiveEquals indicates the feature expression is
	// testing for a feature with a particular value, and only that value.
	FeatureExpressionTypeExclusiveEquals

	// FeatureExpressionTypeNotEquals indicates the feature expression is
	// testing for a feature that doesn't have a particular value.
	FeatureExpressionTypeNotEquals

	// FeatureExpressionTypeWildcard indicates the feature expression that
	// communicates additional features are available, even though they
	// they weren't mentioned.
	FeatureExpressionTypeWildcard
)

// FeatureExpression represents a feature expression communicated in the
// Accept-Features header.
type FeatureExpression string

// IsWildcard indicates if the feature expression is '*'.
func (e FeatureExpression) IsWildcard() bool {
	return e.String() == "*"
}

// Tag retrieves the feature tag for the feature expression, when applicable.
// Wilcard expressions do not contain a feature tag.
func (e FeatureExpression) Tag() FeatureTag {
	return e.parseTag()
}

// Value retrieves the feature tag value for the feature expression, when
// applicable.
func (e FeatureExpression) Value() FeatureTagValue {
	val, _ := e.parseValue()
	return val
}

// HasValue determines if the feature expression has a feature tag value.
// Existance, absence, and wildcard feature expressions do not have tag values.
func (e FeatureExpression) HasValue() bool {
	_, ok := e.parseValue()
	return ok
}

// IsExclusive indicates if the feature expression uses exclusive equality.
func (e FeatureExpression) IsExclusive() bool {
	return e.HasValue() &&
		strings.Contains(e.String(), "{") &&
		strings.Contains(e.String(), "}")
}

// Type indicates the feature expression type.
func (e FeatureExpression) Type() FeatureExpressionType {
	// check for negation ( ! ).
	if strings.HasPrefix(e.String(), "!") {
		return FeatureExpressionTypeNotExists
	}
	// check for inequality ( != ).
	if strings.Contains(e.String(), "!=") {
		return FeatureExpressionTypeNotEquals
	}
	// check for equality ( = ).
	if strings.Contains(e.String(), "=") {
		// check for exclusive equality ( {...} ).
		if e.IsExclusive() {
			return FeatureExpressionTypeExclusiveEquals
		}
		return FeatureExpressionTypeEquals
	}
	if e.IsWildcard() {
		return FeatureExpressionTypeWildcard
	}
	return FeatureExpressionTypeExists
}

// parseTag takes the raw feature expression and parses the tag.
// wildcard expressions do not have tags.
func (e FeatureExpression) parseTag() FeatureTag {
	if strings.HasPrefix(e.String(), "!") {
		return FeatureTag(strings.TrimLeft(e.String(), "!"))
	}
	if s := strings.Split(e.String(), "!="); len(s) > 1 {
		return FeatureTag(s[0])
	}
	if s := strings.Split(e.String(), "="); len(s) > 1 {
		return FeatureTag(s[0])
	}
	return FeatureTag(e.String())
}

// parseValue takes the raw feature expression and parses the tag value.
// the existance, absence, and wildcard expressions do not have tag values.
func (e FeatureExpression) parseValue() (val FeatureTagValue, ok bool) {
	// check for negation ( ! ).
	if strings.HasPrefix(e.String(), "!") {
		return
	}
	// check for inequality ( != ).
	if s := strings.Split(e.String(), "!="); len(s) > 1 {
		val, ok = FeatureTagValue(s[1]), true
		return
	}
	// check for equality ( = ).
	if s := strings.Split(e.String(), "="); len(s) > 1 {
		// check for exclusive equality ( {...} ).
		if strings.HasPrefix(s[1], "{") &&
			strings.HasSuffix(s[1], "}") {

			val, ok = FeatureTagValue(
				strings.ReplaceAll(
					strings.ReplaceAll(s[1], "}", ""), "{", "",
				),
			), true
			return
		}
	}
	return
}

// String provides a textual representation of the feature expression.
func (e FeatureExpression) String() string {
	return string(e)
}
