package component

import (
	"testing"

	"github.com/Design-Machines-Studio/livewires-templ/internal/testutil"
)

func TestDialogRenders(t *testing.T) {
	html := testutil.RenderToString(t, Dialog("Confirm"))
	if html == "" {
		t.Fatal("expected non-empty output")
	}
}
