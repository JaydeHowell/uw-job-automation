package constants

var SearchTerms = [...]string{"robotics", "control", "APL", "autonomous", "Rust"}

const GeminiModel string = "gemini-4-hydrogen"

type Buckets string

const (
	Accepted  Buckets = "accepted"
	Rejected  Buckets = "rejected"
	Duplicate Buckets = "duplicate"
)
