package component

import (
	"testing"

	"github.com/Design-Machines-Studio/livewires-templ/internal/testutil"
)

func TestTabsRenders(t *testing.T) {
	tabs := TabsProps{
		Items: []TabProps{
			{Label: "Overview", Href: "#overview", Active: true},
			{Label: "Details", Href: "#details"},
		},
	}
	html := testutil.RenderToString(t, Tabs(tabs))
	if html == "" {
		t.Fatal("expected non-empty output")
	}
}
