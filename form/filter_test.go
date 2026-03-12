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
