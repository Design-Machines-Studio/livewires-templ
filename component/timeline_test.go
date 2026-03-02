package component

import (
	"testing"

	"github.com/Design-Machines-Studio/livewires-templ/internal/testutil"
)

func TestTimelineItemRenders(t *testing.T) {
	html := testutil.RenderToString(t, TimelineItem(TimelineItemProps{
		Date:  "March 2024",
		Title: "Founded",
	}))
	if html == "" {
		t.Fatal("expected non-empty output")
	}
}

func TestTimelineRenders(t *testing.T) {
	html := testutil.RenderToString(t, Timeline(""))
	if html == "" {
		t.Fatal("expected non-empty output")
	}
}
