package layout

import (
	"testing"

	"github.com/Design-Machines-Studio/livewires-templ/internal/testutil"
)

func TestSectionRenders(t *testing.T) {
	html := testutil.RenderToString(t, Section(SectionProps{Heading: "About", Scheme: "dark"}))
	if html == "" {
		t.Fatal("expected non-empty output")
	}
}
