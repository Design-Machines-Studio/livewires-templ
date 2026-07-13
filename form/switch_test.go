package form

import (
	"strings"
	"testing"

	"github.com/Design-Machines-Studio/livewires-templ/internal/testutil"
	"github.com/a-h/templ"
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

func TestSwitchWithError(t *testing.T) {
	html := testutil.RenderToString(t, SwitchComponent(SwitchProps{
		Name:  "notify",
		Label: "Notifications",
		Error: "Cannot enable",
	}))
	if !strings.Contains(html, `class="error"`) {
		t.Error("expected error class on error message")
	}
	if !strings.Contains(html, `aria-invalid="true"`) {
		t.Error("expected aria-invalid")
	}
	if !strings.Contains(html, `aria-describedby="notify-error"`) {
		t.Error("expected aria-describedby")
	}
	if !strings.Contains(html, `id="notify-error"`) {
		t.Error("expected error message id")
	}
	if !strings.Contains(html, `role="alert"`) {
		t.Error("expected role=alert on error message")
	}
	if !strings.Contains(html, "Cannot enable") {
		t.Error("expected error message text")
	}
}

func TestSwitchWithErrorUsesIDForDescribedby(t *testing.T) {
	html := testutil.RenderToString(t, SwitchComponent(SwitchProps{
		Name:  "notify",
		ID:    "my-switch",
		Label: "Notifications",
		Error: "Nope",
	}))
	if !strings.Contains(html, `aria-describedby="my-switch-error"`) {
		t.Error("expected aria-describedby to use ID when set")
	}
	if !strings.Contains(html, `id="my-switch-error"`) {
		t.Error("expected error message id to use ID when set")
	}
}

func TestSwitchWrapper(t *testing.T) {
	html := testutil.RenderToString(t, Switch("x", "X", false))
	if !strings.Contains(html, "<div>") {
		t.Error("expected div wrapper")
	}
}

func TestSwitchSanitizesErrorID(t *testing.T) {
	html := testutil.RenderToString(t, SwitchComponent(SwitchProps{Name: "email alerts", Label: "Alerts", Error: "Required"}))
	assertDescribedByResolves(t, html, 1)
}

func TestSwitchWithHint(t *testing.T) {
	html := testutil.RenderToString(t, SwitchComponent(SwitchProps{
		Name: "notify", Label: "Notifications", Hint: "Weekly digest",
	}))
	if !strings.Contains(html, `<span id="notify-hint" class="hint">Weekly digest</span>`) {
		t.Errorf("expected hint span, got %s", html)
	}
	if !strings.Contains(html, `aria-labelledby="notify-label"`) {
		t.Errorf("expected aria-labelledby, got %s", html)
	}
	assertDescribedByResolves(t, html, 1)
	// The pill span must stay hidden and ahead of the text.
	if !strings.Contains(html, `<span aria-hidden="true"></span><span>`) {
		t.Errorf("expected hint group after the pill, got %s", html)
	}
}

func TestSwitchWithoutHintRendersUnchanged(t *testing.T) {
	html := testutil.RenderToString(t, Switch("notify", "Notifications", true))
	for _, unwanted := range []string{`class="hint"`, "aria-describedby", "aria-labelledby"} {
		if strings.Contains(html, unwanted) {
			t.Errorf("expected no %s without a hint, got %s", unwanted, html)
		}
	}
	if !strings.Contains(html, `<span aria-hidden="true"></span>Notifications`) {
		t.Errorf("expected unchanged markup, got %s", html)
	}
}

// With no visible label there is no span to name the input from, so
// aria-labelledby must be omitted rather than point at nothing -- the caller
// supplies aria-label via Attrs.
func TestSwitchHintWithoutLabelOmitsLabelledby(t *testing.T) {
	html := testutil.RenderToString(t, SwitchComponent(SwitchProps{
		Name: "notify", Hint: "Weekly digest",
		Attrs: templ.Attributes{"aria-label": "Notifications"},
	}))
	if strings.Contains(html, "aria-labelledby") {
		t.Errorf("expected no aria-labelledby without a label, got %s", html)
	}
	if !strings.Contains(html, `aria-label="Notifications"`) {
		t.Error("expected caller aria-label preserved")
	}
	assertDescribedByResolves(t, html, 1)
}

func TestSwitchHintAndErrorBothAssociated(t *testing.T) {
	html := testutil.RenderToString(t, SwitchComponent(SwitchProps{
		Name: "notify", Label: "Notifications", Hint: "Weekly digest", Error: "Required",
	}))
	if !strings.Contains(html, `aria-describedby="notify-hint notify-error"`) {
		t.Errorf("expected hint then error, got %s", html)
	}
	assertDescribedByResolves(t, html, 2)
	if !strings.Contains(html, `class="switch error"`) {
		t.Errorf("expected error class on the label, got %s", html)
	}
}

// The explicit ID prop wins over Name for id derivation.
func TestSwitchHintUsesExplicitIDAndSanitizesIt(t *testing.T) {
	html := testutil.RenderToString(t, SwitchComponent(SwitchProps{
		Name: "notify", ID: "sw 1", Label: "N", Hint: "H", Error: "E",
	}))
	if !strings.Contains(html, `id="sw_20_1"`) {
		t.Errorf("expected sanitized input id, got %s", html)
	}
	assertDescribedByResolves(t, html, 2)
}

func TestSwitchHintPreservesState(t *testing.T) {
	html := testutil.RenderToString(t, SwitchComponent(SwitchProps{
		Name: "notify", Label: "N", Hint: "H", Checked: true, Disabled: true,
		Variant: "small", Class: "extra", Value: "yes",
		Attrs: templ.Attributes{"data-testid": "sw"},
	}))
	for _, want := range []string{"checked", "disabled", "switch--small", "extra", `value="yes"`, `data-testid="sw"`, `role="switch"`} {
		if !strings.Contains(html, want) {
			t.Errorf("expected %s alongside a hint, got %s", want, html)
		}
	}
}

func TestSwitchHintEscaping(t *testing.T) {
	html := testutil.RenderToString(t, SwitchComponent(SwitchProps{
		Name: "notify", Label: "N", Hint: `a & <script>alert(1)</script>`,
	}))
	if strings.Contains(html, "<script>") {
		t.Errorf("expected escaping, got %s", html)
	}
	if !strings.Contains(html, "&lt;script&gt;") || !strings.Contains(html, "&amp;") {
		t.Errorf("expected escaped entities, got %s", html)
	}
}
