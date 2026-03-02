package component

import (
	"testing"

	"github.com/Design-Machines-Studio/livewires-templ/internal/testutil"
)

func TestBadgeRenders(t *testing.T) {
	html := testutil.RenderToString(t, Badge("Active", "green"))
	if html == "" {
		t.Fatal("expected non-empty output")
	}
}
