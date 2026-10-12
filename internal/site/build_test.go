package site

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func testConfig(t *testing.T, contentDir string) Config {
	return Config{ContentDir: contentDir, TemplatesDir: "../../templates", StaticDir: "../../static", OutDir: t.TempDir()}
}

func TestBuild(t *testing.T) {
	cfg := testConfig(t, "testdata/content")
	if err := Build(cfg); err != nil {
		t.Fatal(err)
	}
	for _, f := range []string{
		"index.html", "blog/index.html", "blog/hello/index.html", "case-studies/sample/index.html",
		"404.html", "sitemap.xml", "robots.txt", "css/style.css", "js/main.js", "resume.pdf",
	} {
		if _, err := os.Stat(filepath.Join(cfg.OutDir, f)); err != nil {
			t.Errorf("missing %s", f)
		}
	}
	sitemap, _ := os.ReadFile(filepath.Join(cfg.OutDir, "sitemap.xml"))
	for _, want := range []string{"https://test.example.com/case-studies/sample/", "https://test.example.com/blog/hello/"} {
		if !strings.Contains(string(sitemap), want) {
			t.Errorf("sitemap missing %s", want)
		}
	}
}

func TestBuildFailureKeepsPreviousOutput(t *testing.T) {
	cfg := testConfig(t, t.TempDir()) // empty content dir: profile.yaml missing
	marker := filepath.Join(cfg.OutDir, "keep.txt")
	writeTestFile(t, marker, "old")
	err := Build(cfg)
	if err == nil || !strings.Contains(err.Error(), "profile.yaml") {
		t.Fatalf("err = %v", err)
	}
	if _, err := os.Stat(marker); err != nil {
		t.Error("output was cleared even though the build failed")
	}
}

func TestResetDirRefusesDangerousPaths(t *testing.T) {
	for _, d := range []string{"", ".", "/"} {
		if err := resetDir(d); err == nil {
			t.Errorf("resetDir(%q) should refuse", d)
		}
	}
}
