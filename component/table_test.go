package component

import (
	"strings"
	"testing"

	"github.com/Design-Machines-Studio/livewires-templ/internal/testutil"
)

func TestTableDefault(t *testing.T) {
	html := testutil.RenderToString(t, Table(TableProps{}))
	if !strings.Contains(html, "<table") {
		t.Error("expected table element")
	}
	if !strings.Contains(html, `class="table"`) {
		t.Error("expected table base class")
	}
}

func TestTableBordered(t *testing.T) {
	html := testutil.RenderToString(t, Table(TableProps{Variant: "bordered"}))
	if !strings.Contains(html, "table--bordered") {
		t.Error("expected table--bordered class")
	}
}

func TestTableStriped(t *testing.T) {
	html := testutil.RenderToString(t, Table(TableProps{Variant: "striped"}))
	if !strings.Contains(html, "table--striped") {
		t.Error("expected table--striped class")
	}
}

func TestTableLined(t *testing.T) {
	html := testutil.RenderToString(t, Table(TableProps{Variant: "lined"}))
	if !strings.Contains(html, "table--lined") {
		t.Error("expected table--lined class")
	}
}
