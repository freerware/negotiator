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

package header

import (
	"errors"
	"fmt"
	"mime"
	"strconv"
	"strings"
)

var (
	// defaultMediaRange is the default media range.
	defaultMediaRange = MediaRange{
		t:      "*",
		subT:   "*",
		params: make(map[string]string),
		qValue: QualityValueMaximum,
	}

	// ErrEmptyMediaRange is an error that indicates that the media
	// range cannot be empty.
	ErrEmptyMediaRange = errors.New("media range cannot be empty")
)

// MediaRange represents a media type matching expression.
type MediaRange struct {
	t      string
	subT   string
	params map[string]string
	qValue QualityValue
}

func NewMediaRange(mediaRange string) (MediaRange, error) {
	return parseMediaRange(mediaRange)
}

// parseMediaRange parses a media range from the provided string.
func parseMediaRange(m string) (MediaRange, error) {
	if len(m) == 0 {
		return MediaRange{}, ErrEmptyMediaRange
	}
	t, subT, params, err := parse(m)
	if err != nil {
		return MediaRange{}, err
	}

	mr := MediaRange{
		t:      t,
		subT:   subT,
		params: params,
		qValue: QualityValue(1.0),
	}
	if q, ok := params["q"]; ok {
		v, err := strconv.ParseFloat(q, 32)
		if err != nil {
			return MediaRange{}, err
		}
		qv, err := NewQualityValue(float32(v))
		if err != nil {
			return MediaRange{}, err
		}
		mr.qValue = qv
	}
	return mr, nil
}

// parse parses the components of a media range from the provided string.
func parse(m string) (string, string, map[string]string, error) {
	// parse the media type.
	mediaType, params, err := mime.ParseMediaType(m)
	if err != nil {
		return "", "", nil, err
	}

	// deconstruct
	var t, subT string
	parts := strings.Split(mediaType, "/")
	t = parts[0]
	if len(parts) > 1 {
		subT = parts[1]
	}
	return t, subT, params, nil
}

// Type retrieves the type of the media range.
func (mr MediaRange) Type() string {
	return mr.t
}

func (mr MediaRange) IsTypeWildcard() bool {
	return mr.t == "*"
}

// SubType retrieves the subtype of the media range.
func (mr MediaRange) SubType() string {
	return mr.subT
}

func (mr MediaRange) IsSubTypeWildcard() bool {
	return mr.subT == "*"
}

// Param retrieves the value for media range parameter provided.
func (mr MediaRange) Param(p string) (string, bool) {
	v, ok := mr.params[p]
	return v, ok
}

func (mr MediaRange) HasParams() (b bool) {
	for k := range mr.params {
		if k != "q" {
			b = true
			return
		}
	}
	return
}

// QualityValue retrieves the quality value of the media range.
func (mr MediaRange) QualityValue() QualityValue {
	return mr.qValue
}

// Compatible determines if the provided media type is compatible with
// the media range.
func (mr MediaRange) Compatible(mediaType string) (bool, error) {
	t, subT, params, err := parse(mediaType)
	if err != nil {
		return false, err
	}

	// TODO(FREER) what if */* is passed in?
	matchedType :=
		strings.ToLower(mr.Type()) == strings.ToLower(t) || mr.Type() == "*"
	matchedSubType :=
		strings.ToLower(mr.SubType()) == strings.ToLower(subT) || mr.SubType() == "*"

	matchedParams := true
	if mr.HasParams() {
		for k, v := range mr.params {
			if vv, ok := params[k]; !ok || v != vv {
				matchedParams = false
				break
			}
		}
	}
	return matchedType && matchedSubType && matchedParams, nil
}

// Precedence determines the specificity of the media range.
func (mr MediaRange) Precedence() int {
	any := mr.t == "*" && mr.subT == "*"
	if any {
		return 0 + len(mr.params)
	}
	if mr.subT == "*" {
		return 1 + len(mr.params)
	}
	return 2 + len(mr.params)
}

// String provides the textual representation of the media range.
func (mr MediaRange) String() string {
	var params []string
	params = append(params, fmt.Sprintf("q=%s", mr.QualityValue().String()))
	for p, v := range mr.params {
		if p == "q" {
			continue
		}
		params = append(params, fmt.Sprintf("%s=%s", p, v))
	}
	t := fmt.Sprintf("%s/%s", mr.Type(), mr.SubType())
	if len(params) > 0 {
		t = fmt.Sprintf("%s;%s", t, strings.Join(params, ";"))
	}
	return t
}
