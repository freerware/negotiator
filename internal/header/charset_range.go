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
)

var (
	// defaultCharsetRange represents the default charset range.
	defaultCharsetRange = CharsetRange{
		r:      "*",
		qValue: QualityValueMaximum,
	}

	// ErrEmptyCharsetRange is an error that indicates that the charset
	// range cannot be empty.
	ErrEmptyCharsetRange = errors.New("charset range cannot be empty")
)

// CharsetRange represents an expression that indicates an encoding
// transformation.
type CharsetRange struct {
	r      string
	qValue QualityValue
}

// NewCharsetRange constructs a charset from the textual representation.
func NewCharsetRange(charset string) (CharsetRange, error) {
	if len(charset) == 0 {
		return CharsetRange{}, ErrEmptyCharsetRange
	}
	parts := strings.Split(charset, ";")
	r := strings.ToLower(parts[0])
	cc := CharsetRange{
		r:      r,
		qValue: QualityValue(1.0),
	}

	if len(parts) > 1 && strings.HasPrefix(strings.Trim(parts[1], " "), "q=") {
		w := strings.Trim(parts[1], " ")
		q := strings.TrimPrefix(w, "q=")
		f, err := strconv.ParseFloat(q, 32)
		if err != nil {
			return CharsetRange{}, err
		}
		qv, err := NewQualityValue(float32(f))
		if err != nil {
			return CharsetRange{}, err
		}
		cc.qValue = qv
	}
	return cc, nil
}

// IsWildcard indicates if the charset range is '*'.
func (c CharsetRange) IsWildcard() bool {
	return c.r == "*"
}

// IsCharset indicates that the charset range is a charset.
func (c CharsetRange) IsCharset() bool {
	return !c.IsWildcard()
}

// Charset retrieves the range value for the charset range.
func (c CharsetRange) Charset() string {
	return c.r
}

// QualityValue retrieves the quality value of the charset range.
func (c CharsetRange) QualityValue() QualityValue {
	return c.qValue
}

// Compatible determines if the provided charset is compatible with the
// charset range.
func (c CharsetRange) Compatible(charset string) bool {
	if c.IsWildcard() {
		return true
	}
	return strings.ToLower(c.r) == strings.ToLower(charset)
}

// String provides a textual representation of the charset range.
func (c CharsetRange) String() string {
	return fmt.Sprintf("%s;q=%s", c.Charset(), c.QualityValue().String())
}
