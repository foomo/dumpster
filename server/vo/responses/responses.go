package responses

import "time"

// Dump info
type Dump struct {
	ID      string `json:"id"`
	Created time.Time
	Report  string
	Errors  string
	Comment string
	Path    string
}

// RestoreReport a report of a dump restoration
type RestoreReport struct {
	Report string
	Errors string
}
