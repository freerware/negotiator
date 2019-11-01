package proactive

import (
	"net/http/httptest"
	"testing"

	_representation "github.com/freerware/negotiator/internal/representation"
	"github.com/freerware/negotiator/internal/test"
	"github.com/freerware/negotiator/representation"
	"github.com/stretchr/testify/suite"
)

type ApacheHTTPDTestSuite struct {
	suite.Suite

	// system under test.
	sut representation.Chooser
}

func TestApacheHTTPDTestSuite(t *testing.T) {
	suite.Run(t, new(ApacheHTTPDTestSuite))
}

func (s *ApacheHTTPDTestSuite) SetupTest() {
	s.sut = ApacheHTTPD()
}

func (s *ApacheHTTPDTestSuite) TestApacheHTTPD_Choose_NewAcceptErr() {
	// arrange.
	request := httptest.NewRequest("GET", "http://freer.ddns.net/thing", nil)
	request.Header.Add("Accept", "invalid/")
	variants := []representation.Representation{}

	// action.
	_, err := s.sut.Choose(request, variants...)

	// assert.
	s.Require().Error(err)
}

func (s *ApacheHTTPDTestSuite) TestApacheHTTPD_Choose_NewAcceptEncodingErr() {
	// arrange.
	request := httptest.NewRequest("GET", "http://freer.ddns.net/thing", nil)
	request.Header.Add("Accept", "application/json")
	request.Header.Add("Accept-Encoding", "invalid")
	variants := []representation.Representation{}

	// action.
	_, err := s.sut.Choose(request, variants...)

	// assert.
	s.Require().Error(err)
}

func (s *ApacheHTTPDTestSuite) TestApacheHTTPD_Choose_NewAcceptLanguageErr() {
	// arrange.
	request := httptest.NewRequest("GET", "http://freer.ddns.net/thing", nil)
	request.Header.Add("Accept", "application/json")
	request.Header.Add("Accept-Encoding", "gzip")
	request.Header.Add("Accept-Language", "invalid")
	variants := []representation.Representation{}

	// action.
	_, err := s.sut.Choose(request, variants...)

	// assert.
	s.Require().Error(err)
}

func (s *ApacheHTTPDTestSuite) TestApacheHTTPD_Choose_NewAcceptCharsetErr() {
	// arrange.
	request := httptest.NewRequest("GET", "http://freer.ddns.net/thing", nil)
	request.Header.Add("Accept", "application/json")
	request.Header.Add("Accept-Encoding", "gzip")
	request.Header.Add("Accept-Language", "en-US")
	request.Header.Add("Accept-Charset", "")
	variants := []representation.Representation{}

	// action.
	_, err := s.sut.Choose(request, variants...)

	// assert.
	s.Require().Error(err)
}

func (s *ApacheHTTPDTestSuite) TestApacheHTTPD_Choose_NoRepresentations() {
	// arrange.
	request := httptest.NewRequest("GET", "http://freer.ddns.net/thing", nil)
	request.Header.Add("Accept", "application/json")
	request.Header.Add("Accept-Encoding", "gzip")
	request.Header.Add("Accept-Language", "en-US")
	request.Header.Add("Accept-Charset", "ascii")
	variants := []representation.Representation{}

	// action.
	chosen, err := s.sut.Choose(request, variants...)

	// assert.
	s.Require().NoError(err)
	s.Require().Nil(chosen)
}

func (s *ApacheHTTPDTestSuite) TestApacheHTTPD_Choose_RFC7231AcceptExample() {
	// arrange.
	english, ascii, gzip := "en-US", "ascii", "gzip"
	request := httptest.NewRequest("GET", "http://freer.ddns.net/thing", nil)
	request.Header.Add("Accept-Language", english)
	request.Header.Add("Accept-Encoding", gzip)
	request.Header.Add("Accept-Charset", ascii)
	request.Header.Add("Accept", "text/*;q=0.3")
	request.Header.Add("Accept", "text/html;q=0.7")
	request.Header.Add("Accept", "text/html;level=1")
	request.Header.Add("Accept", "text/html;level=2;q=0.4")
	request.Header.Add("Accept", "*/*;q=0.5")
	v1 := _representation.NewBuilder().
		WithLocation(*request.URL).
		WithType("text/html").
		WithLanguage(english).
		WithEncoding(gzip).
		WithCharset(ascii).
		WithSourceQuality(1.0).
		Build(test.RepresentationBuilderFunc)
	v2 := _representation.NewBuilder().
		WithLocation(*request.URL).
		WithType("text/html;level=2").
		WithLanguage(english).
		WithEncoding(gzip).
		WithCharset(ascii).
		WithSourceQuality(1.0).
		Build(test.RepresentationBuilderFunc)
	v3 := _representation.NewBuilder().
		WithLocation(*request.URL).
		WithType("text/html;level=3").
		WithLanguage(english).
		WithEncoding(gzip).
		WithCharset(ascii).
		WithSourceQuality(1.0).
		Build(test.RepresentationBuilderFunc)
	v4 := _representation.NewBuilder().
		WithLocation(*request.URL).
		WithType("text/html;level=1").
		WithLanguage(english).
		WithEncoding(gzip).
		WithCharset(ascii).
		WithSourceQuality(1.0).
		Build(test.RepresentationBuilderFunc)
	v5 := _representation.NewBuilder().
		WithLocation(*request.URL).
		WithType("text/plain").
		WithLanguage(english).
		WithEncoding(gzip).
		WithCharset(ascii).
		WithSourceQuality(1.0).
		Build(test.RepresentationBuilderFunc)
	v6 := _representation.NewBuilder().
		WithLocation(*request.URL).
		WithType("image/jpeg").
		WithLanguage(english).
		WithEncoding(gzip).
		WithCharset(ascii).
		WithSourceQuality(1.0).
		Build(test.RepresentationBuilderFunc)
	variants := []representation.Representation{v1, v2, v3, v4, v5, v6}

	// action.
	chosen, err := s.sut.Choose(request, variants...)

	// assert.
	s.Require().NoError(err)
	s.Equal(v4, chosen)
}

func (s *ApacheHTTPDTestSuite) TestApacheHTTPD_Choose_InvalidLevelErr() {
	// arrange.
	htmlLevel2, html, english, ascii, gzip :=
		"text/html;level=2", "text/html", "en-US", "ascii", "gzip"
	request := httptest.NewRequest("GET", "http://freer.ddns.net/thing", nil)
	request.Header.Add("Accept-Language", english)
	request.Header.Add("Accept-Encoding", gzip)
	request.Header.Add("Accept-Charset", ascii)
	request.Header.Add("Accept", htmlLevel2)
	request.Header.Add("Accept", html)
	v1 := _representation.NewBuilder().
		WithLocation(*request.URL).
		WithType(htmlLevel2).
		WithLanguage(english).
		WithEncoding(gzip).
		WithCharset(ascii).
		WithSourceQuality(1.0).
		Build(test.RepresentationBuilderFunc)
	v2 := _representation.NewBuilder().
		WithLocation(*request.URL).
		WithType("text/html;level=x").
		WithLanguage(english).
		WithEncoding(gzip).
		WithCharset(ascii).
		WithSourceQuality(1.0).
		Build(test.RepresentationBuilderFunc)
	variants := []representation.Representation{v1, v2}

	// action.
	_, err := s.sut.Choose(request, variants...)

	// assert.
	s.Require().Error(err)
}

func (s *ApacheHTTPDTestSuite) TestApacheHTTPD_Choose_BestCharset() {
	// arrange.
	htmlLevel2, html, english, asciiRange, gzip, utf8Range, ascii, utf8 :=
		"text/html;level=2", "text/html", "en-US", "ascii;q=0.9", "gzip", "utf8;q=0.8", "ascii", "utf8"
	request := httptest.NewRequest("GET", "http://freer.ddns.net/thing", nil)
	request.Header.Add("Accept-Language", english)
	request.Header.Add("Accept-Encoding", gzip)
	request.Header.Add("Accept-Charset", asciiRange)
	request.Header.Add("Accept-Charset", utf8Range)
	request.Header.Add("Accept", htmlLevel2)
	request.Header.Add("Accept", html)
	v1 := _representation.NewBuilder().
		WithLocation(*request.URL).
		WithType(htmlLevel2).
		WithLanguage(english).
		WithEncoding(gzip).
		WithCharset(ascii).
		WithSourceQuality(1.0).
		Build(test.RepresentationBuilderFunc)
	v2 := _representation.NewBuilder().
		WithLocation(*request.URL).
		WithType(htmlLevel2).
		WithLanguage(english).
		WithEncoding(gzip).
		WithCharset(utf8).
		WithSourceQuality(1.0).
		Build(test.RepresentationBuilderFunc)
	variants := []representation.Representation{v1, v2}

	// action.
	chosen, err := s.sut.Choose(request, variants...)

	// assert.
	s.Require().NoError(err)
	s.Equal(v1, chosen)
}

func (s *ApacheHTTPDTestSuite) TestApacheHTTPD_Choose_NotISO88591() {
	// arrange.
	htmlLevel2, html, english, gzip, ascii, iso88591 :=
		"text/html;level=2", "text/html", "en-US", "gzip", "ascii", "iso-8859-1"
	request := httptest.NewRequest("GET", "http://freer.ddns.net/thing", nil)
	request.Header.Add("Accept-Language", english)
	request.Header.Add("Accept-Encoding", gzip)
	request.Header.Add("Accept-Charset", iso88591)
	request.Header.Add("Accept-Charset", ascii)
	request.Header.Add("Accept", htmlLevel2)
	request.Header.Add("Accept", html)
	v1 := _representation.NewBuilder().
		WithLocation(*request.URL).
		WithType(htmlLevel2).
		WithLanguage(english).
		WithEncoding(gzip).
		WithCharset(ascii).
		WithSourceQuality(1.0).
		Build(test.RepresentationBuilderFunc)
	v2 := _representation.NewBuilder().
		WithLocation(*request.URL).
		WithType(htmlLevel2).
		WithLanguage(english).
		WithEncoding(gzip).
		WithCharset(iso88591).
		WithSourceQuality(1.0).
		Build(test.RepresentationBuilderFunc)
	variants := []representation.Representation{v1, v2}

	// action.
	chosen, err := s.sut.Choose(request, variants...)

	// assert.
	s.Require().NoError(err)
	s.Equal(v1, chosen)
}

func (s *ApacheHTTPDTestSuite) TestApacheHTTPD_Choose_BestEncoding() {
	// arrange.
	htmlLevel2, html, english, gzipRange, ascii, gzip, compressRange, compress :=
		"text/html;level=2", "text/html", "en-US", "gzip;q=0.9", "ascii", "gzip", "compress;q=0.8", "compress"
	request := httptest.NewRequest("GET", "http://freer.ddns.net/thing", nil)
	request.Header.Add("Accept-Language", english)
	request.Header.Add("Accept-Encoding", gzipRange)
	request.Header.Add("Accept-Encoding", compressRange)
	request.Header.Add("Accept-Charset", ascii)
	request.Header.Add("Accept", htmlLevel2)
	request.Header.Add("Accept", html)
	v1 := _representation.NewBuilder().
		WithLocation(*request.URL).
		WithType(htmlLevel2).
		WithLanguage(english).
		WithEncoding(gzip).
		WithCharset(ascii).
		WithSourceQuality(1.0).
		Build(test.RepresentationBuilderFunc)
	v2 := _representation.NewBuilder().
		WithLocation(*request.URL).
		WithType(htmlLevel2).
		WithLanguage(english).
		WithEncoding(compress).
		WithCharset(ascii).
		WithSourceQuality(1.0).
		Build(test.RepresentationBuilderFunc)
	variants := []representation.Representation{v1, v2}

	// action.
	chosen, err := s.sut.Choose(request, variants...)

	// assert.
	s.Require().NoError(err)
	s.Equal(v1, chosen)
}

func (s *ApacheHTTPDTestSuite) TestApacheHTTPD_Choose_SmallestContentLength() {
	// arrange.
	htmlLevel2, html, english, ascii, gzip :=
		"text/html;level=2", "text/html", "en-US", "ascii", "gzip"
	request := httptest.NewRequest("GET", "http://freer.ddns.net/thing", nil)
	request.Header.Add("Accept-Language", english)
	request.Header.Add("Accept-Encoding", gzip)
	request.Header.Add("Accept-Charset", ascii)
	request.Header.Add("Accept", htmlLevel2)
	request.Header.Add("Accept", html)
	v1 := _representation.NewBuilder().
		WithLocation(*request.URL).
		WithType(htmlLevel2).
		WithLanguage(english).
		WithEncoding(gzip).
		WithCharset(ascii).
		WithSourceQuality(1.0).
		Build(test.RepresentationBuilderFunc)
	v2 := _representation.NewBuilder().
		WithLocation(*request.URL).
		WithType(htmlLevel2).
		WithLanguage(english).
		WithEncoding(gzip).
		WithCharset(ascii).
		WithSourceQuality(1.0).
		Build(func(ctx _representation.BuilderContext) representation.Representation {
			r := test.Representation{}
			r.A = "THIS IS REALLY REALLY REALLY LONG!"
			r.B = 100
			r.SetContentType(ctx.ContentType)
			r.SetContentLanguage(ctx.ContentLanguage)
			r.SetContentCharset(ctx.ContentCharset)
			r.SetContentEncoding(ctx.ContentEncoding)
			r.SetContentLocation(ctx.ContentLocation)
			r.SetContentFeatures(ctx.ContentFeatures)
			r.SetSourceQuality(ctx.SourceQuality)
			return r
		})
	variants := []representation.Representation{v1, v2}

	// action.
	chosen, err := s.sut.Choose(request, variants...)

	// assert.
	s.Require().NoError(err)
	s.Equal(v1, chosen)
}

func (s *ApacheHTTPDTestSuite) TearDownTest() {
	s.sut = nil
}
