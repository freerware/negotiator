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

	"golang.org/x/text/language"
)

var (
	// headerAcceptLanguage is the header key for the Accept-Language header.
	headerAcceptLanguage = "Accept-Language"

	// DefaultAcceptLanguage is an Accept-Language header with a
	// single language range of "*".
	DefaultAcceptLanguage = AcceptLanguage([]LanguageRange{defaultLanguageRange})

	// EmptyAcceptLanguage is an empty Accept-Language header.
	EmptyAcceptLanguage = AcceptLanguage([]LanguageRange{})
)

// AcceptLanguage represents the Accept-Language header.
//
// The "Accept-Language" header field can be used by user agents to
// indicate the set of natural languages that are preferred in the
// response.  Language tags are defined in Section 3.1.3.1.
type AcceptLanguage []LanguageRange

// NewAcceptLanguage constructs an Accept-Language header with the provided
// language ranges.
func NewAcceptLanguage(acceptLanguage []string) (AcceptLanguage, error) {
	if len(acceptLanguage) == 0 {
		return EmptyAcceptLanguage, nil
	}
	tags, qValues, err :=
		language.ParseAcceptLanguage(strings.Join(acceptLanguage, ","))
	if err != nil {
		return EmptyAcceptLanguage, err
	}

	var ranges []LanguageRange
	for i := 0; i < len(tags); i++ {
		qv, err := NewQualityValue(qValues[i])
		if err != nil {
			return EmptyAcceptLanguage, err
		}
		ranges = append(ranges, LanguageRange{
			lrange: tags[i].String() + ";q=" + qv.String(),
			tag:    tags[i],
			qValue: qv,
		})
	}
	return AcceptLanguage(ranges), nil
}

// IsEmpty indicates if the Accept-Language header is empty.
func (l AcceptLanguage) IsEmpty() bool {
	return len(l) == len(EmptyAcceptLanguage)
}

// LanguageRanges provides the language ranges sorted on preference from//// highest to lowest.
func (l AcceptLanguage) LanguageRanges() []LanguageRange {
	sort.Slice(l, func(first, second int) bool {
		f := l[first]
		s := l[second]
		return f.QualityValue().GreaterThan(s.QualityValue())
	})
	return l
}

// Compatible determines if the provided language is compatible with any
// of the language ranges within the Accept-Language header value.
func (l AcceptLanguage) Compatible(language string) (c bool, err error) {
	for _, r := range l.LanguageRanges() {
		if c = r.Compatible(language); err != nil || c {
			return
		}
	}
	return
}

// String provides a textual representation of the Accept-Language header.
func (l AcceptLanguage) String() string {
	var languageRanges []string
	for _, lr := range l {
		languageRanges = append(languageRanges, fmt.Sprintf("%s;q=%s", lr.tag.String(), lr.QualityValue().String()))
	}
	return fmt.Sprintf("%s: %s", headerAcceptLanguage, strings.Join(languageRanges, ","))
}
