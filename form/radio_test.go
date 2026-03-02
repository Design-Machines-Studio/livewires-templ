package form

import (
	"testing"

	"github.com/Design-Machines-Studio/livewires-templ/internal/testutil"
)

func TestRadioGroupRenders(t *testing.T) {
	data := RadioGroupProps{
		Name: "choice",
		Options: []RadioOption{
			{Value: "a", Label: "Option A"},
			{Value: "b", Label: "Option B"},
		},
		Selected: "a",
	}
	html := testutil.RenderToString(t, RadioGroup(data))
	if html == "" {
		t.Fatal("expected non-empty output")
	}
}
