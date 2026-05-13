package form

import (
	"strings"
	"testing"

	"github.com/Design-Machines-Studio/livewires-templ/internal/testutil"
)

func TestRadioGroupRenders(t *testing.T) {
	data := RadioGroupProps{
		Name: "choice",
		Options: []RadioOption{
			{Value: "a", Label: "Option A"},
			{Value: "b", Label: "Option B"},
		},
		Selected: "a",
	}
	html := testutil.RenderToString(t, RadioGroup(data))
	if !strings.Contains(html, `class="check-list"`) {
		t.Error("expected check-list class")
	}
	if !strings.Contains(html, `type="radio"`) {
		t.Error("expected radio input type")
	}
	if !strings.Contains(html, "<li>") {
		t.Error("expected list items")
	}
	if !strings.Contains(html, "Option A") {
		t.Error("expected option label text")
	}
}

func TestRadioGroupSelectedChecked(t *testing.T) {
	data := RadioGroupProps{
		Name: "choice",
		Options: []RadioOption{
			{Value: "a", Label: "Option A"},
			{Value: "b", Label: "Option B"},
		},
		Selected: "a",
	}
	html := testutil.RenderToString(t, RadioGroup(data))
	if !strings.Contains(html, "checked") {
		t.Error("expected checked attribute on selected option")
	}
}

func TestRadioGroupInline(t *testing.T) {
	data := RadioGroupProps{
		Name: "choice",
		Options: []RadioOption{
			{Value: "a", Label: "A"},
		},
		Inline: true,
	}
	html := testutil.RenderToString(t, RadioGroup(data))
	if !strings.Contains(html, "check-list--inline") {
		t.Error("expected check-list--inline class")
	}
}

func TestRadioSimple(t *testing.T) {
	html := testutil.RenderToString(t, RadioSimple("color", "red", "Red", false))
	if !strings.Contains(html, `class="radio"`) {
		t.Error("expected radio class")
	}
	if !strings.Contains(html, `name="color"`) {
		t.Error("expected name attribute")
	}
	if !strings.Contains(html, `value="red"`) {
		t.Error("expected value attribute")
	}
}

func TestRadioVariantSuccess(t *testing.T) {
	html := testutil.RenderToString(t, Radio(RadioProps{
		Name:    "status",
		Value:   "ok",
		Label:   "Success",
		Variant: "success",
	}))
	if !strings.Contains(html, "radio--success") {
		t.Error("expected radio--success class")
	}
}

func TestRadioWithError(t *testing.T) {
	html := testutil.RenderToString(t, Radio(RadioProps{
		Name:  "color",
		Value: "red",
		Label: "Red",
		Error: "Pick a color",
	}))
	if !strings.Contains(html, "radio--error") {
		t.Error("expected radio--error variant when Error is set")
	}
	if !strings.Contains(html, `aria-invalid="true"`) {
		t.Error("expected aria-invalid")
	}
	if !strings.Contains(html, `aria-describedby="color-error"`) {
		t.Error("expected aria-describedby")
	}
	if !strings.Contains(html, `id="color-error"`) {
		t.Error("expected error message id")
	}
	if !strings.Contains(html, "Pick a color") {
		t.Error("expected error message text")
	}
}

func TestRadioWrapper(t *testing.T) {
	html := testutil.RenderToString(t, RadioSimple("x", "a", "A", false))
	if !strings.Contains(html, "<div>") {
		t.Error("expected div wrapper")
	}
}

func TestRadioGroupWithError(t *testing.T) {
	html := testutil.RenderToString(t, RadioGroup(RadioGroupProps{
		Name: "choice",
		Options: []RadioOption{
			{Value: "a", Label: "A"},
			{Value: "b", Label: "B"},
		},
		Error: "Please select one",
	}))
	if !strings.Contains(html, "radio--error") {
		t.Error("expected radio--error variant on group radios")
	}
	if !strings.Contains(html, `id="choice-error"`) {
		t.Error("expected group error message id")
	}
	if !strings.Contains(html, `role="alert"`) {
		t.Error("expected role=alert")
	}
	if !strings.Contains(html, "Please select one") {
		t.Error("expected group error message text")
	}
}
