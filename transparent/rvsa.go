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

	"github.com/freerware/negotiator/internal/header"
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

	var (
		a header.Accept
		// TODO(FREER) support encoding extension.
		//ae header.AcceptEncoding
		al  header.AcceptLanguage
		ac  header.AcceptCharset
		af  header.AcceptFeatures
		err error
	)

	accept := r.Header["Accept"]
	if a, err = header.NewAccept(accept); err != nil {
		return nil, err
	}

	// TODO(FREER) support encoding extension.
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

	var variants representation.Set
	for _, rep := range reps {

		qs := rep.SourceQuality()
		qt, twc := c.acceptQuality(rep, a)
		qc, cwc := c.acceptCharsetQuality(rep, ac)
		ql, lwc := c.acceptLanguageQuality(rep, al)
		qf, fwc := c.acceptFeatureQuality(rep, af)

		isDefinite := !twc && !cwc && !lwc && !fwc
		variants = append(variants, representation.RankedRepresentation{
			Representation:        rep,
			SourceQualityValue:    qs,
			MediaTypeQualityValue: qt.Float(),
			CharsetQualityValue:   qc.Float(),
			LanguageQualityValue:  ql.Float(),
			FeatureQualityValue:   qf.Float(),
			IsDefinite:            isDefinite,
		})
	}

	if variants.Empty() {
		return nil, nil
	}

	variants.Sort(func(i, j int) bool {
		v1 := variants[i]
		firstScore := c.overallQuality(v1)
		v2 := variants[j]
		secondScore := c.overallQuality(v2)
		return firstScore > secondScore
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
	if rep.ContentType() == "" {
		return header.QualityValueMaximum, false
	}
	if accept.IsEmpty() {
		return header.QualityValueMaximum, true
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
	if rep.ContentLanguage() == "" {
		return header.QualityValueMaximum, false
	}
	if acceptLanguage.IsEmpty() {
		return header.QualityValueMaximum, true
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
	if rep.ContentCharset() == "" {
		return header.QualityValueMaximum, false
	}
	if acceptCharset.IsEmpty() {
		return header.QualityValueMaximum, true
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
	if len(rep.ContentFeatures()) == 0 {
		return header.QualityValueMaximum, false
	}
	if acceptFeature.IsEmpty() {
		return header.QualityValueMaximum, true
	}
	featureList, err := header.NewFeatureList(rep.ContentFeatures())
	if err != nil {
		panic(err) //TODO(FREER)
	}
	degradation := featureList.QualityDegradation(
		acceptFeature.AsFeatureSets(),
	)
	return header.QualityValue(degradation), usedWildcard //TODO(FREER)
}

func (c rvsa1) overallQuality(v representation.RankedRepresentation) float32 {
	qs := v.SourceQualityValue
	qt := v.MediaTypeQualityValue
	qc := v.CharsetQualityValue
	ql := v.LanguageQualityValue
	qf := v.FeatureQualityValue
	overall := qs * qt * qc * ql * qf
	qv := header.QualityValue(overall)
	return qv.Round(5).Float()
}
