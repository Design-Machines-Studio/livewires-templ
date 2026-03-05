package component

import (
	"strings"
	"testing"

	"github.com/Design-Machines-Studio/livewires-templ/internal/testutil"
)

func TestEmptyStateSimple(t *testing.T) {
	html := testutil.RenderToString(t, EmptyState("No items found"))
	if !strings.Contains(html, "No items found") {
		t.Error("expected message text")
	}
	if !strings.Contains(html, "text-center") {
		t.Error("expected text-center class")
	}
	if !strings.Contains(html, "text-muted") {
		t.Error("expected text-muted class")
	}
}

func TestEmptyStateWithDetail(t *testing.T) {
	html := testutil.RenderToString(t, EmptyStateComponent(EmptyStateProps{
		Message: "No proposals",
		Detail:  "Create one to get started",
	}))
	if !strings.Contains(html, "No proposals") {
		t.Error("expected message text")
	}
	if !strings.Contains(html, "Create one to get started") {
		t.Error("expected detail text")
	}
}

func TestEmptyStateWithAction(t *testing.T) {
	html := testutil.RenderToString(t, EmptyStateComponent(EmptyStateProps{
		Message:     "No proposals yet",
		ActionHref:  "/proposals/new",
		ActionLabel: "Create proposal",
	}))
	if !strings.Contains(html, "/proposals/new") {
		t.Error("expected action link href")
	}
	if !strings.Contains(html, "Create proposal") {
		t.Error("expected action label")
	}
	if !strings.Contains(html, "button") {
		t.Error("expected button class on action link")
	}
}

func TestEmptyStateNoActionWhenPartial(t *testing.T) {
	html := testutil.RenderToString(t, EmptyStateComponent(EmptyStateProps{
		Message:    "No items",
		ActionHref: "/create",
	}))
	if strings.Contains(html, "/create") {
		t.Error("should not render action when ActionLabel is missing")
	}
}
