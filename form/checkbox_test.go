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
