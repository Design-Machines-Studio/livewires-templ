package layout

import (
	"testing"

	"github.com/Design-Machines-Studio/livewires-templ/internal/testutil"
)

func TestBaseRenders(t *testing.T) {
	html := testutil.RenderToString(t, Base(BaseProps{Title: "Test Page", CSSPath: "/css/style.css"}))
	if html == "" {
		t.Fatal("expected non-empty output")
	}
}
