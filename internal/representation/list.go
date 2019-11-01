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
	"bytes"
	"compress/flate"
	"compress/gzip"
	"compress/zlib"
	"encoding/json"
	"encoding/xml"
	"errors"
	"io"
	"net/url"
	"strings"

	rep "github.com/freerware/negotiator/representation"
	"gopkg.in/yaml.v2"
)

var (
	// ErrUnsupportedContentEncoding indicates an error that occurs when the list
	// represention is given an unsupported content encoding.
	ErrUnsupportedContentEncoding = errors.New("represention content encoding is not supported")

	// ErrUnsupportedContentType indicates an error that occurs when the list
	// represention is given an unsupported content type.
	ErrUnsupportedContentType = errors.New("representation content type not supported")
)

// Representation is the core representation.
type Representation struct {
	encoding      []string
	mediaType     string
	charset       string
	language      string
	location      url.URL
	sourceQuality float32
	features      []string
}

// ContentType retrieves the content type of the representation.
func (r Representation) ContentType() string { return r.mediaType }

// SetContentType modifies the content type of the representation.
func (r *Representation) SetContentType(ct string) { r.mediaType = ct }

// ContentLanguage retrieves the content language of the representation.
func (r Representation) ContentLanguage() string { return r.language }

// SetContentLanguage modifies the content language of the representation.
func (r *Representation) SetContentLanguage(cl string) { r.language = cl }

// ContentEncoding retrieves the content encoding of the representation.
func (r Representation) ContentEncoding() []string { return r.encoding }

// SetContentEncoding modifies the content encoding of the representation.
func (r *Representation) SetContentEncoding(ce []string) { r.encoding = ce }

// ContentCharset retrieves the content charset of the representation.
func (r Representation) ContentCharset() string { return r.charset }

// SetContentCharset modifies the content charset of the representation.
func (r *Representation) SetContentCharset(cc string) { r.charset = cc }

// ContentLocation retrieves the content location of the representation.
func (r Representation) ContentLocation() url.URL { return r.location }

// SetContentLocation modifies the content location of the representation.
func (r *Representation) SetContentLocation(l url.URL) { r.location = l }

// ContentFeatures retrieves the content features of the representation.
func (r Representation) ContentFeatures() []string { return r.features }

// SetContentFeatures retrieves the content features of the representation.
func (r *Representation) SetContentFeatures(cf []string) { r.features = cf }

// SourceQuality retrieves the source quality of the representation.
func (r Representation) SourceQuality() float32 { return rep.SourceQualityPerfect }

// SetSourceQuality modifies the source quality of the representation.
func (r *Representation) SetSourceQuality(sq float32) { r.sourceQuality = sq }

// List represents a representation containing a list of descriptions of representations
// for a particular resource.
type List struct {
	Representation

	Representations []rep.Representation
}

// SetRepresentations modifies the represention list within the list representation.
func (l *List) SetRepresentations(reps ...rep.Representation) {
	l.Representations = reps
}

// Bytes retrieves the serialized form of the list representation.
func (l List) Bytes() ([]byte, error) {
	supportedMediaTypes := map[string]marshaller{
		"application/json": json.Marshal,
		"application/xml":  xml.Marshal,
		"application/yaml": yaml.Marshal,
		"text/yaml":        yaml.Marshal,
	}

	var buf bytes.Buffer
	supportedEncodings := map[string]writerConstructor{
		"gzip":       newGzip,
		"x-gzip":     newGzip,
		"compress":   newCompress,
		"x-compress": newCompress,
		"deflate":    newDeflate,
	}

	if _, ok := supportedMediaTypes[strings.ToLower(l.ContentType())]; !ok {
		return []byte{}, ErrUnsupportedContentType
	}

	// serialize.
	b, err := supportedMediaTypes[strings.ToLower(l.ContentType())](l)
	if err != nil {
		return b, err
	}

	if len(l.ContentEncoding()) < 1 {
		return b, nil
	}

	// encode.
	var w io.WriteCloser = &closeableBuffer{&buf}
	for _, e := range l.ContentEncoding() {
		if _, ok := supportedEncodings[strings.ToLower(e)]; !ok {
			return []byte{}, ErrUnsupportedContentEncoding
		}
		if w, err = supportedEncodings[strings.ToLower(e)](w); err != nil {
			return []byte{}, err
		}
		defer w.Close()
	}
	if _, err = w.Write(b); err != nil {
		return []byte{}, err
	}

	// return final bytes.
	return buf.Bytes(), nil
}

// writerConstructor represents a constructor for closeable writers.
type writerConstructor func(io.WriteCloser) (io.WriteCloser, error)

var (
	// newGzip is the constructor for the gzip closeable writer.
	newGzip writerConstructor = func(w io.WriteCloser) (io.WriteCloser, error) {
		return gzip.NewWriter(w), nil
	}

	// newCompress is the constructor for the compress closeable writer.
	newCompress writerConstructor = func(w io.WriteCloser) (io.WriteCloser, error) {
		return zlib.NewWriter(w), nil
	}

	// newDeflate is the constructor for the default closeable writer.
	newDeflate writerConstructor = func(w io.WriteCloser) (io.WriteCloser, error) {
		return flate.NewWriter(w, 1)
	}
)

// closeableBuffer represents a closeable buffer.
type closeableBuffer struct {
	buf *bytes.Buffer
}

// Close closes the buffer.
func (cb closeableBuffer) Close() error {
	return nil
}

// Write writes the provided bytes to the buffer.
func (cb closeableBuffer) Write(b []byte) (int, error) {
	return cb.buf.Write(b)
}

// marshaller represents a marshaling function.
type marshaller func(interface{}) ([]byte, error)
