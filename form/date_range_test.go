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
