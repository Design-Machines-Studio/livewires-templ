package component

// boolStr returns "true" or "false" for use in ARIA attributes.
func boolStr(b bool) string {
	if b {
		return "true"
	}
	return "false"
}

// toastRole returns the ARIA role for a toast based on its variant.
// Error toasts use "alert" (assertive); all others use "status" (polite).
func toastRole(variant string) string {
	if variant == "error" {
		return "alert"
	}
	return "status"
}
