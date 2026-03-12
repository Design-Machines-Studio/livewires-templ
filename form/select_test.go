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
