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
