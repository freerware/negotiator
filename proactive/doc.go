// Package proactive implements proactive content negotiation as defined in
// RFC7231 Section 3.4.1.
//
// Construction
//
// For out of the box proactive negotiation support, use proactive.Default,
// which is the default proactive negotiator.
//  //retrieves the default proactive negotiator.
//  p := proactive.Default
// In situations where more customization is required, use the proactive.New
// constructor function and specify options as arguments.
//  //constructs a proactive negotiator with the provided options.
//  p := proactive.New(
//		proactive.DisableStrictMode(),
//		proactive.DisableNotAcceptableRepresentation(),
//  )
//
// Strict Mode
//
// According to RFC7231, when none of the representations match the values
// provided for a particular proactive content negotiation header, the origin
// server can honor that header and return 406 Not Acceptable, or disregard
// the header field by treating the resource as if it is not subject to
// content negotiation.
//
// The behavior of honoring the header in these scenarios is what we refer to
// as strict mode. It is possible to configure strict mode for each individual
// proactive negotiation header, or disable strict mode for all. Strict mode
// is enabled for all headers by default.
//
// See Also
//
// ➣ https://tools.ietf.org/html/rfc7231#section-3.4.1
//
// ➣ https://httpd.apache.org/docs/2.4/content-negotiation.html
package proactive
