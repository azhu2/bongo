package scorer

import "fmt"

type invalidLetterError struct {
	letter rune
}

func (e invalidLetterError) Error() string {
	return fmt.Sprintf("solution has invalid letter: %c", e.letter)
}

func (e invalidLetterError) Is(target error) bool {
	_, ok := target.(invalidLetterError)
	return ok
}
