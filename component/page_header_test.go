package component

import (
	"context"
	"io"
	"strings"
	"testing"

	"github.com/Design-Machines-Studio/livewires-templ/internal/testutil"
	"github.com/a-h/templ"
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
		Actions: templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
			_, err := io.WriteString(w, `<a href="/members" class="button">View members</a>`)
			return err
		}),
		Supplement: templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
			_, err := io.WriteString(w, `<p class="text-sm">Changes take effect after saving.</p>`)
			return err
		}),
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
		Title: "Notifications",
		Actions: templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
			_, err := io.WriteString(w, `<button type="button">Mark all read</button>`)
			return err
		}),
	}))
	if !strings.Contains(html, `<header class=""><div><h1>Notifications</h1></div><button type="button">Mark all read</button></header>`) {
		t.Fatalf("unexpected actions header: %s", html)
	}
}

func TestPageHeaderCountRemainsInsideHeading(t *testing.T) {
	html := testutil.RenderToString(t, PageHeaderComponent(PageHeaderProps{
		Title: "Members <all>",
		TitleSuffix: templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
			_, err := io.WriteString(w, `<span id="members-count" class="text-muted text-base font-normal">(27)</span>`)
			return err
		}),
	}))
	if !strings.Contains(html, `<h1>Members &lt;all&gt; <span id="members-count" class="text-muted text-base font-normal">(27)</span></h1>`) {
		t.Fatalf("count must remain inside the escaped heading: %s", html)
	}
}

func TestPageHeaderFooterFollowsActionRow(t *testing.T) {
	content := func(value string) templ.Component {
		return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error { _, err := io.WriteString(w, value); return err })
	}
	html := testutil.RenderToString(t, PageHeaderComponent(PageHeaderProps{
		Title:       "Network",
		LayoutClass: "cluster cluster-between items-baseline",
		Supplement:  content(`<p id="context">Title context</p>`),
		Actions:     content(`<button type="button">Refresh</button>`),
		Footer:      content(`<div id="lead-sync"><p class="lead">The connected network.</p><p>Last synced: today</p></div>`),
	}))
	if !strings.Contains(html, `<p id="context">Title context</p></div><button type="button">Refresh</button></div><div id="lead-sync">`) {
		t.Fatalf("supplement must stay with title, footer must follow entire action row: %s", html)
	}
	withoutActions := testutil.RenderToString(t, PageHeaderComponent(PageHeaderProps{Title: "Atlas", Footer: content(`<p>Last synced: today</p>`)}))
	if !strings.Contains(withoutActions, `<h1>Atlas</h1><p>Last synced: today</p></header>`) {
		t.Fatalf("footer must also render without actions: %s", withoutActions)
	}
}
