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
