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
	"strings"
)

var (
	// headerNegotiate is the header key for the Negotiate header.
	headerNegotiate = "Negotiate"

	// EmptyNegotiateHeader is an empty Negotiate header.
	EmptyNegotiateHeader = Negotiate([]NegotiateDirective{})
)

// Negotiate represents the Negotiate header.
type Negotiate []NegotiateDirective

// NewNegotiate constructs a Negotiate header with the provided directives.
func NewNegotiate(directives []string) (Negotiate, error) {
	var dirs []NegotiateDirective
	for _, d := range directives {
		directive, err := NewNegotiateDirective(d)
		if err != nil {
			return EmptyNegotiateHeader, err
		}
		dirs = append(dirs, directive)
	}
	return Negotiate(dirs), nil
}

// Directives provides the negotiation directives.
func (n Negotiate) Directives() []NegotiateDirective {
	return n
}

// Contains determines if the Negotiate header contains at least one of the
// provided directives.
func (n Negotiate) Contains(directives ...NegotiateDirective) (matches bool) {
	for _, dir := range directives {
		for _, d := range n.Directives() {
			if matches = strings.ToLower(d.String()) == strings.ToLower(dir.String()); matches {
				return
			}
		}
	}
	return
}

// ContainsRVSA determines if the Negotiate header contains an RVSA algorithm
// that matches the version provided.
func (n Negotiate) ContainsRVSA(version string) (matches bool) {
	for _, d := range n.Directives() {
		rvsaDir := NegotiateDirective(version)
		if matches = d.IsRVSAVersion() &&
			strings.ToLower(rvsaDir.String()) == strings.ToLower(d.String()); matches {
			return
		}
	}
	return
}

// String provides the textual representation of the Negotiate header.
func (n Negotiate) String() string {
	var s []string
	for _, d := range n.Directives() {
		s = append(s, d.String())
	}
	return fmt.Sprintf("%s: %s", headerNegotiate, strings.Join(s, ","))
}
