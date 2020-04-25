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

package proactive

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/freerware/negotiator"
	"github.com/freerware/negotiator/internal/header"
	"github.com/freerware/negotiator/representation"
	"github.com/uber-go/tally"
	"go.uber.org/zap"
)

// Defines the default proactive negotiator.
//
// The default configuration is as follows:
//
// ➣ Strict mode is enabled for all proactive negotiation headers
//
// ➣ A representation containing information about available representations
// is returned with a 406 Not Acceptable response.
//
// ➣ The algorithm used to choose the 'best' representation is the Apache
// httpd algorithm.
//
// ➣ Representation candidates for 406 Not Acceptable responses support
// JSON (application/json), XML (application/xml), and YAML (application/yaml,
// text/yaml) media types.
//
// ➣ The fallback representation for 406 Not Acceptable responses utilizes
// the JSON (application/json) media type.
var (
	// Default is the default proactive negotiator.
	Default = New()
)

// Defines the tags and scope names for proactive negotiator metrics.
var (
	scopeTagProactive                      = map[string]string{"negotiator": "proactive"}
	scopeNameProactiveTimer                = "negotiate"
	scopeNameProactiveNotAcceptableCounter = "negotiate.not_acceptable"
	scopeNameProactiveAcceptableCounter    = "negotiate.acceptable"
	scopeNameProactiveNoContentCounter     = "negotiate.no_content"
	scopeNameProactiveErrorCounter         = "negotiate.error"
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

	xmlList = func(reps ...representation.Representation) representation.Representation {
		list := representation.List{}
		list.SetContentType("application/xml")
		list.SetContentCharset("ascii")
		list.SetContentEncoding([]string{"identity"})
		list.SetContentLanguage("en-US")
		list.SetRepresentations(reps...)
		return &list
	}

	yamlList = func(reps ...representation.Representation) representation.Representation {
		list := representation.List{}
		list.SetContentType("application/yaml")
		list.SetContentCharset("ascii")
		list.SetContentEncoding([]string{"identity"})
		list.SetContentLanguage("en-US")
		list.SetRepresentations(reps...)
		return &list
	}
)

// Negotiator represents the negotiator responsible for performing
// proactive (server-driven) negotiation.
type Negotiator struct {
	strictAccept                     bool
	strictAcceptLanguage             bool
	strictAcceptCharset              bool
	notAcceptableRepresentation      bool
	defaultRepresentationConstructor representation.ListConstructor
	representationConstructors       []representation.ListConstructor
	chooser                          representation.Chooser
	logger                           *zap.Logger
	scope                            tally.Scope
}

// New constructs a negotiator capable of performing proactive
// (server-driven) negotiation with the options provided.
func New(options ...Option) Negotiator {
	// set defaults.
	o := Options{
		StrictAccept:                     true,
		StrictAcceptLanguage:             true,
		StrictAcceptCharset:              true,
		NotAcceptableRepresentation:      true,
		DefaultRepresentationConstructor: jsonList,
		Chooser:                          ApacheHTTPD(),
		Logger:                           zap.NewNop(),
		Scope:                            tally.NoopScope,
		RepresentationConstructors: []representation.ListConstructor{
			jsonList,
			xmlList,
			yamlList,
		},
	}
	// apply options.
	for _, opt := range options {
		opt(&o)
	}
	n := Negotiator{
		strictAccept:                     o.StrictAccept,
		strictAcceptLanguage:             o.StrictAcceptLanguage,
		strictAcceptCharset:              o.StrictAcceptCharset,
		notAcceptableRepresentation:      o.NotAcceptableRepresentation,
		defaultRepresentationConstructor: o.DefaultRepresentationConstructor,
		representationConstructors:       o.RepresentationConstructors,
		chooser:                          o.Chooser,
		logger:                           o.Logger,
		scope:                            o.Scope.Tagged(scopeTagProactive),
	}
	n.logger.Debug("negotiator configuration",
		zap.String("type", "proactive"),
		zap.Bool("strict-accept", n.strictAccept),
		zap.Bool("strict-accept-language", n.strictAcceptLanguage),
		zap.Bool("strict-accept-charset", n.strictAcceptCharset),
		zap.Bool("not-acceptable-representation", n.notAcceptableRepresentation))
	return n
}

// Negotiate performs proactive (server-driven) content negotiation with the
// representations provided.
func (n Negotiator) Negotiate(
	ctx negotiator.NegotiationContext, reps ...representation.Representation) (err error) {

	defer n.scope.Timer(scopeNameProactiveTimer).Start().Stop()
	defer func() {
		if err != nil {
			n.scope.Counter(scopeNameProactiveErrorCounter).Inc(1)
		}
	}()

	if len(reps) == 0 {
		status := http.StatusNoContent
		ctx.ResponseWriter.WriteHeader(status)
		n.logger.Info("no representations to negotiate", zap.Int("status", status))
		n.scope.Counter(scopeNameProactiveNoContentCounter).Inc(1)
		return
	}

	var (
		accept         = header.DefaultAccept
		acceptLanguage = header.DefaultAcceptLanguage
		acceptCharset  = header.DefaultAcceptCharset
		ac, alc, acc   = 0, 0, 0
		hasHeader      bool
		headerValues   []string
	)
	if headerValues, hasHeader = ctx.Request.Header["Accept"]; hasHeader {
		if accept, err = header.NewAccept(headerValues); err != nil {
			return
		}
	}
	if headerValues, hasHeader = ctx.Request.Header["Accept-Language"]; hasHeader {
		if acceptLanguage, err = header.NewAcceptLanguage(headerValues); err != nil {
			return
		}
	}
	if headerValues, hasHeader = ctx.Request.Header["Accept-Charset"]; hasHeader {
		if acceptCharset, err = header.NewAcceptCharset(headerValues); err != nil {
			return
		}
	}
	for _, r := range reps {
		var c bool
		if c, err = accept.Compatible(r.ContentType()); err != nil {
			return
		} else if !c {
			ac = ac + 1
		}

		if c, err = acceptLanguage.Compatible(r.ContentLanguage()); err != nil {
			return
		} else if !c {
			alc = alc + 1
		}

		if c, err = acceptCharset.Compatible(r.ContentCharset()); err != nil {
			return
		} else if !c {
			acc = acc + 1
		}
	}
	if len(reps) == ac && n.strictAccept {
		n.logger.Debug("failed strict mode for Accept header")
		return n.notAcceptable(ctx, reps...)
	}
	if len(reps) == alc && n.strictAcceptLanguage {
		n.logger.Debug("failed strict mode for Accept-Language header")
		return n.notAcceptable(ctx, reps...)
	}
	if len(reps) == acc && n.strictAcceptCharset {
		n.logger.Debug("failed strict mode for Accept-Charset header")
		return n.notAcceptable(ctx, reps...)
	}

	// choose 'best' representation.
	var rep representation.Representation
	if rep, err = n.chooser.Choose(ctx.Request, reps...); err != nil {
		return
	}

	if rep == nil {
		return n.notAcceptable(ctx, reps...)
	}
	return n.acceptable(ctx, rep)
}

// acceptable is responsible for responding to the user agent with the
// representation chosen by the server-side algorithm.
func (n Negotiator) acceptable(
	ctx negotiator.NegotiationContext, rep representation.Representation) (err error) {

	defer func() {
		if err == nil {
			n.scope.Counter(scopeNameProactiveAcceptableCounter).Inc(1)
		}
	}()

	// serialize.
	var b []byte
	if b, err = rep.Bytes(); err != nil {
		return
	}

	// respond.
	var (
		clen  = len(b)
		ct    = rep.ContentType()
		ce    = strings.Join(rep.ContentEncoding(), ",")
		clang = rep.ContentLanguage()
		cc    = rep.ContentCharset()
		loc   = rep.ContentLocation()
	)
	status := http.StatusOK
	if ctx.IsCreation {
		status = http.StatusCreated
	}
	ctx.ResponseWriter.Header().Add("Content-Location", (&loc).String())
	ctx.ResponseWriter.Header().Add("Content-Length", strconv.Itoa(clen))
	ctx.ResponseWriter.Header().Add("Content-Type", ct)
	ctx.ResponseWriter.Header().Add("Content-Encoding", ce)
	ctx.ResponseWriter.Header().Add("Content-Language", clang)
	ctx.ResponseWriter.Header().Add("Content-Charset", cc)
	ctx.ResponseWriter.WriteHeader(status)
	if _, err = ctx.ResponseWriter.Write(b); err == nil {
		n.logger.Info("acceptable",
			zap.Int("content-length", clen),
			zap.String("content-type", ct),
			zap.String("content-encoding", ce),
			zap.String("content-language", clang),
			zap.String("content-charset", cc),
			zap.String("content-location", (&loc).String()),
			zap.Int("status", status))
	}
	return
}

// notAcceptable is responsible for responding to the user agent with a
// 406 HTTP status code, along with a representation describing the available
// representations and their metadata.
func (n Negotiator) notAcceptable(
	ctx negotiator.NegotiationContext, reps ...representation.Representation) (err error) {

	defer func() {
		if err == nil {
			n.scope.Counter(scopeNameProactiveNotAcceptableCounter).Inc(1)
		}
	}()

	if !n.notAcceptableRepresentation {
		status := http.StatusNotAcceptable
		ctx.ResponseWriter.WriteHeader(status)
		n.logger.Info("not acceptable", zap.Int("status", status))
		return nil
	}

	// perform negotiation on representation.
	var (
		lists  []representation.Representation
		chosen representation.Representation
	)
	for _, c := range n.representationConstructors {
		lists = append(lists, c(reps...))
	}
	if chosen, err = n.chooser.Choose(ctx.Request, lists...); err != nil {
		return
	}
	n.logger.Debug("completed choosing on not acceptable response",
		zap.Int("representation-count", len(lists)))

	// choose default when there are no matches.
	if chosen == nil {
		chosen = n.defaultRepresentationConstructor(reps...)
		n.logger.Debug("chose default representation for not acceptable response")
	}

	// serialize.
	var b []byte
	if b, err = chosen.Bytes(); err != nil {
		return
	}

	// respond.
	var (
		clen   = len(b)
		ct     = chosen.ContentType()
		ce     = strings.Join(chosen.ContentEncoding(), ",")
		clang  = chosen.ContentLanguage()
		cc     = chosen.ContentCharset()
		status = http.StatusNotAcceptable
	)
	ctx.ResponseWriter.Header().Add("Content-Length", strconv.Itoa(clen))
	ctx.ResponseWriter.Header().Add("Content-Type", ct)
	ctx.ResponseWriter.Header().Add("Content-Encoding", ce)
	ctx.ResponseWriter.Header().Add("Content-Language", clang)
	ctx.ResponseWriter.Header().Add("Content-Charset", cc)
	ctx.ResponseWriter.WriteHeader(status)
	if _, err = ctx.ResponseWriter.Write(b); err == nil {
		n.logger.Info("not acceptable",
			zap.Int("content-length", clen),
			zap.String("content-type", ct),
			zap.String("content-encoding", ce),
			zap.String("content-language", clang),
			zap.String("content-charset", cc),
			zap.Int("status", status))
	}
	return
}
