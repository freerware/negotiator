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

package transparent

import (
	"errors"
	"math"
	"net/http"
	"strings"

	"github.com/freerware/negotiator"
	"github.com/freerware/negotiator/internal/header"
	"github.com/freerware/negotiator/internal/representation/json"
	"github.com/freerware/negotiator/representation"
)

var (
	// Default is the default transparent negotiator.
	Default = New()

	// ErrVariantListSizeExceeded represents an error encountered when the
	// provided variant list exceeds the maximum allowed.pr
	ErrVariantListSizeExceeded = errors.New("number of representations exceeds the maximum")
)

// Negotiator represents the negotiator responsible for
// performing transparent negotiation.
type Negotiator struct {
	maximumVariantListSize        int
	listRepresentationConstructor representation.Constructor
	chooser                       representation.Chooser
	guessSmallThreshold           int
}

// New constructs a negotiatior capable of performing transparent
// negotiation with the options provided.
func New(options ...Option) Negotiator {
	// set defaults.
	o := Options{
		MaximumVariantListSize:        10,
		ListRepresentationConstructor: json.List,
		Chooser:                       RVSA1(),
		GuessSmallThreshold:           50,
	}
	// apply options.
	for _, opt := range options {
		opt(&o)
	}
	return Negotiator{
		maximumVariantListSize:        o.MaximumVariantListSize,
		listRepresentationConstructor: o.ListRepresentationConstructor,
		chooser:                       o.Chooser,
		guessSmallThreshold:           o.GuessSmallThreshold,
	}
}

// Negotiate performs transparent content negotiation with the representations
// provided.
func (n Negotiator) Negotiate(
	ctx negotiator.NegotiationContext, reps ...representation.Representation) error {

	if len(reps) > n.maximumVariantListSize {
		return ErrVariantListSizeExceeded
	}

	negotiate, err := header.NewNegotiate(ctx.Request.Header["Negotiate"])
	if err != nil {
		return err
	}

	// determine when the user agent wants the server to choose the best
	// variant on it's behalf.
	shouldChoose :=
		negotiate.Contains(header.NegotiateDirective("*")) ||
			negotiate.ContainsRVSA("1.0") ||
			negotiate.Contains(header.NegotiateDirectiveGuessSmall)

	if !shouldChoose {
		return n.listResponse(ctx, reps...)
	}

	rep, err := n.chooser.Choose(ctx.Request, reps...)
	if err != nil {
		return err
	}

	if rep == nil {
		return n.listResponse(ctx, reps...)
	}

	if negotiate.Contains(header.NegotiateDirectiveGuessSmall) {
		list := n.listRepresentationConstructor(reps...)
		listBytes, err := list.Bytes()
		if err != nil {
			return err
		}

		choiceBytes, err := rep.Bytes()
		if err != nil {
			return err
		}

		less := len(choiceBytes) < len(listBytes)
		diff := math.Abs(float64(len(listBytes) - len(choiceBytes)))
		if !less && diff > float64(n.guessSmallThreshold) {
			return n.listResponse(ctx, reps...)
		}
	}

	//https://tools.ietf.org/html/rfc2296#section-3.5
	loc := rep.ContentLocation()
	neighborURL := strings.TrimRight(loc.String(), "/")
	resourceURL := strings.TrimRight(ctx.Request.URL.String(), "/")
	nLastSlash := strings.LastIndex(neighborURL, "/")
	rLastSlash := strings.LastIndex(resourceURL, "/")
	if nLastSlash == -1 || rLastSlash == -1 || nLastSlash != rLastSlash {
		return n.listResponse(ctx, reps...)
	}
	return n.choiceResponse(ctx, reps, rep)
}

// listResponse is responsible for responding to the user agent with a 'list'
// response, including the representation describing the available
// representations and their metadata.
func (n Negotiator) listResponse(
	ctx negotiator.NegotiationContext, reps ...representation.Representation) error {
	a, err := header.NewAlternates(reps[0], reps...)
	if err != nil {
		return err
	}
	tcn, err := header.NewTCN([]string{header.ResponseTypeList.String()})
	if err != nil {
		return err
	}

	// construct representation.
	list := n.listRepresentationConstructor(reps...)

	// serialize.
	b, err := list.Bytes()
	if err != nil {
		return err
	}

	// respond.
	ctx.ResponseWriter.Header().Add("Alternates", a.ValuesAsString())
	ctx.ResponseWriter.Header().Add("TCN", tcn.ValuesAsString())
	ctx.ResponseWriter.Header().Add("Content-Length", string(len(b)))
	ctx.ResponseWriter.Header().Add("Content-Type", list.ContentType())
	ctx.ResponseWriter.Header().Add("Content-Encoding", strings.Join(list.ContentEncoding(), ","))
	ctx.ResponseWriter.Header().Add("Content-Language", list.ContentLanguage())
	ctx.ResponseWriter.Header().Add("Content-Charset", list.ContentCharset())
	ctx.ResponseWriter.WriteHeader(http.StatusMultipleChoices)
	_, err = ctx.ResponseWriter.Write(b)
	return err
}

// choiceResponse is responsible for responding to the user agent with a
// 'choice' response, including the representation chosen by the remote
// variant selection algorithm.
func (n Negotiator) choiceResponse(
	ctx negotiator.NegotiationContext,
	reps []representation.Representation,
	rep representation.Representation) error {
	// If a response from a transparently negotiable resource includes an
	// Alternates header, this header MUST contain the complete variant list
	// bound to the negotiable resource. Responses from resources which do not
	// support transparent content negotiation MAY also use Alternates headers.
	//
	// https://tools.ietf.org/html/rfc2295#section-8.3
	a, err := header.NewAlternates(nil, reps...)
	if err != nil {
		return err
	}
	tcn, err := header.NewTCN([]string{header.ResponseTypeChoice.String()})
	if err != nil {
		return err
	}
	loc := rep.ContentLocation()

	// serialize.
	b, err := rep.Bytes()
	if err != nil {
		return err
	}

	// respond.
	ctx.ResponseWriter.Header().Add("Alternates", a.ValuesAsString())
	ctx.ResponseWriter.Header().Add("TCN", tcn.ValuesAsString())
	ctx.ResponseWriter.Header().Add("Content-Location", (&loc).String())
	ctx.ResponseWriter.Header().Add("Content-Length", string(len(b)))
	ctx.ResponseWriter.Header().Add("Content-Type", rep.ContentType())
	ctx.ResponseWriter.Header().Add("Content-Encoding", strings.Join(rep.ContentEncoding(), ","))
	ctx.ResponseWriter.Header().Add("Content-Language", rep.ContentLanguage())
	ctx.ResponseWriter.Header().Add("Content-Charset", rep.ContentCharset())
	ctx.ResponseWriter.WriteHeader(http.StatusOK)
	_, err = ctx.ResponseWriter.Write(b)
	return err
}
