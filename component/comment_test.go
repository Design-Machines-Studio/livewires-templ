package component

import (
	"strings"
	"testing"

	"github.com/Design-Machines-Studio/livewires-templ/internal/testutil"
)

func TestCommentRenders(t *testing.T) {
	html := testutil.RenderToString(t, Comment("Jane", "Hello world", "2024-03-15", ""))
	if !strings.Contains(html, "comment") {
		t.Error("expected comment class")
	}
	if !strings.Contains(html, "Jane") {
		t.Error("expected author name")
	}
	if !strings.Contains(html, "Hello world") {
		t.Error("expected body text")
	}
	if !strings.Contains(html, "author") {
		t.Error("expected author class")
	}
}

func TestCommentWithBadge(t *testing.T) {
	html := testutil.RenderToString(t, Comment("Alice", "Great work", "2024-01-01", "Member"))
	if !strings.Contains(html, "Member") {
		t.Error("expected badge label")
	}
	if !strings.Contains(html, "badge") {
		t.Error("expected badge class")
	}
}
