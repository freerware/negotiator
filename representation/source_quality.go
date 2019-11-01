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

package representation

// Source quality guidelines as defined in RFC2295 Section 5.3.
const (
	// SourceQualityPerfect indicates that the representation is perfect
	// quality with no degradation.
	SourceQualityPerfect float32 = 1.0

	// SourceQualityNearlyPerfect indicates the threshold of noticeable loss
	// of quality for the representation.
	SourceQualityNearlyPerfect float32 = 0.9

	// SourceQualityAcceptable indicates that the representation has
	// noticeable but acceptable quality reduction.
	SourceQualityAcceptable float32 = 0.8

	// SourceQualityBarelyAcceptable indicates that the representation
	// has barely acceptable quality.
	SourceQualityBarelyAcceptable float32 = 0.5

	// SourceQualitySeverelyDegraded indicates that the representation
	// has severely degraded quality.
	SourceQualitySeverelyDegraded float32 = 0.3

	// SourceQualityCompletelyDegraded indicates that the representation
	// has completed degraded quality.
	SourceQualityCompletelyDegraded float32 = 0.0
)
