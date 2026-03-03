package component

import (
	"strings"
	"testing"

	"github.com/Design-Machines-Studio/livewires-templ/internal/testutil"
)

func TestToastRenders(t *testing.T) {
	html := testutil.RenderToString(t, Toast("Saved!", "success"))
	if html == "" {
		t.Fatal("expected non-empty output")
	}
}

func TestToastRoleStatus(t *testing.T) {
	html := testutil.RenderToString(t, Toast("Saved!", "success"))
	if !strings.Contains(html, `role="status"`) {
		t.Error("expected role=status for success toast")
	}
}

func TestToastRoleAlert(t *testing.T) {
	html := testutil.RenderToString(t, Toast("Failed!", "error"))
	if !strings.Contains(html, `role="alert"`) {
		t.Error("expected role=alert for error toast")
	}
}

func TestToastDismissible(t *testing.T) {
	html := testutil.RenderToString(t, Toast("Notice", "info"))
	if !strings.Contains(html, `aria-label="Dismiss"`) {
		t.Error("expected dismiss button with aria-label")
	}
}
