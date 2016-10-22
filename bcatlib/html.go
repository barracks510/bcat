// Copyright © 2016 Dennis Chen <barracks510@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package bcatlib

import (
	"html"
	"io"
	"os"
	"strings"
)

// FilterFunc applies some filter on a channel of byte slices.
type FilterFunc func(<-chan []byte) <-chan []byte

// ReaderCollection implements io.ReaderCloser on a collection of files
type ReaderCollection struct {
	files   []io.ReadCloser
	reader  io.Reader
	filters []FilterFunc
}

// NewReaderCollection creates a new ReaderCollection with file descriptors in
// the same nature as cat.
func NewReaderCollection(args []string) (*ReaderCollection, error) {
	if len(args) == 0 {
		return &ReaderCollection{files: []io.ReadCloser{os.Stdin}, reader: os.Stdin}, nil
	}

	files := make([]io.ReadCloser, len(args))
	for i, file := range args {
		var err error
		if file == "-" {
			files[i] = os.Stdin
		}
		files[i], err = os.Open(file)
		if err != nil {
			return nil, err
		}
	}

	readers := make([]io.Reader, len(args))
	for i := range files {
		readers[i] = files[i]
	}

	return &ReaderCollection{files: files, reader: io.MultiReader(readers...)}, nil
}

func (r *ReaderCollection) Read(p []byte) (int, error) {
	return r.reader.Read(p)
}

// Close closes up each of the files and returns the error of the one that
// failed last.
func (r *ReaderCollection) Close() error {
	var err error
	for _, fh := range r.files {
		err = fh.Close()
	}
	return err
}

// MakeFilterableChan returns a byte slice channel appropriate for use in a
// FilterFunc.
func (r *ReaderCollection) MakeFilterableChan() <-chan []byte {
	ch := make(chan []byte)
	go func() {
		for {
			buffer := make([]byte, 4096)
			n, err := r.Read(buffer)
			if err != nil || n == 0 {
				break
			}
			ch <- buffer[:n]
		}
		close(ch)
	}()
	return ch
}

// TextFilter applies a basic filter on the contents of a string chan, escaping
// HTML, and newline characters
func TextFilter(source <-chan []byte) <-chan []byte {
	ch := make(chan []byte)
	var new bool = true
	go func() {
		for content := range source {
			var output string
			if new {
				output += "<pre>"
				new = !new
			}
			output += strings.Replace(html.EscapeString(string(content)), "\n", "<br>", -1)
			ch <- []byte(output)
		}
		ch <- []byte("</pre>")
		close(ch)
	}()
	return ch
}

// TeeFilter applies a basic tee filter on unaltered contents of a string chan.
func TeeFilter(source <-chan []byte, outfile *os.File) <-chan []byte {
	ch := make(chan []byte)
	go func() {
		for output := range source {
			outfile.Write(output)
			ch <- output
		}
		close(ch)
	}()
	return ch
}
