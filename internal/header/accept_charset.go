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
	// headerAcceptCharset is the header key for the Accept-Charset header.
	headerAcceptCharset = "Accept-Charset"

	// DefaultAcceptCharset is an Accept-Charset header with a charset range
	// of "*" and quality value of 1.0.
	DefaultAcceptCharset = AcceptCharset([]CharsetRange{defaultCharsetRange})

	// EmptyAcceptCharset is an empty Accept-Charset header.
	EmptyAcceptCharset = AcceptCharset([]CharsetRange{})
)

// AcceptCharset represents the Accept-Charset header.
//
// The "Accept-Charset" header field can be sent by a user agent to
// indicate what charsets are acceptable in textual response content.
// This field allows user agents capable of understanding more
// comprehensive or special-purpose charsets to signal that capability
// to an origin server that is capable of representing information in
// those charsets.
type AcceptCharset []CharsetRange

// NewAcceptCharset constructs an Accept-Charset header with the provided
// charsets.
func NewAcceptCharset(acceptCharset []string) (AcceptCharset, error) {
	if len(acceptCharset) == 0 {
		return EmptyAcceptCharset, nil
	}
	var charsets []CharsetRange
	for _, c := range acceptCharset {
		charset, err := NewCharsetRange(c)
		if err != nil {
			return EmptyAcceptCharset, err
		}
		charsets = append(charsets, charset)
	}
	return AcceptCharset(charsets), nil
}

// CharsetRanges provides the charsets sorted on preference from highest
// preference to lowest.
func (c AcceptCharset) CharsetRanges() []CharsetRange {
	sort.Slice(c, func(first, second int) bool {
		f := c[first]
		s := c[second]
		return f.QualityValue().GreaterThan(s.QualityValue())
	})
	return c
}

// Compatible determines if the provided charset is compatible with any
// of the charset ranges within the Accept-Charset header value.
func (c AcceptCharset) Compatible(charset string) (cc bool, err error) {
	for _, r := range c.CharsetRanges() {
		if cc = r.Compatible(charset); err != nil || cc {
			return
		}
	}
	return
}

// IsEmpty indicates if the Accept-Charset header is empty.
func (c AcceptCharset) IsEmpty() bool {
	return len(c) == len(EmptyAcceptCharset)
}

// String provides the textual representation of the Accept-Charset header value.
func (c AcceptCharset) String() string {
	var charsets []string
	for _, cr := range c {
		charsets = append(charsets, cr.String())
	}
	return fmt.Sprintf("%s: %s", headerAcceptCharset, strings.Join(charsets, ","))
}
