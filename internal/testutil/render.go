// Package testutil provides shared test helpers for rendering templ components.
package testutil

import (
	"bytes"
	"context"
	"testing"

	"github.com/a-h/templ"
)

// RenderToString renders a templ component to a string, failing the test on error.
func RenderToString(t *testing.T, c templ.Component) string {
	t.Helper()
	var buf bytes.Buffer
	err := c.Render(context.Background(), &buf)
	if err != nil {
		t.Fatalf("render failed: %v", err)
	}
	return buf.String()
}
