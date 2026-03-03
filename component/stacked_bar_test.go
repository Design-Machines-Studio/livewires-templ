package component

import (
	"strings"
	"testing"

	"github.com/Design-Machines-Studio/livewires-templ/internal/testutil"
)

func TestStackedBar(t *testing.T) {
	html := testutil.RenderToString(t, StackedBar(StackedBarProps{
		Segments: []StackedBarSegment{
			{Percent: 40, Variant: "primary", Label: "Primary"},
			{Percent: 35, Variant: "secondary"},
			{Percent: 25, Variant: "muted"},
		},
		Label: "Distribution",
	}))
	if !strings.Contains(html, `class="stacked-bar"`) {
		t.Error("expected stacked-bar class")
	}
	if !strings.Contains(html, "segment--primary") {
		t.Error("expected segment--primary class")
	}
	if !strings.Contains(html, "segment--secondary") {
		t.Error("expected segment--secondary class")
	}
	if !strings.Contains(html, "segment--muted") {
		t.Error("expected segment--muted class")
	}
	if !strings.Contains(html, "width: 40%") {
		t.Errorf("expected 40%% width style, got: %s", html)
	}
	if !strings.Contains(html, `aria-label="Distribution"`) {
		t.Error("expected aria-label on container")
	}
	if !strings.Contains(html, `aria-label="Primary"`) {
		t.Error("expected aria-label on primary segment")
	}
	if !strings.Contains(html, `aria-hidden="true"`) {
		t.Error("expected aria-hidden=true on segments without label")
	}
}
