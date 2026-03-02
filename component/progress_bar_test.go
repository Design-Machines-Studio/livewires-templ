package component

import (
	"testing"

	"github.com/Design-Machines-Studio/livewires-templ/internal/testutil"
)

func TestProgressBarRenders(t *testing.T) {
	html := testutil.RenderToString(t, ProgressBar(65, "", "Progress"))
	if html == "" {
		t.Fatal("expected non-empty output")
	}
}
