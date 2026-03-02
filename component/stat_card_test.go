package component

import (
	"testing"

	"github.com/Design-Machines-Studio/livewires-templ/internal/testutil"
)

func TestStatCardRenders(t *testing.T) {
	html := testutil.RenderToString(t, StatCardSimple("Revenue", "$12,450", "+12%", ""))
	if html == "" {
		t.Fatal("expected non-empty output")
	}
}
