package form

// Option represents a single option in a select dropdown or filter.
type Option struct {
	Value    string // Option value
	Label    string // Display label
	Selected bool   // Whether this option is selected
}

// FilterOption is an alias for Option (backward compatibility).
type FilterOption = Option

// SelectOption is an alias for Option (backward compatibility).
type SelectOption = Option
