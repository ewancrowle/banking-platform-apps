package paymentdecision

type DeclineReason int

const (
	// For `iota` to work correctly, the order of these constants
	// must correspond to the order in which they appear in `payment_decision.proto`.
	DeclineReasonUnspecified DeclineReason = iota
	DeclineReasonInsufficientFunds
	DeclineReasonSuspiciousActivity
)

// Returns the string representation of a DeclineReason.
func (r DeclineReason) String() string {
	switch r {
	case DeclineReasonUnspecified:
		return "unspecified"
	case DeclineReasonInsufficientFunds:
		return "insufficient_funds"
	case DeclineReasonSuspiciousActivity:
		return "suspicious_activity"
	default:
		return "unknown"
	}
}

// FromString returns the DeclineReason corresponding to the given string.
func FromString(s string) DeclineReason {
	switch s {
	case "unspecified":
		return DeclineReasonUnspecified
	case "insufficient_funds":
		return DeclineReasonInsufficientFunds
	case "suspicious_activity":
		return DeclineReasonSuspiciousActivity
	default:
		return DeclineReasonUnspecified
	}
}
