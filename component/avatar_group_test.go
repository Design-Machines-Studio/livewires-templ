package component

import (
	"strings"
	"testing"

	"github.com/Design-Machines-Studio/livewires-templ/internal/testutil"
)

func TestAvatarGroup(t *testing.T) {
	html := testutil.RenderToString(t, AvatarGroup())
	if !strings.Contains(html, "avatar-group") {
		t.Error("expected avatar-group class")
	}
}

func TestAvatarGroupWithClass(t *testing.T) {
	html := testutil.RenderToString(t, AvatarGroupComponent(AvatarGroupProps{Class: "mt-1"}))
	if !strings.Contains(html, "avatar-group mt-1") {
		t.Error("expected custom class appended")
	}
}
