package e

import "fmt"

// Wrap - helper func to beauty wrap.
func Wrap(msg string, err error) error {
	return fmt.Errorf("%s: %w", msg, err)
}

// WrapIfErr - helper func to beauty wrap in the deferred functions.
func WrapIfErr(msg string, err error) error {
	if err == nil {
		return nil
	}

	return Wrap(msg, err)
}
