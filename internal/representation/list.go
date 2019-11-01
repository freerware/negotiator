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
	rep "github.com/freerware/negotiator/representation"
)

// RepresentationMetadata is the metadata about each representation in the
// representation list.
type RepresentationMetadata struct {
	ContentType     string   `json:"contentType,omitempty"`
	ContentLanguage string   `json:"contentLanguage,omitempty"`
	ContentEncoding []string `json:"contentEncoding,omitempty"`
	ContentLocation string   `json:"contentLocation,omitempty"`
	ContentCharset  string   `json:"contentCharset,omitempty"`
	ContentFeatures []string `json:"contentFeatures,omitempty"`
	SourceQuality   float32  `json:"sourceQuality"`
}

// List represents a representation containing a list of descriptions of representations
// for a particular resource.
type List struct {
	rep.Base

	Representations []RepresentationMetadata `json:"representations"`
}

// SetRepresentations modifies the represention list within the list representation.
func (l *List) SetRepresentations(reps ...rep.Representation) {
	for _, rep := range reps {
		loc := rep.ContentLocation()
		l.Representations = append(l.Representations, RepresentationMetadata{
			ContentType:     rep.ContentType(),
			ContentLanguage: rep.ContentLanguage(),
			ContentEncoding: rep.ContentEncoding(),
			ContentLocation: (&loc).String(),
			ContentCharset:  rep.ContentCharset(),
			ContentFeatures: rep.ContentFeatures(),
		})
	}
}

// Bytes retrieves the serialized form of the list representation.
func (l List) Bytes() ([]byte, error) {
	return l.Base.Bytes(&l)
}

// FromBytes constructs the list representation from its serialized form.
func (l List) FromBytes(b []byte) error {
	return l.Base.FromBytes(b, &l)
}
