package semver

import (
	"fmt"
)

type ParseError struct {
	Position int
	Message  string
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("%s (at position %d)", e.Message, e.Position)
}
