package component

import "github.com/a-h/templ"

// dateFilterTriggerIcon returns icon, or the default calendar icon when nil.
func dateFilterTriggerIcon(icon templ.Component) templ.Component {
	if icon != nil {
		return icon
	}
	return dateFilterIcon()
}
