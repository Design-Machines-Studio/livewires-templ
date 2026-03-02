package form

import (
	"testing"

	"github.com/Design-Machines-Studio/livewires-templ/internal/testutil"
)

func TestSearchRenders(t *testing.T) {
	html := testutil.RenderToString(t, SearchSimple("q", "", ""))
	if html == "" {
		t.Fatal("expected non-empty output")
	}
}
