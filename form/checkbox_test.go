package form

import (
	"strings"
	"testing"

	"github.com/Design-Machines-Studio/livewires-templ/internal/testutil"
	"github.com/a-h/templ"
)

func TestCheckboxRenders(t *testing.T) {
	html := testutil.RenderToString(t, CheckboxSimple("agree", "I agree", false))
	if !strings.Contains(html, `class="checkbox"`) {
		t.Error("expected checkbox class")
	}
	if !strings.Contains(html, `type="checkbox"`) {
		t.Error("expected checkbox input type")
	}
	if !strings.Contains(html, `name="agree"`) {
		t.Error("expected name attribute")
	}
	if !strings.Contains(html, "I agree") {
		t.Error("expected label text")
	}
}

func TestCheckboxChecked(t *testing.T) {
	html := testutil.RenderToString(t, CheckboxSimple("agree", "I agree", true))
	if !strings.Contains(html, "checked") {
		t.Error("expected checked attribute")
	}
}

func TestCheckboxVariantSuccess(t *testing.T) {
	html := testutil.RenderToString(t, Checkbox(CheckboxProps{
		Name:    "ok",
		Label:   "Approved",
		Variant: "success",
	}))
	if !strings.Contains(html, "checkbox--success") {
		t.Error("expected checkbox--success class")
	}
}

func TestCheckboxVariantError(t *testing.T) {
	html := testutil.RenderToString(t, Checkbox(CheckboxProps{
		Name:    "err",
		Label:   "Rejected",
		Variant: "error",
	}))
	if !strings.Contains(html, "checkbox--error") {
		t.Error("expected checkbox--error class")
	}
}

func TestCheckboxWithValue(t *testing.T) {
	html := testutil.RenderToString(t, CheckboxWithValue("color", "red", "Red", false))
	if !strings.Contains(html, `value="red"`) {
		t.Error("expected value attribute")
	}
}

func TestCheckboxWithError(t *testing.T) {
	html := testutil.RenderToString(t, Checkbox(CheckboxProps{
		Name:  "agree",
		Label: "I agree",
		Error: "You must agree",
	}))
	if !strings.Contains(html, "checkbox--error") {
		t.Error("expected checkbox--error variant when Error is set")
	}
	if !strings.Contains(html, `aria-invalid="true"`) {
		t.Error("expected aria-invalid")
	}
	if !strings.Contains(html, `aria-describedby="agree-error"`) {
		t.Error("expected aria-describedby")
	}
	if !strings.Contains(html, `id="agree-error"`) {
		t.Error("expected error message id")
	}
	if !strings.Contains(html, `role="alert"`) {
		t.Error("expected role=alert")
	}
	if !strings.Contains(html, "You must agree") {
		t.Error("expected error message text")
	}
}

func TestCheckboxErrorPreservesExplicitVariant(t *testing.T) {
	html := testutil.RenderToString(t, Checkbox(CheckboxProps{
		Name:    "ok",
		Label:   "OK",
		Variant: "warning",
		Error:   "Problem",
	}))
	if !strings.Contains(html, "checkbox--warning") {
		t.Error("expected explicit variant preserved when Error is set")
	}
}

func TestCheckboxWrapper(t *testing.T) {
	html := testutil.RenderToString(t, CheckboxSimple("x", "X", false))
	if !strings.Contains(html, "<div>") {
		t.Error("expected div wrapper")
	}
}

func TestCheckboxSanitizesErrorID(t *testing.T) {
	html := testutil.RenderToString(t, Checkbox(CheckboxProps{Name: "terms accepted", Label: "Agree", Error: "Required"}))
	assertDescribedByResolves(t, html, 1)
}

func TestCheckboxWithHint(t *testing.T) {
	html := testutil.RenderToString(t, Checkbox(CheckboxProps{
		Name: "agree", Label: "I agree", Hint: "You can opt out later",
	}))
	if !strings.Contains(html, `<span id="agree-hint" class="hint block">You can opt out later</span>`) {
		t.Errorf("expected hint span, got %s", html)
	}
	// The hint sits inside the label, so the whole label stays a click target.
	if !strings.Contains(html, `class="hint block">You can opt out later</span></span></label>`) {
		t.Errorf("expected hint nested inside the label, got %s", html)
	}
}

// A wrapping label folds the hint into the accessible name unless the input
// names itself from the label span alone.
func TestCheckboxHintAccessibleAssociation(t *testing.T) {
	html := testutil.RenderToString(t, Checkbox(CheckboxProps{
		Name: "agree", Label: "I agree", Hint: "You can opt out later",
	}))
	if !strings.Contains(html, `aria-labelledby="agree-label"`) {
		t.Errorf("expected aria-labelledby, got %s", html)
	}
	if !strings.Contains(html, `aria-describedby="agree-hint"`) {
		t.Errorf("expected aria-describedby, got %s", html)
	}
	assertDescribedByResolves(t, html, 1)
	if !strings.Contains(html, `id="agree-label"`) {
		t.Error("aria-labelledby resolves to no element")
	}
}

func TestCheckboxWithoutHintRendersUnchanged(t *testing.T) {
	html := testutil.RenderToString(t, CheckboxSimple("agree", "I agree", false))
	for _, unwanted := range []string{`class="hint block"`, "aria-describedby", "aria-labelledby", "<span"} {
		if strings.Contains(html, unwanted) {
			t.Errorf("expected no %s without a hint, got %s", unwanted, html)
		}
	}
	if !strings.Contains(html, `<input type="checkbox" name="agree">I agree`) {
		t.Errorf("expected unchanged markup, got %s", html)
	}
}

func TestCheckboxHintAndErrorBothAssociated(t *testing.T) {
	html := testutil.RenderToString(t, Checkbox(CheckboxProps{
		Name: "agree", Label: "I agree", Hint: "Opt out later", Error: "Required",
	}))
	if !strings.Contains(html, `aria-describedby="agree-hint agree-error"`) {
		t.Errorf("expected hint then error, got %s", html)
	}
	assertDescribedByResolves(t, html, 2)
	if !strings.Contains(html, "checkbox--error") {
		t.Error("expected error variant")
	}
	if !strings.Contains(html, `role="alert"`) {
		t.Error("expected role=alert")
	}
}

// Checkboxes in a set share a name, so their hint ids must be disambiguated
// by value or the wrong hint gets announced.
func TestCheckboxHintIDsUniqueAcrossASharedName(t *testing.T) {
	a := testutil.RenderToString(t, Checkbox(CheckboxProps{Name: "opts", Value: "a", Label: "A", Hint: "First"}))
	b := testutil.RenderToString(t, Checkbox(CheckboxProps{Name: "opts", Value: "b", Label: "B", Hint: "Second"}))
	if !strings.Contains(a, `id="opts-a-hint"`) || !strings.Contains(b, `id="opts-b-hint"`) {
		t.Errorf("expected value-disambiguated hint ids, got %s and %s", a, b)
	}
}

func TestCheckboxHintPreservesState(t *testing.T) {
	html := testutil.RenderToString(t, Checkbox(CheckboxProps{
		Name: "opts", Value: "a", Label: "A", Hint: "H",
		Checked: true, Disabled: true, Variant: "large", Class: "extra",
		Attrs: templ.Attributes{"data-testid": "cb"},
	}))
	for _, want := range []string{"checked", "disabled", "checkbox--large", "extra", `data-testid="cb"`, `value="a"`, `class="hint block"`} {
		if !strings.Contains(html, want) {
			t.Errorf("expected %s alongside a hint, got %s", want, html)
		}
	}
}

func TestCheckboxHintEscaping(t *testing.T) {
	html := testutil.RenderToString(t, Checkbox(CheckboxProps{
		Name: "agree", Label: `<b>L</b>`, Hint: `a & "b" <script>alert(1)</script>`,
	}))
	if strings.Contains(html, "<script>") || strings.Contains(html, "<b>") {
		t.Errorf("expected escaping, got %s", html)
	}
	if !strings.Contains(html, "&lt;script&gt;") || !strings.Contains(html, "&amp;") {
		t.Errorf("expected escaped entities, got %s", html)
	}
}

func TestCheckboxHintIDSanitized(t *testing.T) {
	html := testutil.RenderToString(t, Checkbox(CheckboxProps{
		Name: "opt group", Value: "a b", Label: "A", Hint: "H",
	}))
	assertDescribedByResolves(t, html, 1)
	if strings.Contains(html, `id="opt group-a b-hint"`) {
		t.Error("hint id must not contain whitespace")
	}
}

// With no visible label there is no span to name the input from, so
// aria-labelledby must be omitted rather than point at nothing.
func TestCheckboxHintWithoutLabelOmitsLabelledby(t *testing.T) {
	html := testutil.RenderToString(t, Checkbox(CheckboxProps{
		Name: "agree", Hint: "You can opt out later",
		Attrs: templ.Attributes{"aria-label": "I agree"},
	}))
	if strings.Contains(html, "aria-labelledby") {
		t.Errorf("expected no aria-labelledby without a label, got %s", html)
	}
	if !strings.Contains(html, `aria-label="I agree"`) {
		t.Error("expected caller aria-label preserved")
	}
	assertDescribedByResolves(t, html, 1)
}
