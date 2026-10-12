package site

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTestFile(t *testing.T, path, data string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(data), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestCheckLinks(t *testing.T) {
	dir := t.TempDir()
	writeTestFile(t, filepath.Join(dir, "index.html"), `
<a href="/blog/">blog</a> <a href="/blog">blog no slash</a> <a href="/#about">anchor</a>
<link href="/css/a.css?v=1"> <img src="/missing.png">
<a href="https://example.com/x">ext</a> <a href="mailto:a@b.c">mail</a> <a href="#top">frag</a>
<a href="//cdn.example.com/x.js">protocol-relative</a>`)
	writeTestFile(t, filepath.Join(dir, "blog", "index.html"), `<a href="/">home</a>`)
	writeTestFile(t, filepath.Join(dir, "css", "a.css"), ``)
	// CSS file: one resolvable font, one missing font.
	writeTestFile(t, filepath.Join(dir, "css", "style.css"),
		`@font-face { src: url('/fonts/existing.woff2'); }
@font-face { src: url('/fonts/missing.woff2'); }`)
	writeTestFile(t, filepath.Join(dir, "fonts", "existing.woff2"), ``)

	err := CheckLinks(dir)
	if err == nil {
		t.Fatal("expected broken link error")
	}
	msg := err.Error()
	if !strings.Contains(msg, "index.html: broken link /missing.png") {
		t.Errorf("error = %q", msg)
	}
	if !strings.Contains(msg, "broken link /fonts/missing.woff2") {
		t.Errorf("CSS broken font not reported; error = %q", msg)
	}
	if strings.Count(msg, "broken link") != 2 {
		t.Errorf("expected exactly two broken links, got %q", msg)
	}
}

func TestCheckLinksClean(t *testing.T) {
	dir := t.TempDir()
	writeTestFile(t, filepath.Join(dir, "index.html"), `<a href="/">home</a>`)
	if err := CheckLinks(dir); err != nil {
		t.Fatal(err)
	}
}
