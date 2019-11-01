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
	"strconv"
	"strings"

	"github.com/freerware/negotiator"
	"github.com/freerware/negotiator/internal/header"
	"github.com/freerware/negotiator/internal/representation/json"
	"github.com/freerware/negotiator/representation"
	"github.com/uber-go/tally"
	"go.uber.org/zap"
)

// Defines the default transparent negotiator.
//
// The default configuration is as follows:
//
// ➣ The algorithm used for serving choice responses is the RVSA 1.0 algorithm.
//
// ➣ For guess-small responses, the choice response can be no more than 50
// bytes larger than the list response.
//
// ➣ The representation for list responses utilizes the JSON (application/json)
// media type.
//
// ➣ No more than 10 representations can be used in the negotiation process.
var (
	// Default is the default transparent negotiator.
	Default = New()
)

// Errors that can be thrown from options.
var (
	// ErrVariantListSizeExceeded represents an error encountered when the
	// provided variant list exceeds the maximum allowed.pr
	ErrVariantListSizeExceeded = errors.New("number of representations exceeds the maximum")
)

// Defines the tags and scope names for transparent negotiator metrics.
var (
	scopeTagTransparent              = map[string]string{"negotiator": "transparent"}
	scopeNameTransparentTimer        = "negotiate"
	scopeNameTransparentErrorCounter = "negotiate.error"
)

// Negotiator represents the negotiator responsible for
// performing transparent negotiation.
type Negotiator struct {
	maximumVariantListSize        int
	listRepresentationConstructor representation.ListConstructor
	chooser                       representation.Chooser
	guessSmallThreshold           int
	logger                        *zap.Logger
	scope                         tally.Scope
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
		Logger:                        zap.NewNop(),
		Scope:                         tally.NoopScope,
	}
	// apply options.
	for _, opt := range options {
		opt(&o)
	}
	n := Negotiator{
		maximumVariantListSize:        o.MaximumVariantListSize,
		listRepresentationConstructor: o.ListRepresentationConstructor,
		chooser:                       o.Chooser,
		guessSmallThreshold:           o.GuessSmallThreshold,
		logger:                        o.Logger,
		scope:                         o.Scope.Tagged(scopeTagTransparent),
	}
	n.logger.Debug("negotiator configuration",
		zap.String("type", "transparent"),
		zap.Int("maximum-variant-list-size", n.maximumVariantListSize),
		zap.Int("guess-small-threshold", n.guessSmallThreshold))
	return n
}

// Negotiate performs transparent content negotiation with the representations
// provided.
func (n Negotiator) Negotiate(
	ctx negotiator.NegotiationContext, reps ...representation.Representation) (err error) {

	defer n.scope.Timer(scopeNameTransparentTimer).Start().Stop()
	defer func() {
		if err != nil {
			n.scope.Counter(scopeNameTransparentErrorCounter).Inc(1)
		}
	}()

	if len(reps) > n.maximumVariantListSize {
		return ErrVariantListSizeExceeded
	}

	var negotiate header.Negotiate
	if negotiate, err = header.NewNegotiate(ctx.Request.Header["Negotiate"]); err != nil {
		return
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

	var rep representation.Representation
	if rep, err = n.chooser.Choose(ctx.Request, reps...); err != nil {
		return
	}

	if rep == nil {
		return n.listResponse(ctx, reps...)
	}

	if negotiate.Contains(header.NegotiateDirectiveGuessSmall) {
		list := n.listRepresentationConstructor(reps...)
		var listBytes []byte
		if listBytes, err = list.Bytes(); err != nil {
			return
		}

		var choiceBytes []byte
		if choiceBytes, err = rep.Bytes(); err != nil {
			return
		}

		less := len(choiceBytes) < len(listBytes)
		diff := math.Abs(float64(len(listBytes) - len(choiceBytes)))
		if !less && diff > float64(n.guessSmallThreshold) {
			n.logger.Debug("choice response is not smaller or not much larger than the list response",
				zap.Int("choice-response-size", len(choiceBytes)),
				zap.Int("list-response-size", len(listBytes)),
				zap.Int("guess-small-threshold", n.guessSmallThreshold))
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
		n.logger.Debug("variant resource is not a neighbor of the negotiable resource",
			zap.String("resource-url", resourceURL),
			zap.String("neighbor-url", neighborURL))
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
	t, err := header.NewTCN([]string{header.ResponseTypeList.String()})
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
	var (
		clen       = len(b)
		ct         = list.ContentType()
		ce         = strings.Join(list.ContentEncoding(), ",")
		clang      = list.ContentLanguage()
		cc         = list.ContentCharset()
		alternates = a.ValuesAsString()
		tcn        = t.ValuesAsString()
		status     = http.StatusMultipleChoices
	)
	ctx.ResponseWriter.Header().Add("Alternates", alternates)
	ctx.ResponseWriter.Header().Add("TCN", tcn)
	ctx.ResponseWriter.Header().Add("Content-Length", strconv.Itoa(clen))
	ctx.ResponseWriter.Header().Add("Content-Type", ct)
	ctx.ResponseWriter.Header().Add("Content-Encoding", ce)
	ctx.ResponseWriter.Header().Add("Content-Language", clang)
	ctx.ResponseWriter.Header().Add("Content-Charset", cc)
	ctx.ResponseWriter.WriteHeader(status)
	if _, err = ctx.ResponseWriter.Write(b); err == nil {
		n.logger.Info("list response",
			zap.Int("content-length", clen),
			zap.String("content-type", ct),
			zap.String("content-encoding", ce),
			zap.String("content-language", clang),
			zap.String("content-charset", cc),
			zap.Int("status", status),
			zap.String("alternates", alternates),
			zap.String("tcn", tcn))
	}
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
	t, err := header.NewTCN([]string{header.ResponseTypeChoice.String()})
	if err != nil {
		return err
	}

	// serialize.
	b, err := rep.Bytes()
	if err != nil {
		return err
	}

	// respond.
	var (
		clen       = len(b)
		ct         = rep.ContentType()
		ce         = strings.Join(rep.ContentEncoding(), ",")
		clang      = rep.ContentLanguage()
		cc         = rep.ContentCharset()
		alternates = a.ValuesAsString()
		tcn        = t.ValuesAsString()
		loc        = rep.ContentLocation()
		status     = http.StatusOK
	)
	ctx.ResponseWriter.Header().Add("Alternates", alternates)
	ctx.ResponseWriter.Header().Add("TCN", tcn)
	ctx.ResponseWriter.Header().Add("Content-Location", (&loc).String())
	ctx.ResponseWriter.Header().Add("Content-Length", strconv.Itoa(clen))
	ctx.ResponseWriter.Header().Add("Content-Type", ct)
	ctx.ResponseWriter.Header().Add("Content-Encoding", ce)
	ctx.ResponseWriter.Header().Add("Content-Language", clang)
	ctx.ResponseWriter.Header().Add("Content-Charset", cc)
	ctx.ResponseWriter.WriteHeader(status)
	if _, err = ctx.ResponseWriter.Write(b); err == nil {
		n.logger.Info("choice response",
			zap.Int("content-length", clen),
			zap.String("content-type", ct),
			zap.String("content-encoding", ce),
			zap.String("content-language", clang),
			zap.String("content-charset", cc),
			zap.Int("status", status),
			zap.String("alternates", alternates),
			zap.String("tcn", tcn),
			zap.String("content-location", (&loc).String()))
	}
	return err
}
