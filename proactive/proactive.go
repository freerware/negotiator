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
	"strings"

	"github.com/freerware/negotiator"
	"github.com/freerware/negotiator/internal/header"
	"github.com/freerware/negotiator/internal/representation/json"
	"github.com/freerware/negotiator/internal/representation/xml"
	"github.com/freerware/negotiator/internal/representation/yaml"
	"github.com/freerware/negotiator/representation"
)

var (
	// Default is the default proactive negotiator.
	Default = New(DisableStrictMode())
)

// Negotiator represents the negotiator responsible for performing
// proactive (server-driven) negotiation.
type Negotiator struct {
	strictAccept                     bool
	strictAcceptLanguage             bool
	strictAcceptCharset              bool
	notAcceptableRepresentation      bool
	defaultRepresentationConstructor representation.Constructor
	representationConstructors       []representation.Constructor
	chooser                          representation.Chooser
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
		DefaultRepresentationConstructor: json.List,
		Chooser:                          ApacheHTTPD(),
		RepresentationConstructors: []representation.Constructor{
			json.List,
			xml.List,
			yaml.List,
		},
	}
	// apply options.
	for _, opt := range options {
		opt(&o)
	}
	return Negotiator{
		strictAccept:                     o.StrictAccept,
		strictAcceptLanguage:             o.StrictAcceptLanguage,
		strictAcceptCharset:              o.StrictAcceptCharset,
		notAcceptableRepresentation:      o.NotAcceptableRepresentation,
		defaultRepresentationConstructor: o.DefaultRepresentationConstructor,
		representationConstructors:       o.RepresentationConstructors,
		chooser:                          o.Chooser,
	}
}

// Negotiate performs proactive (server-driven) content negotiation with the
// representations provided.
func (n Negotiator) Negotiate(
	ctx negotiator.NegotiationContext, reps ...representation.Representation) error {

	if len(reps) == 0 {
		ctx.ResponseWriter.WriteHeader(http.StatusNoContent)
		return nil
	}

	/*
		// retrieve headers.
		_, hasAccept := ctx.Request.Header["Accept"]
		_, hasAcceptLanguage := ctx.Request.Header["Accept-Language"]
		_, hasAcceptCharset := ctx.Request.Header["Accept-Charset"]

		// check strict mode for 'Accept'.
		if !hasAccept && n.strictAccept {
			return n.notAcceptable(ctx, reps...)
		}

		// check strict mode for 'Accept-Language'.
		if !hasAcceptLanguage && n.strictAcceptLanguage {
			return n.notAcceptable(ctx, reps...)
		}

		// check strcit mode for 'Accept-Çharset'.
		if !hasAcceptCharset && n.strictAcceptCharset {
			return n.notAcceptable(ctx, reps...)
		}
	*/

	var (
		accept                  = header.DefaultAccept
		acceptLanguage          = header.DefaultAcceptLanguage
		acceptCharset           = header.DefaultAcceptCharset
		ac, alc, acc, hasHeader bool
		headerValues            []string
		err                     error
	)
	if headerValues, hasHeader = ctx.Request.Header["Accept"]; hasHeader {
		if accept, err = header.NewAccept(headerValues); err != nil {
			return err
		}
	}
	if headerValues, hasHeader = ctx.Request.Header["Accept-Language"]; hasHeader {
		if acceptLanguage, err = header.NewAcceptLanguage(headerValues); err != nil {
			return err
		}
	}
	if headerValues, hasHeader = ctx.Request.Header["Accept-Charset"]; hasHeader {
		if acceptCharset, err = header.NewAcceptCharset(headerValues); err != nil {
			return err
		}
	}
	for _, r := range reps {
		if ac, err = accept.Compatible(r.ContentType()); err != nil {
			return err
		}
		if alc, err = acceptLanguage.Compatible(r.ContentLanguage()); err != nil {
			return err
		}
		if acc, err = acceptCharset.Compatible(r.ContentCharset()); err != nil {
			return err
		}
	}
	if !ac && n.strictAccept {
		return n.notAcceptable(ctx, reps...)
	}
	if !alc && n.strictAcceptLanguage {
		return n.notAcceptable(ctx, reps...)
	}
	if !acc && n.strictAcceptCharset {
		return n.notAcceptable(ctx, reps...)
	}

	// choose 'best' representation.
	rep, err := n.chooser.Choose(ctx.Request, reps...)
	if err != nil {
		return err
	}

	if rep == nil {
		return n.notAcceptable(ctx, reps...)
	}
	return n.acceptable(ctx, rep)
}

// acceptable is responsible for responding to the user agent with the
// representation chosen by the server-side algorithm.
func (n Negotiator) acceptable(
	ctx negotiator.NegotiationContext, rep representation.Representation) error {

	// serialize.
	bytes, err := rep.Bytes()
	if err != nil {
		return err
	}

	// respond.
	statusCode := http.StatusOK
	if ctx.IsCreation {
		statusCode = http.StatusCreated
	}
	loc := rep.ContentLocation()
	ctx.ResponseWriter.Header().Add("Content-Location", (&loc).String())
	ctx.ResponseWriter.Header().Add("Content-Length", string(len(bytes)))
	ctx.ResponseWriter.Header().Add("Content-Type", rep.ContentType())
	ctx.ResponseWriter.Header().Add("Content-Encoding", strings.Join(rep.ContentEncoding(), ","))
	ctx.ResponseWriter.Header().Add("Content-Language", rep.ContentLanguage())
	ctx.ResponseWriter.Header().Add("Content-Charset", rep.ContentCharset())
	ctx.ResponseWriter.WriteHeader(statusCode)
	_, err = ctx.ResponseWriter.Write(bytes)
	return err
}

// notAcceptable is responsible for responding to the user agent with a
// 406 HTTP status code, along with a representation describing the available
// representations and their metadata.
func (n Negotiator) notAcceptable(
	ctx negotiator.NegotiationContext, reps ...representation.Representation) error {

	if !n.notAcceptableRepresentation {
		ctx.ResponseWriter.WriteHeader(http.StatusNotAcceptable)
		return nil
	}

	// perform negotiation on representation.
	var lists []representation.Representation
	for _, c := range n.representationConstructors {
		lists = append(lists, c(reps...))
	}
	chosen, err := n.chooser.Choose(ctx.Request, lists...)
	if err != nil {
		return err
	}

	// choose default when there are no matches.
	if chosen == nil {
		chosen = n.defaultRepresentationConstructor(reps...)
	}

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
	ctx.ResponseWriter.WriteHeader(http.StatusNotAcceptable)
	_, err = ctx.ResponseWriter.Write(b)
	return err
}
