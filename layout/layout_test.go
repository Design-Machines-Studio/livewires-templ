package layout

import (
	"bytes"
	"context"
	"testing"

	"github.com/a-h/templ"
)

// renderToString is a test helper that renders a templ component to a string.
func renderToString(t *testing.T, c templ.Component) string {
	t.Helper()
	var buf bytes.Buffer
	err := c.Render(context.Background(), &buf)
	if err != nil {
		t.Fatalf("render failed: %v", err)
	}
	return buf.String()
}

func TestBaseRenders(t *testing.T) {
	html := renderToString(t, Base(BaseProps{Title: "Test Page", CSSPath: "/css/style.css"}))
	if html == "" {
		t.Fatal("expected non-empty output")
	}
}

func TestSectionRenders(t *testing.T) {
	html := renderToString(t, Section(SectionProps{Heading: "About", Scheme: "dark"}))
	if html == "" {
		t.Fatal("expected non-empty output")
	}
}
