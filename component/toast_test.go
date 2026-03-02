package component

import (
	"testing"

	"github.com/Design-Machines-Studio/livewires-templ/internal/testutil"
)

func TestToastRenders(t *testing.T) {
	html := testutil.RenderToString(t, Toast("Saved!", "success"))
	if html == "" {
		t.Fatal("expected non-empty output")
	}
}
