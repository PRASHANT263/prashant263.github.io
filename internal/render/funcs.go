package render

import (
	"html/template"
	"sort"
	"strings"

	"portfolio/internal/content"
)

var funcs = template.FuncMap{
	"take": func(n int, entries []content.BlogEntry) []content.BlogEntry {
		if len(entries) > n {
			return entries[:n]
		}
		return entries
	},
	"join":      strings.Join,
	"hasPrefix": strings.HasPrefix,
	// asset turns a static/-relative path from content into a root URL.
	"asset": func(p string) string { return "/" + strings.TrimPrefix(p, "/") },
}

// blogTags returns every tag used by entries, sorted and de-duplicated.
func blogTags(entries []content.BlogEntry) []string {
	seen := map[string]bool{}
	var tags []string
	for _, e := range entries {
		for _, t := range e.Tags {
			if !seen[t] {
				seen[t] = true
				tags = append(tags, t)
			}
		}
	}
	sort.Strings(tags)
	return tags
}
