package constants

var SearchTerms = [...]string{"robotics", "control", "APL", "autonomous", "Rust"}

type Buckets string

const (
	Accepted  Buckets = "accepted"
	Rejected  Buckets = "rejected"
	Duplicate Buckets = "duplicate"
)
