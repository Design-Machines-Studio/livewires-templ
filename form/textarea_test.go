package form

import (
	"strings"
	"testing"

	"github.com/Design-Machines-Studio/livewires-templ/internal/testutil"
)

func TestTextareaRenders(t *testing.T) {
	html := testutil.RenderToString(t, Textarea(TextareaProps{
		Label:       "Message",
		Name:        "message",
		Placeholder: "Write...",
	}))
	if !strings.Contains(html, "field") {
		t.Error("expected field class")
	}
	if !strings.Contains(html, "Message") {
		t.Error("expected label text")
	}
	if !strings.Contains(html, `name="message"`) {
		t.Error("expected textarea name attribute")
	}
	if !strings.Contains(html, "Write...") {
		t.Error("expected placeholder text")
	}
}

func TestTextareaWithHint(t *testing.T) {
	html := testutil.RenderToString(t, Textarea(TextareaProps{Label: "Bio", Name: "bio", Hint: "Keep it short"}))
	if !strings.Contains(html, "hint") {
		t.Error("expected hint class")
	}
	if !strings.Contains(html, "Keep it short") {
		t.Error("expected hint text")
	}
}

func TestTextareaWithError(t *testing.T) {
	html := testutil.RenderToString(t, Textarea(TextareaProps{
		Label: "Bio",
		Name:  "bio",
		Error: "Too long",
	}))
	if !strings.Contains(html, `class="error"`) {
		t.Error("expected error class")
	}
	if !strings.Contains(html, `aria-invalid="true"`) {
		t.Error("expected aria-invalid")
	}
	if !strings.Contains(html, `aria-describedby="bio-error"`) {
		t.Error("expected aria-describedby")
	}
	if !strings.Contains(html, `id="bio-error"`) {
		t.Error("expected error message id")
	}
	if !strings.Contains(html, "Too long") {
		t.Error("expected error message text")
	}
}

func TestTextareaAriaRequired(t *testing.T) {
	html := testutil.RenderToString(t, Textarea(TextareaProps{
		Label:    "Bio",
		Name:     "bio",
		Required: true,
	}))
	if !strings.Contains(html, `aria-required="true"`) {
		t.Error("expected aria-required when required")
	}
}

func TestTextareaAriaRequiredAbsentWhenNotRequired(t *testing.T) {
	html := testutil.RenderToString(t, Textarea(TextareaProps{
		Label: "Bio",
		Name:  "bio",
	}))
	if strings.Contains(html, `aria-required`) {
		t.Error("expected no aria-required when not required")
	}
}

func TestTextareaHintIsAssociatedWithControl(t *testing.T) {
	html := testutil.RenderToString(t, Textarea(TextareaProps{Label: "Bio", Name: "bio", Hint: "Max 200 characters"}))
	if !strings.Contains(html, `aria-describedby="bio-hint"`) {
		t.Errorf("expected textarea described by its hint, got %s", html)
	}
	if !strings.Contains(html, `<p id="bio-hint" class="hint">`) {
		t.Errorf("expected hint paragraph to carry the referenced id, got %s", html)
	}
	assertDescribedByResolves(t, html, 1)
}

func TestTextareaHintAndErrorBothAssociated(t *testing.T) {
	html := testutil.RenderToString(t, Textarea(TextareaProps{
		Label: "Bio", Name: "bio", Hint: "Max 200 characters", Error: "Too long",
	}))
	assertDescribedByResolves(t, html, 2)
}

func TestTextareaWithoutHintHasNoDescription(t *testing.T) {
	html := testutil.RenderToString(t, Textarea(TextareaProps{Label: "Bio", Name: "bio"}))
	if strings.Contains(html, "aria-describedby") {
		t.Errorf("expected no aria-describedby without hint or error, got %s", html)
	}
}

func TestTextareaSanitizesControlID(t *testing.T) {
	html := testutil.RenderToString(t, Textarea(TextareaProps{Label: "Bio", Name: "about me", Hint: "H"}))
	assertLabelPointsAtControl(t, html, "about me")
	assertDescribedByResolves(t, html, 1)
}
