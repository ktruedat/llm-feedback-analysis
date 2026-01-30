package analysis

// Status represents the status of an analysis.
type Status string

const (
	StatusProcessing Status = "processing"
	StatusSuccess    Status = "success"
	StatusFailed     Status = "failed"
)

// String returns the string representation of the status.
func (s Status) String() string {
	return string(s)
}

// IsValid checks if the status is valid.
func (s Status) IsValid() bool {
	switch s {
	case StatusProcessing, StatusSuccess, StatusFailed:
		return true
	default:
		return false
	}
}
