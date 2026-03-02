package component

import (
	"testing"

	"github.com/Design-Machines-Studio/livewires-templ/internal/testutil"
)

func TestChecklistRenders(t *testing.T) {
	items := []ChecklistItemProps{
		{Label: "Task 1", Completed: true},
		{Label: "Task 2", Completed: false},
	}
	html := testutil.RenderToString(t, Checklist(items))
	if html == "" {
		t.Fatal("expected non-empty output")
	}
}
