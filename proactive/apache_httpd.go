/* Copyright 2020 Freerware
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package proactive

import (
	"mime"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/freerware/negotiator/internal/header"
	"github.com/freerware/negotiator/representation"
)

// httpd represents the proactive (server-driven) content
// negotiation algorithm offered by Apache HTTP server.
// https://httpd.apache.org/docs/2.4/content-negotiation.html
type httpd struct {
	filters []filter
}

// ApacheHTTPD provides the Apache HTTP server proactive content
// negotation algorithm.
func ApacheHTTPD() representation.Chooser {
	filters := []filter{
		// step 2.1
		bestSourceAndType,
		// step 2.2
		bestLanguage,
		// step 2.3
		bestLanguageOrder,
		// step 2.4
		bestLevel,
		// step 2.5
		bestCharset,
		// step 2.6
		notISO88591,
		// step 2.7
		bestEncoding,
		// step 2.8
		smallestContentLength,
	}
	return httpd{filters: filters}
}

// Chooser determines the 'best' representation from the provided set.
func (c httpd) Choose(
	r *http.Request, reps ...representation.Representation) (representation.Representation, error) {

	var (
		a   header.Accept
		ae  header.AcceptEncoding
		al  header.AcceptLanguage
		ac  header.AcceptCharset
		err error
	)

	accept := r.Header["Accept"]
	if a, err = header.NewAccept(accept); err != nil {
		return nil, err
	}

	acceptEncoding := r.Header["Accept-Encoding"]
	if ae, err = header.NewAcceptEncoding(acceptEncoding); err != nil {
		return nil, err
	}

	acceptLanguage := r.Header["Accept-Language"]
	if al, err = header.NewAcceptLanguage(acceptLanguage); err != nil {
		return nil, err
	}

	acceptCharset := r.Header["Accept-Charset"]
	if ac, err = header.NewAcceptCharset(acceptCharset); err != nil {
		return nil, err
	}

	var variants representation.Set
	for _, rp := range reps {
		qt := c.acceptQuality(rp, a)
		qc := c.acceptCharsetQuality(rp, ac)
		ql, los := c.acceptLanguageQuality(rp, al)
		qe := c.acceptEncodingQuality(rp, ae)

		shouldEliminate :=
			qt == header.QualityValueMinimum || qc == header.QualityValueMinimum ||
				qe == header.QualityValueMinimum || ql == header.QualityValueMinimum
		if shouldEliminate {
			continue
		}

		variants = append(variants, representation.RankedRepresentation{
			Representation:        rp,
			SourceQualityValue:    rp.SourceQuality(),
			MediaTypeQualityValue: qt.Float(),
			CharsetQualityValue:   qc.Float(),
			EncodingQualityValue:  qe.Float(),
			LanguageQualityValue:  ql.Float(),
			LanguageOrderScore:    los,
		})
	}

	for _, f := range c.filters {
		// if there are no eligble variants, we are done.
		if variants.Empty() {
			return nil, nil
		}
		// apply filter.
		if variants, err = f(variants); err != nil {
			return nil, err
		}
		// if we are down to one, choose it.
		if variants.Size() == 1 {
			break
		}
	}
	return variants.First().Representation, nil
}

var (
	// smallestContentLength selects the variants with the smallest content length.
	smallestContentLength filter = func(variants representation.Set) (representation.Set, error) {
		variants.Sort(func(i, j int) bool {
			// sort smallest to largest.
			var f, s int
			if bytes, err := variants[i].Representation.Bytes(); err != nil {
				f = len(bytes)
			}
			if bytes, err := variants[j].Representation.Bytes(); err != nil {
				s = len(bytes)
			}
			return f > s
		})
		lowest := variants.First()
		lowestBytes, err := lowest.Representation.Bytes()
		if err != nil {
			return representation.EmptySet, err
		}
		lowestLength := len(lowestBytes)
		return variants.Where(func(v representation.RankedRepresentation) bool {
			var length int
			if bytes, err := v.Representation.Bytes(); err == nil {
				length = len(bytes)
			}
			return length == lowestLength
		}), nil
	}

	// bestEncoding selects the variants with the best encoding.
	bestEncoding filter = func(variants representation.Set) (representation.Set, error) {
		variants.Sort(func(i, j int) bool {
			first := header.QualityValue(variants[i].EncodingQualityValue)
			second := header.QualityValue(variants[j].EncodingQualityValue)
			return first.GreaterThan(second)
		})
		highest := variants.First()
		return variants.Where(func(v representation.RankedRepresentation) bool {
			qv := header.QualityValue(v.EncodingQualityValue)
			highestqv := header.QualityValue(highest.EncodingQualityValue)
			return qv.Equals(highestqv)
		}), nil
	}

	// notISO88591 selects the variants that don't have ISO-8859-1 encoding.
	// if all variants have ISO-8859-1, select all variants instead.
	notISO88591 filter = func(variants representation.Set) (representation.Set, error) {
		notISO88591 := variants.Where(func(v representation.RankedRepresentation) bool {
			return strings.ToLower(v.Representation.ContentCharset()) != "iso-8859-1"
		})
		// only filter for variants that are not ISO8859-1 charset if
		// not all are ISO8859-1.
		if !notISO88591.Empty() && notISO88591.Size() != variants.Size() {
			return notISO88591, nil
		}
		return variants, nil
	}

	// bestCharset selects the variants with the best charset.
	bestCharset filter = func(variants representation.Set) (representation.Set, error) {
		variants.Sort(func(i, j int) bool {
			first := header.QualityValue(variants[i].CharsetQualityValue)
			second := header.QualityValue(variants[j].CharsetQualityValue)
			return first.GreaterThan(second)
		})
		highest := variants.First()
		return variants.Where(func(v representation.RankedRepresentation) bool {
			qv := header.QualityValue(v.CharsetQualityValue)
			highestqv := header.QualityValue(highest.CharsetQualityValue)
			return qv.Equals(highestqv)
		}), nil
	}

	// bestSourceAndType selects the variants with best media type and source quality.
	bestSourceAndType filter = func(variants representation.Set) (representation.Set, error) {
		variants.Sort(func(i, j int) bool {
			fsq, fmtqv :=
				header.QualityValue(variants[i].SourceQualityValue),
				header.QualityValue(variants[i].MediaTypeQualityValue)
			first := fsq.Multiply(fmtqv)
			ssq, smtqv :=
				header.QualityValue(variants[j].SourceQualityValue),
				header.QualityValue(variants[j].MediaTypeQualityValue)
			second := ssq.Multiply(smtqv)
			return first.GreaterThan(second)
		})
		highest := variants.First()
		hsqv, hmtqv :=
			header.QualityValue(highest.SourceQualityValue),
			header.QualityValue(highest.MediaTypeQualityValue)
		highestScore := hsqv.Multiply(hmtqv)
		return variants.Where(func(v representation.RankedRepresentation) bool {
			sqv, mtqv :=
				header.QualityValue(v.SourceQualityValue),
				header.QualityValue(v.MediaTypeQualityValue)
			qv := sqv.Multiply(mtqv)
			highestqv := header.QualityValue(highestScore)
			return qv.Equals(highestqv)
		}), nil
	}

	// bestLanguage selects the variants with the best language.
	bestLanguage filter = func(variants representation.Set) (representation.Set, error) {
		variants.Sort(func(i, j int) bool {
			first := header.QualityValue(variants[i].LanguageQualityValue)
			second := header.QualityValue(variants[j].LanguageQualityValue)
			return first.GreaterThan(second)
		})
		highest := variants.First()
		return variants.Where(func(v representation.RankedRepresentation) bool {
			qv := header.QualityValue(v.LanguageQualityValue)
			highestqv := header.QualityValue(highest.LanguageQualityValue)
			return qv.Equals(highestqv)
		}), nil
	}

	// bestLanguageOrder selects the variants with best language order score.
	bestLanguageOrder filter = func(variants representation.Set) (representation.Set, error) {
		variants.Sort(func(i, j int) bool {
			return variants[i].LanguageOrderScore < variants[j].LanguageOrderScore
		})
		highest := variants.First()
		return variants.Where(func(v representation.RankedRepresentation) bool {
			return v.LanguageOrderScore == highest.LanguageOrderScore
		}), nil
	}

	// bestLevel selects the variants with the highest 'level' media parameter.
	bestLevel filter = func(variants representation.Set) (representation.Set, error) {
		var htmlWithLevel []hwl
		for _, v := range variants {
			mt, p, err := mime.ParseMediaType(v.Representation.ContentType())
			if err != nil {
				return nil, err
			}
			isHTML := mt == "text/html"
			if l, ok := p["level"]; ok && isHTML {
				lNum, err := strconv.ParseInt(l, 0, 32)
				if err != nil {
					return representation.EmptySet, err
				}
				htmlWithLevel = append(htmlWithLevel, hwl{v, int(lNum)})
			}
		}
		var htmlVariants representation.Set
		if len(htmlWithLevel) > 0 {
			sort.Slice(htmlWithLevel, func(i, j int) bool {
				return htmlWithLevel[i].l > htmlWithLevel[j].l
			})
			for _, v := range htmlWithLevel {
				htmlVariants = append(htmlVariants, v.v)
			}
			return htmlVariants, nil
		}
		return variants, nil
	}
)

// acceptQuality determines the quality score for a representations media
// type based on the Accept header.
func (c httpd) acceptQuality(
	rep representation.Representation,
	accept header.Accept,
) header.QualityValue {
	qt := header.QualityValueMinimum
	if rep.ContentType() == "" || accept.IsEmpty() {
		qt = header.QualityValueMaximum
	} else {
		for _, mr := range accept.MediaRanges() {
			compatible, err := mr.Compatible(rep.ContentType())
			if compatible && err == nil {
				qt = mr.QualityValue()
				break
			}
		}
	}
	return qt
}

// acceptCharsetQuality determines the quality score for a representations
// charset based on the Accept-Charset header.
func (c httpd) acceptCharsetQuality(
	rep representation.Representation,
	acceptCharset header.AcceptCharset,
) header.QualityValue {
	qc := header.QualityValueMinimum
	if rep.ContentCharset() == "" || acceptCharset.IsEmpty() {
		qc = header.QualityValueMaximum
	} else {
		for _, c := range acceptCharset.CharsetRanges() {
			if c.Compatible(rep.ContentCharset()) {
				qc = c.QualityValue()
				break
			}
		}
	}
	return qc
}

// acceptLanguageQuality determines the quality score and order score for a
// represenations language based on the Accept-Language header.
func (c httpd) acceptLanguageQuality(
	rep representation.Representation,
	acceptLanguage header.AcceptLanguage,
) (header.QualityValue, int) {
	ql := header.QualityValueMinimum
	var los int
	if rep.ContentLanguage() == "" || acceptLanguage.IsEmpty() {
		ql = header.QualityValueMaximum
	} else {
		for idx, lr := range acceptLanguage {
			if lr.Compatible(rep.ContentLanguage()) {
				ql = lr.QualityValue()
				los = len(acceptLanguage) - idx
				break
			}
		}
	}
	return ql, los
}

// acceptEncodingQuality determines the quality score for a representations
// encoding based on the Accept-Encoding header.
func (c httpd) acceptEncodingQuality(
	rep representation.Representation,
	acceptEncoding header.AcceptEncoding,
) header.QualityValue {
	if len(rep.ContentEncoding()) == 0 || acceptEncoding.IsEmpty() {
		return header.QualityValueMaximum
	}
	for _, c := range acceptEncoding.CodingRanges() {
		for _, e := range rep.ContentEncoding() { // TODO(FREER)
			if c.Compatible(e) {
				return c.QualityValue()
			}
		}
	}
	return header.QualityValueMinimum
}

type hwl struct {
	v representation.RankedRepresentation
	l int
}
