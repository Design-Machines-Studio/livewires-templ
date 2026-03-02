package form

import (
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
	if html == "" {
		t.Fatal("expected non-empty output")
	}
}
