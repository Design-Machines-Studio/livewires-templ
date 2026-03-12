package form

import (
	"strings"
	"testing"

	"github.com/Design-Machines-Studio/livewires-templ/internal/testutil"
)

func TestFieldRenders(t *testing.T) {
	html := testutil.RenderToString(t, FieldText("Name", "name", "", true))
	if !strings.Contains(html, "field") {
		t.Error("expected field class")
	}
	if !strings.Contains(html, "Name") {
		t.Error("expected label text")
	}
	if !strings.Contains(html, `name="name"`) {
		t.Error("expected input name attribute")
	}
	if !strings.Contains(html, "required") {
		t.Error("expected required attribute")
	}
}

func TestFieldWithHint(t *testing.T) {
	html := testutil.RenderToString(t, Field(FieldProps{Label: "Email", Name: "email", Type: "email", Hint: "We won't share it"}))
	if !strings.Contains(html, "hint") {
		t.Error("expected hint class")
	}
	if !strings.Contains(html, "We won&#39;t share it") && !strings.Contains(html, "We won't share it") && !strings.Contains(html, "We won&#x27;t share it") {
		t.Error("expected hint text")
	}
}
