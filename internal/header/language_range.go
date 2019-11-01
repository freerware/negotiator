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
	"strconv"
	"strings"

	"golang.org/x/text/language"
)

var (
	//defaultLanguageRange is the default language range.
	defaultLanguageRange = LanguageRange{
		lrange: "*",
		tag:    language.Und,
		qValue: QualityValueMaximum,
	}

	// ErrEmptyLanguageRange is an error that indicates that the langauge
	// range cannot be empty.
	ErrEmptyLanguageRange = errors.New("language range cannot be empty")
)

// LanguageRange represents a language tag matching expression.
type LanguageRange struct {
	lrange string
	tag    language.Tag
	qValue QualityValue
}

// NewLanguageRange constructs a language range from the textual representation.
func NewLanguageRange(languageRange string) (LanguageRange, error) {
	if len(languageRange) == 0 {
		return LanguageRange{}, ErrEmptyLanguageRange
	}

	parts := strings.Split(languageRange, ";")
	r := strings.ToLower(parts[0])
	t, err := language.Parse(r)
	if r != "*" && err != nil {
		return LanguageRange{}, err
	}
	lr := LanguageRange{
		lrange: r,
		tag:    t,
		qValue: QualityValue(1.0),
	}

	if len(parts) > 1 && strings.HasPrefix(strings.Trim(parts[1], " "), "q=") {
		w := strings.Trim(parts[1], " ")
		q := strings.TrimPrefix(w, "q=")
		f, err := strconv.ParseFloat(q, 32)
		if err != nil {
			return LanguageRange{}, err
		}
		qv, err := NewQualityValue(float32(f))
		if err != nil {
			return LanguageRange{}, err
		}
		lr.qValue = qv
	}
	return lr, nil
}

// IsWildcard indicates if the language range is '*'.
func (lr LanguageRange) IsWildcard() bool {
	return lr.lrange == "*"
}

// IsTag indicates if the language range specifies a language tag.
func (lr LanguageRange) IsTag() bool {
	return !lr.IsWildcard()
}

// Tag retrieves the language tag this language range specifies.
func (lr LanguageRange) Tag() string {
	return lr.tag.String()
}

// QualityValue retrieves the quality value of the language range.
func (lr LanguageRange) QualityValue() QualityValue {
	return lr.qValue
}

// Compatible determines if the provided language tag is compatible with
// the language range.
func (lr LanguageRange) Compatible(tag string) bool {
	_, err := language.Parse(tag)
	if err != nil {
		return false
	}
	if lr.IsWildcard() {
		return true
	}
	m := language.NewMatcher([]language.Tag{language.Und, lr.tag})
	//_, i, _ := m.Match(t)
	_, i := language.MatchStrings(m, tag)
	return i != 0
}

// String provides a textual representation of the language range.
func (lr LanguageRange) String() string {
	return fmt.Sprintf("%s;q=%s", lr.lrange, lr.QualityValue().String())
}
