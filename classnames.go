package livewires

import "strings"

// ClassNames joins non-empty class strings with spaces.
func ClassNames(classes ...string) string {
	result := make([]string, 0, len(classes))
	for _, c := range classes {
		trimmed := strings.TrimSpace(c)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return strings.Join(result, " ")
}

// VariantClass returns "base--variant" if variant is non-empty.
// Used for component modifiers (double-dash convention).
func VariantClass(base, variant string) string {
	if variant == "" {
		return ""
	}
	return base + "--" + variant
}

// ModifierClass returns "base-modifier" if modifier is non-empty.
// Used for layout modifiers (single-dash convention).
func ModifierClass(base, modifier string) string {
	if modifier == "" {
		return ""
	}
	return base + "-" + modifier
}

// SizeClass returns "text-size" if size is non-empty.
func SizeClass(size string) string {
	if size == "" {
		return ""
	}
	return "text-" + size
}

// DefaultStr returns val if non-empty, otherwise def.
func DefaultStr(val, def string) string {
	if val != "" {
		return val
	}
	return def
}

// AlignClass returns "text-align" if align is non-empty.
func AlignClass(align string) string {
	if align == "" {
		return ""
	}
	return "text-" + align
}

// SchemeClass returns "scheme-name" if name is non-empty.
func SchemeClass(name string) string {
	if name == "" {
		return ""
	}
	return "scheme-" + name
}

// AspectRatioClass returns "aspect-{ratio}" if ratio is non-empty (e.g. "aspect-4-3").
func AspectRatioClass(ratio string) string {
	if ratio == "" {
		return ""
	}
	return "aspect-" + ratio
}
