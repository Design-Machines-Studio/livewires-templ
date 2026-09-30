package form

// RequiredIndicatorMarkup is trusted, static label content for native required
// controls and group legends. Its visual text is hidden from assistive technology
// because the control already exposes its required state. Optional fields emit
// nothing. The leading space permits ordinary inline label wrapping.
func RequiredIndicatorMarkup(required bool) string {
	if !required {
		return ""
	}
	return ` <span class="text-xs text-muted" aria-hidden="true">Required</span>`
}
