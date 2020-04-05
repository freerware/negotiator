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

package representation_test

import (
	"bytes"
	"io"
	"net/url"
	"testing"

	"github.com/freerware/negotiator/internal/test"
	"github.com/freerware/negotiator/representation"
	"github.com/stretchr/testify/suite"
)

type BaseTestSuite struct {
	suite.Suite
}

func TestBaseTestSuite(t *testing.T) {
	suite.Run(t, new(BaseTestSuite))
}

func (s *BaseTestSuite) TestBaseRepresentation_Bytes() {
	// json representations.
	identityJSON := representation.Base{}
	identityJSON.SetContentType("application/json")
	identityJSON.SetContentEncoding([]string{"identity"})
	gzippedJSON := representation.Base{}
	gzippedJSON.SetContentType("application/json")
	gzippedJSON.SetContentEncoding([]string{"gzip"})
	compressedJSON := representation.Base{}
	compressedJSON.SetContentType("application/json")
	compressedJSON.SetContentEncoding([]string{"compress"})
	deflatedJSON := representation.Base{}
	deflatedJSON.SetContentType("application/json")
	deflatedJSON.SetContentEncoding([]string{"deflate"})

	// yaml representations.
	identityYAML := representation.Base{}
	identityYAML.SetContentType("application/yaml")
	identityYAML.SetContentEncoding([]string{"identity"})
	gzippedYAML := representation.Base{}
	gzippedYAML.SetContentType("application/yaml")
	gzippedYAML.SetContentEncoding([]string{"gzip"})
	compressedYAML := representation.Base{}
	compressedYAML.SetContentType("application/yaml")
	compressedYAML.SetContentEncoding([]string{"compress"})
	deflatedYAML := representation.Base{}
	deflatedYAML.SetContentType("application/yaml")
	deflatedYAML.SetContentEncoding([]string{"deflate"})

	// xml representations.
	identityXML := representation.Base{}
	identityXML.SetContentType("application/xml")
	identityXML.SetContentEncoding([]string{"identity"})
	gzippedXML := representation.Base{}
	gzippedXML.SetContentType("application/xml")
	gzippedXML.SetContentEncoding([]string{"gzip"})
	compressedXML := representation.Base{}
	compressedXML.SetContentType("application/xml")
	compressedXML.SetContentEncoding([]string{"compress"})
	deflatedXML := representation.Base{}
	deflatedXML.SetContentType("application/xml")
	deflatedXML.SetContentEncoding([]string{"deflate"})

	// error causing representations.
	unsupportedMediaType := representation.Base{}
	unsupportedMediaType.SetContentType("application/beeboop")

	unsupportedContentEncoding := representation.Base{}
	unsupportedContentEncoding.SetContentType("application/json")
	unsupportedContentEncoding.SetContentEncoding([]string{"beeboop"})

	tests := []struct {
		name string
		in   representation.Base
		err  error
	}{
		{"IdentityJSON", identityJSON, nil},
		{"GzippedJSON", gzippedJSON, nil},
		{"CompressedJSON", compressedJSON, nil},
		{"DeflatedJSON", deflatedJSON, nil},
		{"IdentityYAML", identityYAML, nil},
		{"GzippedYAML", gzippedYAML, nil},
		{"CompressedYAML", compressedYAML, nil},
		{"DeflatedYAML", deflatedYAML, nil},
		{"IdentityXML", identityXML, nil},
		{"GzippedXML", gzippedXML, nil},
		{"CompressedXML", compressedXML, nil},
		{"DeflatedXML", deflatedXML, nil},
		{
			"UnsupportedMediaType",
			unsupportedMediaType,
			representation.ErrUnsupportedContentType,
		},
		{
			"UnsupportedContentEncoding",
			unsupportedContentEncoding,
			representation.ErrUnsupportedContentEncoding,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			// arrange.
			rep := test.Representation{A: "TEST", B: 28}

			// action.
			b, err := tt.in.Bytes(&rep)

			// assert.
			if tt.err != nil {
				s.Require().EqualError(err, tt.err.Error())
			} else {
				s.Require().NoError(err)
				actualRep := test.Representation{}
				actualRep.SetContentEncoding(rep.ContentEncoding())
				actualRep.SetContentType(rep.ContentType())
				err = tt.in.FromBytes(b, &actualRep)
				s.Require().NoError(err)
				s.Equal(rep, actualRep)
			}
		})
	}
}

func (s *BaseTestSuite) TestBaseRepresentation_FromBytes() {

	// json representations.
	identityJSON := representation.Base{}
	identityJSON.SetContentType("application/json")
	identityJSON.SetContentEncoding([]string{"identity"})
	gzippedJSON := representation.Base{}
	gzippedJSON.SetContentType("application/json")
	gzippedJSON.SetContentEncoding([]string{"gzip"})
	compressedJSON := representation.Base{}
	compressedJSON.SetContentType("application/json")
	compressedJSON.SetContentEncoding([]string{"compress"})
	deflatedJSON := representation.Base{}
	deflatedJSON.SetContentType("application/json")
	deflatedJSON.SetContentEncoding([]string{"deflate"})

	// yaml representations.
	identityYAML := representation.Base{}
	identityYAML.SetContentType("application/yaml")
	identityYAML.SetContentEncoding([]string{"identity"})
	gzippedYAML := representation.Base{}
	gzippedYAML.SetContentType("application/yaml")
	gzippedYAML.SetContentEncoding([]string{"gzip"})
	compressedYAML := representation.Base{}
	compressedYAML.SetContentType("application/yaml")
	compressedYAML.SetContentEncoding([]string{"compress"})
	deflatedYAML := representation.Base{}
	deflatedYAML.SetContentType("application/yaml")
	deflatedYAML.SetContentEncoding([]string{"deflate"})

	// xml representations.
	identityXML := representation.Base{}
	identityXML.SetContentType("application/xml")
	identityXML.SetContentEncoding([]string{"identity"})
	gzippedXML := representation.Base{}
	gzippedXML.SetContentType("application/xml")
	gzippedXML.SetContentEncoding([]string{"gzip"})
	compressedXML := representation.Base{}
	compressedXML.SetContentType("application/xml")
	compressedXML.SetContentEncoding([]string{"compress"})
	deflatedXML := representation.Base{}
	deflatedXML.SetContentType("application/xml")
	deflatedXML.SetContentEncoding([]string{"deflate"})

	// error causing representations.
	unsupportedMediaType := representation.Base{}
	unsupportedMediaType.SetContentType("application/beeboop")

	unsupportedContentEncoding := representation.Base{}
	unsupportedContentEncoding.SetContentType("application/json")
	unsupportedContentEncoding.SetContentEncoding([]string{"beeboop"})

	tests := []struct {
		name string
		in   representation.Base
		err  error
	}{
		{"IdentityJSON", identityJSON, nil},
		{"GzippedJSON", gzippedJSON, nil},
		{"CompressedJSON", compressedJSON, nil},
		{"DeflatedJSON", deflatedJSON, nil},
		{"IdentityYAML", identityYAML, nil},
		{"GzippedYAML", gzippedYAML, nil},
		{"CompressedYAML", compressedYAML, nil},
		{"DeflatedYAML", deflatedYAML, nil},
		{"IdentityXML", identityXML, nil},
		{"GzippedXML", gzippedXML, nil},
		{"CompressedXML", compressedXML, nil},
		{"DeflatedXML", deflatedXML, nil},
		{
			"UnsupportedMediaType",
			unsupportedMediaType,
			representation.ErrUnsupportedContentType,
		},
		{
			"UnsupportedContentEncoding",
			unsupportedContentEncoding,
			representation.ErrUnsupportedContentEncoding,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			// arrange - calculate input bytes.
			rep := test.Representation{A: "TEST", B: 28}
			rep.SetContentEncoding(tt.in.ContentEncoding())
			rep.SetContentType(tt.in.ContentType())
			b, _ := tt.in.Bytes(&rep)

			// arrange - prepare result.
			actualRep := test.Representation{}
			actualRep.SetContentEncoding(tt.in.ContentEncoding())
			actualRep.SetContentType(tt.in.ContentType())

			// action.
			err := tt.in.FromBytes(b, &actualRep)

			// assert.
			if tt.err != nil {
				s.Require().EqualError(err, tt.err.Error())
			} else {
				s.Require().NoError(err)
				s.Equal(rep, actualRep)
			}
		})
	}
}

func (s BaseTestSuite) TestBaseRepresentation_ContentType() {

	// arrange.
	json := "application/json"
	rep := test.Representation{A: "TEST", B: 28}
	rep.SetContentType(json)

	// action.
	ct := rep.ContentType()

	// assert.
	s.Equal(json, ct)
}

func (s BaseTestSuite) TestBaseRepresentation_SetContentType() {

	// arrange.
	json := "application/json"
	yaml := "application/yaml"
	rep := test.Representation{A: "TEST", B: 28}
	rep.SetContentType(json)

	// action.
	rep.SetContentType(yaml)

	// assert.
	s.Equal(yaml, rep.ContentType())
}

func (s BaseTestSuite) TestBaseRepresentation_ContentLanguage() {

	// arrange.
	english := "en-US"
	rep := test.Representation{A: "TEST", B: 28}
	rep.SetContentLanguage(english)

	// action.
	cl := rep.ContentLanguage()

	// assert.
	s.Equal(english, cl)
}

func (s BaseTestSuite) TestBaseRepresentation_SetContentLanguage() {

	// arrange.
	english := "en-US"
	british := "en-GB"
	rep := test.Representation{A: "TEST", B: 28}
	rep.SetContentLanguage(english)

	// action.
	rep.SetContentLanguage(british)

	// assert.
	s.Equal(british, rep.ContentLanguage())
}

func (s BaseTestSuite) TestBaseRepresentation_ContentEncoding() {

	// arrange.
	gzip := "gzip"
	rep := test.Representation{A: "TEST", B: 28}
	rep.SetContentEncoding([]string{gzip})

	// action.
	ce := rep.ContentEncoding()

	// assert.
	s.Equal(gzip, ce[0])
}

func (s BaseTestSuite) TestBaseRepresentation_SetContentEncoding() {

	// arrange.
	gzip := "gzip"
	compress := "compress"
	rep := test.Representation{A: "TEST", B: 28}
	rep.SetContentEncoding([]string{gzip})

	// action.
	rep.SetContentEncoding([]string{compress})

	// assert.
	s.Equal(compress, rep.ContentEncoding()[0])
}

func (s BaseTestSuite) TestBaseRepresentation_ContentCharset() {

	// arrange.
	ascii := "ascii"
	rep := test.Representation{A: "TEST", B: 28}
	rep.SetContentCharset(ascii)

	// action.
	cc := rep.ContentCharset()

	// assert.
	s.Equal(ascii, cc)
}

func (s BaseTestSuite) TestBaseRepresentation_SetContentCharset() {

	// arrange.
	ascii := "ascii"
	utf8 := "utf8"
	rep := test.Representation{A: "TEST", B: 28}
	rep.SetContentCharset(ascii)

	// action.
	rep.SetContentCharset(utf8)

	// assert.
	s.Equal(utf8, rep.ContentCharset())
}

func (s BaseTestSuite) TestBaseRepresentation_ContentLocation() {

	// arrange.
	loc, _ := url.Parse("http://www.example.com/")
	rep := test.Representation{A: "TEST", B: 28}
	rep.SetContentLocation(*loc)

	// action.
	cc := rep.ContentLocation()

	// assert.
	s.Equal(*loc, cc)
}

func (s BaseTestSuite) TestBaseRepresentation_SetContentLocation() {

	// arrange.
	loc, _ := url.Parse("http://www.example.com/")
	newLoc, _ := url.Parse("http://www.example.com/another")
	rep := test.Representation{A: "TEST", B: 28}
	rep.SetContentLocation(*loc)

	// action.
	rep.SetContentLocation(*newLoc)

	// assert.
	s.Equal(*newLoc, rep.ContentLocation())
}

func (s BaseTestSuite) TestBaseRepresentation_ContentFeatures() {

	// arrange.
	feature := "foo"
	rep := test.Representation{A: "TEST", B: 28}
	rep.SetContentFeatures([]string{feature})

	// action.
	cf := rep.ContentFeatures()

	// assert.
	s.Equal(feature, cf[0])
}

func (s BaseTestSuite) TestBaseRepresentation_SetContentFeatures() {

	// arrange.
	feature := "foo"
	newFeature := "bar"
	rep := test.Representation{A: "TEST", B: 28}
	rep.SetContentFeatures([]string{feature})

	// action.
	rep.SetContentFeatures([]string{newFeature})

	// assert.
	s.Equal(newFeature, rep.ContentFeatures()[0])
}

func (s BaseTestSuite) TestBaseRepresentation_SourceQuality() {

	// arrange.
	perfect := representation.SourceQualityPerfect
	rep := test.Representation{A: "TEST", B: 28}
	rep.SetSourceQuality(perfect)

	// action.
	sq := rep.SourceQuality()

	// assert.
	s.Equal(perfect, sq)
}

func (s BaseTestSuite) TestBaseRepresentation_SetSourceQuality() {

	// arrange.
	perfect := representation.SourceQualityPerfect
	nearlyPerfect := representation.SourceQualityNearlyPerfect
	rep := test.Representation{A: "TEST", B: 28}
	rep.SetSourceQuality(perfect)

	// action.
	rep.SetSourceQuality(nearlyPerfect)

	// assert.
	s.Equal(nearlyPerfect, rep.SourceQuality())
}

func (s BaseTestSuite) TestBaseRepresentation_SetMarshallers() {

	// arrange.
	testContentType := "text/test"
	marshaller := func(in interface{}) ([]byte, error) {
		return []byte{}, nil
	}
	rep := test.Representation{A: "TEST", B: 28}
	rep.SetContentType(testContentType)

	// action.
	rep.SetMarshallers(map[string]representation.Marshaller{
		testContentType: marshaller,
	})

	// assert.
	b, err := rep.Bytes()
	s.Require().NoError(err)
	s.Equal([]byte{}, b)
}

func (s BaseTestSuite) TestBaseRepresentation_SetUnmarshallers() {

	// arrange.
	testContentType := "text/test"
	unmarshaller := func(b []byte, in interface{}) error {
		return nil
	}
	rep := test.Representation{A: "TEST", B: 28}
	rep.SetContentType(testContentType)

	// action.
	rep.SetUnmarshallers(map[string]representation.Unmarshaller{
		testContentType: unmarshaller,
	})

	// assert.
	err := rep.FromBytes([]byte{})
	s.Require().NoError(err)
}

func (s BaseTestSuite) TestBaseRepresentation_SetEncodingReaders() {

	// arrange.
	testContentEncoding := "test"
	reader := func(r io.Reader) (io.ReadCloser, error) {
		cb := closeableBuffer{}
		return &cb, nil
	}
	rep := test.Representation{A: "TEST", B: 28}
	rep.SetContentType("application/json")
	rep.SetContentEncoding([]string{testContentEncoding})

	// action.
	rep.SetEncodingReaders(map[string]representation.EncodingReaderConstructor{
		testContentEncoding: reader,
	})
}

func (s BaseTestSuite) TestBaseRepresentation_SetEncodingWriters() {

	// arrange.
	testContentEncoding := "test"
	writer := func(r io.WriteCloser) (io.WriteCloser, error) {
		cb := closeableBuffer{}
		return &cb, nil
	}
	rep := test.Representation{A: "TEST", B: 28}
	rep.SetContentType("application/json")
	rep.SetContentEncoding([]string{testContentEncoding})

	// action.
	rep.SetEncodingWriters(map[string]representation.EncodingWriterConstructor{
		testContentEncoding: writer,
	})
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
