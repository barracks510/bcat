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
	"fmt"
	"net"
	"net/http"
)

// Server allows content to be served on a random unused port.
type Server struct {
	url         string
	netListener *net.Listener
	mux         *http.ServeMux
}

// NewServer creates a new server instance and mounts a single http.HandleFunc
// to "/". NewServer does not use the net/http singleton ServeMux.
func NewServer(h http.HandlerFunc) (*Server, error) {
	l, err := net.Listen("tcp4", "127.0.0.1:0")
	if err != nil {
		return nil, err
	}
	m := http.NewServeMux()
	m.HandleFunc("/", h)
	return &Server{
		url:         fmt.Sprintf("http://%s", l.Addr().String()),
		netListener: &l,
		mux:         m,
	}, nil
}

// Serve accepts incoming connections and responds to them in their own
// goroutine.
func (s *Server) Serve() error {
	return http.Serve(*s.netListener, s.mux)
}

// Returns the URL that the server is listening on.
func (s *Server) Url() string {
	return s.url
}
