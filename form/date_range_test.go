package form

import (
	"strings"
	"testing"

	"github.com/Design-Machines-Studio/livewires-templ/internal/testutil"
)

func TestDateRangeRenders(t *testing.T) {
	data := DateRangeProps{
		StartName: "from",
		EndName:   "to",
		Label:     "Date Range",
	}
	html := testutil.RenderToString(t, DateRange(data))
	if !strings.Contains(html, `<div class="date-range"><input type="date" id="from" name="from" value="" aria-label="Start date"> <span aria-hidden="true">–</span> <input type="date" id="to" name="to" value="" aria-label="End date"></div>`) {
		t.Errorf("expected Live Wires .date-range markup, got %s", html)
	}
	if strings.Contains(html, "stack") || strings.Contains(html, "cluster") || strings.Contains(html, "<label") {
		t.Error("expected no layout classes or visible labels")
	}
	if strings.Contains(html, "<fieldset class") || strings.Contains(html, "<legend class") {
		t.Error("expected no empty class attributes")
	}
	if !strings.Contains(html, "Date Range") {
		t.Error("expected legend text")
	}
	if !strings.Contains(html, `name="from"`) {
		t.Error("expected start date field name")
	}
	if !strings.Contains(html, `name="to"`) {
		t.Error("expected end date field name")
	}
}

func TestDateRangeWithError(t *testing.T) {
	html := testutil.RenderToString(t, DateRange(DateRangeProps{
		StartName: "from",
		EndName:   "to",
		Label:     "Date Range",
		Error:     "End date must be after start",
	}))
	if !strings.Contains(html, `class="error"`) {
		t.Error("expected error class")
	}
	if !strings.Contains(html, `aria-invalid="true"`) {
		t.Error("expected aria-invalid on inputs")
	}
	if !strings.Contains(html, `id="from-error"`) {
		t.Error("expected error message id")
	}
	if !strings.Contains(html, `role="alert"`) {
		t.Error("expected role=alert")
	}
	if !strings.Contains(html, "End date must be after start") {
		t.Error("expected error message text")
	}
}

// Both inputs point at the single shared error paragraph, so both IDREFs must
// resolve after sanitization.
func TestDateRangeSanitizesIDs(t *testing.T) {
	html := testutil.RenderToString(t, DateRange(DateRangeProps{
		StartName: "from date", EndName: "to date", Label: "Range", Error: "Invalid range",
	}))
	if strings.Count(html, `aria-describedby="`+errorID("from date")+`"`) != 2 {
		t.Errorf("expected both inputs to reference the sanitized error id, got %s", html)
	}
	if !strings.Contains(html, `<p id="`+errorID("from date")+`"`) {
		t.Errorf("expected sanitized error paragraph id, got %s", html)
	}
	if strings.Contains(html, `id="from date"`) || strings.Contains(html, `id="to date"`) {
		t.Error("ids must not contain whitespace")
	}
}
