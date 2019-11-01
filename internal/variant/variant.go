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
package variant

import (
	"github.com/freerware/negotiator/internal/header"
	"github.com/freerware/negotiator/representation"
)

// Variant describes a single variation (representation) that can be
// served by a resource. Contains all of the parameters leveraged in
// server-side selection algorithms.
type Variant struct {
	Representation        representation.Representation
	SourceQualityValue    header.QualityValue
	MediaTypeQualityValue header.QualityValue
	CharsetQualityValue   header.QualityValue
	LanguageQualityValue  header.QualityValue
	EncodingQualityValue  header.QualityValue
	FeatureQualityValue   header.QualityValue
	IsDefinite            bool
	LanguageOrderScore    int
}
