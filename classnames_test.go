package livewires

import "testing"

func TestClassNames(t *testing.T) {
	tests := []struct {
		name   string
		input  []string
		expect string
	}{
		{"empty", nil, ""},
		{"single", []string{"button"}, "button"},
		{"multiple", []string{"button", "button--accent"}, "button button--accent"},
		{"skips empty", []string{"button", "", "active"}, "button active"},
		{"trims whitespace", []string{" button ", "  ", "active"}, "button active"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ClassNames(tt.input...)
			if got != tt.expect {
				t.Errorf("ClassNames(%v) = %q, want %q", tt.input, got, tt.expect)
			}
		})
	}
}

func TestVariantClass(t *testing.T) {
	if got := VariantClass("button", "accent"); got != "button--accent" {
		t.Errorf("got %q, want button--accent", got)
	}
	if got := VariantClass("button", ""); got != "" {
		t.Errorf("got %q, want empty", got)
	}
}

func TestModifierClass(t *testing.T) {
	if got := ModifierClass("stack", "compact"); got != "stack-compact" {
		t.Errorf("got %q, want stack-compact", got)
	}
	if got := ModifierClass("stack", ""); got != "" {
		t.Errorf("got %q, want empty", got)
	}
}

func TestSchemeClass(t *testing.T) {
	if got := SchemeClass("subtle"); got != "scheme-subtle" {
		t.Errorf("got %q, want scheme-subtle", got)
	}
	if got := SchemeClass(""); got != "" {
		t.Errorf("got %q, want empty", got)
	}
}

func TestAspectRatioClass(t *testing.T) {
	if got := AspectRatioClass("4-3"); got != "aspect-4-3" {
		t.Errorf("got %q, want aspect-4-3", got)
	}
	if got := AspectRatioClass(""); got != "" {
		t.Errorf("got %q, want empty", got)
	}
}
