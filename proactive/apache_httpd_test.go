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
	htmlLevel2, english, ascii, gzip :=
		"text/html;level=2", "en-US", "ascii", "gzip"
	request := httptest.NewRequest("GET", "http://freer.ddns.net/thing", nil)
	request.Header.Add("Accept-Language", english)
	request.Header.Add("Accept-Encoding", gzip)
	request.Header.Add("Accept-Charset", ascii)
	request.Header.Add("Accept", htmlLevel2)
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

func (s *ApacheHTTPDTestSuite) TearDownTest() {
	s.sut = nil
}
