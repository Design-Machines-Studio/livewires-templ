package form

import (
	"strings"
	"testing"

	"github.com/Design-Machines-Studio/livewires-templ/internal/testutil"
)

func TestSwitchRenders(t *testing.T) {
	html := testutil.RenderToString(t, Switch("notifications", "Enable notifications", true))
	if html == "" {
		t.Fatal("expected non-empty output")
	}
	if !strings.Contains(html, `role="switch"`) {
		t.Error("expected role=switch on input")
	}
	if !strings.Contains(html, `checked`) {
		t.Error("expected checked attribute")
	}
	if !strings.Contains(html, "Enable notifications") {
		t.Error("expected label text")
	}
}

func TestSwitchUnchecked(t *testing.T) {
	html := testutil.RenderToString(t, Switch("dark-mode", "Dark mode", false))
	if strings.Contains(html, `checked`) {
		t.Error("expected no checked attribute when unchecked")
	}
}

func TestSwitchSmallVariant(t *testing.T) {
	html := testutil.RenderToString(t, SwitchSmall("compact", "Compact", false))
	if !strings.Contains(html, "switch--small") {
		t.Error("expected switch--small class")
	}
}

func TestSwitchDisabled(t *testing.T) {
	html := testutil.RenderToString(t, SwitchComponent(SwitchProps{
		Name:     "feature",
		Label:    "Feature",
		Disabled: true,
	}))
	if !strings.Contains(html, `disabled`) {
		t.Error("expected disabled attribute")
	}
}

func TestSwitchDefaultValue(t *testing.T) {
	html := testutil.RenderToString(t, Switch("toggle", "Toggle", false))
	if !strings.Contains(html, `value="1"`) {
		t.Error("expected default value of 1")
	}
}

func TestSwitchCustomValue(t *testing.T) {
	html := testutil.RenderToString(t, SwitchComponent(SwitchProps{
		Name:  "toggle",
		Label: "Toggle",
		Value: "on",
	}))
	if !strings.Contains(html, `value="on"`) {
		t.Error("expected custom value")
	}
}

func TestSwitchID(t *testing.T) {
	html := testutil.RenderToString(t, SwitchComponent(SwitchProps{
		Name:  "x",
		ID:    "my-switch",
		Label: "X",
	}))
	if !strings.Contains(html, `id="my-switch"`) {
		t.Error("expected id attribute")
	}
}

func TestSwitchNoIDWhenEmpty(t *testing.T) {
	html := testutil.RenderToString(t, SwitchComponent(SwitchProps{Name: "x", Label: "X"}))
	if strings.Contains(html, `id=`) {
		t.Error("expected no id attribute when ID is empty")
	}
}

func TestSwitchCustomClass(t *testing.T) {
	html := testutil.RenderToString(t, SwitchComponent(SwitchProps{
		Name:  "x",
		Label: "X",
		Class: "my-extra-class",
	}))
	if !strings.Contains(html, "my-extra-class") {
		t.Error("expected custom class in output")
	}
}

func TestSwitchAriaHidden(t *testing.T) {
	html := testutil.RenderToString(t, Switch("x", "X", false))
	if !strings.Contains(html, `aria-hidden="true"`) {
		t.Error("expected aria-hidden on span")
	}
}
