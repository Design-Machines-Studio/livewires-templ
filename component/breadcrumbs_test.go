package component

import (
	"strings"
	"testing"

	"github.com/Design-Machines-Studio/livewires-templ/internal/testutil"
)

func TestBreadcrumbs(t *testing.T) {
	html := testutil.RenderToString(t, Breadcrumbs([]BreadcrumbItem{
		{Label: "Home", Href: "/"},
		{Label: "Docs", Href: "/docs/"},
		{Label: "Current Page"},
	}))
	if !strings.Contains(html, `class="breadcrumbs"`) {
		t.Error("expected breadcrumbs class")
	}
	if !strings.Contains(html, `aria-label="Breadcrumb"`) {
		t.Error("expected aria-label")
	}
	if !strings.Contains(html, `<a href="/"`) {
		t.Error("expected link for Home")
	}
	if !strings.Contains(html, `aria-current="page"`) {
		t.Error("expected aria-current on last item")
	}
	if !strings.Contains(html, "Current Page") {
		t.Error("expected current page text")
	}
	if !strings.Contains(html, "<ol>") {
		t.Error("expected ol element for breadcrumb list")
	}
	if !strings.Contains(html, "<li>") {
		t.Error("expected li elements for breadcrumb items")
	}
}
