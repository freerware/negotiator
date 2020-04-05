<p align="center"><img src="https://user-images.githubusercontent.com/5921929/73627486-ba232b00-4601-11ea-9c45-26e9b31da69d.jpg" width="360"></p>

# negotiator

> A compact library for handling HTTP content negotiation for RESTful APIs.

[![GoDoc][doc-img]][doc] [![Build Status][ci-img]][ci]
[![Coverage Status][coverage-img]][coverage] [![Release][release-img]][release]
[![License][license-img]][license]

## What is it?

`negotiator` enhances the interoperability of your HTTP API by equipping it
with capabilities to facilitate [proactive][rfc7231-3.4.1], 
[reactive][rfc7231-3.4.2], and [transparent][rfc2295] content negotiation. 
This is accomplished by providing customizable and extendable functionality 
that adheres to RFC specifications, as well as industry adopted algorithms. 
With `negotiator`, your API no longer needs to take on the burden of
implementing content negotiation, allowing you to focus on simply defining
your representations and letting us do the rest.

## Why use it?

There are many reasons why you should leave your HTTP content negotiation to us:

- content negotiation algorithms are not trivial, with some algorithms 
  detailed in lengthy RFC documentation while others lacking any
  standardization at all.
- allows you to focus purely on defining your representations.
- maximizes your APIs interoperability, lowering friction for client adoption.
- unlocks all forms of content negotiation, allowing your API leverage
  different kinds of negotiation algorithms to support all of your flows.
- customization allowing you to disable or enable particular features.
- extensibility allowing you to define your own algorithms.

## How to use it?

### Quickstart

```go
http.HandleFunc("/foo", func(rw http.ResponseWriter, r *http.Request) {

	// gather representations.
	representations := []representation.Representation { Foo { ID: 1 }  }

	// choose a negotiator.
	p := proactive.Default

	// negotiate.
	ctx := negotiator.NegotiationContext { Request: r, ResponseWriter: rw }
	if err := p.Negotiate(ctx, representations...); err != nil {
		http.Error(rw, "oops!", 500)
	}
})
```

#### Examples

If you are looking for hands-on examples, we've created a sample RESTful API,
called [`tutor`][tutor], that leverages this library.

### Proactive

#### Construction

For out of the box proactive negotiation support, use
[`proactive.Default`][proactive-default-doc], which is the default proactive 
negotiator. 

```go
// retrieves the default proactive negotiator.
p := proactive.Default
```

In situations where more customization is required, use the 
[`proactive.New`][proactive-new-doc] constructor function and specify options
as arguments.

```go
// constructs a proactive negotiator with the provided options.
p := proactive.New(
	proactive.DisableStrictMode(),
	proactive.DisableNotAcceptableRepresentation(),
)
```

#### Strict Mode

According to [RFC7231][rfc7231], when none of the representations match the
values provided for a particular proactive content negotiation header, the
origin server can honor that header and return `406 Not Acceptable`, or
disregard the header field by treating the resource as if it is not subject
to content negotiation.

The behavior of honoring the header in these scenarios is what we refer to as
**strict mode**. It is possible to configure strict mode for each individual
proactive negotiation header, or disable strict mode for all. Strict mode is
enabled for all headers by default.

### Reactive

#### Construction

For out of the box reactive negotiation support, use
[`reactive.Default`][reactive-default-doc], which is the default reactive 
negotiator. 

```go
// retrieves the default reactive negotiator.
p := reactive.Default
```

In situations where more customization is required, use the 
[`reactive.New`][reactive-new-doc] constructor function and specify options
as arguments.

```go
// constructs a reactive negotiator with the provided options.
p := reactive.New(
	reactive.Logger(l),
)
```

### Transparent

#### Construction

For out of the box transparent negotiation support, use
[`transparent.Default`][transparent-default-doc], which is the default transparent 
negotiator. 

```go
// retrieves the default transparent negotiator.
p := transparent.Default
```

In situations where more customization is required, use the 
[`transparent.New`][transparent-new-doc] constructor function and specify options
as arguments.

```go
// constructs a transparent negotiator with the provided options.
p := transparent.New(
	transparent.MaximumVariantListSize(5),
)
```

### Logging

We use [`zap`][zap] as our logging library of choice. To leverage the logs
emitted from the negotiator, utilize the [`proactive.Logger`][proactive-logger-doc],
[`reactive.Logger`][reactive-logger-doc], or [`transparent.Logger`][transparent-logger-doc] 
option with a [`*zap.Logger`][logger-doc] upon creation.

```go
l, _ := zap.NewDevelopment()

// create a proactive negotiator with logging.
pn := proactive.New(
	proactive.Logger(l),
)

// create a reactive negotiator with logging.
rn := reactive.New(
	reactive.Logger(l),
)

// create a transparent negotiator with logging.
tn := transparent.New(
	transparent.Logger(l),
)
```

### Metrics

For emitting metrics, we use [`tally`][tally]. To utilize the metrics emitted
from the negotiator, leverage the [`proactive.Scope`][proactive-scope-doc],
[`reactive.Scope`][reactive-scope-doc], or [`transparent.Scope`][transparent-scope-doc] 
option with a [`tally.Scope`][scope-doc] upon creation. 

```go
s := tally.NewTestScope("example", map[string]string{}) 

// create a reactive negotiator with metrics.
rn := reactive.New(
	reactive.Scope(s),
)

// create a proactive negotiator with metrics.
pn := proactive.New(
	proactive.Scope(s),
)

// create a transparent negotiator with metrics.
tn := transparent.New(
	transparent.Scope(s),
)
```

#### Emitted Metrics

| Name                                   | Tag                     | Type      | Description                                                  |
| -------------------------------------- | ----------------------- | --------- | ------------------------------------------------------------ |
| [_PREFIX._].negotiate                  | negotiator: proactive   | timer     | The time spent during proactive negotiation.                 |
| [_PREFIX._].negotiate                  | negotiator: reactive    | timer     | The time spent during reactive negotiation.                  |
| [_PREFIX._].negotiate                  | negotiator: transparent | timer     | The time spent during transparent negotiation.               |
| [_PREFIX._].negotiate.no_content       | negotiator: proactive   | counter   | The count of proactive negotiation resulting in HTTP 204.    |
| [_PREFIX._].negotiate.no_content       | negotiator: reactive    | counter   | The count of reactive negotiation resulting in HTTP 204.     |
| [_PREFIX._].negotiate.error            | negotiator: proactive   | counter   | The count of proactive negotiation resulting in an error.    |
| [_PREFIX._].negotiate.error            | negotiator: reactive    | counter   | The count of reactive negotiation resulting in an error.     |
| [_PREFIX._].negotiate.error            | negotiator: transparent | counter   | The count of transparent negotiation resulting in an error.  |
| [_PREFIX._].negotiate.multiple_choices | negotiator: reactive    | counter   | The count of reactive negotiation resulting in HTTP 302.     |
| [_PREFIX._].negotiate.acceptable       | negotiator: proactive   | counter   | The count of reactive negotiation resulting in HTTP 200.     |
| [_PREFIX._].negotiate.not_acceptable   | negotiator: proactive   | counter   | The count of reactive negotiation resulting in HTTP 406.     |

## Contribute

Want to lend us a hand? Check out our guidelines for
[contributing][contributing].

## License

We are rocking an [Apache 2.0 license][apache-license] for this project.

## Code of Conduct

Please check out our [code of conduct][code-of-conduct] to get up to speed
how we do things.

[zap]: https://github.com/uber-go/zap
[tally]: https://github.com/uber-go/tally
[tutor]: https://github.com/freerware/tutor
[logger-doc]: https://godoc.org/go.uber.org/zap#Logger
[scope-doc]: https://godoc.org/github.com/uber-go/tally#Scope
[contributing]: https://github.com/freerware/negotiator/blob/master/CONTRIBUTING.md
[apache-license]: https://github.com/freerware/negotiator/blob/master/LICENSE.txt
[code-of-conduct]: https://github.com/freerware/negotiator/blob/master/CODE_OF_CONDUCT.md
[gophercises]: https://gophercises.com
[doc-img]: https://godoc.org/github.com/freerware/negotiator?status.svg
[doc]: https://godoc.org/github.com/freerware/negotiator
[ci-img]: https://travis-ci.com/freerware/negotiator.svg?branch=master
[ci]: https://travis-ci.com/freerware/negotiator
[coverage-img]: https://coveralls.io/repos/github/freerware/negotiator/badge.svg?branch=master
[coverage]: https://coveralls.io/github/freerware/negotiator?branch=master
[license]: https://opensource.org/licenses/Apache-2.0
[license-img]: https://img.shields.io/badge/License-Apache%202.0-blue.svg
[release]: https://github.com/freerware/negotiator/releases
[release-img]: https://img.shields.io/github/tag/freerware/negotiator.svg?label=version
[proactive-default-doc]: https://godoc.org/github.com/freerware/negotiator/proactive#pkg-variables
[proactive-new-doc]: https://godoc.org/github.com/freerware/negotiator/proactive#New
[proactive-logger-doc]: https://godoc.org/github.com/freerware/negotiator/proactive#pkg-variables
[proactive-scope-doc]: https://godoc.org/github.com/freerware/negotiator/proactive#pkg-variables
[reactive-default-doc]: https://godoc.org/github.com/freerware/negotiator/reactive#pkg-variables
[reactive-new-doc]: https://godoc.org/github.com/freerware/negotiator/reactive#New
[reactive-logger-doc]: https://godoc.org/github.com/freerware/negotiator/reactive#pkg-variables
[reactive-scope-doc]: https://godoc.org/github.com/freerware/negotiator/reactive#pkg-variables
[transparent-default-doc]: https://godoc.org/github.com/freerware/negotiator/transparent#pkg-variables
[transparent-new-doc]: https://godoc.org/github.com/freerware/negotiator/transparent#New
[transparent-logger-doc]: https://godoc.org/github.com/freerware/negotiator/transparent#pkg-variables
[transparent-scope-doc]: https://godoc.org/github.com/freerware/negotiator/transparent#pkg-variables
[rfc7231]: https://tools.ietf.org/html/rfc7231
[rfc7231-5.3]: https://tools.ietf.org/html/rfc7231#section-5.3
[rfc2295]: https://tools.ietf.org/html/rfc2295
[rfc7231-3.4.1]: https://tools.ietf.org/html/rfc7231#section-3.4.1
[rfc7231-3.4.2]: https://tools.ietf.org/html/rfc7231#section-3.4.2
