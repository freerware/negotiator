package reactive_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/freerware/negotiator"
	_representation "github.com/freerware/negotiator/internal/representation"
	"github.com/freerware/negotiator/internal/test"
	"github.com/freerware/negotiator/reactive"
	"github.com/freerware/negotiator/representation"
	"github.com/stretchr/testify/suite"
	"github.com/uber-go/tally"
	"go.uber.org/zap"
)

var (
	jsonList = func(reps ...representation.Representation) representation.Representation {
		list := representation.List{}
		list.SetContentType("application/json")
		list.SetContentCharset("ascii")
		list.SetContentEncoding([]string{"identity"})
		list.SetContentLanguage("en-US")
		list.SetRepresentations(reps...)
		return &list
	}
)

type ReactiveTestSuite struct {
	suite.Suite

	// system under test.
	sut negotiator.Negotiator
}

func TestReactiveTestSuite(t *testing.T) {
	suite.Run(t, new(ReactiveTestSuite))
}

func (s *ReactiveTestSuite) SetupTest() {
	s.sut = reactive.New(
		reactive.Scope(tally.NoopScope),
		reactive.Logger(zap.NewNop()),
	)
}

func (s ReactiveTestSuite) TestReactive_NoRepresentations() {
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

func (s ReactiveTestSuite) TestReactive() {
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
	jList := jsonList(variants...)
	expectedBytes, _ := jList.Bytes()
	expectedLen := len(expectedBytes)

	// action.
	err := s.sut.Negotiate(ctx, variants...)

	// assert.
	s.Require().NoError(err)
	response := responseWriter.Result()
	s.Equal(http.StatusMultipleChoices, response.StatusCode)
	s.Equal(expectedLen, responseWriter.Body.Len())
	s.Equal(jList.ContentType(), response.Header.Get("Content-Type"))
}

func (s ReactiveTestSuite) TestReactive_ListConstructor() {
	// arrange.
	s.sut = reactive.New(
		reactive.Representation(jsonList),
		reactive.Scope(tally.NoopScope),
		reactive.Logger(zap.NewNop()),
	)
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
	jList := jsonList(variants...)
	expectedBytes, _ := jList.Bytes()
	expectedLen := len(expectedBytes)

	// action.
	err := s.sut.Negotiate(ctx, variants...)

	// assert.
	s.Require().NoError(err)
	response := responseWriter.Result()
	s.Equal(http.StatusMultipleChoices, response.StatusCode)
	s.Equal(expectedLen, responseWriter.Body.Len())
	s.Equal(jList.ContentType(), response.Header.Get("Content-Type"))
}

func (s *ReactiveTestSuite) TearDownTest() {
	s.sut = nil
}
