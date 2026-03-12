package component

import (
	"strings"
	"testing"

	"github.com/Design-Machines-Studio/livewires-templ/internal/testutil"
)

func TestActivityCardRenders(t *testing.T) {
	html := testutil.RenderToString(t, ActivityCardComponent(ActivityCardProps{Title: "Proposal A", BadgeLabel: "Active", BadgeVariant: "active"}))
	if !strings.Contains(html, "activity-card") {
		t.Error("expected activity-card class")
	}
	if !strings.Contains(html, "Proposal A") {
		t.Error("expected title text")
	}
	if !strings.Contains(html, "Active") {
		t.Error("expected badge label")
	}
}

func TestActivityCardWithHref(t *testing.T) {
	html := testutil.RenderToString(t, ActivityCardComponent(ActivityCardProps{Title: "Link Card", Href: "/proposals/1"}))
	if !strings.Contains(html, "/proposals/1") {
		t.Error("expected href in output")
	}
	if !strings.Contains(html, "title") {
		t.Error("expected title class")
	}
}
