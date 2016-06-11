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

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
	defer cancel()
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

		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
		defer cancel()
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
		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
		defer cancel()
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

		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
		defer cancel()
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

		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
		defer cancel()
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
