package livewires

import "strings"

// ClassNames joins non-empty class strings with spaces.
func ClassNames(classes ...string) string {
	var result []string
	for _, c := range classes {
		if strings.TrimSpace(c) != "" {
			result = append(result, strings.TrimSpace(c))
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

// SchemeClass returns "scheme-name" if name is non-empty.
func SchemeClass(name string) string {
	if name == "" {
		return ""
	}
	return "scheme-" + name
}
