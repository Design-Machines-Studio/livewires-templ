package component

import (
	"strings"
	"testing"

	"github.com/Design-Machines-Studio/livewires-templ/internal/testutil"
)

func TestStatusIndicatorDefault(t *testing.T) {
	html := testutil.RenderToString(t, StatusIndicator("Online", ""))
	if !strings.Contains(html, "status-indicator") {
		t.Error("expected status-indicator class")
	}
	if !strings.Contains(html, "Online") {
		t.Error("expected label text")
	}
}

func TestStatusIndicatorSuccess(t *testing.T) {
	html := testutil.RenderToString(t, StatusIndicator("Active", "success"))
	if !strings.Contains(html, "status-indicator--success") {
		t.Error("expected status-indicator--success class")
	}
}

func TestStatusIndicatorColor(t *testing.T) {
	html := testutil.RenderToString(t, StatusIndicator("Green", "green"))
	if !strings.Contains(html, "status-indicator--green") {
		t.Error("expected status-indicator--green class")
	}
}

func TestStatusIndicatorDotOnly(t *testing.T) {
	html := testutil.RenderToString(t, StatusIndicatorComponent(StatusIndicatorProps{
		Variant: "success",
		Attrs:   map[string]any{"aria-label": "Online"},
	}))
	if !strings.Contains(html, "status-indicator") {
		t.Error("expected status-indicator class")
	}
	if !strings.Contains(html, `aria-hidden="true"`) {
		t.Error("expected aria-hidden=true when label is empty")
	}
}

func TestStatusIndicatorSize(t *testing.T) {
	html := testutil.RenderToString(t, StatusIndicatorComponent(StatusIndicatorProps{
		Label: "Small", Variant: "success", Size: "small",
	}))
	if !strings.Contains(html, "status-indicator--small") {
		t.Error("expected status-indicator--small class")
	}
}
