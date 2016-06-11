package httpstub_test

import (
	"encoding/json"
	"net/http"

	"github.com/jharlap/httpstub"
)

const (
	ctJSON = "application/json"
	ctXML  = "application/xml"
)

func Example() {
	ts := httpstub.New().WithDefaultContentType(ctJSON)
	defer ts.Close()

	// the default status for name requests will be 204 no content, this will match PUT, DELETE, ...
	nameEndpoint := ts.Path("/user/*/name").WithStatus(http.StatusNoContent)

	// GET overrides the status and body
	nameEndpoint.WithMethod("GET").WithBody(`{"id":"a1","name":"Alice"}`).WithStatus(http.StatusOK)

	ts.Path("/user/*/xml").WithContentType(ctXML).WithBody(`<user id="a1"><name>Alice</name></user>`)
	ts.Path("/user").WithBody(`{"id":"a1","name":"Alice","gender":"f"}`)

	resp, err := http.Get(ts.URL + "/user/a1/meep")
	if err != nil {
		panic("httpstub server misbehaved?")
	}

	if resp.Header.Get("content-type") != ctJSON {
		panic("won't happen: the server respects your content type directions")
	}

	if resp.StatusCode != http.StatusOK {
		panic("nor this: the default status code is OK")
	}

	var alice struct {
		Gender string
	}
	json.NewDecoder(resp.Body).Decode(&alice)
	if alice.Gender != "f" {
		panic("note that we requested .../meep. the first matching path prefix was /user")
	}
}
