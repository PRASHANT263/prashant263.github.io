# Part 1 — Scaffold and content layer

### Task 1: Scaffold

**Files:** Create `go.mod`, `.gitignore`; Delete `main.go`

- [ ] **Step 1: Init module and deps**

```bash
cd /Users/prashant/Desktop/temp
rm main.go
go mod init portfolio
go get gopkg.in/yaml.v3@v3.0.1 github.com/yuin/goldmark@latest
mkdir -p internal/content internal/render internal/site cmd/build cmd/serve templates/partials templates/pages static/css static/js static/fonts static/img content/case-studies content/posts assets
```

- [ ] **Step 2: `.gitignore`**

```
dist/
.DS_Store
/Prashants_a25.pdf
```

- [ ] **Step 3: Commit** — `git add -A && git commit -m "chore: scaffold Go module"`

---

### Task 2: Date type

**Files:** Create `internal/content/date.go`, `internal/content/date_test.go`

- [ ] **Step 1: Failing test**

```go
package content

import (
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestDateUnmarshal(t *testing.T) {
	var v struct {
		D Date `yaml:"d"`
	}
	if err := yaml.Unmarshal([]byte("d: 2026-09-12"), &v); err != nil {
		t.Fatal(err)
	}
	if got := v.D.ISO(); got != "2026-09-12" {
		t.Errorf("ISO() = %q", got)
	}
	if got := v.D.Human(); got != "Sep 12, 2026" {
		t.Errorf("Human() = %q", got)
	}
}

func TestDateUnmarshalRejectsBadFormat(t *testing.T) {
	var v struct {
		D Date `yaml:"d"`
	}
	err := yaml.Unmarshal([]byte("d: 12/09/2026"), &v)
	if err == nil || !strings.Contains(err.Error(), "YYYY-MM-DD") {
		t.Fatalf("err = %v, want YYYY-MM-DD hint", err)
	}
}

func TestZeroDateFormatsEmpty(t *testing.T) {
	var d Date
	if d.ISO() != "" || d.Human() != "" {
		t.Errorf("zero date should format as empty")
	}
}
```

- [ ] **Step 2: Run** `go test ./internal/content/` → FAIL (undefined: Date)

- [ ] **Step 3: Implement `date.go`**

```go
package content

import (
	"fmt"
	"time"

	"gopkg.in/yaml.v3"
)

const dateLayout = "2006-01-02"

// Date is a calendar date written as YYYY-MM-DD in content files.
type Date struct{ time.Time }

func (d *Date) UnmarshalYAML(n *yaml.Node) error {
	t, err := time.Parse(dateLayout, n.Value)
	if err != nil {
		return fmt.Errorf("line %d: date %q must be YYYY-MM-DD", n.Line, n.Value)
	}
	d.Time = t
	return nil
}

// ISO returns the date as YYYY-MM-DD, or "" for the zero date.
func (d Date) ISO() string {
	if d.IsZero() {
		return ""
	}
	return d.Format(dateLayout)
}

// Human returns the date as e.g. "Sep 12, 2026", or "" for the zero date.
func (d Date) Human() string {
	if d.IsZero() {
		return ""
	}
	return d.Format("Jan 2, 2006")
}
```

- [ ] **Step 4: Run** `go test ./internal/content/` → PASS
- [ ] **Step 5: Commit** — `git commit -am "feat(content): add Date type"` (add new files first)
