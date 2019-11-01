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

package reactive

import (
	"net/http"
	"strings"

	"github.com/freerware/negotiator"
	"github.com/freerware/negotiator/internal/representation/json"
	"github.com/freerware/negotiator/representation"
)

var (
	// Default is the default reactive negotiator.
	Default = New()
)

// Negotiator represents the negotiator responsible for performing
// reactive (agent-driven) negotiation.
type Negotiator struct {
	representationConstructor representation.Constructor
}

// New constructs a negotiator capable of performing reactive
// (agent-driven) negotiation with the options provided.
func New(options ...Option) Negotiator {
	// set defaults.
	o := Options{
		RepresentationConstructor: json.List,
	}
	// apply options.
	for _, opt := range options {
		opt(&o)
	}
	return Negotiator{
		representationConstructor: o.RepresentationConstructor,
	}
}

// Negotiate performs reactive (agent-driven) content negotiation with the
// representations provided.
func (n Negotiator) Negotiate(
	ctx negotiator.NegotiationContext,
	reps ...representation.Representation) error {

	if len(reps) == 0 {
		ctx.ResponseWriter.WriteHeader(http.StatusNoContent)
		return nil
	}

	// construct representation.
	chosen := n.representationConstructor(reps...)

	// serialize.
	b, err := chosen.Bytes()
	if err != nil {
		return err
	}

	// respond.
	ctx.ResponseWriter.Header().Add("Content-Length", string(len(b)))
	ctx.ResponseWriter.Header().Add("Content-Type", chosen.ContentType())
	ctx.ResponseWriter.Header().Add("Content-Encoding", strings.Join(chosen.ContentEncoding(), ","))
	ctx.ResponseWriter.Header().Add("Content-Language", chosen.ContentLanguage())
	ctx.ResponseWriter.Header().Add("Content-Charset", chosen.ContentCharset())
	ctx.ResponseWriter.WriteHeader(http.StatusMultipleChoices)
	_, err = ctx.ResponseWriter.Write(b)
	return err
}
