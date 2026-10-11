package utilityfeedback

import (
	"strings"
	"testing"

	"github.com/Design-Machines-Studio/livewires-templ/internal/testutil"
)

func TestPinnedUtilityExampleRendersFallbackAndManagedHooks(t *testing.T) {
	html := testutil.RenderToString(t, Page())
	for _, want := range []string{
		relativeTimeAsset, copyAsset, loadingAsset, toastAsset,
		liveWiresCSS,
		`data-relative-time-fallback>October 10, 2026`, `data-copy-source`, `readonly`,
		`data-loading-region`, `data-loading-content`, `data-loading-status`,
		`data-toast-region`, `data-toast-list`, `data-toast-id="save-42:confirmed"`,
		`data-toast-kind="success"`, `data-toast-message`, `type="button"`,
	} {
		if !strings.Contains(html, want) {
			t.Errorf("rendered example missing %q: %s", want, html)
		}
	}
}
