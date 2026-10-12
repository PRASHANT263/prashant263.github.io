package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestHandlerServesFilesAnd404Page(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "index.html"), []byte("home"), 0o644)
	os.WriteFile(filepath.Join(dir, "404.html"), []byte("not found page"), 0o644)
	srv := httptest.NewServer(handler(dir))
	defer srv.Close()

	get := func(path string) (int, string) {
		res, err := http.Get(srv.URL + path)
		if err != nil {
			t.Fatal(err)
		}
		defer res.Body.Close()
		b, _ := io.ReadAll(res.Body)
		return res.StatusCode, string(b)
	}
	if code, body := get("/"); code != 200 || body != "home" {
		t.Errorf("/ = %d %q", code, body)
	}
	if code, body := get("/missing/"); code != 404 || body != "not found page" {
		t.Errorf("/missing/ = %d %q", code, body)
	}
}

func TestSnapshotChangesWhenFileChanges(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "a.txt")
	os.WriteFile(f, []byte("one"), 0o644)
	before := snapshot([]string{dir})
	os.WriteFile(f, []byte("three"), 0o644)
	if snapshot([]string{dir}) == before {
		t.Error("snapshot did not change after edit")
	}
}
