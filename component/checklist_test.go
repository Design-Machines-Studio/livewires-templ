package component

import (
	"strings"
	"testing"

	"github.com/Design-Machines-Studio/livewires-templ/internal/testutil"
)

func TestChecklistRenders(t *testing.T) {
	items := []ChecklistItemProps{
		{Label: "Task 1", Completed: true},
		{Label: "Task 2", Completed: false},
	}
	html := testutil.RenderToString(t, Checklist(items))
	if !strings.Contains(html, "checklist") {
		t.Error("expected checklist class")
	}
	if !strings.Contains(html, "Task 1") {
		t.Error("expected completed item label")
	}
	if !strings.Contains(html, "Task 2") {
		t.Error("expected pending item label")
	}
	if !strings.Contains(html, "check") {
		t.Error("expected check class on completed item")
	}
}

func TestProgressChecklistZeroTotal(t *testing.T) {
	html := testutil.RenderToString(t, ProgressChecklistComponent(ProgressChecklistProps{
		Title:     "Empty checklist",
		Completed: 0,
		Total:     0,
	}))
	if !strings.Contains(html, "progress-checklist") {
		t.Error("expected progress-checklist class")
	}
	if !strings.Contains(html, "Empty checklist") {
		t.Error("expected title text")
	}
	if !strings.Contains(html, "0 of 0 done") {
		t.Error("expected progress count")
	}
}

func TestProgressChecklistHeadingLevel(t *testing.T) {
	html := testutil.RenderToString(t, ProgressChecklistComponent(ProgressChecklistProps{
		Title:        "Setup",
		HeadingLevel: "h3",
		Total:        2,
		Completed:    1,
	}))
	if !strings.Contains(html, "<h3>") {
		t.Error("expected h3 heading")
	}
}
