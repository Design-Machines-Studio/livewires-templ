package form

import (
	"testing"

	"github.com/Design-Machines-Studio/livewires-templ/internal/testutil"
)

func TestSwitchRenders(t *testing.T) {
	html := testutil.RenderToString(t, Switch("notifications", "notif", true))
	if html == "" {
		t.Fatal("expected non-empty output")
	}
}
