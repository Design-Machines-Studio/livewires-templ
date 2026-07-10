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

// The hint must be announced, not merely displayed.
func TestFieldHintIsAssociatedWithInput(t *testing.T) {
	html := testutil.RenderToString(t, Field(FieldProps{Label: "Email", Name: "email", Hint: "We never share it"}))
	if !strings.Contains(html, `aria-describedby="email-hint"`) {
		t.Errorf("expected input described by its hint, got %s", html)
	}
	if !strings.Contains(html, `<p id="email-hint" class="hint">`) {
		t.Errorf("expected hint paragraph to carry the referenced id, got %s", html)
	}
	assertDescribedByResolves(t, html, 1)
}

func TestFieldHintAndErrorBothAssociated(t *testing.T) {
	html := testutil.RenderToString(t, Field(FieldProps{
		Label: "Email", Name: "email", Hint: "We never share it", Error: "Invalid email",
	}))
	// Hint first, then error: guidance before failure.
	if !strings.Contains(html, `aria-describedby="email-hint email-error"`) {
		t.Errorf("expected hint then error in aria-describedby, got %s", html)
	}
	assertDescribedByResolves(t, html, 2)
}

func TestFieldWithoutHintHasNoDescription(t *testing.T) {
	html := testutil.RenderToString(t, FieldText("Name", "name", "", false))
	if strings.Contains(html, "aria-describedby") {
		t.Errorf("expected no aria-describedby without hint or error, got %s", html)
	}
	if strings.Contains(html, `class="hint"`) {
		t.Error("expected no empty hint element")
	}
}

// A field name need not be id-safe. The id is sanitized; the submitted name is not.
func TestFieldSanitizesIDButPreservesName(t *testing.T) {
	html := testutil.RenderToString(t, Field(FieldProps{
		Label: "Email", Name: "user[email]", Hint: "Hint", Error: "Bad",
	}))
	assertLabelPointsAtControl(t, html, "user[email]")
	assertDescribedByResolves(t, html, 2)
	if !strings.Contains(html, `name="user[email]"`) {
		t.Errorf("submitted name must not be sanitized, got %s", html)
	}
	if strings.Contains(html, `id="user[email]"`) {
		t.Error("id must not contain brackets")
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
