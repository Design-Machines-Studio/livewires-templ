package form

import (
	"strings"
	"testing"

	"github.com/Design-Machines-Studio/livewires-templ/internal/testutil"
)

func TestSelectRenders(t *testing.T) {
	opts := []SelectOption{
		{Value: "a", Label: "Option A"},
		{Value: "b", Label: "Option B"},
	}
	html := testutil.RenderToString(t, Select(SelectProps{
		Label:       "Country",
		Name:        "country",
		Placeholder: "Choose...",
		Options:     opts,
	}))
	if !strings.Contains(html, "field") {
		t.Error("expected field class")
	}
	if !strings.Contains(html, "Country") {
		t.Error("expected label text")
	}
	if !strings.Contains(html, "Choose...") {
		t.Error("expected placeholder option")
	}
	if !strings.Contains(html, "Option A") {
		t.Error("expected option label")
	}
}

func TestSelectWithError(t *testing.T) {
	html := testutil.RenderToString(t, Select(SelectProps{
		Label: "Country",
		Name:  "country",
		Error: "Required field",
	}))
	if !strings.Contains(html, `class="error"`) {
		t.Error("expected error class")
	}
	if !strings.Contains(html, `aria-invalid="true"`) {
		t.Error("expected aria-invalid on select")
	}
	if !strings.Contains(html, `aria-describedby="country-error"`) {
		t.Error("expected aria-describedby")
	}
	if !strings.Contains(html, `id="country-error"`) {
		t.Error("expected error message id")
	}
	if !strings.Contains(html, `role="alert"`) {
		t.Error("expected role=alert")
	}
	if !strings.Contains(html, "Required field") {
		t.Error("expected error message text")
	}
}
