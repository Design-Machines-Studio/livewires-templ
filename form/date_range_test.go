package form

import (
	"strings"
	"testing"

	"github.com/Design-Machines-Studio/livewires-templ/internal/testutil"
	"github.com/a-h/templ"
)

func TestDateRangeRenders(t *testing.T) {
	data := DateRangeProps{
		StartName: "from",
		EndName:   "to",
		Label:     "Date Range",
	}
	html := testutil.RenderToString(t, DateRange(data))
	if !strings.Contains(html, `<fieldset><legend>Date Range</legend><div class="date-range">`) {
		t.Errorf("expected fieldset with legend around div.date-range, got %s", html)
	}
	if !strings.Contains(html, `name="from"`) {
		t.Error("expected start date field name")
	}
	if !strings.Contains(html, `name="to"`) {
		t.Error("expected end date field name")
	}
	if !strings.Contains(html, `aria-label="Start date"`) || !strings.Contains(html, `aria-label="End date"`) {
		t.Error("expected default aria-labels on both inputs")
	}
	if !strings.Contains(html, `<span aria-hidden="true">–</span>`) {
		t.Error("expected hidden separator")
	}
	if strings.Contains(html, "stack") || strings.Contains(html, "cluster") {
		t.Error("must not bake layout primitives into markup")
	}
}

// Without a label or error the root is the reference div.date-range and carries Class and Attrs.
func TestDateRangeBare(t *testing.T) {
	html := testutil.RenderToString(t, DateRange(DateRangeProps{
		StartName:  "r-from",
		EndName:    "r-to",
		StartLabel: "From",
		Small:      true,
		Class:      "extra",
		Attrs:      templ.Attributes{"data-test": "x"},
		EndAttrs:   templ.Attributes{"data-bind": "end"},
	}))
	if strings.Contains(html, "<fieldset") {
		t.Error("expected no fieldset without label or error")
	}
	if !strings.Contains(html, `<div class="date-range date-range--small extra" data-test="x">`) {
		t.Errorf("expected small date-range root with class and attrs, got %s", html)
	}
	if !strings.Contains(html, `aria-label="From"`) {
		t.Error("expected custom start aria-label")
	}
	if !strings.Contains(html, `data-bind="end"`) {
		t.Error("expected EndAttrs on the end input")
	}
}

func TestDateRangeWithError(t *testing.T) {
	html := testutil.RenderToString(t, DateRange(DateRangeProps{
		StartName: "from",
		EndName:   "to",
		Label:     "Date Range",
		Error:     "End date must be after start",
	}))
	if !strings.Contains(html, `<legend class="error">`) {
		t.Error("expected error legend")
	}
	if strings.Count(html, `aria-invalid="true"`) != 2 {
		t.Error("expected aria-invalid on both inputs")
	}
	if strings.Count(html, `aria-describedby="from-error"`) != 2 {
		t.Error("expected both inputs to describe the shared error")
	}
	if !strings.Contains(html, `<p id="from-error" class="error" role="alert">End date must be after start</p>`) {
		t.Error("expected shared error paragraph")
	}
}

// Both inputs point at the single shared error paragraph, so the IDREF must
// resolve after sanitization.
func TestDateRangeSanitizesIDs(t *testing.T) {
	html := testutil.RenderToString(t, DateRange(DateRangeProps{
		StartName: "from date", EndName: "to date", Error: "Invalid range",
	}))
	id := errorID("from date")
	if !strings.Contains(html, `<p id="`+id+`"`) {
		t.Errorf("expected sanitized error paragraph id, got %s", html)
	}
	if strings.Count(html, `aria-describedby="`+id+`"`) != 2 {
		t.Error("expected both inputs to reference the sanitized error id")
	}
	if strings.Contains(html, `id="from date"`) {
		t.Error("ids must not contain whitespace")
	}
}
