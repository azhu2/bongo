package scorer

import "fmt"

type InvalidLetterError struct {
	letter rune
}

func (e InvalidLetterError) Error() string {
	return fmt.Sprintf("solution has invalid letter: %c", e.letter)
}

func (e InvalidLetterError) Is(target error) bool {
	_, ok := target.(InvalidLetterError)
	return ok
}
