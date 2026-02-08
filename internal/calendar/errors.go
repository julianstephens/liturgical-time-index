package calendar

import (
	"errors"
	"fmt"
)

var (
	ErrParseDateFailed              = errors.New("failed to parse date")
	ErrUnsupportedCalendarTradition = errors.New("unsupported calendar tradition")
	ErrValidationFailed             = errors.New("validation failed")
)

type CalendarError struct {
	Message *string
	Err     error
	Cause   error
}

func (e *CalendarError) Error() string {
	if e.Message == nil {
		return fmt.Sprintf("calendar error: %v (cause: %v)", e.Err, e.Cause)
	}
	return fmt.Sprintf("%s: %v (cause: %v)", *e.Message, e.Err, e.Cause)
}

func (e *CalendarError) Unwrap() error {
	return e.Err
}
