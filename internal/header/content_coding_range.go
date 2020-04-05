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
	"regexp"
	"strconv"
	"strings"

	"github.com/stretchr/stew/slice"
)

var (
	contentCodingRangeRegex = regexp.MustCompile(`^([A-Za-z0-9-]+|\*)(;\s?q=(\d(\.\d{1,3})?))?$`)
)

var (
	gzip      = "gzip"
	xgzip     = "x-gzip"
	deflate   = "deflate"
	compress  = "compress"
	xcompress = "x-compress"
	identity  = "identity"

	contentCodings = []string{
		gzip,
		xgzip,
		deflate,
		compress,
		xcompress,
		identity,
		"*",
	}

	// defaultContentCodingRange represents the default content coding range.
	defaultContentCodingRange = ContentCodingRange{
		coding: "*",
		qValue: QualityValueMaximum,
	}

	// ErrEmptyContentCodingRange is an error that indicates that the content
	// coding range cannot be empty.
	ErrEmptyContentCodingRange = errors.New("content coding range cannot be empty")

	// ErrInvalidContentCodingRange is an error that indicates that the content
	// coding range is invalid.
	ErrInvalidContentCodingRange = errors.New("content coding range is invalid")
)

// ContentCodingRange represents an expression that indicates an encoding
// transformation.
//
// Content coding values indicate an encoding transformation that has
// been or can be applied to a representation.  Content codings are
// primarily used to allow a representation to be compressed or
// otherwise usefully transformed without losing the identity of its
// underlying media type and without loss of information.
type ContentCodingRange struct {
	coding string
	qValue QualityValue
}

// NewContentCodingRange constructs a content coding from the textual representation.
func NewContentCodingRange(contentCoding string) (ContentCodingRange, error) {
	if len(contentCoding) == 0 {
		return ContentCodingRange{}, ErrEmptyContentCodingRange
	}

	if ok := contentCodingRangeRegex.MatchString(contentCoding); !ok {
		return ContentCodingRange{}, ErrInvalidContentCodingRange
	}
	groups := contentCodingRangeRegex.FindStringSubmatch(contentCoding)

	var valid []string
	valid = append(valid, contentCodings...)
	if !slice.ContainsString(contentCodings, groups[1]) {
		return ContentCodingRange{}, ErrInvalidContentCodingRange
	}
	cc := ContentCodingRange{coding: groups[1], qValue: QualityValue(1.0)}

	if len(groups[3]) > 0 {
		q, _ := strconv.ParseFloat(groups[4], 32)
		qv, err := NewQualityValue(float32(q))
		if err != nil {
			return ContentCodingRange{}, err
		}
		cc.qValue = qv
	}
	return cc, nil
}

// IsWildcard indicates if the specified coding range is '*'.
func (cc ContentCodingRange) IsWildcard() bool {
	return cc.coding == "*"
}

// IsIdentity indicates if the specified coding range is 'identi'.
func (cc ContentCodingRange) IsIdentity() bool {
	return strings.ToLower(cc.coding) == strings.ToLower(identity)
}

// IsCoding indicates if the specified coding range is a content coding.
func (cc ContentCodingRange) IsCoding() bool {
	return !cc.IsWildcard() && !cc.IsIdentity()
}

// Coding retrieves the content coding.
func (cc ContentCodingRange) CodingRange() string {
	return cc.coding
}

// Compatible determines if the provided content coding is compatible with the
// content coding range.
func (cc ContentCodingRange) Compatible(coding string) bool {
	if !slice.ContainsString(contentCodings, strings.ToLower(coding)) {
		return false
	}
	if cc.IsWildcard() {
		return true
	}
	return strings.ToLower(cc.CodingRange()) == strings.ToLower(coding)
}

// QualityValue retrieves the quality value of the content coding.
//
// Each codings value MAY be given an associated quality value
// representing the preference for that encoding, as defined in
// Section 5.3.1.
func (cc ContentCodingRange) QualityValue() QualityValue {
	return cc.qValue
}

// String provides a textual representation of the content coding.
func (cc ContentCodingRange) String() string {
	return fmt.Sprintf("%s;q=%s", cc.CodingRange(), cc.QualityValue().String())
}
