package component

import (
	"testing"

	"github.com/Design-Machines-Studio/livewires-templ/internal/testutil"
)

func TestActivityCardRenders(t *testing.T) {
	html := testutil.RenderToString(t, ActivityCard(ActivityCardProps{Title: "Proposal A", BadgeLabel: "Active", BadgeVariant: "active"}))
	if html == "" {
		t.Fatal("expected non-empty output")
	}
}
