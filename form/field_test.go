package form

import (
	"strings"
	"testing"

	"github.com/Design-Machines-Studio/livewires-templ/internal/testutil"
)

func TestFieldRenders(t *testing.T) {
	html := testutil.RenderToString(t, FieldText("Name", "name", "", true))
	if !strings.Contains(html, "field") {
		t.Error("expected field class")
	}
	if !strings.Contains(html, "Name") {
		t.Error("expected label text")
	}
	if !strings.Contains(html, `name="name"`) {
		t.Error("expected input name attribute")
	}
	if !strings.Contains(html, "required") {
		t.Error("expected required attribute")
	}
}

func TestFieldWithHint(t *testing.T) {
	html := testutil.RenderToString(t, Field(FieldProps{Label: "Email", Name: "email", Type: "email", Hint: "We won't share it"}))
	if !strings.Contains(html, "hint") {
		t.Error("expected hint class")
	}
	if !strings.Contains(html, "We won&#39;t share it") && !strings.Contains(html, "We won't share it") && !strings.Contains(html, "We won&#x27;t share it") {
		t.Error("expected hint text")
	}
}

func TestFieldWithError(t *testing.T) {
	html := testutil.RenderToString(t, Field(FieldProps{
		Label: "Email",
		Name:  "email",
		Type:  "email",
		Error: "Invalid email",
	}))
	if !strings.Contains(html, `class="error"`) {
		t.Error("expected error class on label or input")
	}
	if !strings.Contains(html, `aria-invalid="true"`) {
		t.Error("expected aria-invalid on input")
	}
	if !strings.Contains(html, `aria-describedby="email-error"`) {
		t.Error("expected aria-describedby linking to error message")
	}
	if !strings.Contains(html, `id="email-error"`) {
		t.Error("expected error message id")
	}
	if !strings.Contains(html, `role="alert"`) {
		t.Error("expected role=alert on error message")
	}
	if !strings.Contains(html, "Invalid email") {
		t.Error("expected error message text")
	}
}

func TestFieldWithoutErrorNoErrorMarkup(t *testing.T) {
	html := testutil.RenderToString(t, FieldText("Name", "name", "", false))
	if strings.Contains(html, `aria-invalid`) {
		t.Error("expected no aria-invalid when no error")
	}
	if strings.Contains(html, `role="alert"`) {
		t.Error("expected no role=alert when no error")
	}
}

func TestFieldAriaRequired(t *testing.T) {
	html := testutil.RenderToString(t, FieldText("Name", "name", "", true))
	if !strings.Contains(html, `aria-required="true"`) {
		t.Error("expected aria-required when required")
	}
}

func TestFieldAriaRequiredAbsentWhenNotRequired(t *testing.T) {
	html := testutil.RenderToString(t, FieldText("Name", "name", "", false))
	if strings.Contains(html, `aria-required`) {
		t.Error("expected no aria-required when not required")
	}
}

func TestFieldAriaRequiredWithError(t *testing.T) {
	html := testutil.RenderToString(t, Field(FieldProps{
		Label:    "Email",
		Name:     "email",
		Required: true,
		Error:    "Required",
	}))
	if !strings.Contains(html, `aria-required="true"`) {
		t.Error("expected aria-required in error branch")
	}
}
