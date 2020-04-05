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
	// newGzipWriter is the constructor for the gzip closeable writer.
	newGzipWriter EncodingWriterConstructor = func(w io.WriteCloser) (io.WriteCloser, error) {
		return gzip.NewWriter(w), nil
	}

	// newCompressWriter is the constructor for the compress closeable writer.
	newCompressWriter EncodingWriterConstructor = func(w io.WriteCloser) (io.WriteCloser, error) {
		return zlib.NewWriter(w), nil
	}

	// newDeflateWriter is the constructor for the default closeable writer.
	newDeflateWriter EncodingWriterConstructor = func(w io.WriteCloser) (io.WriteCloser, error) {
		return flate.NewWriter(w, 1)
	}
)

// EncodingWriterConstructor represents a constructor for closeable encoding writers.
type EncodingWriterConstructor func(io.WriteCloser) (io.WriteCloser, error)
