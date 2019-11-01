// Package transparent implements transparent content negotiation as defined in
// RFC2295.
//
//Construction
//
// For out of the box transparent negotiation support, use
// transparent.Default, which is the default transparent negotiator.
//	//retrieves the default transparent negotiator.
//	p := transparent.Default
// In situations where more customization is required, use the
// transparent.New constructor function and specify options as arguments.
//	constructs a transparent negotiator with the provided options.
//	p := transparent.New(
//		transparent.MaximumVariantListSize(5),
//	)
//
// See Also
//
// ➣ https://tools.ietf.org/html/rfc2295
package transparent
