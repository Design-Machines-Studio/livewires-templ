package component

import (
	"strings"
	"testing"

	"github.com/Design-Machines-Studio/livewires-templ/internal/testutil"
	"github.com/a-h/templ"
)

func TestRelativeTimeKeepsNativeDateFallback(t *testing.T) {
	html := testutil.RenderToString(t, RelativeTimeComponent(RelativeTimeProps{
		DateTime: "2026-10-10", Fallback: "October 10, 2026", Lang: "en", TimeZone: "America/New_York",
		Class: "updated-at", Attrs: templ.Attributes{"data-owner": "activity"},
	}))
	nodes := testutil.ParseFragment(t, html)
	host := testutil.FindElement(nodes, "lw-relative-time")
	if host == nil {
		t.Fatal("relative-time custom element missing")
	}
	if value, ok := testutil.AttrVal(host, "class"); !ok || value != "updated-at" {
		t.Fatalf("class belongs on custom element, got %q", value)
	}
	if value, ok := testutil.AttrVal(host, "data-owner"); !ok || value != "activity" {
		t.Fatalf("Attrs belong on custom element, got %q", value)
	}
	time := testutil.FindElement(nodes, "time")
	if value, ok := testutil.AttrVal(time, "datetime"); !ok || value != "2026-10-10" {
		t.Fatalf("datetime missing: %q", value)
	}
	fallback := testutil.FindElement(nodes, "span")
	if fallback == nil || !strings.Contains(html, `data-relative-time-fallback>October 10, 2026`) || !strings.Contains(html, `data-relative-time-output hidden`) {
		t.Fatalf("native fallback/output contract missing: %s", html)
	}
}

func TestRelativeTimeReferenceMarkup(t *testing.T) {
	assertMatchesReference(t, "relative-time.html", RelativeTimeComponent(RelativeTimeProps{
		DateTime: "2026-10-10T15:30:00-04:00", Fallback: "October 10, 2026 at 3:30 PM EDT",
		Lang: "en", TimeZone: "America/New_York",
	}))
}

func TestCopyKeepsSelectableManualFallbackAndHooks(t *testing.T) {
	html := testutil.RenderToString(t, CopyComponent(CopyProps{
		ID: "share-url", Label: "Share address", Value: "https://example.test/meeting", ButtonLabel: "Copy address",
		Class: "share-copy", Attrs: templ.Attributes{"data-scope": "meeting"},
	}))
	nodes := testutil.ParseFragment(t, html)
	host := testutil.FindElement(nodes, "lw-copy")
	if value, ok := testutil.AttrVal(host, "data-scope"); !ok || value != "meeting" {
		t.Fatalf("Attrs belong on copy host, got %q", value)
	}
	textarea := testutil.FindElement(nodes, "textarea")
	if value, ok := testutil.AttrVal(textarea, "id"); !ok || value != "share-url" {
		t.Fatalf("source ID mismatch: %q", value)
	}
	if _, ok := testutil.AttrVal(textarea, "readonly"); !ok || !strings.Contains(html, "https://example.test/meeting") {
		t.Fatalf("copy value must remain readable/selectable: %s", html)
	}
	button := testutil.FindElement(nodes, "button")
	if value, ok := testutil.AttrVal(button, "type"); !ok || value != "button" {
		t.Fatalf("copy action must not submit a form: %s", html)
	}
	if !strings.Contains(html, "Copy address") || !strings.Contains(html, `data-copy-status role="status" aria-live="polite" aria-atomic="true"`) {
		t.Fatalf("copy control or truthful status hook missing: %s", html)
	}
	if strings.Contains(html, `data-state="copied"`) || strings.Contains(html, `>Copied.<`) {
		t.Fatalf("server markup must not claim a copy before clipboard confirmation: %s", html)
	}
}

func TestCopyReferenceMarkup(t *testing.T) {
	assertMatchesReference(t, "copy-button.html", CopyComponent(CopyProps{
		ID: "share-address", Label: "Share address", Value: "https://example.test/meeting/annual-planning",
		ButtonLabel: "Copy address", Class: "stack stack-compact",
	}))
}

func TestToastManagedRegionSemanticsAndExistingWrappers(t *testing.T) {
	standalone := testutil.RenderToString(t, Toast("Saved!", "success"))
	if !strings.Contains(standalone, `role="status"`) || !strings.Contains(standalone, `type="button"`) {
		t.Fatalf("existing standalone toast semantics changed: %s", standalone)
	}
	managed := withChildren(ToastRegionComponent(ToastRegionProps{Label: "Editor feedback", Attrs: templ.Attributes{"data-view": "editor"}}),
		withChildren(ToastComponent(ToastProps{ID: "save-42:confirmed", Message: "Saved changes", Variant: "success", Duration: ptr(6000), Dismissible: true, DismissLabel: "Dismiss saved message", Managed: true}),
			staticHTML(`<button type="button">Undo</button>`),
		),
	)
	html := testutil.RenderToString(t, managed)
	nodes := testutil.ParseFragment(t, html)
	region := testutil.FindElement(nodes, "section")
	list := testutil.FindElement(nodes, "div")
	if _, ok := testutil.AttrVal(region, "data-toast-region"); !ok || testutil.FindElement(nodes, "p") == nil {
		t.Fatalf("region/list/message structure missing: %s", html)
	}
	if value, ok := testutil.AttrVal(region, "aria-label"); !ok || value != "Editor feedback" {
		t.Fatalf("region label missing: %q", value)
	}
	if _, ok := testutil.AttrVal(list, "data-toast-list"); !ok {
		t.Fatalf("stable list boundary missing: %s", html)
	}
	toast := testutil.FindElementByClass(nodes, "toast")
	if id, ok := testutil.AttrVal(toast, "data-toast-id"); !ok || id != "save-42:confirmed" {
		t.Fatalf("occurrence ID missing: %q", id)
	}
	if kind, _ := testutil.AttrVal(toast, "data-toast-kind"); kind != "success" {
		t.Fatalf("toast kind missing: %q", kind)
	}
	if duration, _ := testutil.AttrVal(toast, "data-toast-duration"); duration != "6000" {
		t.Fatalf("optional duration missing: %q", duration)
	}
	if _, ok := testutil.AttrVal(toast, "role"); ok {
		t.Fatalf("managed toast must not add a duplicate live role: %s", html)
	}
	if _, ok := testutil.AttrVal(toast, "aria-live"); ok {
		t.Fatalf("managed toast must not add a nested live region: %s", html)
	}
	if _, ok := testutil.AttrVal(region, "role"); ok {
		t.Fatalf("region wrapper must not duplicate the controller announcement: %s", html)
	}
	button := testutil.FindElementByClass(nodes, "close")
	if typ, _ := testutil.AttrVal(button, "type"); typ != "button" {
		t.Fatalf("dismiss button type is not safe: %s", html)
	}
	if label, _ := testutil.AttrVal(button, "aria-label"); label != "Dismiss saved message" {
		t.Fatalf("dismiss button label missing: %q", label)
	}
	if !strings.Contains(html, `type="button">Undo</button>`) {
		t.Fatalf("toast action content must be available for sticky actions: %s", html)
	}
}

func TestLoadingRegionKeepsAnnouncementsOutsideBusyContent(t *testing.T) {
	content := withChildren(LoadingRegionComponent(LoadingRegionProps{ID: "editor-load", Label: "Editor update", Attrs: templ.Attributes{"data-owner": "editor"}}),
		staticHTML(`<input aria-label="Still usable">`),
	)
	html := testutil.RenderToString(t, content)
	nodes := testutil.ParseFragment(t, html)
	root := testutil.FindElement(nodes, "section")
	if id, _ := testutil.AttrVal(root, "id"); id != "editor-load" {
		t.Fatalf("stable root ID missing: %q", id)
	}
	busy := testutil.FindElement(nodes, "div")
	if _, ok := testutil.AttrVal(busy, "data-loading-content"); !ok {
		t.Fatalf("operation content boundary missing: %s", html)
	}
	status := testutil.FindElement(nodes, "p")
	if status == nil || status.Parent != root || busy == status {
		t.Fatalf("status announcement must remain outside busy content: %s", html)
	}
	indicator := testutil.FindElementByClass(nodes, "loading")
	if indicator == nil {
		t.Fatalf("decorative indicator missing: %s", html)
	}
	if _, ok := testutil.AttrVal(indicator, "hidden"); !ok {
		t.Fatalf("indicator should be hidden until controlled: %s", html)
	}
	if indicator.Parent != busy {
		t.Fatalf("indicator should stay within its operation content: %s", html)
	}
	if _, ok := testutil.AttrVal(root, "data-owner"); !ok {
		t.Fatalf("developer Attrs should be on the region root: %s", html)
	}
}

func ptr(value int) *int { return &value }
