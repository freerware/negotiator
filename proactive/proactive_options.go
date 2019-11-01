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
	"github.com/freerware/negotiator/representation"
	"github.com/uber-go/tally"
	"go.uber.org/zap"
)

// Options represents the configuration options for proactive
// (server-driven) content negotiation.
type Options struct {
	StrictAccept                     bool
	StrictAcceptLanguage             bool
	StrictAcceptCharset              bool
	NotAcceptableRepresentation      bool
	DefaultRepresentationConstructor representation.ListConstructor
	RepresentationConstructors       []representation.ListConstructor
	Chooser                          representation.Chooser
	Logger                           *zap.Logger
	Scope                            tally.Scope
}

// Option represents a configurable option for proactive
// (server-driven) content negotiation.
type Option func(*Options)

// Options that can be used to configure and extend proactive negotiators.
var (
	// DisableStrictAccept deactivates strict mode for the Accept header,
	// meaning that a 406 HTTP status code is not returned if none of the
	// representations have a media type that is acceptable.
	DisableStrictAccept = func() Option {
		return func(o *Options) {
			o.StrictAccept = false
		}
	}

	// DisableStrictAcceptLanguage deactivates strict mode for the
	// Accept-Language header, meaning that a 406 HTTP status code is not
	// returned if none of the representations have a language that is
	// acceptable.
	DisableStrictAcceptLanguage = func() Option {
		return func(o *Options) {
			o.StrictAcceptLanguage = false
		}
	}

	// DisableStrictAcceptCharset deactivates strict mode for the
	// Accept-Charset header, meaning that a 406 HTTP status code is not
	// returned if none of the representations have a charset that is
	// acceptable.
	DisableStrictAcceptCharset = func() Option {
		return func(o *Options) {
			o.StrictAcceptCharset = false
		}
	}

	// DisableStrictMode deactivates strict mode for all Accept-* headers,
	// meaning that a 406 HTTP status code will not be returned if none of the
	// representations are acceptable for any of the headers.
	DisableStrictMode = func() Option {
		return func(o *Options) {
			o.StrictAccept = false
			o.StrictAcceptLanguage = false
			o.StrictAcceptCharset = false
		}
	}

	// Algorithm defines the algorithm to leverage for proactive negotiation.
	Algorithm = func(c representation.Chooser) Option {
		return func(o *Options) {
			o.Chooser = c
		}
	}

	// DisableNotAcceptableRepresentation deactivates functionality that
	// provides a representation in body of responses having a 406 HTTP
	// status code.
	DisableNotAcceptableRepresentation = func() Option {
		return func(o *Options) {
			o.NotAcceptableRepresentation = false
		}
	}

	// DefaultRepresentation defines the representation to utilize when
	// returning responses with the 406 HTTP status code.
	DefaultRepresentation = func(constructor representation.ListConstructor) Option {
		return func(o *Options) {
			o.DefaultRepresentationConstructor = constructor
		}
	}

	// Representations defines the set of representations to consider when
	// returning responses with the 406 HTTP status code. These the algorithm
	// that the proactive negotiator is configured with is utilized to select
	// the 'best' representation.
	Representations = func(constructors ...representation.ListConstructor) Option {
		return func(o *Options) {
			o.RepresentationConstructors = constructors
		}
	}

	// Logger specifies the logger for the proactive negotiator.
	Logger = func(l *zap.Logger) Option {
		return func(o *Options) {
			o.Logger = l
		}
	}

	// Scope specifies the metric scope to leverage for the proactive
	// negotiator.
	Scope = func(s tally.Scope) Option {
		return func(o *Options) {
			o.Scope = s
		}
	}
)
