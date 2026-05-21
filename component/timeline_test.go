package component

import (
	"strings"
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

func TestTimelineItemTitleHrefEmpty(t *testing.T) {
	html := testutil.RenderToString(t, TimelineItem(TimelineItemProps{
		Date:  "March 2024",
		Title: "Founded",
	}))
	if !strings.Contains(html, "Founded") {
		t.Errorf("expected literal title, got: %s", html)
	}
	if strings.Contains(html, "<a") {
		t.Errorf("expected no anchor when TitleHref empty, got: %s", html)
	}
}

func TestTimelineItemTitleHrefSet(t *testing.T) {
	html := testutil.RenderToString(t, TimelineItem(TimelineItemProps{
		Date:      "March 2024",
		Title:     "Proposal",
		TitleHref: "/proposals/1",
	}))
	if !strings.Contains(html, `<a href="/proposals/1">Proposal</a>`) {
		t.Errorf("expected anchored title, got: %s", html)
	}
}

func TestTimelineItemIndicatorUnread(t *testing.T) {
	html := testutil.RenderToString(t, TimelineItem(TimelineItemProps{
		Date:      "March 2024",
		Title:     "Founded",
		Indicator: "unread",
	}))
	if !strings.Contains(html, `class="timeline-indicator timeline-indicator--unread"`) {
		t.Errorf("expected indicator class, got: %s", html)
	}
	if !strings.Contains(html, `<span class="sr-only">Unread.</span>`) {
		t.Errorf("expected sr-only word, got: %s", html)
	}
}

func TestTimelineItemIndicatorSlugThrough(t *testing.T) {
	html := testutil.RenderToString(t, TimelineItem(TimelineItemProps{
		Date:      "March 2024",
		Title:     "Founded",
		Indicator: "alert",
	}))
	if !strings.Contains(html, "timeline-indicator--alert") {
		t.Errorf("expected slugged variant class, got: %s", html)
	}
	if !strings.Contains(html, `<span class="sr-only">Alert.</span>`) {
		t.Errorf("expected sr-only word, got: %s", html)
	}
}

func TestTimelineItemIndicatorEmpty(t *testing.T) {
	html := testutil.RenderToString(t, TimelineItem(TimelineItemProps{
		Date:  "March 2024",
		Title: "Founded",
	}))
	if strings.Contains(html, "timeline-indicator") {
		t.Errorf("expected no indicator element, got: %s", html)
	}
}

func TestTimelineItemCombined(t *testing.T) {
	html := testutil.RenderToString(t, TimelineItem(TimelineItemProps{
		Date:      "March 2024",
		Title:     "Proposal",
		TitleHref: "/proposals/1",
		Indicator: "unread",
		Badge:     "New",
	}))
	indicator := strings.Index(html, "timeline-indicator--unread")
	anchor := strings.Index(html, `<a href="/proposals/1">`)
	badge := strings.Index(html, "New")
	if indicator < 0 || anchor < 0 || badge < 0 {
		t.Fatalf("expected indicator, anchor, and badge present, got: %s", html)
	}
	if !(indicator < anchor && anchor < badge) {
		t.Errorf("expected order indicator < anchor < badge, got %d/%d/%d: %s", indicator, anchor, badge, html)
	}
}

func TestTimelineRenders(t *testing.T) {
	html := testutil.RenderToString(t, Timeline(""))
	if html == "" {
		t.Fatal("expected non-empty output")
	}
}
