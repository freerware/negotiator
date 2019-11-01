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
	"net/url"
	"sort"
	"strings"

	"github.com/freerware/negotiator/representation"
)

var (
	// headerAlternates is the header key for the Alternates header.
	headerAlternates = "Alternates"

	// variantAttributeType represents the attribute key that communicates
	// the media type of the variant.
	variantAttributeType = "type"

	// variantAttributeCharset represents the attribute key that communicates
	// the charset of the variant.
	variantAttributeCharset = "charset"

	// variantAttributeLanguage represents the attribute key that communicates
	// the language(s) of the variant.
	variantAttributeLanguage = "language"

	// variantAttributeLength represents the attribute that communicates the
	// variant size in bytes.
	variantAttributeLength = "length"

	// variantAttributeFeatures represents the attribute that communicates the
	// feature list of the variant.
	variantAttributeFeatures = "features"

	// variantAttributeDescription represents the attribute that communicates
	// the textual description of the variant.
	variantAttributeDescription = "description"
)

// variantAttributes represents the various dimensions of a variant described
// within a variant description.
type variantAttributes map[string]interface{}

// String provides the textual representation of the variant attributes.
func (a variantAttributes) String() string {
	// sort the keys for deterministic output.
	var keys []string
	for key := range a {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	var s []string
	for _, key := range keys {
		value := a[key]
		s = append(s, fmt.Sprintf("{ %s %v }", key, value))
	}
	return strings.Join(s, " ")
}

// variantDescription represents the complete description of a variant,
// including it's URL, source quality, and attributes.
type variantDescription struct {
	uri           url.URL
	sourceQuality float32
	attributes    variantAttributes
}

// String provides the textual representation of the variant description.
func (a variantDescription) String() string {
	return fmt.Sprintf("{ %q %.3f %s }", a.uri.String(), a.sourceQuality, a.attributes)
}

// variantFallback represents the URL of the variant that user agents should
// use when it finds all variants listed in the variant description as unacceptable.
type variantFallback url.URL

// String provides the textual representation of the variant fallback.
func (f variantFallback) String() string {
	u := url.URL(f)
	return fmt.Sprintf("{ %q }", u.String())
}

// Alternates represents the Alternates header.
type Alternates struct {
	descriptions []variantDescription
	fallback     *variantFallback
}

// NewAlternates constructs an Alternates header with the provided representations.
func NewAlternates(
	fb representation.Representation, reps ...representation.Representation) (Alternates, error) {
	var descriptions []variantDescription
	for _, rep := range reps {
		bytes, err := rep.Bytes()
		if err != nil {
			return Alternates{}, err
		}
		descriptions = append(descriptions, variantDescription{
			uri:           rep.ContentLocation(),
			sourceQuality: rep.SourceQuality(),
			attributes: map[string]interface{}{
				variantAttributeType:     rep.ContentType(),
				variantAttributeCharset:  rep.ContentCharset(),
				variantAttributeLanguage: rep.ContentLanguage(),
				variantAttributeFeatures: strings.Join(rep.ContentFeatures(), " "),
				variantAttributeLength:   len(bytes),
			},
		})
	}
	var fallback *variantFallback
	if fb != nil {
		vf := variantFallback(fb.ContentLocation())
		fallback = &vf
	}
	return Alternates{descriptions: descriptions, fallback: fallback}, nil
}

// HasFallback indicates if a fallback variant has been specified.
func (a Alternates) HasFallback() bool {
	return a.fallback != nil
}

// String provides the textual representation of the Alternates header value.
func (a Alternates) String() string {
	return fmt.Sprintf("%s: %s", headerAlternates, a.ValuesAsString())
}

// ValuesAsStrings provides the string representation for each value of
// for the Alternates header.
func (a Alternates) ValuesAsStrings() []string {
	s := []string{}
	for _, d := range a.descriptions {
		s = append(s, d.String())
	}
	if a.HasFallback() {
		s = append(s, a.fallback.String())
	}
	return s
}

// ValuesAsString provides a single string containing all of the values for// the Alternates header.
func (a Alternates) ValuesAsString() string {
	return strings.Join(a.ValuesAsStrings(), ",")
}
