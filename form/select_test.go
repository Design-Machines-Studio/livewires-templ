package form

import (
	"testing"

	"github.com/Design-Machines-Studio/livewires-templ/internal/testutil"
)

func TestSelectRenders(t *testing.T) {
	opts := []SelectOption{
		{Value: "a", Label: "Option A"},
		{Value: "b", Label: "Option B"},
	}
	html := testutil.RenderToString(t, SelectSimple("Country", "country", "Choose...", opts, false))
	if html == "" {
		t.Fatal("expected non-empty output")
	}
}
