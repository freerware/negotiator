package representation

import (
	"net/url"

	rep "github.com/freerware/negotiator/representation"
)

// Builder constructs a new representation builder.
func NewBuilder() Builder {
	return Builder{}
}

// Builder represents a representation builder.
type Builder struct {
	ct  string
	cl  string
	ce  []string
	cc  string
	cf  []string
	loc url.URL
	sq  float32
}

// WithType associates the provided content type with the representation to be built.
func (b Builder) WithType(ct string) Builder {
	b.ct = ct
	return b
}

// WithLanguage associates the provided language with the representation to be built.
func (b Builder) WithLanguage(cl string) Builder {
	b.cl = cl
	return b
}

// WithEncoding associates the provided encoding with the representation to be built.
func (b Builder) WithEncoding(ce string) Builder {
	b.ce = append(b.ce, ce)
	return b
}

// WithCharset associates the provided charset with the representation to be built.
func (b Builder) WithCharset(cc string) Builder {
	b.cc = cc
	return b
}

// WithLocation associates the provided content location with the
// representation to be built
func (b Builder) WithLocation(loc url.URL) Builder {
	b.loc = loc
	return b
}

// WithSourceQuality associates the provided source quality with the
// representation to be built.
func (b Builder) WithSourceQuality(sq float32) Builder {
	b.sq = sq
	return b
}

// WithFeature associates the provided feature with the representation to be built.
func (b Builder) WithFeature(cf string) Builder {
	b.cf = append(b.cf, cf)
	return b
}

// Build builds the representation.
func (b Builder) Build(bf BuilderFunc) rep.Representation {
	ctx := BuilderContext{
		ContentType:     b.ct,
		ContentLanguage: b.cl,
		ContentEncoding: b.ce,
		ContentCharset:  b.cc,
		ContentLocation: b.loc,
		ContentFeatures: b.cf,
		SourceQuality:   b.sq,
	}
	return bf(ctx)
}

// BuilderContext represents the context used to build the representation.
type BuilderContext struct {
	ContentType     string
	ContentLanguage string
	ContentEncoding []string
	ContentCharset  string
	ContentFeatures []string
	ContentLocation url.URL
	SourceQuality   float32
}

// BuilderFunc represents the function for creating the representation using
// the builder context.
type BuilderFunc func(ctx BuilderContext) rep.Representation
