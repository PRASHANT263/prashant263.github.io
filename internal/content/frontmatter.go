package content

import (
	"errors"
	"strings"
)

// splitFrontMatter separates a leading "---" YAML block from a Markdown body.
func splitFrontMatter(data []byte) (front []byte, body string, err error) {
	s := strings.ReplaceAll(string(data), "\r\n", "\n")
	if !strings.HasPrefix(s, "---\n") {
		return nil, "", errors.New("must start with a --- front matter block")
	}
	s = s[len("---\n"):]
	// empty front matter with body: ---\n---\n...
	if strings.HasPrefix(s, "---\n") {
		return nil, s[len("---\n"):], nil
	}
	// empty front matter at EOF: ---\n---
	if s == "---" {
		return nil, "", nil
	}
	// closing --- followed by newline (body may follow)
	end := strings.Index(s, "\n---\n")
	if end >= 0 {
		return []byte(s[:end+1]), s[end+len("\n---\n"):], nil
	}
	// closing --- at EOF with no trailing newline
	if strings.HasSuffix(s, "\n---") {
		end = len(s) - len("\n---")
		return []byte(s[:end+1]), "", nil
	}
	return nil, "", errors.New("front matter is missing its closing ---")
}
