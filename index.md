---
# Feel free to add content and custom Front Matter to this file.
# To modify the layout, see https://jekyllrb.com/docs/themes/#overriding-theme-defaults

layout: home
---

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
[rfc2295]: https://tools.ietf.org/html/rfc2295
[rfc7231-3.4.1]: https://tools.ietf.org/html/rfc7231#section-3.4.1
[rfc7231-3.4.2]: https://tools.ietf.org/html/rfc7231#section-3.4.2
