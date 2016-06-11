package httpstub_test

import (
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/jharlap/httpstub"

	"golang.org/x/net/context"
	"golang.org/x/net/context/ctxhttp"
)

func TestIsAnHTTPServer(t *testing.T) {
	ts := httpstub.New()
	defer ts.Close()

	ctx, _ := context.WithTimeout(context.Background(), time.Millisecond)
	_, err := ctxhttp.Get(ctx, nil, ts.URL)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
	}
}

func TestEndpointWithStatus(t *testing.T) {
	cases := []int{http.StatusOK, http.StatusNoContent, http.StatusNotFound, http.StatusInternalServerError}

	for _, ex := range cases {
		ts := httpstub.New()
		defer ts.Close()

		ts.Path("/").WithStatus(ex)

		ctx, _ := context.WithTimeout(context.Background(), time.Millisecond)
		resp, err := ctxhttp.Get(ctx, nil, ts.URL)
		if err != nil {
			t.Errorf("%d: unexpected err: %s", ex, err)
		}
		if resp.StatusCode != ex {
			t.Errorf("status: expected: %d got: %d", ex, resp.StatusCode)
		}
	}
}

func TestPathMatching(t *testing.T) {
	ts := httpstub.New()
	defer ts.Close()

	ts.Path("/nocontent").WithStatus(http.StatusNoContent)
	ts.Path("/err").WithStatus(http.StatusInternalServerError)
	ts.Path("/user/*/name").WithStatus(http.StatusOK)
	ts.Path("/user/*").WithStatus(http.StatusNotFound)

	cases := []struct {
		path   string
		status int
	}{
		{"/", http.StatusOK},
		{"/err", http.StatusInternalServerError},
		{"/nocontent", http.StatusNoContent},
		{"/user/1", http.StatusNotFound},
		{"/user/1/name", http.StatusOK},
		{"/user/hello/name", http.StatusOK},
	}

	for _, tc := range cases {
		ctx, _ := context.WithTimeout(context.Background(), time.Millisecond)
		resp, err := ctxhttp.Get(ctx, nil, ts.URL+tc.path)
		if err != nil {
			t.Errorf("%s: unexpected err: %s", tc.path, err)
		}
		if resp.StatusCode != tc.status {
			t.Errorf("%s: status expected: %d got: %d", tc.path, tc.status, resp.StatusCode)
		}
	}
}

func TestEndpointWithContentType(t *testing.T) {
	cases := []string{"text/plain", "application/json", "text/html", "application/json; charset=utf-8"}

	for _, ex := range cases {
		ts := httpstub.New()
		defer ts.Close()

		ts.Path("/").WithContentType(ex)

		ctx, _ := context.WithTimeout(context.Background(), time.Millisecond)
		resp, err := ctxhttp.Get(ctx, nil, ts.URL)
		if err != nil {
			t.Errorf("%s: unexpected err: %s", ex, err)
		}

		ct := resp.Header.Get("content-type")
		if ct != ex {
			t.Errorf("expected: %s got: %s", ex, ct)
		}
	}
}

func TestEndpointWithBody(t *testing.T) {
	cases := []string{"", "hello world"}

	for _, ex := range cases {
		ts := httpstub.New()
		defer ts.Close()

		ts.Path("/").WithBody(ex)

		ctx, _ := context.WithTimeout(context.Background(), time.Millisecond)
		resp, err := ctxhttp.Get(ctx, nil, ts.URL)
		if err != nil {
			t.Errorf("%s: unexpected err: %s", ex, err)
		}

		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Errorf("%s: unexpected err: %s", ex, err)
		}

		if string(b) != ex {
			t.Errorf("expected: %s got: %s", ex, string(b))
		}
	}

}

func TestEndpointWithMethods(t *testing.T) {
	cases := []struct {
		method string
		status int
		body   string
	}{
		{"GET", http.StatusOK, "hello"},
		{"PUT", http.StatusNoContent, ""},
		{"TEAPOT", http.StatusTeapot, "ceylon?"},
	}

	ts := httpstub.New()
	defer ts.Close()

	e := ts.Path("/").WithContentType("text/plain").WithStatus(http.StatusTeapot)
	e.WithMethod("GET").WithStatus(200).WithBody("hello")
	e.WithMethod("PUT").WithStatus(204)
	e.WithMethod("TEAPOT").WithBody("ceylon?") // inherits the default status

	for _, tc := range cases {

		ctx, _ := context.WithTimeout(context.Background(), time.Millisecond)
		req, err := http.NewRequest(tc.method, ts.URL, nil)
		if err != nil {
			t.Errorf("%s: unexpected err: %s", tc.method, err)
		}

		resp, err := ctxhttp.Do(ctx, nil, req)
		if err != nil {
			t.Errorf("%s: unexpected err: %s", tc.method, err)
		}

		if resp.StatusCode != tc.status {
			t.Errorf("%s: status expected: %d, got: %d", tc.method, tc.status, resp.StatusCode)
		}

		if resp.Header.Get("content-type") != "text/plain" {
			t.Errorf("%s: content-type expected: text/plain, got: %s", tc.method, resp.Header.Get("content-type"))
		}

		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Errorf("%s: unexpected err: %s", tc.method, err)
		}

		if string(b) != tc.body {
			t.Errorf("%s: body expected: %s got: %s", tc.method, tc.body, string(b))
		}
	}

}
