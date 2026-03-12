package component

import (
	"strings"
	"testing"

	"github.com/Design-Machines-Studio/livewires-templ/internal/testutil"
)

func TestInfoCardUsesCardClass(t *testing.T) {
	html := testutil.RenderToString(t, InfoCardComponent(InfoCardProps{Title: "Board meeting", Header: "Upcoming"}))
	if strings.Contains(html, "info-card") {
		t.Error("should not use info-card class, should use card")
	}
	if !strings.Contains(html, "card") {
		t.Error("expected card base class")
	}
}

func TestInfoCardSizeDefault(t *testing.T) {
	html := testutil.RenderToString(t, InfoCardComponent(InfoCardProps{Title: "Test"}))
	if !strings.Contains(html, "text-sm") {
		t.Error("expected default text-sm size class")
	}
}

func TestInfoCardSizeOverride(t *testing.T) {
	html := testutil.RenderToString(t, InfoCardComponent(InfoCardProps{Title: "Test", Size: "base"}))
	if !strings.Contains(html, "text-base") {
		t.Error("expected text-base size class")
	}
	if strings.Contains(html, "text-sm") {
		t.Error("should not contain default text-sm when overridden")
	}
}

func TestInfoCardCategoryClass(t *testing.T) {
	html := testutil.RenderToString(t, InfoCardComponent(InfoCardProps{Title: "Test", Header: "Upcoming"}))
	if !strings.Contains(html, `class="category"`) {
		t.Error("expected category class on header element")
	}
}

func TestInfoCardNoUtilityClassesOnInnerElements(t *testing.T) {
	html := testutil.RenderToString(t, InfoCardComponent(InfoCardProps{Title: "Test", Header: "Cat", Subtitle: "Sub"}))
	if strings.Contains(html, `class="category text-muted"`) || strings.Contains(html, `class="category font-`) {
		t.Error("inner elements should not have utility classes for theme styling")
	}
	if strings.Contains(html, `class="card-heading font-bold"`) {
		t.Error("card-heading should not have font-bold utility class")
	}
}

func TestInfoCardWithAvatar(t *testing.T) {
	html := testutil.RenderToString(t, InfoCardComponent(InfoCardProps{
		Title:      "Board meeting",
		AvatarName: "John Doe",
		AvatarText: "Facilitator",
	}))
	if !strings.Contains(html, "avatar") {
		t.Error("expected avatar when AvatarName provided")
	}
}

func TestInfoCardWithHref(t *testing.T) {
	html := testutil.RenderToString(t, InfoCardComponent(InfoCardProps{
		Title: "Board meeting",
		Href:  "/meetings/1",
	}))
	if !strings.Contains(html, `<a href="/meetings/1"`) {
		t.Error("expected link when Href provided")
	}
}

func TestInfoCardWithSubtitle(t *testing.T) {
	html := testutil.RenderToString(t, InfoCardComponent(InfoCardProps{
		Title:    "Board meeting",
		Subtitle: "March 15, 2026",
	}))
	if !strings.Contains(html, "March 15, 2026") {
		t.Error("expected subtitle text")
	}
	if !strings.Contains(html, `subtitle`) {
		t.Error("expected subtitle class")
	}
}
