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
