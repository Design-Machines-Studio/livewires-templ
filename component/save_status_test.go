package component

import (
	"strings"
	"testing"

	"github.com/Design-Machines-Studio/livewires-templ/internal/testutil"
	"github.com/a-h/templ"
)

func TestSaveStatusStateDefaults(t *testing.T) {
	tests := []struct {
		state string
		label string
	}{
		{"ready", "Ready"},
		{"dirty", "Unsaved changes"},
		{"saving", "Saving…"},
		{"saved", "Draft saved"},
		{"failed", "Couldn’t save draft"},
		{"unknown", "Save not confirmed"},
	}
	for _, test := range tests {
		t.Run(test.state, func(t *testing.T) {
			html := testutil.RenderToString(t, SaveStatus(test.state))
			nodes := testutil.ParseFragment(t, html)
			status := testutil.FindElement(nodes, "div")
			line := testutil.FindElementByClass(nodes, "save-status")
			if status == nil || line == nil {
				t.Fatalf("status region or text missing: %s", html)
			}
			if role, _ := testutil.AttrVal(status, "role"); role != "status" {
				t.Fatalf("role=status missing: %s", html)
			}
			if atomic, _ := testutil.AttrVal(status, "aria-atomic"); atomic != "true" {
				t.Fatalf("aria-atomic=true missing: %s", html)
			}
			if got, _ := testutil.AttrVal(line, "data-save-state"); got != test.state {
				t.Fatalf("state = %q, want %q: %s", got, test.state, html)
			}
			if !strings.Contains(testutil.NodeText(line), test.label) {
				t.Fatalf("label = %q, want %q: %s", html, test.label, html)
			}
		})
	}
}

func TestSaveStatusInvalidStateFailsClosed(t *testing.T) {
	for _, state := range []string{"", "confirmed", " SAVED ", "<script>"} {
		t.Run(state, func(t *testing.T) {
			html := testutil.RenderToString(t, SaveStatusComponent(SaveStatusProps{
				State: state, Label: "Draft saved",
			}))
			if !strings.Contains(html, `data-save-state="unknown"`) || !strings.Contains(html, ">Save not confirmed</p>") {
				t.Fatalf("invalid state must render the safe unknown default: %s", html)
			}
		})
	}
}

func TestSaveStatusScopeLabelAndSupportingText(t *testing.T) {
	html := testutil.RenderToString(t, SaveStatusComponent(SaveStatusProps{
		State: "saved", Label: "Text saved", SupportingText: "Rules and collaborators still need saving.",
	}))
	nodes := testutil.ParseFragment(t, html)
	status := testutil.FindElement(nodes, "div")
	line := testutil.FindElementByClass(nodes, "save-status")
	if line == nil || !strings.Contains(testutil.NodeText(line), "Text saved") {
		t.Fatalf("scope-aware label missing: %s", html)
	}
	if status == nil || !strings.Contains(testutil.NodeText(status), "Rules and collaborators still need saving.") {
		t.Fatalf("supporting text must remain inside the status region: %s", html)
	}
}

func TestSaveStatusEscapesTextAndPlacesAttributes(t *testing.T) {
	html := testutil.RenderToString(t, SaveStatusComponent(SaveStatusProps{
		State: "saved", Label: `<b>Text saved</b>`, SupportingText: `<img src=x onerror=alert(1)>`,
		ID: "editor-save-status", Class: "editor-status", Attrs: templ.Attributes{"data-editor": "proposal"},
	}))
	nodes := testutil.ParseFragment(t, html)
	status := testutil.FindElement(nodes, "div")
	line := testutil.FindElementByClass(nodes, "save-status")
	if id, _ := testutil.AttrVal(status, "id"); id != "editor-save-status" {
		t.Fatalf("stable ID belongs to the status region: %s", html)
	}
	if editor, _ := testutil.AttrVal(status, "data-editor"); editor != "proposal" {
		t.Fatalf("Attrs belong to the status region: %s", html)
	}
	if class, _ := testutil.AttrVal(line, "class"); class != "save-status editor-status" {
		t.Fatalf("Class belongs to the short status text: %s", html)
	}
	if strings.Contains(html, "<b>") || strings.Contains(html, "<img") || !strings.Contains(html, "&lt;") {
		t.Fatalf("label and supporting text must be HTML escaped: %s", html)
	}
}

func TestSaveStatusUnknownCannotBeOverriddenToSaved(t *testing.T) {
	for _, state := range []string{"unknown", "invalid"} {
		html := testutil.RenderToString(t, SaveStatusComponent(SaveStatusProps{State: state, Label: "Draft saved"}))
		if !strings.Contains(html, ">Save not confirmed</p>") {
			t.Fatalf("unknown status must not imply a save: %s", html)
		}
	}
}

func TestSaveStatusAccessibleCompositionAndUniqueInstances(t *testing.T) {
	first := SaveStatusComponent(SaveStatusProps{ID: "text-status", State: "saved", Label: "Text saved", SupportingText: "Rules still need saving."})
	second := SaveStatusComponent(SaveStatusProps{ID: "rules-status", State: "saving"})
	html := testutil.RenderToString(t, templ.Join(first, second))
	nodes := testutil.ParseFragment(t, html)
	if strings.Count(html, `role="status"`) != 2 || strings.Count(html, `aria-atomic="true"`) != 2 || strings.Contains(html, `aria-live=`) {
		t.Fatalf("each component must have exactly one status region without a duplicate live role: %s", html)
	}
	if !strings.Contains(html, `id="text-status"`) || !strings.Contains(html, `id="rules-status"`) {
		t.Fatalf("status region IDs must be unique and stable: %s", html)
	}
	textRegion := testutil.FindElement(nodes, "div")
	if textRegion == nil || !strings.Contains(html, `id="text-status"`) || !strings.Contains(testutil.NodeText(textRegion), "Text saved") || !strings.Contains(testutil.NodeText(textRegion), "Rules still need saving.") {
		t.Fatalf("scope label and support must share their atomic status region: %s", html)
	}
}
