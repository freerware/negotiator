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
	"strconv"
	"strings"

	"github.com/freerware/negotiator"
	"github.com/freerware/negotiator/representation"
	"github.com/uber-go/tally"
	"go.uber.org/zap"
)

// Defines the default reactive negotiator.
//
// The default configuration is as follows:
//
// ➣ The representation for 406 Not Acceptable responses utilizes the JSON
// (application/json) media type.
var (
	// Default is the default reactive negotiator.
	Default = New()
)

// Defines the tags and scope names for reactive negotiator metrics.
var (
	scopeTagReactive                        = map[string]string{"negotiator": "reactive"}
	scopeNameReactiveTimer                  = "negotiate"
	scopeNameReactiveMultipleChoicesCounter = "negotiate.multiple_choices"
	scopeNameReactiveNoContentCounter       = "negotiate.no_content"
	scopeNameReactiveErrorCounter           = "negotiate.error"
)

var (
	jsonList = func(reps ...representation.Representation) representation.Representation {
		list := representation.List{}
		list.SetContentType("application/json")
		list.SetContentCharset("ascii")
		list.SetContentEncoding([]string{"identity"})
		list.SetContentLanguage("en-US")
		list.SetRepresentations(reps...)
		return &list
	}
)

// Negotiator represents the negotiator responsible for performing
// reactive (agent-driven) negotiation.
type Negotiator struct {
	representationConstructor representation.ListConstructor
	logger                    *zap.Logger
	scope                     tally.Scope
}

// New constructs a negotiator capable of performing reactive
// (agent-driven) negotiation with the options provided.
func New(options ...Option) Negotiator {
	// set defaults.
	o := Options{
		RepresentationConstructor: jsonList,
		Logger:                    zap.NewNop(),
		Scope:                     tally.NoopScope,
	}
	// apply options.
	for _, opt := range options {
		opt(&o)
	}
	n := Negotiator{
		representationConstructor: o.RepresentationConstructor,
		logger:                    o.Logger,
		scope:                     o.Scope.Tagged(scopeTagReactive),
	}
	n.logger.Debug("negotiator configuration", zap.String("type", "reactive"))
	return n
}

// Negotiate performs reactive (agent-driven) content negotiation with the
// representations provided.
func (n Negotiator) Negotiate(
	ctx negotiator.NegotiationContext,
	reps ...representation.Representation) (err error) {

	defer n.scope.Timer(scopeNameReactiveTimer).Start().Stop()
	defer func() {
		if err != nil {
			n.scope.Counter(scopeNameReactiveErrorCounter).Inc(1)
		} else {
			n.scope.Counter(scopeNameReactiveMultipleChoicesCounter).Inc(1)
		}
	}()

	if len(reps) == 0 {
		status := http.StatusNoContent
		ctx.ResponseWriter.WriteHeader(status)
		n.logger.Info("no representations to negotiate", zap.Int("status", status))
		n.scope.Counter(scopeNameReactiveNoContentCounter).Inc(1)
		return
	}

	// construct representation.
	rep := n.representationConstructor(reps...)

	// serialize.
	var b []byte
	if b, err = rep.Bytes(); err != nil {
		return
	}

	// respond.
	var (
		clen   = len(b)
		ct     = rep.ContentType()
		ce     = strings.Join(rep.ContentEncoding(), ",")
		clang  = rep.ContentLanguage()
		cc     = rep.ContentCharset()
		status = http.StatusMultipleChoices
	)
	ctx.ResponseWriter.Header().Add("Content-Length", strconv.Itoa(clen))
	ctx.ResponseWriter.Header().Add("Content-Type", ct)
	ctx.ResponseWriter.Header().Add("Content-Encoding", ce)
	ctx.ResponseWriter.Header().Add("Content-Language", clang)
	ctx.ResponseWriter.Header().Add("Content-Charset", cc)
	ctx.ResponseWriter.WriteHeader(status)
	if _, err = ctx.ResponseWriter.Write(b); err == nil {
		n.logger.Info("multiple choices",
			zap.Int("content-length", clen),
			zap.String("content-type", ct),
			zap.String("content-encoding", ce),
			zap.String("content-language", clang),
			zap.String("content-charset", cc),
			zap.Int("status", status))
	}
	return
}
