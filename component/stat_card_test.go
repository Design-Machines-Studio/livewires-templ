package component

import (
	"strings"
	"testing"

	"github.com/Design-Machines-Studio/livewires-templ/internal/testutil"
)

func TestStatCardUsesCardClass(t *testing.T) {
	html := testutil.RenderToString(t, StatCardSimple("Revenue", "$12,450", "+12%", ""))
	if strings.Contains(html, "stat-card") {
		t.Error("should not use stat-card class, should use card")
	}
	if !strings.Contains(html, `class="card `) {
		t.Error("expected card base class")
	}
}

func TestStatCardChildClasses(t *testing.T) {
	html := testutil.RenderToString(t, StatCardSimple("Revenue", "$12,450", "+12%", ""))
	if !strings.Contains(html, `class="title"`) {
		t.Error("expected title class on label")
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
	html := testutil.RenderToString(t, StatCard(StatCardProps{
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
	html := testutil.RenderToString(t, StatCard(StatCardProps{
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
