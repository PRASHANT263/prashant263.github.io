package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// writeFile is a test helper that writes content to path, fataling on error.
func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

// get performs an HTTP GET and returns (statusCode, body).
func get(t *testing.T, srv *httptest.Server, urlPath string) (int, string) {
	t.Helper()
	res, err := http.Get(srv.URL + urlPath)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	b, _ := io.ReadAll(res.Body)
	return res.StatusCode, string(b)
}

func TestHandlerServesFilesAnd404Page(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "index.html"), "home")
	writeFile(t, filepath.Join(dir, "404.html"), "not found page")
	srv := httptest.NewServer(handler(dir))
	defer srv.Close()

	if code, body := get(t, srv, "/"); code != 200 || body != "home" {
		t.Errorf("/ = %d %q", code, body)
	}
	if code, body := get(t, srv, "/missing/"); code != 404 || body != "not found page" {
		t.Errorf("/missing/ = %d %q", code, body)
	}
}

func TestHandlerTraversalGets404(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "index.html"), "home")
	writeFile(t, filepath.Join(dir, "404.html"), "not found page")
	srv := httptest.NewServer(handler(dir))
	defer srv.Close()

	// Go's http server redirects /../x → /x; the client follows the redirect.
	// Either way, the non-existent path must produce a 404.
	if code, _ := get(t, srv, "/../x"); code != 404 {
		t.Errorf("/../x = %d, want 404", code)
	}
}

func TestHandlerDirectoryWithoutIndexGets404(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "index.html"), "home")
	writeFile(t, filepath.Join(dir, "404.html"), "not found page")
	if err := os.MkdirAll(filepath.Join(dir, "emptydir"), 0o755); err != nil {
		t.Fatal(err)
	}
	srv := httptest.NewServer(handler(dir))
	defer srv.Close()

	// Directory without index.html must not produce a directory listing.
	if code, body := get(t, srv, "/emptydir/"); code != 404 || body != "not found page" {
		t.Errorf("/emptydir/ = %d %q, want 404 %q", code, body, "not found page")
	}
}

func TestHandlerFallbackWhenNo404File(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "index.html"), "home")
	// Deliberately no 404.html.
	srv := httptest.NewServer(handler(dir))
	defer srv.Close()

	code, body := get(t, srv, "/missing/")
	if code != 404 {
		t.Errorf("status = %d, want 404", code)
	}
	if !strings.Contains(body, "no 404.html") {
		t.Errorf("fallback text not present; body = %q", body)
	}
}

func TestSnapshotChangesWhenFileChanges(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "a.txt")
	writeFile(t, f, "one")
	before := snapshot([]string{dir})
	writeFile(t, f, "three")
	if snapshot([]string{dir}) == before {
		t.Error("snapshot did not change after edit")
	}
}
