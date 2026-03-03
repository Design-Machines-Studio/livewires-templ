package component

import (
	"strings"
	"testing"

	"github.com/Design-Machines-Studio/livewires-templ/internal/testutil"
)

func TestDividerDefault(t *testing.T) {
	html := testutil.RenderToString(t, Divider(""))
	if !strings.Contains(html, "<hr") {
		t.Error("expected hr element")
	}
	if !strings.Contains(html, "divider") {
		t.Error("expected divider base class")
	}
}

func TestDividerVariant(t *testing.T) {
	html := testutil.RenderToString(t, Divider("hairline"))
	if !strings.Contains(html, "divider--hairline") {
		t.Error("expected divider--hairline class")
	}
}

func TestDividerAccent(t *testing.T) {
	html := testutil.RenderToString(t, Divider("accent"))
	if !strings.Contains(html, "divider--accent") {
		t.Error("expected divider--accent class")
	}
}
