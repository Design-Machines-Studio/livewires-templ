package component

import (
	"strings"
	"testing"

	"github.com/Design-Machines-Studio/livewires-templ/internal/testutil"
)

func TestLoaderRenders(t *testing.T) {
	html := testutil.RenderToString(t, Loader("members", ""))
	if html == "" {
		t.Fatal("expected non-empty output")
	}
}

func TestLoaderRoleStatus(t *testing.T) {
	html := testutil.RenderToString(t, Loader("Loading data", ""))
	if !strings.Contains(html, `role="status"`) {
		t.Error("expected role=status on loader")
	}
}

func TestLoaderAriaLabel(t *testing.T) {
	html := testutil.RenderToString(t, Loader("Loading members", ""))
	if !strings.Contains(html, `aria-label="Loading members"`) {
		t.Error("expected aria-label with loader text")
	}
}
