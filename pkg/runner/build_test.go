//go:build builder
// +build builder

package runner

import (
	"testing"
)

func TestAppPath(t *testing.T) {
	path := execPath()
	if path == "" {
		t.Fatalf("path cannot be empty")
	}
}
