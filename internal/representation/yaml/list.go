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

package yaml

import (
	_representation "github.com/freerware/negotiator/internal/representation"
	"github.com/freerware/negotiator/representation"
)

// List constructs an YAML representation containing a list describing available
// representations for a particular resource.
func List(reps ...representation.Representation) representation.Representation {
	list := _representation.List{}
	list.SetContentType("application/yaml")
	list.SetRepresentations(reps...)
	return list
}