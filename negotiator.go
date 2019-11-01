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

package negotiator

import (
	"net/http"

	"github.com/freerware/negotiator/representation"
)

// NegotiationContext represents the context, such as the nature of the request
// and request itself, in which content negotiation is occuring.
type NegotiationContext struct {
	ResponseWriter http.ResponseWriter
	Request        *http.Request
	IsCreation     bool
}

// Negotiator represents a content negotiator.
type Negotiator interface {
	Negotiate(NegotiationContext, ...representation.Representation) error
}
