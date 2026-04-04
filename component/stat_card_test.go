package component

import (
	"strings"
	"testing"

	"github.com/Design-Machines-Studio/livewires-templ/internal/testutil"
)

func TestStatCardUsesCardClass(t *testing.T) {
	html := testutil.RenderToString(t, StatCardSimple("Revenue", "$12,450", "+12%", ""))
	if !strings.Contains(html, "card--stat") {
		t.Error("expected card--stat class")
	}
	if !strings.Contains(html, "card") {
		t.Error("expected card base class")
	}
}

func TestStatCardChildClasses(t *testing.T) {
	html := testutil.RenderToString(t, StatCardSimple("Revenue", "$12,450", "+12%", ""))
	if !strings.Contains(html, `class="title"`) {
		t.Error("expected title class on label")
	}
	if !strings.Contains(html, `class="value"`) {
		t.Error("expected value class on value")
	}
	if !strings.Contains(html, `class="description"`) {
		t.Error("expected description class on detail")
	}
}

func TestStatCardLink(t *testing.T) {
	html := testutil.RenderToString(t, StatCardLink("Revenue", "$12,450", "+12%", "/stats", ""))
	if !strings.Contains(html, "<a") {
		t.Error("expected anchor tag when href provided")
	}
	if !strings.Contains(html, `<dt class="title"><a`) {
		t.Error("expected link inside dt element, not wrapping dt/dd")
	}
}

func TestStatCardWithStatus(t *testing.T) {
	html := testutil.RenderToString(t, StatCardComponent(StatCardProps{
		Label:  "Compliance",
		Value:  "98%",
		Status: "success",
	}))
	if !strings.Contains(html, "status-indicator") {
		t.Error("expected status-indicator when Status provided")
	}
	if !strings.Contains(html, "status-indicator--success") {
		t.Error("expected status-indicator--success variant")
	}
}

func TestStatCardSubDetail(t *testing.T) {
	html := testutil.RenderToString(t, StatCardComponent(StatCardProps{
		Label:     "Revenue",
		Value:     "$12,450",
		Detail:    "+12%",
		SubDetail: "vs last quarter",
	}))
	if !strings.Contains(html, `class="sub-detail"`) {
		t.Error("expected sub-detail class when SubDetail provided")
	}
	if !strings.Contains(html, "vs last quarter") {
		t.Error("expected sub-detail text")
	}
}

func TestStatCardSizeDefault(t *testing.T) {
	html := testutil.RenderToString(t, StatCardSimple("Revenue", "$12,450", "+12%", ""))
	if !strings.Contains(html, "text-sm") {
		t.Error("expected default text-sm size class")
	}
}

func TestStatCardSizeOverride(t *testing.T) {
	html := testutil.RenderToString(t, StatCardComponent(StatCardProps{
		Label: "Revenue",
		Value: "$12,450",
		Size:  "base",
	}))
	if !strings.Contains(html, "text-base") {
		t.Error("expected text-base size class")
	}
	if strings.Contains(html, "text-sm") {
		t.Error("should not contain default text-sm when overridden")
	}
}

func TestStatCardValueSize(t *testing.T) {
	html := testutil.RenderToString(t, StatCardComponent(StatCardProps{
		Label:     "Revenue",
		Value:     "$12,450",
		ValueSize: "6xl",
	}))
	if !strings.Contains(html, "text-6xl") {
		t.Error("expected text-6xl class on value element")
	}
}

func TestStatCardValueSizeDefault(t *testing.T) {
	html := testutil.RenderToString(t, StatCardSimple("Revenue", "$12,450", "+12%", ""))
	if strings.Contains(html, "text-6xl") {
		t.Error("should not contain text-6xl when ValueSize not set")
	}
}

func TestStatCardWithProgress(t *testing.T) {
	html := testutil.RenderToString(t, StatCardComponent(StatCardProps{
		Label:    "Completion",
		Value:    "75%",
		Progress: 75,
	}))
	if !strings.Contains(html, "progress-bar") {
		t.Error("expected progress-bar class when Progress > 0")
	}
	if !strings.Contains(html, `aria-valuenow="75"`) {
		t.Error("expected aria-valuenow to reflect Progress value")
	}
	if !strings.Contains(html, `role="progressbar"`) {
		t.Error("expected role=progressbar from ProgressBar component")
	}
	if !strings.Contains(html, `class="progress"`) {
		t.Error("expected progress class on wrapping dd")
	}
	if !strings.Contains(html, "progress-bar--thin") {
		t.Error("expected progress-bar--thin variant for stat card progress bar")
	}
}

func TestStatCardNoProgressBarByDefault(t *testing.T) {
	html := testutil.RenderToString(t, StatCardSimple("Revenue", "$12,450", "+12%", ""))
	if strings.Contains(html, "progress-bar") {
		t.Error("should not render progress-bar when Progress is 0")
	}
}

func TestStatCardProgressAriaLabel(t *testing.T) {
	html := testutil.RenderToString(t, StatCardComponent(StatCardProps{
		Label:    "Completion",
		Value:    "75%",
		Progress: 75,
	}))
	if !strings.Contains(html, `aria-label="Completion progress"`) {
		t.Error(`expected aria-label="Completion progress" on progress bar`)
	}
}

func TestStatCardProgressEmptyLabel(t *testing.T) {
	html := testutil.RenderToString(t, StatCardComponent(StatCardProps{
		Value:    "75%",
		Progress: 75,
	}))
	if strings.Contains(html, `aria-label=" progress"`) {
		t.Error("should not produce aria-label with leading space when Label is empty")
	}
}

func TestStatCardProgressClamped(t *testing.T) {
	html := testutil.RenderToString(t, StatCardComponent(StatCardProps{
		Label:    "Over",
		Value:    "150%",
		Progress: 150,
	}))
	if !strings.Contains(html, `aria-valuenow="100"`) {
		t.Error("expected Progress clamped to 100")
	}
}

func TestStatCardGroupComponent(t *testing.T) {
	html := testutil.RenderToString(t, StatCardGroupComponent(StatCardGroupProps{
		Class: "grid",
	}))
	if !strings.Contains(html, "<dl") {
		t.Error("expected dl element")
	}
	if !strings.Contains(html, `class="grid"`) {
		t.Error("expected class on dl element")
	}
}
