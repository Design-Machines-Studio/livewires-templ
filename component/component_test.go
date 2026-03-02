package component

import (
	"bytes"
	"context"
	"testing"

	"github.com/a-h/templ"
)

// renderToString is a test helper that renders a templ component to a string.
func renderToString(t *testing.T, c templ.Component) string {
	t.Helper()
	var buf bytes.Buffer
	err := c.Render(context.Background(), &buf)
	if err != nil {
		t.Fatalf("render failed: %v", err)
	}
	return buf.String()
}

func TestButtonRenders(t *testing.T) {
	html := renderToString(t, Button("Click me", "accent"))
	if html == "" {
		t.Fatal("expected non-empty output")
	}
}

func TestBadgeRenders(t *testing.T) {
	html := renderToString(t, Badge("Active", "green"))
	if html == "" {
		t.Fatal("expected non-empty output")
	}
}

func TestAvatarRenders(t *testing.T) {
	html := renderToString(t, Avatar("John Doe", ""))
	if html == "" {
		t.Fatal("expected non-empty output")
	}
}

func TestStatCardRenders(t *testing.T) {
	html := renderToString(t, StatCardSimple("Revenue", "$12,450", "+12%", ""))
	if html == "" {
		t.Fatal("expected non-empty output")
	}
}

func TestInfoCardRenders(t *testing.T) {
	html := renderToString(t, InfoCard(InfoCardProps{Title: "Board meeting", Header: "Upcoming"}))
	if html == "" {
		t.Fatal("expected non-empty output")
	}
}

func TestActivityCardRenders(t *testing.T) {
	html := renderToString(t, ActivityCard(ActivityCardProps{Title: "Proposal A", BadgeLabel: "Active", BadgeVariant: "active"}))
	if html == "" {
		t.Fatal("expected non-empty output")
	}
}

func TestCommentRenders(t *testing.T) {
	html := renderToString(t, Comment("Jane", "Hello world", "2024-03-15", ""))
	if html == "" {
		t.Fatal("expected non-empty output")
	}
}

func TestChecklistRenders(t *testing.T) {
	items := []ChecklistItemProps{
		{Label: "Task 1", Completed: true},
		{Label: "Task 2", Completed: false},
	}
	html := renderToString(t, Checklist(items))
	if html == "" {
		t.Fatal("expected non-empty output")
	}
}

func TestProgressBarRenders(t *testing.T) {
	html := renderToString(t, ProgressBar(65, "", "Progress"))
	if html == "" {
		t.Fatal("expected non-empty output")
	}
}

func TestLoaderRenders(t *testing.T) {
	html := renderToString(t, Loader("members", ""))
	if html == "" {
		t.Fatal("expected non-empty output")
	}
}

func TestDialogRenders(t *testing.T) {
	html := renderToString(t, Dialog("Confirm"))
	if html == "" {
		t.Fatal("expected non-empty output")
	}
}

func TestPaginationRenders(t *testing.T) {
	html := renderToString(t, Pagination(PaginationProps{CurrentPage: 2, TotalPages: 5, BaseURL: "/items"}))
	if html == "" {
		t.Fatal("expected non-empty output")
	}
}

func TestTabsRenders(t *testing.T) {
	tabs := TabsProps{
		Items: []TabProps{
			{Label: "Overview", Href: "#overview", Active: true},
			{Label: "Details", Href: "#details"},
		},
	}
	html := renderToString(t, Tabs(tabs))
	if html == "" {
		t.Fatal("expected non-empty output")
	}
}

func TestToastRenders(t *testing.T) {
	html := renderToString(t, Toast("Saved!", "success"))
	if html == "" {
		t.Fatal("expected non-empty output")
	}
}
