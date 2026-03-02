package form

import (
	"testing"

	"github.com/Design-Machines-Studio/livewires-templ/internal/testutil"
)

func TestCheckboxRenders(t *testing.T) {
	html := testutil.RenderToString(t, CheckboxSimple("agree", "I agree", false))
	if html == "" {
		t.Fatal("expected non-empty output")
	}
}
