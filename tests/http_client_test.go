package utilities_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/teoba/any-media-convertor/Internal/client"
)

// newTestServer spins up an httptest.Server that always responds with the given
// status code and body, then returns the server and a client pointed at it.
func newTestServer(t *testing.T, status int, body string, checkReq func(r *http.Request)) (*httptest.Server, *client.Client) {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if checkReq != nil {
			checkReq(r)
		}
		w.WriteHeader(status)
		_, _ = w.Write([]byte(body))
	}))
	t.Cleanup(srv.Close)
	return srv, client.New()
}

func TestClientGet_Success(t *testing.T) {
	want := `{"ok":true}`
	srv, c := newTestServer(t, http.StatusOK, want, nil)

	got, err := c.Get(srv.URL+"/test", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(got) != want {
		t.Errorf("body = %q, want %q", got, want)
	}
}

func TestClientGet_NonOKStatusReturnsError(t *testing.T) {
	for _, status := range []int{400, 403, 404, 500, 503} {
		status := status
		t.Run(http.StatusText(status), func(t *testing.T) {
			srv, c := newTestServer(t, status, "error body", nil)
			_, err := c.Get(srv.URL, nil)
			if err == nil {
				t.Errorf("status %d: expected error, got nil", status)
			}
		})
	}
}

func TestClientGet_HeadersForwarded(t *testing.T) {
	headers := map[string]string{
		"Authorization": "Bearer token123",
		"X-Custom":      "hello",
	}

	srv, c := newTestServer(t, http.StatusOK, "ok", func(r *http.Request) {
		for k, v := range headers {
			if got := r.Header.Get(k); got != v {
				// We can't call t.Errorf here because this runs inside the server goroutine.
				// Use a panic that the test framework will capture via the deferred recover.
				panic("header " + k + " = " + got + ", want " + v)
			}
		}
	})

	_, err := c.Get(srv.URL, headers)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestClientPost_Success(t *testing.T) {
	want := `{"guest_token":"abc123"}`
	srv, c := newTestServer(t, http.StatusOK, want, func(r *http.Request) {
		if r.Method != http.MethodPost {
			panic("expected POST, got " + r.Method)
		}
	})

	got, err := c.Post(srv.URL, nil, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(got) != want {
		t.Errorf("body = %q, want %q", got, want)
	}
}

func TestClientPost_NonOKStatusReturnsError(t *testing.T) {
	srv, c := newTestServer(t, http.StatusUnauthorized, "unauthorized", nil)
	_, err := c.Post(srv.URL, nil, "")
	if err == nil {
		t.Error("expected error on 401, got nil")
	}
}

func TestClientGet_InvalidURLReturnsError(t *testing.T) {
	c := client.New()
	_, err := c.Get("://bad-url", nil)
	if err == nil {
		t.Error("expected error for invalid URL, got nil")
	}
}

func TestClientPost_InvalidURLReturnsError(t *testing.T) {
	c := client.New()
	_, err := c.Post("://bad-url", nil, "body")
	if err == nil {
		t.Error("expected error for invalid URL, got nil")
	}
}
