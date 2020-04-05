package proactive_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/freerware/negotiator"
	_representation "github.com/freerware/negotiator/internal/representation"
	"github.com/freerware/negotiator/internal/representation/json"
	"github.com/freerware/negotiator/internal/representation/xml"
	"github.com/freerware/negotiator/internal/test"
	"github.com/freerware/negotiator/internal/test/mock"
	"github.com/freerware/negotiator/proactive"
	"github.com/freerware/negotiator/representation"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"github.com/uber-go/tally"
	"go.uber.org/zap"
)

type ProactiveTestSuite struct {
	suite.Suite

	// system under test.
	sut negotiator.Negotiator

	// mocks.
	mc      *gomock.Controller
	chooser *mock.Chooser
}

func TestProactiveTestSuite(t *testing.T) {
	suite.Run(t, new(ProactiveTestSuite))
}

func (s *ProactiveTestSuite) SetupTest() {
	s.mc = gomock.NewController(s.T())
	s.chooser = mock.NewChooser(s.mc)
	s.sut = proactive.New(
		proactive.Algorithm(s.chooser),
		proactive.Representations(json.List),
		proactive.Scope(tally.NoopScope),
		proactive.Logger(zap.NewNop()),
	)
}

func (s ProactiveTestSuite) TestProactive_NoRepresentations() {
	// arrange.
	request := httptest.NewRequest("GET", "http://freer.ddns.net/thing", nil)
	responseWriter := httptest.NewRecorder()
	request.Header.Add("Accept", "application/json")
	ctx := negotiator.NegotiationContext{Request: request, ResponseWriter: responseWriter}
	variants := []representation.Representation{}

	// action.
	err := s.sut.Negotiate(ctx, variants...)

	// assert.
	s.Require().NoError(err)
	response := responseWriter.Result()
	s.Equal(http.StatusNoContent, response.StatusCode)
	s.Equal(0, responseWriter.Body.Len())
}

func (s ProactiveTestSuite) TestProactive_InvalidAccept() {
	// arrange.
	_json, english, ascii, gzip := "application/json", "en-US", "ascii", "gzip"
	request := httptest.NewRequest("GET", "http://freer.ddns.net/thing", nil)
	responseWriter := httptest.NewRecorder()
	request.Header.Add("Accept", "invalid/")
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
	s.Require().Error(err)
	s.Equal(0, responseWriter.Body.Len())
}

func (s ProactiveTestSuite) TestProactive_StrictMode_MissingAccept() {
	// arrange.
	s.sut = proactive.New(proactive.Algorithm(s.chooser))
	_json, english, ascii, gzip := "application/json", "en-US", "ascii", "gzip"
	request := httptest.NewRequest("GET", "http://freer.ddns.net/thing", nil)
	responseWriter := httptest.NewRecorder()
	request.Header.Add("Accept-Language", english)
	request.Header.Add("Accept-Encoding", gzip)
	request.Header.Add("Accept-Charset", ascii)
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
	jsonList := json.List(variants...)
	expectedBytes, _ := jsonList.Bytes()
	expectedLen := len(expectedBytes)
	s.chooser.EXPECT().Choose(ctx.Request, gomock.Any()).Return(jsonList, nil)

	// action.
	err := s.sut.Negotiate(ctx, variants...)

	// assert.
	s.Require().NoError(err)
	response := responseWriter.Result()
	s.Equal(http.StatusOK, response.StatusCode)
	s.Equal(expectedLen, responseWriter.Body.Len())
	s.Equal(_json, response.Header.Get("Content-Type"))
}

func (s ProactiveTestSuite) TestProactive_AcceptStrictModeDisabled_MissingAccept() {
	// arrange.
	s.sut = proactive.New(proactive.Algorithm(s.chooser), proactive.DisableStrictAccept())
	_json, english, ascii, gzip := "application/json", "en-US", "ascii", "gzip"
	request := httptest.NewRequest("GET", "http://freer.ddns.net/thing", nil)
	responseWriter := httptest.NewRecorder()
	request.Header.Add("Accept-Language", english)
	request.Header.Add("Accept-Encoding", gzip)
	request.Header.Add("Accept-Charset", ascii)
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
	jsonList := json.List(variants...)
	expectedBytes, _ := jsonList.Bytes()
	expectedLen := len(expectedBytes)
	s.chooser.EXPECT().Choose(ctx.Request, gomock.Any()).Return(jsonList, nil)

	// action.
	err := s.sut.Negotiate(ctx, variants...)

	// assert.
	s.Require().NoError(err)
	response := responseWriter.Result()
	s.Equal(http.StatusOK, response.StatusCode)
	s.Equal(expectedLen, responseWriter.Body.Len())
	s.Equal(_json, response.Header.Get("Content-Type"))
}

func (s ProactiveTestSuite) TestProactive_StrictMode_NoMatchesForAccept() {
	// arrange.
	s.sut = proactive.New(proactive.Algorithm(s.chooser))
	_json, _xml, english, ascii, gzip :=
		"application/json", "application/xml", "en-US", "ascii", "gzip"
	request := httptest.NewRequest("GET", "http://freer.ddns.net/thing", nil)
	responseWriter := httptest.NewRecorder()
	request.Header.Add("Accept", _json)
	request.Header.Add("Accept-Encoding", gzip)
	request.Header.Add("Accept-Language", english)
	request.Header.Add("Accept-Charset", ascii)
	ctx := negotiator.NegotiationContext{Request: request, ResponseWriter: responseWriter}
	v := _representation.NewBuilder().
		WithLocation(*request.URL).
		WithType(_xml).
		WithLanguage(english).
		WithEncoding(gzip).
		WithCharset(ascii).
		WithSourceQuality(1.0).
		Build(test.RepresentationBuilderFunc)
	variants := []representation.Representation{v}
	jsonList := json.List(variants...)
	expectedBytes, _ := jsonList.Bytes()
	expectedLen := len(expectedBytes)
	s.chooser.EXPECT().Choose(ctx.Request, gomock.Any()).Return(jsonList, nil)

	// action.
	err := s.sut.Negotiate(ctx, variants...)

	// assert.
	s.Require().NoError(err)
	response := responseWriter.Result()
	s.Equal(http.StatusNotAcceptable, response.StatusCode)
	s.Equal(expectedLen, responseWriter.Body.Len())
	s.Equal(_json, response.Header.Get("Content-Type"))
}

func (s ProactiveTestSuite) TestProactive_AcceptStrictModeDisabled_NoMatchesForAccept() {
	// arrange.
	s.sut = proactive.New(proactive.Algorithm(s.chooser), proactive.DisableStrictAccept())
	_json, _xml, english, ascii, gzip :=
		"application/json", "application/xml", "en-US", "ascii", "gzip"
	request := httptest.NewRequest("GET", "http://freer.ddns.net/thing", nil)
	responseWriter := httptest.NewRecorder()
	request.Header.Add("Accept", _json)
	request.Header.Add("Accept-Encoding", gzip)
	request.Header.Add("Accept-Language", english)
	request.Header.Add("Accept-Charset", ascii)
	ctx := negotiator.NegotiationContext{Request: request, ResponseWriter: responseWriter}
	v := _representation.NewBuilder().
		WithLocation(*request.URL).
		WithType(_xml).
		WithLanguage(english).
		WithEncoding(gzip).
		WithCharset(ascii).
		WithSourceQuality(1.0).
		Build(test.RepresentationBuilderFunc)
	variants := []representation.Representation{v}
	jsonList := json.List(variants...)
	expectedBytes, _ := jsonList.Bytes()
	expectedLen := len(expectedBytes)
	s.chooser.EXPECT().Choose(ctx.Request, gomock.Any()).Return(jsonList, nil)

	// action.
	err := s.sut.Negotiate(ctx, variants...)

	// assert.
	s.Require().NoError(err)
	response := responseWriter.Result()
	s.Equal(http.StatusOK, response.StatusCode)
	s.Equal(expectedLen, responseWriter.Body.Len())
	s.Equal(_json, response.Header.Get("Content-Type"))
}

func (s ProactiveTestSuite) TestProactive_InvalidAcceptLanguage() {
	// arrange.
	_json, english, ascii, gzip := "application/json", "en-US", "ascii", "gzip"
	request := httptest.NewRequest("GET", "http://freer.ddns.net/thing", nil)
	responseWriter := httptest.NewRecorder()
	request.Header.Add("Accept", "application/json")
	request.Header.Add("Accept-Encoding", "gzip")
	request.Header.Add("Accept-Language", "invalid")
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
	s.Require().Error(err)
	s.Equal(0, responseWriter.Body.Len())
}

func (s ProactiveTestSuite) TestProactive_StrictMode_MissingAcceptLanguage() {
	// arrange.
	s.sut = proactive.New(proactive.Algorithm(s.chooser))
	_json, english, ascii, gzip := "application/json", "en-US", "ascii", "gzip"
	request := httptest.NewRequest("GET", "http://freer.ddns.net/thing", nil)
	responseWriter := httptest.NewRecorder()
	request.Header.Add("Accept", _json)
	request.Header.Add("Accept-Encoding", gzip)
	request.Header.Add("Accept-Charset", ascii)
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
	jsonList := json.List(variants...)
	expectedBytes, _ := jsonList.Bytes()
	expectedLen := len(expectedBytes)
	s.chooser.EXPECT().Choose(ctx.Request, gomock.Any()).Return(jsonList, nil)

	// action.
	err := s.sut.Negotiate(ctx, variants...)

	// assert.
	s.Require().NoError(err)
	response := responseWriter.Result()
	s.Equal(http.StatusOK, response.StatusCode)
	s.Equal(expectedLen, responseWriter.Body.Len())
	s.Equal(_json, response.Header.Get("Content-Type"))
}

func (s ProactiveTestSuite) TestProactive_AcceptLanguageStrictModeDisabled_MissingAcceptLanguage() {
	// arrange.
	s.sut = proactive.New(proactive.Algorithm(s.chooser), proactive.DisableStrictAcceptLanguage())
	_json, english, ascii, gzip := "application/json", "en-US", "ascii", "gzip"
	request := httptest.NewRequest("GET", "http://freer.ddns.net/thing", nil)
	responseWriter := httptest.NewRecorder()
	request.Header.Add("Accept", _json)
	request.Header.Add("Accept-Encoding", gzip)
	request.Header.Add("Accept-Charset", ascii)
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
	jsonList := json.List(variants...)
	expectedBytes, _ := jsonList.Bytes()
	expectedLen := len(expectedBytes)
	s.chooser.EXPECT().Choose(ctx.Request, gomock.Any()).Return(jsonList, nil)

	// action.
	err := s.sut.Negotiate(ctx, variants...)

	// assert.
	s.Require().NoError(err)
	response := responseWriter.Result()
	s.Equal(http.StatusOK, response.StatusCode)
	s.Equal(expectedLen, responseWriter.Body.Len())
	s.Equal(_json, response.Header.Get("Content-Type"))
}

func (s ProactiveTestSuite) TestProactive_StrictMode_NoMatchesForAcceptLanguage() {
	// arrange.
	s.sut = proactive.New(proactive.Algorithm(s.chooser))
	_json, english, french, ascii, gzip := "application/json", "en-US", "fr", "ascii", "gzip"
	request := httptest.NewRequest("GET", "http://freer.ddns.net/thing", nil)
	responseWriter := httptest.NewRecorder()
	request.Header.Add("Accept", _json)
	request.Header.Add("Accept-Language", french)
	request.Header.Add("Accept-Encoding", gzip)
	request.Header.Add("Accept-Charset", ascii)
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
	jsonList := json.List(variants...)
	expectedBytes, _ := jsonList.Bytes()
	expectedLen := len(expectedBytes)
	s.chooser.EXPECT().Choose(ctx.Request, gomock.Any()).Return(jsonList, nil)

	// action.
	err := s.sut.Negotiate(ctx, variants...)

	// assert.
	s.Require().NoError(err)
	response := responseWriter.Result()
	s.Equal(http.StatusNotAcceptable, response.StatusCode)
	s.Equal(expectedLen, responseWriter.Body.Len())
	s.Equal(_json, response.Header.Get("Content-Type"))
}

func (s ProactiveTestSuite) TestProactive_AcceptLanguageStrictModeDisabled_NoMatchesForAcceptLanguage() {
	// arrange.
	s.sut = proactive.New(proactive.Algorithm(s.chooser), proactive.DisableStrictAcceptLanguage())
	_json, english, french, ascii, gzip := "application/json", "en-US", "fr", "ascii", "gzip"
	request := httptest.NewRequest("GET", "http://freer.ddns.net/thing", nil)
	responseWriter := httptest.NewRecorder()
	request.Header.Add("Accept", _json)
	request.Header.Add("Accept-Language", french)
	request.Header.Add("Accept-Encoding", gzip)
	request.Header.Add("Accept-Charset", ascii)
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
	jsonList := json.List(variants...)
	expectedBytes, _ := jsonList.Bytes()
	expectedLen := len(expectedBytes)
	s.chooser.EXPECT().Choose(ctx.Request, gomock.Any()).Return(jsonList, nil)

	// action.
	err := s.sut.Negotiate(ctx, variants...)

	// assert.
	s.Require().NoError(err)
	response := responseWriter.Result()
	s.Equal(http.StatusOK, response.StatusCode)
	s.Equal(expectedLen, responseWriter.Body.Len())
	s.Equal(_json, response.Header.Get("Content-Type"))
}

func (s ProactiveTestSuite) TestProactive_InvalidAcceptCharset() {
	// arrange.
	_json, english, ascii, gzip := "application/json", "en-US", "ascii", "gzip"
	request := httptest.NewRequest("GET", "http://freer.ddns.net/thing", nil)
	responseWriter := httptest.NewRecorder()
	request.Header.Add("Accept", "application/json")
	request.Header.Add("Accept-Encoding", "gzip")
	request.Header.Add("Accept-Language", "en-US")
	request.Header.Add("Accept-Charset", "")
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
	s.Require().Error(err)
	s.Equal(0, responseWriter.Body.Len())
}

func (s ProactiveTestSuite) TestProactive_StrictMode_MissingAcceptCharset() {
	// arrange.
	s.sut = proactive.New(proactive.Algorithm(s.chooser))
	_json, english, ascii, gzip := "application/json", "en-US", "ascii", "gzip"
	request := httptest.NewRequest("GET", "http://freer.ddns.net/thing", nil)
	responseWriter := httptest.NewRecorder()
	request.Header.Add("Accept", _json)
	request.Header.Add("Accept-Encoding", gzip)
	request.Header.Add("Accept-Language", english)
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
	jsonList := json.List(variants...)
	expectedBytes, _ := jsonList.Bytes()
	expectedLen := len(expectedBytes)
	s.chooser.EXPECT().Choose(ctx.Request, gomock.Any()).Return(jsonList, nil)

	// action.
	err := s.sut.Negotiate(ctx, variants...)

	// assert.
	s.Require().NoError(err)
	response := responseWriter.Result()
	s.Equal(http.StatusOK, response.StatusCode)
	s.Equal(expectedLen, responseWriter.Body.Len())
	s.Equal(_json, response.Header.Get("Content-Type"))
}

func (s ProactiveTestSuite) TestProactive_AcceptCharsetStrictModeDisableld_MissingAcceptCharset() {
	// arrange.
	s.sut = proactive.New(proactive.Algorithm(s.chooser), proactive.DisableStrictAcceptCharset())
	_json, english, ascii, gzip := "application/json", "en-US", "ascii", "gzip"
	request := httptest.NewRequest("GET", "http://freer.ddns.net/thing", nil)
	responseWriter := httptest.NewRecorder()
	request.Header.Add("Accept", _json)
	request.Header.Add("Accept-Encoding", gzip)
	request.Header.Add("Accept-Language", english)
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
	jsonList := json.List(variants...)
	expectedBytes, _ := jsonList.Bytes()
	expectedLen := len(expectedBytes)
	s.chooser.EXPECT().Choose(ctx.Request, gomock.Any()).Return(jsonList, nil)

	// action.
	err := s.sut.Negotiate(ctx, variants...)

	// assert.
	s.Require().NoError(err)
	response := responseWriter.Result()
	s.Equal(http.StatusOK, response.StatusCode)
	s.Equal(expectedLen, responseWriter.Body.Len())
	s.Equal(_json, response.Header.Get("Content-Type"))
}

func (s ProactiveTestSuite) TestProactive_StrictMode_NoMatchesForAcceptCharset() {
	// arrange.
	s.sut = proactive.New(proactive.Algorithm(s.chooser))
	_json, english, ascii, utf8, gzip :=
		"application/json", "en-US", "ascii", "utf8", "gzip"
	request := httptest.NewRequest("GET", "http://freer.ddns.net/thing", nil)
	responseWriter := httptest.NewRecorder()
	request.Header.Add("Accept", _json)
	request.Header.Add("Accept-Encoding", gzip)
	request.Header.Add("Accept-Language", english)
	request.Header.Add("Accept-Charset", ascii)
	ctx := negotiator.NegotiationContext{Request: request, ResponseWriter: responseWriter}
	v := _representation.NewBuilder().
		WithLocation(*request.URL).
		WithType(_json).
		WithLanguage(english).
		WithEncoding(gzip).
		WithCharset(utf8).
		WithSourceQuality(1.0).
		Build(test.RepresentationBuilderFunc)
	variants := []representation.Representation{v}
	jsonList := json.List(variants...)
	expectedBytes, _ := jsonList.Bytes()
	expectedLen := len(expectedBytes)
	s.chooser.EXPECT().Choose(ctx.Request, gomock.Any()).Return(jsonList, nil)

	// action.
	err := s.sut.Negotiate(ctx, variants...)

	// assert.
	s.Require().NoError(err)
	response := responseWriter.Result()
	s.Equal(http.StatusNotAcceptable, response.StatusCode)
	s.Equal(expectedLen, responseWriter.Body.Len())
	s.Equal(_json, response.Header.Get("Content-Type"))
}

func (s ProactiveTestSuite) TestProactive_AcceptCharsetStrictModeDisabled_NoMatchesForAcceptCharset() {
	// arrange.
	s.sut = proactive.New(proactive.Algorithm(s.chooser), proactive.DisableStrictAcceptCharset())
	_json, english, ascii, utf8, gzip :=
		"application/json", "en-US", "ascii", "utf8", "gzip"
	request := httptest.NewRequest("GET", "http://freer.ddns.net/thing", nil)
	responseWriter := httptest.NewRecorder()
	request.Header.Add("Accept", _json)
	request.Header.Add("Accept-Encoding", gzip)
	request.Header.Add("Accept-Language", english)
	request.Header.Add("Accept-Charset", ascii)
	ctx := negotiator.NegotiationContext{Request: request, ResponseWriter: responseWriter}
	v := _representation.NewBuilder().
		WithLocation(*request.URL).
		WithType(_json).
		WithLanguage(english).
		WithEncoding(gzip).
		WithCharset(utf8).
		WithSourceQuality(1.0).
		Build(test.RepresentationBuilderFunc)
	variants := []representation.Representation{v}
	jsonList := json.List(variants...)
	expectedBytes, _ := jsonList.Bytes()
	expectedLen := len(expectedBytes)
	s.chooser.EXPECT().Choose(ctx.Request, gomock.Any()).Return(jsonList, nil)

	// action.
	err := s.sut.Negotiate(ctx, variants...)

	// assert.
	s.Require().NoError(err)
	response := responseWriter.Result()
	s.Equal(http.StatusOK, response.StatusCode)
	s.Equal(expectedLen, responseWriter.Body.Len())
	s.Equal(_json, response.Header.Get("Content-Type"))
}

func (s ProactiveTestSuite) TestProactive_NotAcceptable() {
	// arrange.
	_json, english, ascii, gzip := "application/json", "en-US", "ascii", "gzip"
	request := httptest.NewRequest("GET", "http://freer.ddns.net/thing", nil)
	responseWriter := httptest.NewRecorder()
	request.Header.Add("Accept", _json)
	request.Header.Add("Accept-Language", english)
	request.Header.Add("Accept-Encoding", gzip)
	request.Header.Add("Accept-Language", english)
	request.Header.Add("Accept-Charset", ascii)
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
	jsonList := json.List(variants...)
	expectedBytes, _ := jsonList.Bytes()
	expectedLen := len(expectedBytes)
	s.chooser.
		EXPECT().
		Choose(ctx.Request, gomock.Any()).
		Return(nil, nil).
		Times(1)
	s.chooser.
		EXPECT().
		Choose(ctx.Request, gomock.Any()).
		Return(jsonList, nil).
		Times(1)

	// action.
	err := s.sut.Negotiate(ctx, variants...)

	// assert.
	s.Require().NoError(err)
	response := responseWriter.Result()
	s.Equal(http.StatusNotAcceptable, response.StatusCode)
	s.Equal(expectedLen, responseWriter.Body.Len())
	s.Equal(_json, response.Header.Get("Content-Type"))
}

func (s ProactiveTestSuite) TestProactive_NotAcceptable_NoRepresentation() {
	// arrange.
	s.sut = proactive.New(
		proactive.Algorithm(s.chooser),
		proactive.DisableNotAcceptableRepresentation(),
		proactive.DisableStrictMode(),
	)
	_json, english, ascii, gzip := "application/json", "en-US", "ascii", "gzip"
	request := httptest.NewRequest("GET", "http://freer.ddns.net/thing", nil)
	responseWriter := httptest.NewRecorder()
	request.Header.Add("Accept", _json)
	request.Header.Add("Accept-Language", english)
	request.Header.Add("Accept-Encoding", gzip)
	request.Header.Add("Accept-Language", english)
	request.Header.Add("Accept-Charset", ascii)
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
	s.chooser.
		EXPECT().
		Choose(ctx.Request, gomock.Any()).
		Return(nil, nil).
		Times(1)

	// action.
	err := s.sut.Negotiate(ctx, variants...)

	// assert.
	s.Require().NoError(err)
	response := responseWriter.Result()
	s.Equal(http.StatusNotAcceptable, response.StatusCode)
	s.Zero(responseWriter.Body.Len())
}

func (s ProactiveTestSuite) TestProactive_NotAcceptable_DefaultRepresentation() {
	// arrange.
	s.sut = proactive.New(
		proactive.Algorithm(s.chooser),
		proactive.DisableStrictMode(),
		proactive.DefaultRepresentation(xml.List),
	)
	_json, english, ascii, gzip := "application/json", "en-US", "ascii", "gzip"
	request := httptest.NewRequest("GET", "http://freer.ddns.net/thing", nil)
	responseWriter := httptest.NewRecorder()
	request.Header.Add("Accept", _json)
	request.Header.Add("Accept-Language", english)
	request.Header.Add("Accept-Encoding", gzip)
	request.Header.Add("Accept-Language", english)
	request.Header.Add("Accept-Charset", ascii)
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
	xmlList := xml.List(variants...)
	expectedBytes, _ := xmlList.Bytes()
	expectedLen := len(expectedBytes)
	s.chooser.
		EXPECT().
		Choose(ctx.Request, gomock.Any()).
		Return(nil, nil).
		Times(2)

	// action.
	err := s.sut.Negotiate(ctx, variants...)

	// assert.
	s.Require().NoError(err)
	response := responseWriter.Result()
	s.Equal(http.StatusNotAcceptable, response.StatusCode)
	s.Equal(expectedLen, responseWriter.Body.Len())
	s.Equal(xmlList.ContentType(), response.Header.Get("Content-Type"))
}

func (s ProactiveTestSuite) TestProactive_NotAcceptable_ChooserError() {
	// arrange.
	_json, english, ascii, gzip := "application/json", "en-US", "ascii", "gzip"
	request := httptest.NewRequest("GET", "http://freer.ddns.net/thing", nil)
	responseWriter := httptest.NewRecorder()
	request.Header.Add("Accept", _json)
	request.Header.Add("Accept-Language", english)
	request.Header.Add("Accept-Encoding", gzip)
	request.Header.Add("Accept-Language", english)
	request.Header.Add("Accept-Charset", ascii)
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
	errExpected := errors.New("whoa")
	s.chooser.
		EXPECT().
		Choose(ctx.Request, gomock.Any()).
		Return(nil, nil).
		Times(1)
	s.chooser.
		EXPECT().
		Choose(ctx.Request, gomock.Any()).
		Return(nil, errExpected).
		Times(1)

	// action.
	err := s.sut.Negotiate(ctx, variants...)

	// assert.
	s.Require().Error(err)
}

func (s ProactiveTestSuite) TestProactive() {
	// arrange.
	_json, english, ascii, gzip := "application/json", "en-US", "ascii", "gzip"
	request := httptest.NewRequest("GET", "http://freer.ddns.net/thing", nil)
	responseWriter := httptest.NewRecorder()
	request.Header.Add("Accept", _json)
	request.Header.Add("Accept-Language", english)
	request.Header.Add("Accept-Encoding", gzip)
	request.Header.Add("Accept-Language", english)
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
	expectedBytes, _ := v.Bytes()
	expectedLen := len(expectedBytes)
	s.chooser.
		EXPECT().
		Choose(ctx.Request, gomock.Any()).
		Return(v, nil).
		Times(1)

	// action.
	err := s.sut.Negotiate(ctx, variants...)

	// assert.
	s.Require().NoError(err)
	response := responseWriter.Result()
	s.Equal(http.StatusOK, response.StatusCode)
	s.Equal(expectedLen, responseWriter.Body.Len())
	s.Equal(_json, response.Header.Get("Content-Type"))
}

func (s ProactiveTestSuite) TestProactive_ChooseError() {
	// arrange.
	_json, english, ascii, gzip := "application/json", "en-US", "ascii", "gzip"
	request := httptest.NewRequest("GET", "http://freer.ddns.net/thing", nil)
	responseWriter := httptest.NewRecorder()
	request.Header.Add("Accept", _json)
	request.Header.Add("Accept-Language", english)
	request.Header.Add("Accept-Encoding", gzip)
	request.Header.Add("Accept-Language", english)
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
	errExpected := errors.New("whoa")
	s.chooser.
		EXPECT().
		Choose(ctx.Request, gomock.Any()).
		Return(nil, errExpected).
		Times(1)

	// action.
	err := s.sut.Negotiate(ctx, variants...)

	// assert.
	s.Require().Error(err)
}

func (s ProactiveTestSuite) TestProactive_WriteError() {
	// arrange.
	_json, english, ascii, gzip := "application/json", "en-US", "ascii", "gzip"
	request := httptest.NewRequest("GET", "http://freer.ddns.net/thing", nil)
	responseWriter := mock.NewMockResponseWriter(s.mc)
	request.Header.Add("Accept", _json)
	request.Header.Add("Accept-Language", english)
	request.Header.Add("Accept-Encoding", gzip)
	request.Header.Add("Accept-Language", english)
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
	bytes, _ := v.Bytes()
	s.chooser.
		EXPECT().
		Choose(ctx.Request, gomock.Any()).
		Return(v, nil).
		Times(1)
	errExpected := errors.New("whoa")
	responseWriter.EXPECT().Write(bytes).Return(0, errExpected).Times(1)
	responseWriter.EXPECT().WriteHeader(http.StatusOK).Times(1)
	responseWriter.EXPECT().Header().Return(make(http.Header)).AnyTimes()

	// action.
	err := s.sut.Negotiate(ctx, variants...)

	// assert.
	s.Require().Error(err)
}

func (s ProactiveTestSuite) TestProactive_IsCreation() {
	// arrange.
	_json, english, ascii, gzip := "application/json", "en-US", "ascii", "gzip"
	request := httptest.NewRequest("GET", "http://freer.ddns.net/thing", nil)
	responseWriter := httptest.NewRecorder()
	request.Header.Add("Accept", _json)
	request.Header.Add("Accept-Language", english)
	request.Header.Add("Accept-Encoding", gzip)
	request.Header.Add("Accept-Language", english)
	ctx := negotiator.NegotiationContext{
		Request:        request,
		ResponseWriter: responseWriter,
		IsCreation:     true,
	}
	v := _representation.NewBuilder().
		WithLocation(*request.URL).
		WithType(_json).
		WithLanguage(english).
		WithEncoding(gzip).
		WithCharset(ascii).
		WithSourceQuality(1.0).
		Build(test.RepresentationBuilderFunc)
	variants := []representation.Representation{v}
	expectedBytes, _ := v.Bytes()
	expectedLen := len(expectedBytes)
	s.chooser.
		EXPECT().
		Choose(ctx.Request, gomock.Any()).
		Return(v, nil).
		Times(1)

	// action.
	err := s.sut.Negotiate(ctx, variants...)

	// assert.
	s.Require().NoError(err)
	response := responseWriter.Result()
	s.Equal(http.StatusCreated, response.StatusCode)
	s.Equal(expectedLen, responseWriter.Body.Len())
	s.Equal(_json, response.Header.Get("Content-Type"))
}

func (s *ProactiveTestSuite) TearDownTest() {
	s.mc.Finish()
	s.mc = nil
	s.sut = nil
	s.chooser = nil
}
