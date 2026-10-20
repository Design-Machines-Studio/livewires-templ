package component

import (
	"strings"
	"testing"

	"github.com/Design-Machines-Studio/livewires-templ/internal/testutil"
)

func TestBadgeRenders(t *testing.T) {
	html := testutil.RenderToString(t, Badge("Active", "green"))
	if html == "" {
		t.Fatal("expected non-empty output")
	}
}

func TestBadgePreservesLabelCase(t *testing.T) {
	tests := []struct {
		name   string
		render func(*testing.T) string
	}{
		{name: "Badge", render: func(t *testing.T) string { return testutil.RenderToString(t, Badge("New comments", "accent")) }},
		{name: "BadgeSmall", render: func(t *testing.T) string { return testutil.RenderToString(t, BadgeSmall("New comments", "accent")) }},
		{name: "BadgeComponent", render: func(t *testing.T) string { return testutil.RenderToString(t, BadgeComponent(BadgeProps{Label: "New comments", Variant: "accent"})) }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			html := tt.render(t)
			if !strings.Contains(html, ">New comments</span>") {
				t.Fatalf("expected sentence-case label to be preserved, got %q", html)
			}
		})
	}
}
