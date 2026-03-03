package component

import (
	"strings"
	"testing"

	"github.com/Design-Machines-Studio/livewires-templ/internal/testutil"
)

func TestButtonGroup(t *testing.T) {
	html := testutil.RenderToString(t, ButtonGroup())
	if !strings.Contains(html, "button-group") {
		t.Error("expected button-group class")
	}
}

func TestButtonGroupWithClass(t *testing.T) {
	html := testutil.RenderToString(t, ButtonGroupComponent(ButtonGroupProps{Class: "mt-1"}))
	if !strings.Contains(html, "button-group mt-1") {
		t.Error("expected custom class appended")
	}
}

func TestButtonGroupRole(t *testing.T) {
	html := testutil.RenderToString(t, ButtonGroup())
	if !strings.Contains(html, `role="group"`) {
		t.Error("expected role=group on button group")
	}
}

func TestButtonGroupLabel(t *testing.T) {
	html := testutil.RenderToString(t, ButtonGroupComponent(ButtonGroupProps{Label: "Actions"}))
	if !strings.Contains(html, `aria-label="Actions"`) {
		t.Error("expected aria-label when Label provided")
	}
}
