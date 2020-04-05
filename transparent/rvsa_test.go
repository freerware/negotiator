package transparent

import (
	"net/http/httptest"
	"testing"

	_representation "github.com/freerware/negotiator/internal/representation"
	"github.com/freerware/negotiator/internal/test"
	"github.com/freerware/negotiator/representation"
	"github.com/stretchr/testify/suite"
)

type RVSATestSuite struct {
	suite.Suite

	// system under test.
	sut representation.Chooser
}

func TestRVSATestSuite(t *testing.T) {
	suite.Run(t, new(RVSATestSuite))
}

func (s *RVSATestSuite) SetupTest() {
	s.sut = RVSA1()
}

func (s *RVSATestSuite) TestRVSA_Choose_NewAcceptErr() {
	// arrange.
	request := httptest.NewRequest("GET", "http://freer.ddns.net/thing", nil)
	request.Header.Add("Accept", "invalid/")
	variants := []representation.Representation{}

	// action.
	_, err := s.sut.Choose(request, variants...)

	// assert.
	s.Require().Error(err)
}

func (s *RVSATestSuite) TestRVSA_Choose_NewAcceptLanguageErr() {
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

func (s *RVSATestSuite) TestRVSA_Choose_NewAcceptCharsetErr() {
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

func (s *RVSATestSuite) TestRVSA_Choose_NoRepresentation() {
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

func (s *RVSATestSuite) TestRVSA_Choose() {

	// arrange.
	var (
		html       = "text/html"
		english    = "en-US"
		asciiRange = "ascii;q=0.9"
		gzip       = "gzip"
		utf8Range  = "utf8;q=0.8"
		ascii      = "ascii"
		utf8       = "utf8"
		features   = "foo"
	)
	request := httptest.NewRequest("GET", "http://freer.ddns.net/thing", nil)
	request.Header.Add("Accept-Language", english)
	request.Header.Add("Accept-Encoding", gzip)
	request.Header.Add("Accept-Charset", asciiRange)
	request.Header.Add("Accept-Charset", utf8Range)
	request.Header.Add("Accept-Features", features)
	request.Header.Add("Accept", html)
	v1 := _representation.NewBuilder().
		WithLocation(*request.URL).
		WithType(html).
		WithLanguage(english).
		WithEncoding(gzip).
		WithCharset(ascii).
		WithSourceQuality(1.0).
		WithFeature(features).
		Build(test.RepresentationBuilderFunc)
	v2 := _representation.NewBuilder().
		WithLocation(*request.URL).
		WithType(html).
		WithLanguage(english).
		WithEncoding(gzip).
		WithCharset(utf8).
		WithSourceQuality(1.0).
		WithFeature(features).
		Build(test.RepresentationBuilderFunc)
	variants := []representation.Representation{v1, v2}

	// action.
	chosen, err := s.sut.Choose(request, variants...)

	// assert.
	s.Require().NoError(err)
	s.Equal(v1, chosen)
}

func (s *RVSATestSuite) TestRVSA_Choose_NotDefinite_MissingAccept() {

	// arrange.
	var (
		html       = "text/html"
		english    = "en-US"
		asciiRange = "ascii;q=0.9"
		gzip       = "gzip"
		utf8Range  = "utf8;q=0.8"
		ascii      = "ascii"
		utf8       = "utf8"
		features   = "foo"
	)
	request := httptest.NewRequest("GET", "http://freer.ddns.net/thing", nil)
	request.Header.Add("Accept-Language", english)
	request.Header.Add("Accept-Encoding", gzip)
	request.Header.Add("Accept-Charset", asciiRange)
	request.Header.Add("Accept-Charset", utf8Range)
	v1 := _representation.NewBuilder().
		WithLocation(*request.URL).
		WithType(html).
		WithLanguage(english).
		WithEncoding(gzip).
		WithCharset(ascii).
		WithSourceQuality(1.0).
		WithFeature(features).
		Build(test.RepresentationBuilderFunc)
	v2 := _representation.NewBuilder().
		WithLocation(*request.URL).
		WithType(html).
		WithLanguage(english).
		WithEncoding(gzip).
		WithCharset(utf8).
		WithSourceQuality(1.0).
		WithFeature(features).
		Build(test.RepresentationBuilderFunc)
	variants := []representation.Representation{v1, v2}

	// action.
	chosen, err := s.sut.Choose(request, variants...)

	// assert.
	s.Require().NoError(err)
	s.Require().Nil(chosen)
}

func (s *RVSATestSuite) TestRVSA_Choose_NotDefinite_MissingAcceptLanguage() {

	// arrange.
	var (
		html       = "text/html"
		english    = "en-US"
		asciiRange = "ascii;q=0.9"
		gzip       = "gzip"
		utf8Range  = "utf8;q=0.8"
		ascii      = "ascii"
		utf8       = "utf8"
		features   = "foo"
	)
	request := httptest.NewRequest("GET", "http://freer.ddns.net/thing", nil)
	request.Header.Add("Accept-Encoding", gzip)
	request.Header.Add("Accept-Charset", asciiRange)
	request.Header.Add("Accept-Charset", utf8Range)
	request.Header.Add("Accept-Features", features)
	request.Header.Add("Accept", html)
	v1 := _representation.NewBuilder().
		WithLocation(*request.URL).
		WithType(html).
		WithLanguage(english).
		WithEncoding(gzip).
		WithCharset(ascii).
		WithSourceQuality(1.0).
		WithFeature(features).
		Build(test.RepresentationBuilderFunc)
	v2 := _representation.NewBuilder().
		WithLocation(*request.URL).
		WithType(html).
		WithLanguage(english).
		WithEncoding(gzip).
		WithCharset(utf8).
		WithSourceQuality(1.0).
		WithFeature(features).
		Build(test.RepresentationBuilderFunc)
	variants := []representation.Representation{v1, v2}

	// action.
	chosen, err := s.sut.Choose(request, variants...)

	// assert.
	s.Require().NoError(err)
	s.Require().Nil(chosen)
}

func (s *RVSATestSuite) TestRVSA_Choose_NotDefinite_MissingAcceptCharset() {

	// arrange.
	var (
		html     = "text/html"
		english  = "en-US"
		gzip     = "gzip"
		ascii    = "ascii"
		utf8     = "utf8"
		features = "foo"
	)
	request := httptest.NewRequest("GET", "http://freer.ddns.net/thing", nil)
	request.Header.Add("Accept-Language", english)
	request.Header.Add("Accept-Encoding", gzip)
	request.Header.Add("Accept-Features", features)
	request.Header.Add("Accept", html)
	v1 := _representation.NewBuilder().
		WithLocation(*request.URL).
		WithType(html).
		WithLanguage(english).
		WithEncoding(gzip).
		WithCharset(ascii).
		WithSourceQuality(1.0).
		WithFeature(features).
		Build(test.RepresentationBuilderFunc)
	v2 := _representation.NewBuilder().
		WithLocation(*request.URL).
		WithType(html).
		WithLanguage(english).
		WithEncoding(gzip).
		WithCharset(utf8).
		WithSourceQuality(1.0).
		WithFeature(features).
		Build(test.RepresentationBuilderFunc)
	variants := []representation.Representation{v1, v2}

	// action.
	chosen, err := s.sut.Choose(request, variants...)

	// assert.
	s.Require().NoError(err)
	s.Require().Nil(chosen)
}

func (s *RVSATestSuite) TestRVSA_Choose_ZeroQualityFactor() {

	// arrange.
	var (
		html      = "text/html"
		english   = "en-US"
		gzip      = "gzip"
		utf8Range = "utf8;q=0.8"
		ascii     = "ascii"
		features  = "foo"
	)
	request := httptest.NewRequest("GET", "http://freer.ddns.net/thing", nil)
	request.Header.Add("Accept-Language", english)
	request.Header.Add("Accept-Encoding", gzip)
	request.Header.Add("Accept-Charset", utf8Range)
	request.Header.Add("Accept-Features", features)
	request.Header.Add("Accept", html)
	v1 := _representation.NewBuilder().
		WithLocation(*request.URL).
		WithType(html).
		WithLanguage(english).
		WithEncoding(gzip).
		WithCharset(ascii).
		WithSourceQuality(1.0).
		WithFeature(features).
		Build(test.RepresentationBuilderFunc)
	v2 := _representation.NewBuilder().
		WithLocation(*request.URL).
		WithType(html).
		WithLanguage(english).
		WithEncoding(gzip).
		WithCharset(ascii).
		WithSourceQuality(1.0).
		WithFeature(features).
		Build(test.RepresentationBuilderFunc)
	variants := []representation.Representation{v1, v2}

	// action.
	chosen, err := s.sut.Choose(request, variants...)

	// assert.
	s.Require().NoError(err)
	s.Require().Nil(chosen)
}

func (s *RVSATestSuite) TearDownTest() {
	s.sut = nil
}
