package helpers

import (
	"errors"
	"io"
)

var IgnoreEOF = IgnoreErrorType(io.EOF)

func IgnoreErrorType(target error) func(error) error {
	return func(err error) error {
		if errors.Is(err, target) {
			return nil
		}
		return err
	}
}
