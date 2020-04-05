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
	"github.com/freerware/negotiator/representation"
	"github.com/uber-go/tally"
	"go.uber.org/zap"
)

// Options represents the configuration options for reactive
// (agent-driven) content negotiation.
type Options struct {
	RepresentationConstructor representation.ListConstructor
	Logger                    *zap.Logger
	Scope                     tally.Scope
}

// Option represents a configurable option for reactive
// (agent-driven) content negotiation.
type Option = func(*Options)

// Options that can be used to configure and extend reactive negotiators.
var (
	// Representation defines representation to utilize when returning
	// responses to user agents engaging in reactive negotiation.
	Representation = func(constructor representation.ListConstructor) Option {
		return func(o *Options) {
			o.RepresentationConstructor = constructor
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
