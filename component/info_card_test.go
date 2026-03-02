package component

import (
	"testing"

	"github.com/Design-Machines-Studio/livewires-templ/internal/testutil"
)

func TestInfoCardRenders(t *testing.T) {
	html := testutil.RenderToString(t, InfoCard(InfoCardProps{Title: "Board meeting", Header: "Upcoming"}))
	if html == "" {
		t.Fatal("expected non-empty output")
	}
}
