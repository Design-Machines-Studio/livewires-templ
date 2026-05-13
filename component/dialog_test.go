package component

import (
	"strings"
	"testing"

	"github.com/Design-Machines-Studio/livewires-templ/internal/testutil"
)

func TestDialogRenders(t *testing.T) {
	html := testutil.RenderToString(t, Dialog("Confirm"))
	if html == "" {
		t.Fatal("expected non-empty output")
	}
}

func TestDialogDefaultH2(t *testing.T) {
	html := testutil.RenderToString(t, DialogComponent(DialogProps{Title: "Confirm"}))
	if !strings.Contains(html, "<h2>") {
		t.Error("expected default h2 heading")
	}
}

func TestDialogHeadingLevelH3(t *testing.T) {
	html := testutil.RenderToString(t, DialogComponent(DialogProps{
		Title:        "Confirm",
		HeadingLevel: "h3",
	}))
	if !strings.Contains(html, "<h3>") {
		t.Error("expected h3 heading")
	}
	if strings.Contains(html, "<h2>") {
		t.Error("expected no h2 when HeadingLevel is h3")
	}
}

func TestDialogHeadingLevelH4(t *testing.T) {
	html := testutil.RenderToString(t, DialogComponent(DialogProps{
		Title:        "Confirm",
		HeadingLevel: "h4",
	}))
	if !strings.Contains(html, "<h4>") {
		t.Error("expected h4 heading")
	}
}
