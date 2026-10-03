package link

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func newTestServer() http.Handler {
	return NewHandler(NewMemoryStore()).Routes()
}

func TestHealth(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	newTestServer().ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
}

func TestCreateAndRedirect(t *testing.T) {
	srv := newTestServer()
	target := "https://example.com/some/long/path"

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/links", strings.NewReader(`{"url":"`+target+`"}`))
	srv.ServeHTTP(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("create status = %d, want %d", rec.Code, http.StatusCreated)
	}
	var resp createResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if len(resp.Code) != 7 {
		t.Fatalf("code length = %d, want 7", len(resp.Code))
	}

	rec = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/"+resp.Code, nil)
	srv.ServeHTTP(rec, req)
	if rec.Code != http.StatusFound {
		t.Fatalf("redirect status = %d, want %d", rec.Code, http.StatusFound)
	}
	if got := rec.Header().Get("Location"); got != target {
		t.Fatalf("location = %q, want %q", got, target)
	}
}

func TestCreateRejectsInvalidInput(t *testing.T) {
	cases := []string{
		`{"url":"not-a-url"}`,
		`{"url":"ftp://example.com/file"}`,
		`{"url":""}`,
		`not json`,
	}
	for _, body := range cases {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/api/links", strings.NewReader(body))
		newTestServer().ServeHTTP(rec, req)
		if rec.Code != http.StatusBadRequest {
			t.Errorf("body %q: status = %d, want %d", body, rec.Code, http.StatusBadRequest)
		}
	}
}

func TestRedirectUnknownCode(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/doesnotexist", nil)
	newTestServer().ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusNotFound)
	}
}
