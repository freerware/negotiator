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

import (
	"net/url"
)

// Representation is an HTTP resource representation.
//
// For the purposes of HTTP, a "representation" is information that is
// intended to reflect a past, current, or desired state of a given
// resource, in a format that can be readily communicated via the
// protocol, and that consists of a set of representation metadata and a
// potentially unbounded stream of representation data.
type Representation interface {
	ContentLocation() url.URL
	ContentType() string
	ContentEncoding() []string
	ContentCharset() string
	ContentLanguage() string
	ContentFeatures() []string
	SourceQuality() float32
	Bytes() ([]byte, error)
	FromBytes([]byte) error
}
