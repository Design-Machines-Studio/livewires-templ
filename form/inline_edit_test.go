package form

import (
	"bytes"
	"context"
	"io"
	"strings"
	"testing"

	"github.com/a-h/templ"
	"github.com/Design-Machines-Studio/livewires-templ/internal/testutil"
)

func TestInlineEditEscapesTextAndAttributes(t *testing.T) {
	html := testutil.RenderToString(t, InlineEditDisclosure(InlineEditDisclosureProps{
		ID: "edit-one", ContextID: `record&"one`, Value: `<script>title</script>`,
		Label: `Title & owner`, Name: `title"name`, Action: `/save?title="x"`, Method: "post",
		TriggerLabel: "Edit title", SaveLabel: "Save", CancelLabel: "Cancel",
		Error: `<img src=x onerror=alert(1)>`, Status: `<script>status</script>`,
	}))
	if strings.Contains(html, "<script>title</script>") || strings.Contains(html, `<img src=x`) || strings.Contains(html, "<script>status</script>") {
		t.Fatalf("text was not escaped: %s", html)
	}
	if !strings.Contains(html, `data-context-id="record&amp;&#34;one"`) {
		t.Errorf("context ID was not escaped: %s", html)
	}
	if !strings.Contains(html, `action="/save?title=&#34;x&#34;"`) {
		t.Errorf("form action was not escaped: %s", html)
	}
	if !strings.Contains(html, "&lt;script&gt;title&lt;/script&gt;") || !strings.Contains(html, "&lt;img src=x onerror=alert(1)&gt;") || !strings.Contains(html, "&lt;script&gt;status&lt;/script&gt;") {
		t.Errorf("text content was not escaped: %s", html)
	}
}

func TestInlineEditInstancesHaveDistinctIDsAndNativeComposition(t *testing.T) {
	var html bytes.Buffer
	components := []templ.Component{
		InlineEditCompact(InlineEditCompactProps{
			ID: "first", ContextID: "record-first", Value: "First", Label: "First value",
			Name: "first", Action: "/first", Method: "post", TriggerLabel: "Edit", EditLabel: "Edit first",
			SaveLabel: "Save", CancelLabel: "Cancel",
		}),
		InlineEditDisclosure(InlineEditDisclosureProps{
			ID: "second", ContextID: "record-second", Value: "Second", Label: "Second value",
			Name: "second", Action: "/second", Method: "post", TriggerLabel: "Edit second",
			SaveLabel: "Save", CancelLabel: "Cancel",
		}),
	}
	for _, component := range components {
		if err := component.Render(context.Background(), &html); err != nil {
			t.Fatal(err)
		}
	}
	rendered := html.String()
	for _, id := range []string{"first", "first-panel", "first-input", "first-error", "second", "second-panel", "second-input", "second-error"} {
		if count := strings.Count(rendered, `id="`+id+`"`); count != 1 {
			t.Errorf("instance ID %q occurs %d times, want once: %s", id, count, rendered)
		}
	}
	if strings.Count(rendered, `data-inline-edit-status`) != 2 {
		t.Fatalf("want one status per editor: %s", rendered)
	}
	if !strings.Contains(rendered, `<details`) || !strings.Contains(rendered, `<summary`) || !strings.Contains(rendered, `<form`) {
		t.Fatalf("disclosure is not composed from native elements: %s", rendered)
	}
	if !strings.Contains(rendered, `href="#first-panel"`) || !strings.Contains(rendered, `data-inline-edit="compact"`) {
		t.Fatalf("compact native link/panel hooks are missing: %s", rendered)
	}
	if !strings.Contains(rendered, `aria-label="Edit first"`) || !strings.Contains(rendered, `>Edit</a>`) {
		t.Fatalf("compact visible and accessible labels are not distinct: %s", rendered)
	}
	if !strings.Contains(rendered, `<p data-inline-edit-value>Second</p>`) {
		t.Fatalf("disclosure value should use neutral heading-independent markup: %s", rendered)
	}
}

func TestInlineEditAttributesLandOnTheirNativeElements(t *testing.T) {
	csrf := templ.ComponentFunc(func(_ context.Context, w io.Writer) error {
		_, err := io.WriteString(w, `<input type="hidden" name="csrf" value="app-token">`)
		return err
	})
	html := testutil.RenderToString(t, InlineEditCompact(InlineEditCompactProps{
		ID: "title", ContextID: "record", Value: "A title", Label: "Title", Name: "title",
		Action: "/records/title", Method: "post", EditLabel: "Edit", SaveLabel: "Save",
		CancelLabel: "Cancel", CSRF: csrf,
		Attrs:      templ.Attributes{"data-root": "root-value"},
		FormAttrs:  templ.Attributes{"data-form": "form-value", "autocomplete": "off"},
		InputAttrs: templ.Attributes{"data-input": "input-value", "aria-label": "Title editor"},
	}))
	nodes := testutil.ParseFragment(t, html)
	root := testutil.FindElement(nodes, "section")
	form := testutil.FindElement(nodes, "form")
	input := testutil.FindElement(nodes, "input")
	if got, ok := testutil.AttrVal(root, "data-root"); !ok || got != "root-value" {
		t.Errorf("root attribute = %q, %v", got, ok)
	}
	for key, want := range map[string]string{"data-form": "form-value", "autocomplete": "off", "action": "/records/title", "method": "post"} {
		if got, ok := testutil.AttrVal(form, key); !ok || got != want {
			t.Errorf("form %s = %q, %v; want %q", key, got, ok, want)
		}
		if _, ok := testutil.AttrVal(input, key); ok {
			t.Errorf("form-only attribute %s appeared on input", key)
		}
	}
	for key, want := range map[string]string{"data-input": "input-value", "aria-label": "Title editor", "id": "title-input", "name": "title"} {
		if got, ok := testutil.AttrVal(input, key); !ok || got != want {
			t.Errorf("input %s = %q, %v; want %q", key, got, ok, want)
		}
		if _, ok := testutil.AttrVal(form, key); ok {
			t.Errorf("input-only attribute %s appeared on form", key)
		}
	}
	csrfField := strings.Index(html, `name="csrf" value="app-token"`)
	formStart := strings.Index(html, "<form")
	formEnd := -1
	if formStart >= 0 {
		formEnd = strings.Index(html[formStart:], "</form>")
	}
	if csrfField < 0 || formStart < 0 || formEnd < 0 || csrfField < formStart || csrfField >= formStart+formEnd {
		t.Errorf("application CSRF field missing from native form: %s", html)
	}
}
