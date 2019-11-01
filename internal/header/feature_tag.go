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
	"strconv"
	"strings"
)

// FeatureTag represents a feature tag.
// https://tools.ietf.org/html/rfc2295#section-6.1
type FeatureTag string

// Equals performs a case-insensitive comparison with the provided feature tag.
// The US-ASCII charset is used for feature tag.
func (t FeatureTag) Equals(tag FeatureTag) bool {
	u, err := t.unquotedASCII()
	if err != nil {
		return false
	}
	q, err := t.quotedASCII()
	if err != nil {
		return false
	}
	uu, err := tag.unquotedASCII()
	if err != nil {
		return false
	}
	qq, err := tag.quotedASCII()
	if err != nil {
		return false
	}
	return strings.ToLower(u) == strings.ToLower(uu) ||
		strings.ToLower(q) == strings.ToLower(qq)
}

// quotedASCII encodes the tag to ASCII and surrounds it with double quotes.
func (t FeatureTag) quotedASCII() (string, error) {
	// remove quotes if they exist.
	u := strings.Trim(t.String(), "\"")
	// convert to ASCII and quote.
	return strconv.QuoteToASCII(u), nil
}

// unquotedASCII encodes the tag to ASCII and removes surrounding double quotes.
func (t FeatureTag) unquotedASCII() (string, error) {
	var q string
	var err error
	// convert ASCII and quote.
	if q, err = t.quotedASCII(); err != nil {
		return q, err
	}
	// unquote.
	return strconv.Unquote(q)
}

// String provides the textual representation of the feature tag.
func (t FeatureTag) String() string {
	return string(t)
}

// FeatureTagValue represents a feature tag value.
// https://tools.ietf.org/html/rfc2295#section-6.1.1
type FeatureTagValue string

// Equals performs a case-sensitive, octect-by-octet comparison with the
// provided feature tag value. The US-ASCII charset is used for feature tag values.
func (t FeatureTagValue) Equals(val FeatureTagValue) bool {
	u, err := t.unquotedASCII()
	if err != nil {
		return false
	}
	q, err := t.quotedASCII()
	if err != nil {
		return false
	}
	uu, err := val.unquotedASCII()
	if err != nil {
		return false
	}
	qq, err := val.quotedASCII()
	if err != nil {
		return false
	}
	return u == uu || q == qq
}

// quotedASCII encodes the tag value to ASCII and surrounds it with double quotes.
func (t FeatureTagValue) quotedASCII() (string, error) {
	// remove quotes if they exist.
	u := strings.Trim(t.String(), "\"")
	// convert to ASCII and quote.
	return strconv.QuoteToASCII(u), nil
}

// unquotedASCII encodes the tag value to ASCII and removes surrounding double quotes.
func (t FeatureTagValue) unquotedASCII() (string, error) {
	var q string
	var err error
	// convert ASCII and quote.
	if q, err = t.quotedASCII(); err != nil {
		return q, err
	}
	// unquote.
	return strconv.Unquote(q)
}

// IsNumeric indicates if the feature tag value is numeric.
func (t FeatureTagValue) IsNumeric() bool {
	_, err := strconv.ParseFloat(t.String(), 32)
	return err == nil
}

// String provides the textual representation of the feature tag value.
func (t FeatureTagValue) String() string {
	return string(t)
}
