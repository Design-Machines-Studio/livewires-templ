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
	html := testutil.RenderToString(t, Avatar("Test", "lg"))
	if !strings.Contains(html, "avatar--lg") {
		t.Error("expected avatar--lg class")
	}
}

func TestAvatarShowName(t *testing.T) {
	html := testutil.RenderToString(t, AvatarComponent(AvatarProps{Name: "John Doe", ShowName: true}))
	if !strings.Contains(html, `class="cluster items-center"`) {
		t.Error("expected cluster items-center wrapper when ShowName is true")
	}
	if !strings.Contains(html, "John Doe") {
		t.Error("expected name text")
	}
	if !strings.Contains(html, `role="img"`) {
		t.Error("initials avatar should have role=img when ShowName is true")
	}
	if !strings.Contains(html, `aria-label="John Doe"`) {
		t.Error("initials avatar should keep aria-label when ShowName is true (role=img requires accessible name)")
	}
}

func TestAvatarShowNameImageAltEmpty(t *testing.T) {
	html := testutil.RenderToString(t, AvatarComponent(AvatarProps{
		Name: "Jane Doe", Src: "/img/test.jpg", ShowName: true,
	}))
	if !strings.Contains(html, `alt=""`) {
		t.Error("expected empty alt when ShowName is true (visible text provides name)")
	}
	if !strings.Contains(html, `class="cluster items-center"`) {
		t.Error("expected cluster items-center wrapper")
	}
}

func TestAvatarShowNameNoClusterWhenFalse(t *testing.T) {
	html := testutil.RenderToString(t, Avatar("John Doe", ""))
	if strings.Contains(html, `class="cluster"`) {
		t.Error("should not have cluster wrapper when ShowName is false")
	}
}

func TestAvatarAriaLabel(t *testing.T) {
	html := testutil.RenderToString(t, Avatar("Jane Doe", ""))
	if !strings.Contains(html, `aria-label="Jane Doe"`) {
		t.Error("expected aria-label on initials avatar when ShowName is false")
	}
	imgHtml := testutil.RenderToString(t, AvatarImage("/img/test.jpg", "Jane Doe", ""))
	if strings.Contains(imgHtml, `aria-label`) {
		t.Error("should not have aria-label on image avatar (img alt is sufficient)")
	}
}

func TestAvatarSmall(t *testing.T) {
	html := testutil.RenderToString(t, AvatarSmall("Test"))
	if !strings.Contains(html, "avatar--sm") {
		t.Error("expected avatar--sm class")
	}
}

func TestAvatarSquare(t *testing.T) {
	html := testutil.RenderToString(t, AvatarComponent(AvatarProps{Name: "Test", Size: "sm", Square: true}))
	if !strings.Contains(html, "avatar--sm") {
		t.Error("expected avatar--sm class")
	}
	if !strings.Contains(html, "avatar--square") {
		t.Error("expected avatar--square class")
	}
}

func TestAvatarRoleImg(t *testing.T) {
	html := testutil.RenderToString(t, Avatar("Jane Doe", ""))
	if !strings.Contains(html, `role="img"`) {
		t.Error("expected role=img on initials avatar")
	}
	imgHtml := testutil.RenderToString(t, AvatarImage("/img/test.jpg", "Jane Doe", ""))
	if strings.Contains(imgHtml, `role="img"`) {
		t.Error("should not have role=img on image avatar")
	}
}

func TestAvatarInitialsAriaHidden(t *testing.T) {
	html := testutil.RenderToString(t, Avatar("Jane Doe", ""))
	if !strings.Contains(html, `aria-hidden="true"`) {
		t.Error("expected aria-hidden on initials span")
	}
}

func TestAvatarSingleInitialSmallSizes(t *testing.T) {
	xsHtml := testutil.RenderToString(t, Avatar("John Doe", "xs"))
	if strings.Contains(xsHtml, "JD") {
		t.Error("xs size should use single initial, not JD")
	}
	if !strings.Contains(xsHtml, ">J<") {
		t.Error("xs size should show single initial J")
	}

	smHtml := testutil.RenderToString(t, AvatarSmall("John Doe"))
	if strings.Contains(smHtml, "JD") {
		t.Error("sm size should use single initial, not JD")
	}
	if !strings.Contains(smHtml, ">J<") {
		t.Error("sm size should show single initial J")
	}
}
