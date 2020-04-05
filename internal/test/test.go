package test

import (
	_representation "github.com/freerware/negotiator/internal/representation"
	"github.com/freerware/negotiator/representation"
)

// Representation represents a representation to utilize in unit tests.
type Representation struct {
	representation.Base

	A string
	B int
}

// Bytes serializes the test representation.
func (r Representation) Bytes() ([]byte, error) {
	return r.Base.Bytes(&r)
}

// FromBytes deserializes the test representation.
func (r Representation) FromBytes(b []byte) error {
	return r.Base.FromBytes(b, &r)
}

var (
	RepresentationBuilderFunc = func(ctx _representation.BuilderContext) representation.Representation {
		r := Representation{}
		r.SetContentType(ctx.ContentType)
		r.SetContentLanguage(ctx.ContentLanguage)
		r.SetContentCharset(ctx.ContentCharset)
		r.SetContentEncoding(ctx.ContentEncoding)
		r.SetContentLocation(ctx.ContentLocation)
		r.SetContentFeatures(ctx.ContentFeatures)
		r.SetSourceQuality(ctx.SourceQuality)
		return r
	}
)
