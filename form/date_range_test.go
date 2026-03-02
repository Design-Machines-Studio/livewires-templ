package form

import (
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
	if html == "" {
		t.Fatal("expected non-empty output")
	}
}
