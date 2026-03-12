package component

import (
	"strings"
	"testing"

	"github.com/Design-Machines-Studio/livewires-templ/internal/testutil"
)

func TestKanbanBoardRenders(t *testing.T) {
	html := testutil.RenderToString(t, KanbanBoard("board-1"))
	if !strings.Contains(html, "reel") {
		t.Error("expected reel class")
	}
	if !strings.Contains(html, `id="board-1"`) {
		t.Error("expected board ID attribute")
	}
}

func TestKanbanColumnRenders(t *testing.T) {
	html := testutil.RenderToString(t, KanbanColumn("To Do"))
	if !strings.Contains(html, "kanban-column") {
		t.Error("expected kanban-column class")
	}
	if !strings.Contains(html, "column-title") {
		t.Error("expected column-title class")
	}
	if !strings.Contains(html, "To Do") {
		t.Error("expected column title text")
	}
}

func TestKanbanColumnHeadingLevel(t *testing.T) {
	html := testutil.RenderToString(t, KanbanColumnComponent(KanbanColumnProps{Title: "Done", HeadingLevel: "h3"}))
	if !strings.Contains(html, "<h3") {
		t.Error("expected h3 heading")
	}
}

func TestKanbanColumnTitleSizeDefault(t *testing.T) {
	html := testutil.RenderToString(t, KanbanColumn("To Do"))
	if !strings.Contains(html, "text-sm") {
		t.Error("expected default text-sm on column title")
	}
}

func TestKanbanColumnTitleSizeOverride(t *testing.T) {
	html := testutil.RenderToString(t, KanbanColumnComponent(KanbanColumnProps{Title: "Done", TitleSize: "base"}))
	if !strings.Contains(html, "text-base") {
		t.Error("expected text-base on column title")
	}
}

func TestKanbanColumnNoThemeUtilities(t *testing.T) {
	html := testutil.RenderToString(t, KanbanColumn("To Do"))
	if strings.Contains(html, `text-muted`) {
		t.Error("column title should not have text-muted utility class (handled by CSS)")
	}
}

func TestKanbanCardRenders(t *testing.T) {
	html := testutil.RenderToString(t, KanbanCard("/items/1"))
	if !strings.Contains(html, "/items/1") {
		t.Error("expected card href")
	}
	if !strings.Contains(html, "kanban-card") {
		t.Error("expected kanban-card class")
	}
}
