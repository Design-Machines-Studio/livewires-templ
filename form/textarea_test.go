package form

import (
	"testing"

	"github.com/Design-Machines-Studio/livewires-templ/internal/testutil"
)

func TestTextareaRenders(t *testing.T) {
	html := testutil.RenderToString(t, TextareaSimple("Message", "message", "", "Write...", ""))
	if html == "" {
		t.Fatal("expected non-empty output")
	}
}
