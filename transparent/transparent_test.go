package transparent_test

import (
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/freerware/negotiator"
	"github.com/freerware/negotiator/internal/header"
	_representation "github.com/freerware/negotiator/internal/representation"
	"github.com/freerware/negotiator/internal/test"
	"github.com/freerware/negotiator/internal/test/mock"
	"github.com/freerware/negotiator/representation"
	"github.com/freerware/negotiator/transparent"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

type TransparentTestSuite struct {
	suite.Suite

	// system under test.
	sut negotiator.Negotiator

	// mocks.
	mc      *gomock.Controller
	chooser *mock.Chooser
}

func TestTransparentTestSuite(t *testing.T) {
	suite.Run(t, new(TransparentTestSuite))
}

func (s *TransparentTestSuite) SetupTest() {
	s.mc = gomock.NewController(s.T())
	s.chooser = mock.NewChooser(s.mc)
	s.sut = transparent.New(
		transparent.RVSA(s.chooser),
		transparent.MaximumVariantListSize(3),
	)
}

func (s TransparentTestSuite) TestTransparent_WildcardNegotiateHeader() {
	// arrange.
	_json, english, ascii, gzip := "application/json", "en-US", "ascii", "gzip"
	request := httptest.NewRequest("GET", "http://freer.ddns.net/thing", nil)
	request.Header.Add("Negotiate", "*")
	responseWriter := httptest.NewRecorder()
	ctx := negotiator.NegotiationContext{Request: request, ResponseWriter: responseWriter}
	v := _representation.NewBuilder().
		WithLocation(*request.URL).
		WithType(_json).
		WithLanguage(english).
		WithEncoding(gzip).
		WithCharset(ascii).
		WithSourceQuality(1.0).
		Build(test.RepresentationBuilderFunc)
	variants := []representation.Representation{v}
	s.chooser.EXPECT().Choose(ctx.Request, gomock.Any()).Return(v, nil)

	// action.
	err := s.sut.Negotiate(ctx, variants...)

	// assert.
	s.Require().NoError(err)
	response := responseWriter.Result()
	s.Equal(v.ContentType(), response.Header.Get("Content-Type"))
	s.ElementsMatch(v.ContentEncoding(), response.Header["Content-Encoding"])
	s.Equal(v.ContentCharset(), response.Header.Get("Content-Charset"))
	loc := v.ContentLocation()
	s.Equal(loc.String(), response.Header.Get("Content-Location"))
	s.NotEmpty(response.Header.Get("Alternates"))
	s.Equal(header.ResponseTypeChoice.String(), response.Header.Get("TCN"))
}

func (s TransparentTestSuite) TestTransparent_GuessSmallNegotiateHeader() {
	// arrange.
	_json, english, ascii, gzip := "application/json", "en-US", "ascii", "gzip"

	request := httptest.NewRequest("GET", "http://freer.ddns.net/thing", nil)
	request.Header.Add("Negotiate", "guess-small")
	responseWriter := httptest.NewRecorder()
	ctx := negotiator.NegotiationContext{Request: request, ResponseWriter: responseWriter}
	v := _representation.NewBuilder().
		WithLocation(*request.URL).
		WithType(_json).
		WithLanguage(english).
		WithEncoding(gzip).
		WithCharset(ascii).
		WithSourceQuality(1.0).
		Build(test.RepresentationBuilderFunc)
	variants := []representation.Representation{v}
	s.chooser.EXPECT().Choose(ctx.Request, gomock.Any()).Return(v, nil)

	// action.
	err := s.sut.Negotiate(ctx, variants...)

	// assert.
	s.Require().NoError(err)
	response := responseWriter.Result()
	s.NotEmpty(response.Header.Get("Alternates"))
	s.Equal(header.ResponseTypeChoice.String(), response.Header.Get("TCN"))
}

func (s TransparentTestSuite) TestTransparent_RSVA1NegotiateHeader() {
	// arrange.
	_json, english, ascii, gzip := "application/json", "en-US", "ascii", "gzip"

	request := httptest.NewRequest("GET", "http://freer.ddns.net/thing", nil)
	request.Header.Add("Negotiate", "1.0")
	responseWriter := httptest.NewRecorder()
	ctx := negotiator.NegotiationContext{Request: request, ResponseWriter: responseWriter}
	v := _representation.NewBuilder().
		WithLocation(*request.URL).
		WithType(_json).
		WithLanguage(english).
		WithEncoding(gzip).
		WithCharset(ascii).
		WithSourceQuality(1.0).
		Build(test.RepresentationBuilderFunc)
	variants := []representation.Representation{v}
	s.chooser.EXPECT().Choose(ctx.Request, gomock.Any()).Return(v, nil)

	// action.
	err := s.sut.Negotiate(ctx, variants...)

	// assert.
	s.Require().NoError(err)
	response := responseWriter.Result()
	s.Equal(v.ContentType(), response.Header.Get("Content-Type"))
	s.ElementsMatch(v.ContentEncoding(), response.Header["Content-Encoding"])
	s.Equal(v.ContentCharset(), response.Header.Get("Content-Charset"))
	loc := v.ContentLocation()
	s.Equal(loc.String(), response.Header.Get("Content-Location"))
	s.NotEmpty(response.Header.Get("Alternates"))
	s.Equal("choice", response.Header.Get("TCN"))
}

func (s TransparentTestSuite) TestTransparent_UnrecognizedRSVANegotiateHeader() {
	// arrange.
	_json, english, ascii, gzip := "application/json", "en-US", "ascii", "gzip"

	request := httptest.NewRequest("GET", "http://freer.ddns.net/thing", nil)
	request.Header.Add("Negotiate", "2.0")
	responseWriter := httptest.NewRecorder()
	ctx := negotiator.NegotiationContext{Request: request, ResponseWriter: responseWriter}
	v := _representation.NewBuilder().
		WithLocation(*request.URL).
		WithType(_json).
		WithLanguage(english).
		WithEncoding(gzip).
		WithCharset(ascii).
		WithSourceQuality(1.0).
		Build(test.RepresentationBuilderFunc)
	variants := []representation.Representation{v}

	// action.
	err := s.sut.Negotiate(ctx, variants...)

	// assert.
	s.Require().NoError(err)
	response := responseWriter.Result()
	s.NotEmpty(response.Header.Get("Alternates"))
	s.Equal("list", response.Header.Get("TCN"))
}

func (s TransparentTestSuite) TestTransparent_NoMatches() {
	// arrange.
	_json, english, ascii, gzip := "application/json", "en-US", "ascii", "gzip"

	request := httptest.NewRequest("GET", "http://freer.ddns.net/thing", nil)
	request.Header.Add("Negotiate", "1.0")
	responseWriter := httptest.NewRecorder()
	ctx := negotiator.NegotiationContext{Request: request, ResponseWriter: responseWriter}
	v := _representation.NewBuilder().
		WithLocation(*request.URL).
		WithType(_json).
		WithLanguage(english).
		WithEncoding(gzip).
		WithCharset(ascii).
		WithSourceQuality(1.0).
		Build(test.RepresentationBuilderFunc)
	variants := []representation.Representation{v}
	s.chooser.EXPECT().Choose(ctx.Request, gomock.Any()).Return(nil, nil)

	// action.
	err := s.sut.Negotiate(ctx, variants...)

	// assert.
	s.Require().NoError(err)
	response := responseWriter.Result()
	s.NotEmpty(response.Header.Get("Alternates"))
	s.Equal("list", response.Header.Get("TCN"))
}

func (s TransparentTestSuite) TestTransparent_ChooseError() {
	// arrange.
	_json, english, ascii, gzip := "application/json", "en-US", "ascii", "gzip"
	request := httptest.NewRequest("GET", "http://freer.ddns.net/thing", nil)
	request.Header.Add("Negotiate", "1.0")
	responseWriter := httptest.NewRecorder()
	ctx := negotiator.NegotiationContext{Request: request, ResponseWriter: responseWriter}
	v := _representation.NewBuilder().
		WithLocation(*request.URL).
		WithType(_json).
		WithLanguage(english).
		WithEncoding(gzip).
		WithCharset(ascii).
		WithSourceQuality(1.0).
		Build(test.RepresentationBuilderFunc)
	variants := []representation.Representation{v}
	expectedErr := errors.New("oh no")
	s.chooser.EXPECT().Choose(ctx.Request, gomock.Any()).Return(nil, expectedErr)

	// action.
	err := s.sut.Negotiate(ctx, variants...)

	// assert.
	s.Require().EqualError(err, expectedErr.Error())
}

func (s TransparentTestSuite) TestTransparent_ChoiceResponse() {
	// arrange.
	_json, english, ascii, gzip := "application/json", "en-US", "ascii", "gzip"
	request := httptest.NewRequest("GET", "http://freer.ddns.net/thing", nil)
	responseWriter := httptest.NewRecorder()
	request.Header.Add("Negotiate", "1.0")
	ctx := negotiator.NegotiationContext{Request: request, ResponseWriter: responseWriter}
	v := _representation.NewBuilder().
		WithLocation(*request.URL).
		WithType(_json).
		WithLanguage(english).
		WithEncoding(gzip).
		WithCharset(ascii).
		WithSourceQuality(1.0).
		Build(test.RepresentationBuilderFunc)
	variants := []representation.Representation{v}
	s.chooser.EXPECT().Choose(ctx.Request, gomock.Any()).Return(v, nil)

	// action.
	err := s.sut.Negotiate(ctx, variants...)

	// assert.
	s.Require().NoError(err)
	response := responseWriter.Result()
	s.Equal(v.ContentType(), response.Header.Get("Content-Type"))
	s.ElementsMatch(v.ContentEncoding(), response.Header["Content-Encoding"])
	s.Equal(v.ContentCharset(), response.Header.Get("Content-Charset"))
	loc := v.ContentLocation()
	s.Equal(loc.String(), response.Header.Get("Content-Location"))
	s.NotEmpty(response.Header.Get("Alternates"))
	s.Equal("choice", response.Header.Get("TCN"))
}

func (s TransparentTestSuite) TestTransparent_VariantListSizeExceeded() {
	// arrange.
	_json, english, ascii, gzip := "application/json", "en-US", "ascii", "gzip"

	request := httptest.NewRequest("GET", "http://freer.ddns.net/thing", nil)
	responseWriter := httptest.NewRecorder()
	request.Header.Add("Negotiate", "1.0")
	ctx := negotiator.NegotiationContext{Request: request, ResponseWriter: responseWriter}
	v := _representation.NewBuilder().
		WithLocation(*request.URL).
		WithType(_json).
		WithLanguage(english).
		WithEncoding(gzip).
		WithCharset(ascii).
		WithSourceQuality(1.0).
		Build(test.RepresentationBuilderFunc)
	variants := []representation.Representation{v, v, v, v}

	// action.
	err := s.sut.Negotiate(ctx, variants...)

	// assert.
	s.Require().EqualError(err, transparent.ErrVariantListSizeExceeded.Error())
}

func (s *TransparentTestSuite) TearDownTest() {
	s.mc.Finish()
	s.mc = nil
	s.sut = nil
	s.chooser = nil
}
