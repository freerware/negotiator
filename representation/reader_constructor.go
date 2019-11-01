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

package representation

import (
	"compress/flate"
	"compress/gzip"
	"compress/zlib"
	"io"
)

var (
	// newGzipReader is the constructor for the gzip closeable reader.
	newGzipReader EncodingReaderConstructor = func(r io.Reader) (io.ReadCloser, error) {
		return gzip.NewReader(r)
	}

	// newCompressReader is the constructor for the compress closeable reader.
	newCompressReader EncodingReaderConstructor = func(r io.Reader) (io.ReadCloser, error) {
		return zlib.NewReader(r)
	}

	// newDeflateReader is the constructor for the deflate closeable reader.
	newDeflateReader EncodingReaderConstructor = func(r io.Reader) (io.ReadCloser, error) {
		return flate.NewReader(r), nil
	}
)

// EncodingReaderConstructor represents a constructor for closeable encoding readers.
type EncodingReaderConstructor func(io.Reader) (io.ReadCloser, error)
