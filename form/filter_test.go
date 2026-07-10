package form

import (
	"strings"
	"testing"

	"github.com/Design-Machines-Studio/livewires-templ/internal/testutil"
)

func TestFilterRenders(t *testing.T) {
	data := FilterProps{
		Title: "Status",
		Name:  "status",
		Options: []FilterOption{
			{Value: "active", Label: "Active"},
			{Value: "archived", Label: "Archived"},
		},
	}
	html := testutil.RenderToString(t, Filter(data))
	if !strings.Contains(html, "filter") {
		t.Error("expected filter class")
	}
	if !strings.Contains(html, "Status") {
		t.Error("expected filter title")
	}
	if !strings.Contains(html, "Active") {
		t.Error("expected option label")
	}
	if !strings.Contains(html, `value="active"`) {
		t.Error("expected option value")
	}
}

func TestFilterWithError(t *testing.T) {
	html := testutil.RenderToString(t, Filter(FilterProps{
		Title: "Status",
		Name:  "status",
		Error: "Selection required",
	}))
	if !strings.Contains(html, `class="error"`) {
		t.Error("expected error class")
	}
	if !strings.Contains(html, `aria-invalid="true"`) {
		t.Error("expected aria-invalid")
	}
	if !strings.Contains(html, `id="status-error"`) {
		t.Error("expected error message id")
	}
	if !strings.Contains(html, "Selection required") {
		t.Error("expected error message text")
	}
}

func TestFilterSanitizesControlID(t *testing.T) {
	html := testutil.RenderToString(t, Filter(FilterProps{Title: "Status", Name: "order status", Error: "Bad"}))
	assertLabelPointsAtControl(t, html, "order status")
	assertDescribedByResolves(t, html, 1)
}
