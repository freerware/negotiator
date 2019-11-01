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
	"strconv"
)

var (

	// ErrEmptyNegotiateDirective is an error that indicates that the
	// negotiate directive cannot be empty.
	ErrEmptyNegotiateDirective = errors.New("negotiate directive cannot be empty")
)

// NegotiateDirective represents a directive specified within the Negotiate
// header.
type NegotiateDirective string

const (
	// NegotiateDirectiveTrans indicates the user agent supports transparent
	// content negotiation for the current request.
	NegotiateDirectiveTrans NegotiateDirective = "trans"

	// NegotiateDirectiveVList indicates the user agent requests that any
	// transparently negotiated response for the current request includes an
	// Alternates header with the variant list bound to the negotiable resource.
	NegotiateDirectiveVList NegotiateDirective = "vlist"

	// NegotiateDirectiveGuessSmall indicates the user agent allows origin
	// servers to run a custom algorithm which guesses the best variant for
	// the request, and to return this variant in a choice response, if the
	// resulting choice response is smaller than or not much larger than a list
	// response.
	NegotiateDirectiveGuessSmall NegotiateDirective = "guess-small"
)

// NewNegotiateDirective constructs a new directive for the Negotiate header.
func NewNegotiateDirective(directive string) (NegotiateDirective, error) {
	if len(directive) == 0 {
		return NegotiateDirective(""), ErrEmptyNegotiateDirective
	}
	return NegotiateDirective(directive), nil
}

// IsWildcard indicates if the negotiate directive is a wildcard.
func (d NegotiateDirective) IsWildcard() bool {
	return d.String() == "*"
}

// IsRVSAVersion indicates if the negotiate directive is an RSVA version.
func (d NegotiateDirective) IsRVSAVersion() bool {
	_, err := strconv.ParseFloat(string(d), 32)
	return err == nil
}

// IsExtension indicates if the negotiate directive is an extension.
func (d NegotiateDirective) IsExtension() bool {
	directive := map[NegotiateDirective]bool{
		NegotiateDirectiveTrans:      true,
		NegotiateDirectiveVList:      true,
		NegotiateDirectiveGuessSmall: true,
	}
	return !directive[d] && !d.IsWildcard() && !d.IsRVSAVersion()
}

// String provides the textual representation of the TCN value.
func (d NegotiateDirective) String() string {
	return string(d)
}
