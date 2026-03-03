package component

import (
	"strings"
	"testing"

	"github.com/Design-Machines-Studio/livewires-templ/internal/testutil"
)

func TestInfoCardUsesCardClass(t *testing.T) {
	html := testutil.RenderToString(t, InfoCard(InfoCardProps{Title: "Board meeting", Header: "Upcoming"}))
	if strings.Contains(html, "info-card") {
		t.Error("should not use info-card class, should use card")
	}
	if !strings.Contains(html, `class="card `) {
		t.Error("expected card base class")
	}
}

func TestInfoCardWithAvatar(t *testing.T) {
	html := testutil.RenderToString(t, InfoCard(InfoCardProps{
		Title:      "Board meeting",
		AvatarName: "John Doe",
		AvatarText: "Facilitator",
	}))
	if !strings.Contains(html, "avatar") {
		t.Error("expected avatar when AvatarName provided")
	}
}

func TestInfoCardWithHref(t *testing.T) {
	html := testutil.RenderToString(t, InfoCard(InfoCardProps{
		Title: "Board meeting",
		Href:  "/meetings/1",
	}))
	if !strings.Contains(html, `<a href="/meetings/1"`) {
		t.Error("expected link when Href provided")
	}
}

func TestInfoCardWithSubtitle(t *testing.T) {
	html := testutil.RenderToString(t, InfoCard(InfoCardProps{
		Title:    "Board meeting",
		Subtitle: "March 15, 2026",
	}))
	if !strings.Contains(html, "March 15, 2026") {
		t.Error("expected subtitle text")
	}
	if !strings.Contains(html, "text-muted") {
		t.Error("expected text-muted class on subtitle")
	}
}
