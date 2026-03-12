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

func TestActivityCardSizeDefault(t *testing.T) {
	html := testutil.RenderToString(t, ActivityCardComponent(ActivityCardProps{Title: "Test"}))
	if !strings.Contains(html, "text-sm") {
		t.Error("expected default text-sm size class")
	}
}

func TestActivityCardNoThemeUtilities(t *testing.T) {
	html := testutil.RenderToString(t, ActivityCardComponent(ActivityCardProps{
		Title:     "Test",
		Timestamp: "2 days ago",
	}))
	if strings.Contains(html, `class="title font-bold"`) {
		t.Error("title should not have font-bold utility class")
	}
	if strings.Contains(html, `class="meta text-sm text-muted"`) {
		t.Error("meta should not have text-sm/text-muted utility classes")
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
