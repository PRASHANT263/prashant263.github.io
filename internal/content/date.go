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
