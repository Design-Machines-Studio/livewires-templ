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

func TestCommentSizeDefault(t *testing.T) {
	html := testutil.RenderToString(t, Comment("Jane", "Hello", "2024-03-15", ""))
	if !strings.Contains(html, "text-sm") {
		t.Error("expected default text-sm size class")
	}
}

func TestCommentStructuralUtilities(t *testing.T) {
	html := testutil.RenderToString(t, Comment("Jane", "Hello", "2024-03-15", ""))
	if !strings.Contains(html, "items-start") {
		t.Error("expected items-start structural utility")
	}
	if !strings.Contains(html, "flex-1") {
		t.Error("expected flex-1 structural utility on comment-body")
	}
}

func TestCommentNoThemeUtilities(t *testing.T) {
	html := testutil.RenderToString(t, Comment("Jane", "Hello", "2024-03-15", ""))
	if strings.Contains(html, `class="author font-semibold"`) {
		t.Error("author should not have font-semibold utility class")
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
