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
	"math"
)

// QualityValue represents a relative weight associated with proactive
// negotiation headers.
//
// Many of the request header fields for proactive negotiation use a
// common parameter, named "q" (case-insensitive), to assign a relative
// "weight" to the preference for that associated kind of content.  This
// weight is referred to as a "quality value" (or "qvalue")
type QualityValue float32

const (
	// QualityValueMinimum represents the minimum allowed quality value.
	QualityValueMinimum QualityValue = 0.0

	// QualityValueMaximum represents the maximum allowed quality value.
	QualityValueMaximum QualityValue = 1.0

	// QualityValueDefault represents the default quality value when not
	// specified explicitly.
	QualityValueDefault QualityValue = QualityValueMaximum
)

var (
	// ErrInvalidQualityValue is an error that indicates that the quality value must be between 0.0 and 1.0.
	ErrInvalidQualityValue = errors.New("quality range must be between 0.0 and 1.0")
)

// NewQualityValue constructs a new quality value using the provided number.
//
// The weight is normalized to a real number in the range 0 through 1,
// where 0.001 is the least preferred and 1 is the most preferred; a
// value of 0 means "not acceptable".
func NewQualityValue(qvalue float32) (QualityValue, error) {
	tooLow, tooHigh :=
		qvalue < QualityValueMinimum.Float(),
		qvalue > QualityValueMaximum.Float()
	if tooLow || tooHigh {
		return QualityValue(0), ErrInvalidQualityValue
	}
	return QualityValue(qvalue), nil
}

// Equals determines if the provided quality value is equivalent.
func (qv QualityValue) Equals(q QualityValue) bool {
	// Per https://tools.ietf.org/html/rfc7231#section-5.3.1, quality values
	// have a precision to the thousands place.
	return qv.Round(3) == q.Round(3)
}

// Round rounds the quality value with the scale provided
// and returns the rounded result.
func (qv QualityValue) Round(scale int) QualityValue {
	if scale < 1 {
		return QualityValue(float32(math.Round(float64(qv))))
	}
	factor := float32(math.Pow10(scale))
	rounded := float32(math.Round(float64(qv.Float()*factor)) / float64(factor))
	return QualityValue(rounded)
}

// LessThan determines if the provided quality value is lesser.
func (qv QualityValue) LessThan(q QualityValue) bool {
	// Per https://tools.ietf.org/html/rfc7231#section-5.3.1, quality values
	// have a precision to the thousands place.
	return qv.Round(3) < q.Round(3)
}

// GreaterThan determines if the provided quality value is greater.
func (qv QualityValue) GreaterThan(q QualityValue) bool {
	// Per https://tools.ietf.org/html/rfc7231#section-5.3.1, quality values
	// have a precision to the thousands place.
	return qv.Round(3) > q.Round(3)
}

// Multiply multiplies the quality value with the one provided.
func (qv QualityValue) Multiply(q QualityValue) QualityValue {
	return qv * q
}

// String provides the textual representation of the quality value.
func (qv QualityValue) String() string {
	return fmt.Sprintf("%.3f", qv)
}

// Float provides the floating point representation of the quality value.
func (qv QualityValue) Float() float32 {
	return float32(qv)
}
