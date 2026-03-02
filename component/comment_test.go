package component

import (
	"testing"

	"github.com/Design-Machines-Studio/livewires-templ/internal/testutil"
)

func TestCommentRenders(t *testing.T) {
	html := testutil.RenderToString(t, Comment("Jane", "Hello world", "2024-03-15", ""))
	if html == "" {
		t.Fatal("expected non-empty output")
	}
}
