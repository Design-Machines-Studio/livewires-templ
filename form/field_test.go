package form

import (
	"testing"

	"github.com/Design-Machines-Studio/livewires-templ/internal/testutil"
)

func TestFieldRenders(t *testing.T) {
	html := testutil.RenderToString(t, FieldText("Name", "name", "", true))
	if html == "" {
		t.Fatal("expected non-empty output")
	}
}
