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
