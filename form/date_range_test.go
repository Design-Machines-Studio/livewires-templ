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
	if !strings.Contains(html, "date-range") {
		t.Error("expected date-range class")
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
	assertLabelPointsAtControl(t, html, "from date")
	assertLabelPointsAtControl(t, html, "to date")
	if !strings.Contains(html, `<p id="`+errorID("from date")+`"`) {
		t.Errorf("expected sanitized error paragraph id, got %s", html)
	}
	if strings.Contains(html, `id="from date"`) || strings.Contains(html, `id="to date"`) {
		t.Error("ids must not contain whitespace")
	}
}
