---
layout: page
title: Releases
permalink: /releases/
---

The official releases for `negotiator`.

> We use [Semantic Versioning v2.0.0][semantic-versioning-docs] for our releases.

## [v0.2.0][v0.2.0]
- Removes the `github.com/freerware/negotiator/internal/representation/json` package.
- Removes the `github.com/freerware/negotiator/internal/representation/xml` package.
- Removes the `github.com/freerware/negotiator/internal/representation/yaml` package.
- Exposes the `representation.List` type.

## [v0.1.1][v0.1.1]
- Fixes the issue where `sourceQuality` was not being populated within the
representation metadata in the list representations. [ [#8][issue-8] ]
- Fixes the issue where `Content-Encoding`, `Content-Language`, and 
`Content-Charset` were not being set for the list representations. [ [#7][issue-7] ]

## [v0.1.0][v0.1.0]
- Fixes neighbor resource logic for transparent negotiation such that it adheres to [`RFC2296 Section 3.5`][rfc2296-3.5-docs] and [`RFC2068 Section 3.2.3`][rfc2068-3.2.3-docs].

## [v0.0.1][v0.0.1]
- Initial draft.

[semantic-versioning-docs]: https://semver.org/
[v0.2.0]: https://github.com/freerware/negotiator/releases/tag/v0.2.0
[v0.1.1]: https://github.com/freerware/negotiator/releases/tag/v0.1.1
[v0.1.0]: https://github.com/freerware/negotiator/releases/tag/v0.1.0
[v0.0.1]: https://github.com/freerware/negotiator/releases/tag/v0.0.1
[issue-8]: https://github.com/freerware/negotiator/issues/8
[issue-7]: https://github.com/freerware/negotiator/issues/7
[rfc2296-3.5-docs]: https://tools.ietf.org/html/rfc2296#section-3.5
[rfc2068-3.2.3-docs]: https://tools.ietf.org/html/rfc2068#section-3.2.3
