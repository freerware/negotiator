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
	"strings"
)

var (
	// ErrEmptyTCNValue is an error that indicates that the TCN value cannot
	// be empty.
	ErrEmptyTCNValue = errors.New("TCN value cannot be empty")
)

// ResponseType represents the type of transparent negotiation response type.
type ResponseType string

const (

	// ResponseTypeList indicates the response to the transparent negotiation
	// request contains a list of the available representations.
	ResponseTypeList ResponseType = "list"

	// ResponseTypeChoice indicates the response to the transparent negotiation
	// request contains a chosen representation using a server-side algorithm.
	ResponseTypeChoice ResponseType = "choice"

	// ResponseTypeAdhoc indicates the response to the transparent negotiation
	// request is acting in the interest of achieving compatibility with a
	// non-negotiation or buggy client.
	ResponseTypeAdhoc ResponseType = "adhoc"
)

// String provides the textual representation of the response type.
func (rt ResponseType) String() string {
	return string(rt)
}

// OverrideDirective represents a server-side override performed when producting
// a response during transparent negotiation.
type OverrideDirective string

const (

	// OverrideDirectiveReChoose indicates to the user agent it SHOULD use its
	// internal variant selection algorithm to choose, retrieve, and display
	// the best variant from this list.
	OverrideDirectiveReChoose OverrideDirective = "re-choose"

	// OverrideDirectiveKeep indicates to the user agent it should not renegotiation
	// on the response to the transparent negotiation request and use it directly.
	OverrideDirectiveKeep OverrideDirective = "keep"
)

// String provides the textual representation of the override directive.
func (od OverrideDirective) String() string {
	return string(od)
}

// TCNValue represents a value specified within the TCN header.
type TCNValue string

// NewTCNValue constructs a new value for the TCN header.
func NewTCNValue(value string) (TCNValue, error) {
	if len(value) == 0 {
		return TCNValue(""), ErrEmptyTCNValue
	}
	return TCNValue(value), nil
}

// IsExtension indicates if the TCN value is an extension.
func (v TCNValue) IsExtension() bool {
	override := map[OverrideDirective]bool{
		OverrideDirectiveReChoose: true,
		OverrideDirectiveKeep:     true,
	}
	responseType := map[ResponseType]bool{
		ResponseTypeList:   true,
		ResponseTypeChoice: true,
		ResponseTypeAdhoc:  true,
	}
	return !override[OverrideDirective(v)] && !responseType[ResponseType(v)]
}

// String provides the textual representation of the TCN value.
func (v TCNValue) String() string {
	return string(v)
}

var (

	// headerTCN is the header key for the TCN header.
	headerTCN = "TCN"

	// EmptyTCN is an empty TCN header.
	EmptyTCN = TCN([]TCNValue{})
)

// TCN represents the TCN header.
type TCN []TCNValue

// NewTCN constructs a new TCN header with the value provided.
func NewTCN(values []string) (TCN, error) {
	if len(values) == 0 {
		return EmptyTCN, nil
	}
	var vals []TCNValue
	for _, value := range values {
		val, err := NewTCNValue(value)
		if err != nil {
			return EmptyTCN, err
		}
		vals = append(vals, val)
	}
	return TCN(vals), nil
}

// String provides the textual representation of the TCN header value.
func (t TCN) String() string {
	return fmt.Sprintf("%s: %s", headerTCN, t.ValuesAsString())
}

// ValuesAsStrings provides the string representation for each value of
// for the TCN header.
func (t TCN) ValuesAsStrings() []string {
	var s []string
	for _, v := range t {
		s = append(s, v.String())
	}
	return s
}

// ValuesAsString provides a single string containing all of the values for
// the TCN header.
func (t TCN) ValuesAsString() string {
	return strings.Join(t.ValuesAsStrings(), ",")
}
