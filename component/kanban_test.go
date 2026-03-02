package component

import (
	"testing"

	"github.com/Design-Machines-Studio/livewires-templ/internal/testutil"
)

func TestKanbanBoardRenders(t *testing.T) {
	html := testutil.RenderToString(t, KanbanBoard("board-1"))
	if html == "" {
		t.Fatal("expected non-empty output")
	}
}

func TestKanbanColumnRenders(t *testing.T) {
	html := testutil.RenderToString(t, KanbanColumn("To Do"))
	if html == "" {
		t.Fatal("expected non-empty output")
	}
}

func TestKanbanCardRenders(t *testing.T) {
	html := testutil.RenderToString(t, KanbanCard("/items/1"))
	if html == "" {
		t.Fatal("expected non-empty output")
	}
}
