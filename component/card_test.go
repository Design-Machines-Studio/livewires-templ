package component

import (
	"strings"
	"testing"

	"github.com/a-h/templ"

	"github.com/Design-Machines-Studio/livewires-templ/internal/testutil"
)

func TestCardDefault(t *testing.T) {
	html := testutil.RenderToString(t, Card(""))
	if !strings.Contains(html, "card box") {
		t.Error("expected card box stack classes")
	}
}

func TestCardSubtleScheme(t *testing.T) {
	html := testutil.RenderToString(t, Card("subtle"))
	if !strings.Contains(html, "scheme-subtle") {
		t.Error("expected scheme-subtle class")
	}
}

func TestCardLink(t *testing.T) {
	html := testutil.RenderToString(t, CardLink("/test", ""))
	if !strings.Contains(html, "<a") {
		t.Error("expected anchor tag when href provided")
	}
	if !strings.Contains(html, "card box") {
		t.Error("expected card classes on link")
	}
}

func TestCardWithHeader(t *testing.T) {
	html := testutil.RenderToString(t, CardWithHeader(CardProps{Title: "Panel Title"}))
	if !strings.Contains(html, `class="header box"`) {
		t.Error("expected header box class")
	}
	if !strings.Contains(html, "Panel Title") {
		t.Error("expected header title text")
	}
	if !strings.Contains(html, "<h3>") {
		t.Error("expected default h3 heading")
	}
}

func TestCardWithHeaderHeadingLevel(t *testing.T) {
	html := testutil.RenderToString(t, CardWithHeader(CardProps{Title: "Section", HeadingLevel: 2}))
	if !strings.Contains(html, "<h2>") {
		t.Error("expected h2 heading for HeadingLevel 2")
	}
	if strings.Contains(html, "<h3>") {
		t.Error("should not render h3 when HeadingLevel is 2")
	}
}

func TestCardWithMedia(t *testing.T) {
	media := templ.Raw(`<img src="/photo.jpg" alt="Profile">`)
	html := testutil.RenderToString(t, CardWithMedia(CardProps{Scheme: "subtle"}, media))
	if !strings.Contains(html, "card card--media") {
		t.Error("expected card--media class")
	}
	if !strings.Contains(html, "scheme-subtle") {
		t.Error("expected scheme-subtle class")
	}
	if !strings.Contains(html, "<figure>") {
		t.Error("expected figure wrapper for media")
	}
	if !strings.Contains(html, `<img src="/photo.jpg"`) {
		t.Error("expected media content inside figure")
	}
}
