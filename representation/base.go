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
	"encoding/json"
	"encoding/xml"
	"errors"
	"io"
	"io/ioutil"
	"net/url"
	"strings"

	"gopkg.in/yaml.v2"
)

// Errors that can be encountered when serializing and deserializing representations.
var (
	// ErrUnsupportedContentEncoding indicates an error that occurs when the
	// represention is given an unsupported content encoding.
	ErrUnsupportedContentEncoding = errors.New("representation content encoding is not supported")

	// ErrUnsupportedContentType indicates an error that occurs when the
	// represention is given an unsupported content type.
	ErrUnsupportedContentType = errors.New("representation content type is not supported")
)

var (
	defaultUnmarshallers = map[string]Unmarshaller{
		"application/json": json.Unmarshal,
		"application/xml":  xml.Unmarshal,
		"application/yaml": yaml.Unmarshal,
		"text/yaml":        yaml.Unmarshal,
		"text/html":        xml.Unmarshal,
	}

	defaultEncodingReaders = map[string]EncodingReaderConstructor{
		"gzip":       newGzipReader,
		"x-gzip":     newGzipReader,
		"compress":   newCompressReader,
		"x-compress": newCompressReader,
		"deflate":    newDeflateReader,
	}

	defaultMarshallers = map[string]Marshaller{
		"application/json": json.Marshal,
		"application/xml":  xml.Marshal,
		"application/yaml": yaml.Marshal,
		"text/yaml":        yaml.Marshal,
		"text/html":        xml.Marshal,
	}

	defaultEncodingWriters = map[string]EncodingWriterConstructor{
		"gzip":       newGzipWriter,
		"x-gzip":     newGzipWriter,
		"compress":   newCompressWriter,
		"x-compress": newCompressWriter,
		"deflate":    newDeflateWriter,
	}
)

// Base is the base representation.
type Base struct {
	encoding        []string
	mediaType       string
	charset         string
	language        string
	location        url.URL
	sourceQuality   float32
	features        []string
	marshallers     map[string]Marshaller
	unmarshallers   map[string]Unmarshaller
	encodingReaders map[string]EncodingReaderConstructor
	encodingWriters map[string]EncodingWriterConstructor
}

// ContentType retrieves the content type of the representation.
func (r Base) ContentType() string { return r.mediaType }

// SetContentType modifies the content type of the representation.
func (r *Base) SetContentType(ct string) { r.mediaType = ct }

// ContentLanguage retrieves the content language of the representation.
func (r Base) ContentLanguage() string { return r.language }

// SetContentLanguage modifies the content language of the representation.
func (r *Base) SetContentLanguage(cl string) { r.language = cl }

// ContentEncoding retrieves the content encoding of the representation.
func (r Base) ContentEncoding() []string { return r.encoding }

// SetContentEncoding modifies the content encoding of the representation.
func (r *Base) SetContentEncoding(ce []string) { r.encoding = ce }

// ContentCharset retrieves the content charset of the representation.
func (r Base) ContentCharset() string { return r.charset }

// SetContentCharset modifies the content charset of the representation.
func (r *Base) SetContentCharset(cc string) { r.charset = cc }

// ContentLocation retrieves the content location of the representation.
func (r Base) ContentLocation() url.URL { return r.location }

// SetContentLocation modifies the content location of the representation.
func (r *Base) SetContentLocation(l url.URL) { r.location = l }

// ContentFeatures retrieves the content features of the representation.
func (r Base) ContentFeatures() []string { return r.features }

// SetContentFeatures retrieves the content features of the representation.
func (r *Base) SetContentFeatures(cf []string) { r.features = cf }

// SourceQuality retrieves the source quality of the representation.
func (r Base) SourceQuality() float32 { return r.sourceQuality }

// SetSourceQuality modifies the source quality of the representation.
func (r *Base) SetSourceQuality(sq float32) { r.sourceQuality = sq }

// SetMarshallers modifies the marshallers for the representation.
func (r *Base) SetMarshallers(m map[string]Marshaller) {
	r.marshallers = m
}

// SetUnmarshallers modifies the unmarshallers for the representation.
func (r *Base) SetUnmarshallers(u map[string]Unmarshaller) {
	r.unmarshallers = u
}

// SetEncodingReaders modifies the encoding readers for the representation.
func (r *Base) SetEncodingReaders(e map[string]EncodingReaderConstructor) {
	r.encodingReaders = e
}

// SetEncodingWriters modifies the encoding writers for the representation.
func (r *Base) SetEncodingWriters(e map[string]EncodingWriterConstructor) {
	r.encodingWriters = e
}

// Bytes retrieves the serialized form of the representation.
func (r Base) Bytes(out interface{}) ([]byte, error) {
	marshallers := defaultMarshallers
	if len(r.marshallers) > 0 {
		marshallers = r.marshallers
	}

	ct := strings.Split(r.ContentType(), ";")[0]
	if _, ok := marshallers[strings.ToLower(ct)]; !ok {
		return []byte{}, ErrUnsupportedContentType
	}

	// serialize.
	b, err := marshallers[strings.ToLower(ct)](out)
	if err != nil {
		return b, err
	}

	encodings := r.ContentEncoding()
	if len(encodings) < 1 || strings.ToLower(encodings[0]) == "identity" {
		return b, nil
	}

	// encode.
	return r.encode(b)
}

func (r *Base) encode(b []byte) (bb []byte, err error) {
	var (
		buf                bytes.Buffer
		writer             io.WriteCloser = &closeableBuffer{&buf}
		encodings                         = r.ContentEncoding()
		writerConstructors                = defaultEncodingWriters
	)

	if len(r.encodingWriters) > 0 {
		writerConstructors = r.encodingWriters
	}

	for _, e := range encodings {
		if _, ok := writerConstructors[strings.ToLower(e)]; !ok {
			err = ErrUnsupportedContentEncoding
			return
		}
		if writer, err = writerConstructors[strings.ToLower(e)](writer); err != nil {
			return
		}
	}
	if _, err = writer.Write(b); err != nil {
		return
	}
	if err = writer.Close(); err != nil {
		return
	}
	bb = buf.Bytes()
	return
}

// FromBytes constructs the representation from its serialized form.
func (r Base) FromBytes(b []byte, in interface{}) (err error) {
	unmarshallers := defaultUnmarshallers
	if len(r.unmarshallers) > 0 {
		unmarshallers = r.unmarshallers
	}

	ct := strings.Split(r.ContentType(), ";")[0]
	if _, ok := unmarshallers[strings.ToLower(ct)]; !ok {
		err = ErrUnsupportedContentType
		return
	}

	// decode.
	encodings := r.ContentEncoding()
	if len(encodings) > 0 && strings.ToLower(encodings[0]) != "identity" {
		if b, err = r.decode(b); err != nil {
			return
		}
	}

	// deserialize.
	return unmarshallers[strings.ToLower(ct)](b, in)
}

func (r *Base) decode(b []byte) (bb []byte, err error) {
	var (
		buf                              = bytes.NewBuffer(b)
		reader             io.ReadCloser = &closeableBuffer{buf}
		encodings                        = r.ContentEncoding()
		readerConstructors               = defaultEncodingReaders
	)

	if len(r.encodingReaders) > 0 {
		readerConstructors = r.encodingReaders
	}

	for idx := len(encodings) - 1; idx >= 0; idx-- {
		e := encodings[idx]
		if _, ok := readerConstructors[strings.ToLower(e)]; !ok {
			err = ErrUnsupportedContentEncoding
			return
		}
		if reader, err = readerConstructors[strings.ToLower(e)](reader); err != nil {
			return
		}
		defer func() {
			err = reader.Close()
		}()
	}
	if bb, err = ioutil.ReadAll(reader); err != nil {
		return
	}
	return
}

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

func (cb closeableBuffer) Read(b []byte) (int, error) {
	return cb.buf.Read(b)
}
