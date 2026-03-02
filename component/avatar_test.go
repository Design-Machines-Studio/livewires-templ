package component

import (
	"testing"

	"github.com/Design-Machines-Studio/livewires-templ/internal/testutil"
)

func TestAvatarRenders(t *testing.T) {
	html := testutil.RenderToString(t, Avatar("John Doe", ""))
	if html == "" {
		t.Fatal("expected non-empty output")
	}
}
