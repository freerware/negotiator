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

// Package reactive implements reactive content negotiation as defined in
// RFC7231 Section 3.4.2.
//
// Construction
//
// For out of the box reactive negotiation support, use reactive.Default,
// which is the default reactive negotiator.
//	//retrieves the default reactive negotiator.
//	p := reactive.Default
// In situations where more customization is required, use the reactive.New
// constructor function and specify options as arguments.
//	// constructs a reactive negotiator with the provided options.
//	p := reactive.New(
//		reactive.Logger(l),
//	)
//
// See Also
//
// ➣ https://tools.ietf.org/html/rfc7231#section-3.4.2
package reactive
