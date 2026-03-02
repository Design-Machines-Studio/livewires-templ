package component

import (
	"testing"

	"github.com/Design-Machines-Studio/livewires-templ/internal/testutil"
)

func TestButtonRenders(t *testing.T) {
	html := testutil.RenderToString(t, Button("Click me", "accent"))
	if html == "" {
		t.Fatal("expected non-empty output")
	}
}
