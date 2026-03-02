package component

import (
	"testing"

	"github.com/Design-Machines-Studio/livewires-templ/internal/testutil"
)

func TestPaginationRenders(t *testing.T) {
	html := testutil.RenderToString(t, Pagination(PaginationProps{CurrentPage: 2, TotalPages: 5, BaseURL: "/items"}))
	if html == "" {
		t.Fatal("expected non-empty output")
	}
}
