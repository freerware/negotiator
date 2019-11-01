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
	"github.com/freerware/negotiator/representation"
	"github.com/uber-go/tally"
	"go.uber.org/zap"
)

// Options represents the configuration options for transparent
// content negotiation.
type Options struct {
	MaximumVariantListSize        int
	Chooser                       representation.Chooser
	ListRepresentationConstructor representation.ListConstructor
	Logger                        *zap.Logger
	Scope                         tally.Scope
	GuessSmallThreshold           int
}

// Option represents a configurable option for transparent
// content negotiation.
type Option func(*Options)

// Options that can be used to configure and extend transparent negotiators.
var (
	// MaximumVariantListSize specifies the maximum allowable size of the
	// variant list used within transparent negotiation. Variant lists larger
	// than this number will result in an error. If a value less than 1 is
	// provided, it will be set to 1.
	MaximumVariantListSize = func(size int) Option {
		if size < 1 {
			size = 1
		}
		return func(o *Options) {
			o.MaximumVariantListSize = size
		}
	}

	// RVSA defines the remove variant selection algorithm (RVSA) to leverage
	// for transparant negotiation.
	RVSA = func(c representation.Chooser) Option {
		return func(o *Options) {
			o.Chooser = c
		}
	}

	// ListRepresentation defines the representation to utilize when
	// returning list responses.
	ListRepresentation = func(c representation.ListConstructor) Option {
		return func(o *Options) {
			o.ListRepresentationConstructor = c
		}
	}

	// GuessSmallThreshold specifies the threshold in bytes that the choice
	// response can exceed the list response in a 'guess-small' request.
	GuessSmallThreshold = func(threshold int) Option {
		return func(o *Options) {
			o.GuessSmallThreshold = threshold
		}
	}

	// Logger specifies the logger for the reactive negotiator.
	Logger = func(l *zap.Logger) Option {
		return func(o *Options) {
			o.Logger = l
		}
	}

	// Scope specifies the metric scope to leverage for the reactive
	// negotiator.
	Scope = func(s tally.Scope) Option {
		return func(o *Options) {
			o.Scope = s
		}
	}
)
