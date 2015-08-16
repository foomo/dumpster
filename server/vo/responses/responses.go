package responses

import "time"

// Dump info
type Dump struct {
	ID      string    `json:"id"`
	Created time.Time `json:"created"`
	Report  string    `json:"report"`
	Errors  string    `json:"errors"`
	Comment string    `json:"comment"`
	Path    string    `json:"path"`
}

// RestoreReport a report of a dump restoration
type RestoreReport struct {
	Report string `json:"report"`
	Errors string `json:"errors"`
}
