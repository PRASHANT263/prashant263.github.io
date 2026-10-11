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
	if strings.HasPrefix(s, "---\n") {
		return nil, s[len("---\n"):], nil
	}
	end := strings.Index(s, "\n---\n")
	if end < 0 {
		return nil, "", errors.New("front matter is missing its closing ---")
	}
	return []byte(s[:end+1]), s[end+len("\n---\n"):], nil
}
