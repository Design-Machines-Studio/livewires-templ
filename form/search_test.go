package form

import (
	"strings"
	"testing"

	"github.com/Design-Machines-Studio/livewires-templ/internal/testutil"
)

func TestSearchRenders(t *testing.T) {
	html := testutil.RenderToString(t, SearchSimple("q", "", ""))
	if html == "" {
		t.Fatal("expected non-empty output")
	}
}

func TestSearchWithError(t *testing.T) {
	html := testutil.RenderToString(t, Search(SearchProps{
		Name:  "q",
		Error: "Search failed",
	}))
	if !strings.Contains(html, `class="error"`) {
		t.Error("expected error class")
	}
	if !strings.Contains(html, `aria-invalid="true"`) {
		t.Error("expected aria-invalid")
	}
	if !strings.Contains(html, `id="q-error"`) {
		t.Error("expected error message id")
	}
	if !strings.Contains(html, "Search failed") {
		t.Error("expected error message text")
	}
}

func TestSearchSanitizesControlID(t *testing.T) {
	html := testutil.RenderToString(t, Search(SearchProps{Name: "q term", Error: "Bad"}))
	assertLabelPointsAtControl(t, html, "q term")
	assertDescribedByResolves(t, html, 1)
}
