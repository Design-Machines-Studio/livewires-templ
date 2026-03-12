package component

import (
	"strings"
	"testing"

	"github.com/Design-Machines-Studio/livewires-templ/internal/testutil"
)

func TestAvatarInitials(t *testing.T) {
	html := testutil.RenderToString(t, Avatar("John Doe", ""))
	if !strings.Contains(html, `class="avatar"`) {
		t.Error("expected avatar class")
	}
	if !strings.Contains(html, `class="initials"`) {
		t.Error("expected initials span when no Src provided")
	}
	if !strings.Contains(html, "JD") {
		t.Error("expected initials JD")
	}
}

func TestAvatarImage(t *testing.T) {
	html := testutil.RenderToString(t, AvatarImage("/img/test.jpg", "Jane Doe", ""))
	if !strings.Contains(html, "<img") {
		t.Error("expected img tag when Src provided")
	}
	if !strings.Contains(html, `alt="Jane Doe"`) {
		t.Error("expected alt text")
	}
	if strings.Contains(html, "initials") {
		t.Error("should not render initials when Src provided")
	}
}

func TestAvatarSizeVariant(t *testing.T) {
	html := testutil.RenderToString(t, Avatar("Test", "large"))
	if !strings.Contains(html, "avatar--large") {
		t.Error("expected avatar--large class")
	}
}

func TestAvatarShowName(t *testing.T) {
	html := testutil.RenderToString(t, AvatarComponent(AvatarProps{Name: "John Doe", ShowName: true}))
	if !strings.Contains(html, `class="name"`) {
		t.Error("expected name span when ShowName is true")
	}
	if !strings.Contains(html, "John Doe") {
		t.Error("expected name text")
	}
	if strings.Contains(html, `aria-label=`) {
		t.Error("should not have aria-label when ShowName is true")
	}
}

func TestAvatarAriaLabel(t *testing.T) {
	html := testutil.RenderToString(t, Avatar("Jane Doe", ""))
	if !strings.Contains(html, `aria-label="Jane Doe"`) {
		t.Error("expected aria-label when ShowName is false")
	}
}

func TestAvatarSmall(t *testing.T) {
	html := testutil.RenderToString(t, AvatarSmall("Test"))
	if !strings.Contains(html, "avatar--small") {
		t.Error("expected avatar--small class")
	}
}

func TestAvatarLarge(t *testing.T) {
	html := testutil.RenderToString(t, Avatar("Test", "large"))
	if !strings.Contains(html, "avatar--large") {
		t.Error("expected avatar--large class")
	}
}
