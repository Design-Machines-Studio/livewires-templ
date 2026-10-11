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

func toastKind(props ToastProps) string {
	if props.Kind != "" {
		return props.Kind
	}
	if props.Variant != "" {
		return props.Variant
	}
	return "info"
}

func toastDismissLabel(label string) string {
	if label == "" {
		return "Dismiss notification"
	}
	return label
}

func copyButtonLabel(label string) string {
	if label == "" {
		return "Copy"
	}
	return label
}
