package site

import (
	"errors"
	"fmt"
	"html"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var linkAttr = regexp.MustCompile(`(?:href|src)="([^"]*)"`)

// CheckLinks reports every root-relative href or src in dir's HTML files that
// doesn't resolve to a file in dir. External, mailto and fragment-only links
// are skipped.
func CheckLinks(dir string) error {
	var broken []error
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || filepath.Ext(path) != ".html" {
			return err
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(dir, path)
		for _, m := range linkAttr.FindAllSubmatch(data, -1) {
			link := html.UnescapeString(string(m[1]))
			if !strings.HasPrefix(link, "/") || strings.HasPrefix(link, "//") {
				continue
			}
			if !resolves(dir, link) {
				broken = append(broken, fmt.Errorf("%s: broken link %s", rel, link))
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	return errors.Join(broken...)
}

func resolves(dir, link string) bool {
	if i := strings.IndexAny(link, "?#"); i >= 0 {
		link = link[:i]
	}
	p := filepath.Join(dir, filepath.FromSlash(link))
	info, err := os.Stat(p)
	if err != nil {
		return false
	}
	if info.IsDir() {
		_, err = os.Stat(filepath.Join(p, "index.html"))
		return err == nil
	}
	return true
}
