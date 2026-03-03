package component

import (
	"strings"
	"testing"

	"github.com/Design-Machines-Studio/livewires-templ/internal/testutil"
)

func TestDropdownDefault(t *testing.T) {
	html := testutil.RenderToString(t, Dropdown(""))
	if !strings.Contains(html, `class="dropdown"`) {
		t.Error("expected dropdown class")
	}
}

func TestDropdownEnd(t *testing.T) {
	html := testutil.RenderToString(t, Dropdown("end"))
	if !strings.Contains(html, "dropdown--end") {
		t.Error("expected dropdown--end class")
	}
}

func TestDropdownTrigger(t *testing.T) {
	html := testutil.RenderToString(t, DropdownTrigger("Menu", ""))
	if !strings.Contains(html, "button trigger") {
		t.Error("expected button trigger classes")
	}
	if !strings.Contains(html, "Menu") {
		t.Error("expected trigger label")
	}
}

func TestDropdownMenu(t *testing.T) {
	html := testutil.RenderToString(t, DropdownMenu())
	if !strings.Contains(html, `class="menu"`) {
		t.Error("expected menu class")
	}
}

func TestDropdownItemLink(t *testing.T) {
	html := testutil.RenderToString(t, DropdownItem(DropdownItemProps{Label: "Edit", Href: "/edit"}))
	if !strings.Contains(html, "<a") {
		t.Error("expected anchor tag")
	}
	if !strings.Contains(html, "Edit") {
		t.Error("expected item label")
	}
}

func TestDropdownItemButton(t *testing.T) {
	html := testutil.RenderToString(t, DropdownItem(DropdownItemProps{Label: "Delete"}))
	if !strings.Contains(html, "<button") {
		t.Error("expected button tag when no href")
	}
}

func TestDropdownDivider(t *testing.T) {
	html := testutil.RenderToString(t, DropdownDivider())
	if !strings.Contains(html, `role="separator"`) {
		t.Error("expected role=separator on divider li")
	}
	if !strings.Contains(html, "<hr") {
		t.Error("expected hr element")
	}
}

func TestDropdownTriggerARIA(t *testing.T) {
	html := testutil.RenderToString(t, DropdownTriggerComponent(DropdownTriggerProps{
		Label:    "Actions",
		MenuID:   "menu-1",
		Expanded: true,
	}))
	if !strings.Contains(html, `aria-haspopup="menu"`) {
		t.Error("expected aria-haspopup=menu")
	}
	if !strings.Contains(html, `aria-expanded="true"`) {
		t.Error("expected aria-expanded=true")
	}
	if !strings.Contains(html, `aria-controls="menu-1"`) {
		t.Error("expected aria-controls referencing menu ID")
	}
}

func TestDropdownMenuRole(t *testing.T) {
	html := testutil.RenderToString(t, DropdownMenuComponent(DropdownMenuProps{ID: "menu-1"}))
	if !strings.Contains(html, `role="menu"`) {
		t.Error("expected role=menu on dropdown menu")
	}
	if !strings.Contains(html, `id="menu-1"`) {
		t.Error("expected id on dropdown menu")
	}
}

func TestDropdownItemRoles(t *testing.T) {
	html := testutil.RenderToString(t, DropdownItem(DropdownItemProps{Label: "Edit", Href: "/edit"}))
	if !strings.Contains(html, `role="none"`) {
		t.Error("expected role=none on li")
	}
	if !strings.Contains(html, `role="menuitem"`) {
		t.Error("expected role=menuitem on interactive element")
	}
}

func TestDropdownItemDisabled(t *testing.T) {
	html := testutil.RenderToString(t, DropdownItem(DropdownItemProps{Label: "Delete", Href: "/delete", Disabled: true}))
	if !strings.Contains(html, "<button") {
		t.Error("expected button element when disabled with href")
	}
	if strings.Contains(html, "<a") {
		t.Error("should not render anchor when disabled")
	}
	if !strings.Contains(html, "disabled") {
		t.Error("expected disabled attribute")
	}
}
