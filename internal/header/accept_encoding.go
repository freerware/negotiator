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
)

var (
	// headerAcceptEncoding is the header key for the Accept-Encoding header.
	headerAcceptEncoding = "Accept-Encoding"

	// DefaultAcceptEncoding is an Accept-Encoding header value with a
	// single content coding of "*".
	DefaultAcceptEncoding = AcceptEncoding([]ContentCodingRange{defaultContentCodingRange})

	// EmptyAcceptEncoding is an empty Accept-Encoding header.
	EmptyAcceptEncoding = AcceptEncoding([]ContentCodingRange{})
)

// AcceptEncoding represents the Accept-Encoding header.
//
// The "Accept-Encoding" header field can be used by user agents to
// indicate what response content-codings (Section 3.1.2.1) are
// acceptable in the response.
type AcceptEncoding []ContentCodingRange

// NewAcceptEncoding constructs an Accept-Encoding header with the provided
// content codings.
func NewAcceptEncoding(acceptEncoding []string) (AcceptEncoding, error) {
	if len(acceptEncoding) == 0 {
		return EmptyAcceptEncoding, nil
	}
	var contentCodings []ContentCodingRange
	for _, cc := range acceptEncoding {
		coding, err := NewContentCodingRange(cc)
		if err != nil {
			return EmptyAcceptEncoding, err
		}
		contentCodings = append(contentCodings, coding)
	}
	return AcceptEncoding(contentCodings), nil
}

// Codings provides the content codings sorted on preference from highest
// preference to lowest.
func (e AcceptEncoding) CodingRanges() []ContentCodingRange {
	sort.Slice(e, func(first, second int) bool {
		f := e[first]
		s := e[second]
		return f.QualityValue().GreaterThan(s.QualityValue())
	})
	return e
}

// IsEmpty indicates if the Accept-Encoding header is empty.
func (e AcceptEncoding) IsEmpty() bool {
	return len(e) == len(EmptyAcceptEncoding)
}

// String provides a textual representation of the Accept-Encoding header.
func (e AcceptEncoding) String() string {
	var codings []string
	for _, c := range e {
		codings = append(codings, c.String())
	}
	return fmt.Sprintf("%s: %s", headerAcceptEncoding, strings.Join(codings, ","))
}
