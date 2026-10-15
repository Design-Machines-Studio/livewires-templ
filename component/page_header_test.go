package component

import (
	"context"
	"io"
	"strings"
	"testing"

	"github.com/a-h/templ"
	"github.com/Design-Machines-Studio/livewires-templ/internal/testutil"
)

func TestPageHeaderTitleOnly(t *testing.T) {
	html := testutil.RenderToString(t, PageHeader("Organization settings"))
	if !strings.Contains(html, `<header class=""><h1>Organization settings</h1></header>`) {
		t.Fatalf("unexpected title-only page header: %s", html)
	}
}

func TestPageHeaderOptionalContentAndAttributes(t *testing.T) {
	html := testutil.RenderToString(t, PageHeaderComponent(PageHeaderProps{
		Title:        "Membership settings",
		Lead:         "A long description that remains readable below the page title.",
		Class:        "section section-snug",
		Attrs:        templ.Attributes{"data-page-header": "membership", "aria-label": "Membership settings"},
		HeadingAttrs: templ.Attributes{"id": "membership-heading", "tabindex": "-1"},
		LeadClass:    "lead text-sm",
		LayoutClass:  "cluster cluster-between items-end",
		ContentClass: "prose",
		Actions:      templ.ComponentFunc(func(ctx context.Context, w io.Writer) error { _, err := io.WriteString(w, `<a href="/members" class="button">View members</a>`); return err }),
		Supplement:   templ.ComponentFunc(func(ctx context.Context, w io.Writer) error { _, err := io.WriteString(w, `<p class="text-sm">Changes take effect after saving.</p>`); return err }),
	}))
	wants := []string{
		`<header class="section section-snug"`,
		`data-page-header="membership"`,
		`aria-label="Membership settings"`,
		`<div class="cluster cluster-between items-end"><div class="prose"><h1 id="membership-heading" tabindex="-1">Membership settings</h1>`,
		`<p class="lead text-sm">A long description that remains readable below the page title.</p>`,
		`<p class="text-sm">Changes take effect after saving.</p>`,
		`<a href="/members" class="button">View members</a>`,
	}
	for _, want := range wants {
		if !strings.Contains(html, want) {
			t.Errorf("page header missing %q in %s", want, html)
		}
	}
}

func TestPageHeaderActionsWithoutLayoutClassKeepsContentWrapper(t *testing.T) {
	html := testutil.RenderToString(t, PageHeaderComponent(PageHeaderProps{
		Title:   "Notifications",
		Actions: templ.ComponentFunc(func(ctx context.Context, w io.Writer) error { _, err := io.WriteString(w, `<button type="button">Mark all read</button>`); return err }),
	}))
	if !strings.Contains(html, `<header class=""><div><h1>Notifications</h1></div><button type="button">Mark all read</button></header>`) {
		t.Fatalf("unexpected actions header: %s", html)
	}
}
