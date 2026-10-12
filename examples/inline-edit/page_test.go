package inlineedit

import (
	"strings"
	"testing"

	"github.com/Design-Machines-Studio/livewires-templ/internal/testutil"
)

func TestPinnedInlineEditExampleRendersNativeVariants(t *testing.T) {
	html := testutil.RenderToString(t, Page())
	for _, want := range []string{
		producerCommit, producerBaseCSS, producerInlineEditCSS, producerInlineEditJS,
		`data-inline-edit="compact"`, `data-inline-edit-panel`, `<details`, `<summary`,
		`name="title-a"`, `name="title-b"`, `method="get"`,
		`href="#title-a-panel"`, `data-inline-edit-status`,
		`aria-label="Edit first title"`, `>Edit</a>`, `Without JavaScript`,
		`query parameters; use fictional text`, `persists nothing`,
	} {
		if !strings.Contains(html, want) {
			t.Errorf("rendered example missing %q: %s", want, html)
		}
	}
}
