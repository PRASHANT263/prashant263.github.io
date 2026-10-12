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

// TestCheckOutDir verifies that Build refuses when OutDir would overlap with
// a source directory or the working directory.
func TestCheckOutDir(t *testing.T) {
	tmp := t.TempDir()
	content := filepath.Join(tmp, "content")
	templates := filepath.Join(tmp, "templates")
	static := filepath.Join(tmp, "static")
	out := filepath.Join(tmp, "out")
	for _, p := range []string{content, templates, static, out} {
		if err := os.MkdirAll(p, 0o755); err != nil {
			t.Fatal(err)
		}
	}

	good := Config{ContentDir: content, TemplatesDir: templates, StaticDir: static, OutDir: out}
	if err := checkOutDir(good); err != nil {
		t.Errorf("good config: unexpected error: %v", err)
	}

	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	parentOfCwd := filepath.Dir(cwd)

	cases := []struct {
		name   string
		outDir string
	}{
		{"out == content", content},
		{"out == static", static},
		{"out == templates", templates},
		{"out inside static", filepath.Join(static, "sub")},
		{"out == parent of cwd (..)", parentOfCwd},
		{"out == cwd", cwd},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cfg := Config{ContentDir: content, TemplatesDir: templates, StaticDir: static, OutDir: tc.outDir}
			if err := checkOutDir(cfg); err == nil {
				t.Errorf("checkOutDir with out=%q: expected error, got nil", tc.outDir)
			}
		})
	}
}

// TestBuildRenderErrorKeepsPreviousOutput confirms that a template
// execution error does not wipe the previous OutDir contents.
func TestBuildRenderErrorKeepsPreviousOutput(t *testing.T) {
	// Copy the real templates to a temp dir and overwrite pages/404.html with
	// a template that fails at execution time (non-existent field on Page).
	tmpTemplates := t.TempDir()
	if err := copyDir("../../templates", tmpTemplates); err != nil {
		t.Fatalf("copy templates: %v", err)
	}
	writeTestFile(t, filepath.Join(tmpTemplates, "pages", "404.html"),
		`{{define "main"}}{{.NoSuchField}}{{end}}`)

	cfg := Config{
		ContentDir:   "testdata/content",
		TemplatesDir: tmpTemplates,
		StaticDir:    "../../static",
		OutDir:       t.TempDir(),
	}
	marker := filepath.Join(cfg.OutDir, "keep.txt")
	writeTestFile(t, marker, "old")

	err := Build(cfg)
	if err == nil {
		t.Fatal("expected render error, got nil")
	}
	if _, err := os.Stat(marker); err != nil {
		t.Error("output was cleared even though render failed")
	}
}
