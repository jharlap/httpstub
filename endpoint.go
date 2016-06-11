package httpstub

import "net/http"

// An Endpoint is added to a stub server with server.Path(), and implements a handler that will be matched against a path prefix.
type Endpoint struct {
	path        string
	status      int
	contentType string
	body        []byte
}

// WithStatus sets the response status code for the endpoint
func (e *Endpoint) WithStatus(s int) *Endpoint {
	e.status = s
	return e
}

// WithContentType overrides the server's default content type for the endpoint
func (e *Endpoint) WithContentType(t string) *Endpoint {
	e.contentType = t
	return e
}

// WithBody sets the body the endpoint should return
func (e *Endpoint) WithBody(b string) *Endpoint {
	e.body = []byte(b)
	return e
}

// ServeHTTP lets Endpoint implement http.Handler
func (e *Endpoint) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if len(e.contentType) > 0 {
		w.Header().Set("Content-Type", e.contentType)
	}

	if e.status > 0 {
		w.WriteHeader(e.status)
	}

	if len(e.body) > 0 {
		w.Write(e.body)
	}
}
