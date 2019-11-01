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
	// headerAccept is the header key for the Accept header.
	headerAccept = "Accept"

	// DefaultAccept is an Accept header with a single media range of "*/*".
	DefaultAccept = Accept([]MediaRange{defaultMediaRange})

	// EmptyAccept is an empty Accept header.
	EmptyAccept = Accept([]MediaRange{})
)

// Accept represents the Accept header.
//
// The "Accept" header field can be used by user agents to specify
// response media types that are acceptable.  Accept header fields can
// be used to indicate that the request is specifically limited to a
// small set of desired types, as in the case of a request for an
// in-line image.
type Accept []MediaRange

// NewAccept constructs an Accept header with the provided media ranges.
func NewAccept(accept []string) (Accept, error) {
	if len(accept) == 0 {
		return EmptyAccept, nil
	}

	// parse media ranges
	var mediaRanges []MediaRange
	for _, m := range accept {
		mediaRange, err := parseMediaRange(m)
		if err != nil {
			return EmptyAccept, err
		}
		mediaRanges = append(mediaRanges, mediaRange)
	}
	return Accept(mediaRanges), nil
}

// MediaRanges provides the media ranges sorted on preference and precedence,
// from highest preference and precedence to lowest.
func (a Accept) MediaRanges() []MediaRange {
	sort.Slice(a, func(first, second int) bool {
		f := a[first]
		s := a[second]

		if f.QualityValue().Equals(s.QualityValue()) {
			return f.Precedence() > s.Precedence()
		}
		return f.QualityValue().GreaterThan(s.QualityValue())
	})
	return a
}

// Compatible determines if the provided media type is compatible with any
// of the media ranges within the Accept header value.
func (a Accept) Compatible(mediaType string) (c bool, err error) {
	for _, r := range a.MediaRanges() {
		if c, err = r.Compatible(mediaType); err != nil || c {
			return
		}
	}
	return
}

// IsEmpty indicates if the Accept header is empty.
func (a Accept) IsEmpty() bool {
	return len(a) == len(EmptyAccept)
}

// String provides a textual representation of the Accept header.
func (a Accept) String() string {
	var mediaRanges []string
	for _, mr := range a {
		mediaRanges = append(mediaRanges, mr.String())
	}
	return fmt.Sprintf("%s: %s", headerAccept, strings.Join(mediaRanges, ","))
}
