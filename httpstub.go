package httpstub

import (
	"net/http"
	"net/http/httptest"
	"strings"
)

// Server is a configurable stub server.
type Server struct {
	URL string

	srv                *httptest.Server
	endpoints          []*Endpoint
	defaultContentType string
}

// New starts a new stub server.
func New() *Server {
	s := &Server{}
	ts := httptest.NewServer(s)
	s.srv = ts
	s.URL = ts.URL
	return s
}

// Close shuts down the server and releases the port it was listening on.
func (s *Server) Close() {
	s.srv.Close()
}

// Path creates an endpoint to respond to a request URL path. Paths can be static prefixes or may contain * to signify a path component wildcard, so /u/*/n matches /u/2/n.
func (s *Server) Path(p string) *Endpoint {
	e := &Endpoint{
		path:        p,
		contentType: s.defaultContentType,
	}
	s.endpoints = append(s.endpoints, e)
	return e
}

// WithDefaultContentType sets the default content type for the server. Evaluated when creating endpoints, so must be set first.
func (s *Server) WithDefaultContentType(t string) *Server {
	s.defaultContentType = t
	return s
}

// ServeHTTP sets *Server implement http.Handler.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// compare each path component one by one, treating * as a wildcard
	pc := strings.Split(r.URL.Path, "/")

ENDPOINT:
	for _, e := range s.endpoints {
		ec := strings.Split(e.path, "/")
		if len(ec) > len(pc) {
			continue
		}

		for i, c := range ec {
			if c == "*" {
				continue
			}
			if c != pc[i] {
				// mismatching path, so stop comparing this endpoint
				continue ENDPOINT
			}
		}

		// the pattern must match to get here
		e.ServeHTTP(w, r)
		return
	}
}
