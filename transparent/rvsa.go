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

package transparent

import (
	"net/http"

	"github.com/freerware/negotiator/internal"
	"github.com/freerware/negotiator/internal/header"
	"github.com/freerware/negotiator/internal/variant"
	"github.com/freerware/negotiator/representation"
)

// rvsa1 represents the Remote Variant Selection Algorithm 1.0 as
// defined in RFC2296. This algorithm is leveraged in remote variant
// selection within transparent content negotiation.
type rvsa1 struct{}

// RVSA1 provides the Remote Variant Selection Algorithm 1.0 as
// defined in RFC2296.
func RVSA1() representation.Chooser {
	return rvsa1{}
}

// Choose determines the 'best' representation from the provided set.
func (c rvsa1) Choose(
	r *http.Request, reps ...representation.Representation) (representation.Representation, error) {

	var a header.Accept
	//var ae header.AcceptEncoding
	var al header.AcceptLanguage
	var ac header.AcceptCharset
	var af header.AcceptFeatures
	var err error

	accept := r.Header["Accept"]
	if a, err = header.NewAccept(accept); err != nil {
		return nil, err
	}

	//acceptEncodingEncoding := r.Header["Accept-Encoding"]
	//if ae, err = header.NewAcceptEncoding(acceptEncoding); err != nil {
	//	return nil, err
	//}

	acceptLanguage := r.Header["Accept-Language"]
	if al, err = header.NewAcceptLanguage(acceptLanguage); err != nil {
		return nil, err
	}

	acceptCharset := r.Header["Accept-Charset"]
	if ac, err = header.NewAcceptCharset(acceptCharset); err != nil {
		return nil, err
	}

	acceptFeatures := r.Header["Accept-Features"]
	if af, err = header.NewAcceptFeatures(acceptFeatures); err != nil {
		return nil, err
	}

	var variants variant.Set
	for _, rep := range reps {

		qs := rep.SourceQuality()
		qt, twc := c.acceptQuality(rep, a)
		qc, cwc := c.acceptCharsetQuality(rep, ac)
		ql, lwc := c.acceptLanguageQuality(rep, al)
		qf, fwc := c.acceptFeatureQuality(rep, af)

		isDefinite := !twc && !cwc && !lwc && !fwc
		variants = append(variants, variant.Variant{
			Representation:        rep,
			SourceQualityValue:    header.QualityValue(qs),
			MediaTypeQualityValue: qt,
			CharsetQualityValue:   qc,
			LanguageQualityValue:  ql,
			FeatureQualityValue:   qf,
			IsDefinite:            isDefinite,
		})
	}

	variants.Sort(func(i, j int) bool {
		v1 := variants[i]
		firstScore := c.overallQuality(v1)
		v2 := variants[j]
		secondScore := c.overallQuality(v2)
		return firstScore < secondScore
	})

	highest := variants.First()
	score := c.overallQuality(highest)
	//https://tools.ietf.org/html/rfc2296#section-3.5 accomplishes #1 and #2
	if score > 0.0 && highest.IsDefinite {
		return highest.Representation, nil
	}
	return nil, nil
}

// acceptQuality determines the quality score for a
// represenations media type based on the Accept header.
func (c rvsa1) acceptQuality(
	rep representation.Representation,
	accept header.Accept,
) (header.QualityValue, bool) {
	var usedWildcard bool
	if rep.ContentType() == "" || accept.IsEmpty() {
		return header.QualityValueMaximum, usedWildcard
	}
	qt := header.QualityValueMinimum
	for _, mr := range accept.MediaRanges() {
		compatible, err := mr.Compatible(rep.ContentType())
		if compatible && err == nil {
			qt = mr.QualityValue()
			if mr.Type() == "*" || mr.SubType() == "*" {
				usedWildcard = true
			}
			break
		}
	}
	return qt, usedWildcard
}

// acceptLanguageQuality determines the quality score for a
// represenations language based on the Accept-Language header.
func (c rvsa1) acceptLanguageQuality(
	rep representation.Representation,
	acceptLanguage header.AcceptLanguage,
) (header.QualityValue, bool) {
	var usedWildcard bool
	if rep.ContentLanguage() == "" || acceptLanguage.IsEmpty() {
		return header.QualityValueMaximum, usedWildcard
	}
	ql := header.QualityValueMinimum
	for _, lr := range acceptLanguage {
		if lr.Compatible(rep.ContentLanguage()) {
			ql = lr.QualityValue()
			usedWildcard = lr.IsWildcard()
			break
		}
	}
	return ql, usedWildcard
}

// acceptCharsetQuality determines the quality score for a
// represenations language based on the Accept-Charset header.
func (c rvsa1) acceptCharsetQuality(
	rep representation.Representation,
	acceptCharset header.AcceptCharset,
) (header.QualityValue, bool) {
	var usedWildcard bool
	if rep.ContentCharset() == "" || acceptCharset.IsEmpty() {
		return header.QualityValueMaximum, usedWildcard
	}
	qc := header.QualityValueMinimum
	for _, c := range acceptCharset.CharsetRanges() {
		if c.Compatible(rep.ContentCharset()) {
			qc = c.QualityValue()
			usedWildcard = c.IsWildcard()
			break
		}
	}
	return qc, usedWildcard
}

// acceptFeatureQuality determines the quality score for a
// represenations language based on the Accept-Feature header.
func (c rvsa1) acceptFeatureQuality(
	rep representation.Representation,
	acceptFeature header.AcceptFeatures,
) (header.QualityValue, bool) {
	var usedWildcard bool
	if len(rep.ContentFeatures()) == 0 || acceptFeature.IsEmpty() {
		return header.QualityValueMaximum, usedWildcard
	}
	featureList := header.NewFeatureList(rep.ContentFeatures())
	degradation := featureList.QualityDegradation(
		acceptFeature.AsFeatureSets(),
	)
	return header.QualityValue(degradation), usedWildcard //TODO(FREER)
}

func (c rvsa1) overallQuality(v variant.Variant) float64 {
	qs := float32(v.SourceQualityValue)
	qt := float32(v.MediaTypeQualityValue)
	qc := float32(v.CharsetQualityValue)
	ql := float32(v.LanguageQualityValue)
	qf := float32(v.FeatureQualityValue)
	overall := qs * qt * qc * ql * qf
	return internal.Round5(float64(overall))
}
