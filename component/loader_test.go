package component

import (
	"testing"

	"github.com/Design-Machines-Studio/livewires-templ/internal/testutil"
)

func TestLoaderRenders(t *testing.T) {
	html := testutil.RenderToString(t, Loader("members", ""))
	if html == "" {
		t.Fatal("expected non-empty output")
	}
}
