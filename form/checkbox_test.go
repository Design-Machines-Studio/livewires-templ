package form

import (
	"strings"
	"testing"

	"github.com/Design-Machines-Studio/livewires-templ/internal/testutil"
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
