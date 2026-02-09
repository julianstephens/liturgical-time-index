package output

import "fmt"

var (
	ErrInvalidEntry          = fmt.Errorf("invalid entry")
	ErrInvalidOutputPath     = fmt.Errorf("invalid output path")
	ErrOpenOutputFileFailed  = fmt.Errorf("failed to open output file")
	ErrCloseOutputFileFailed = fmt.Errorf("failed to close output file")
	ErrSerializationFailed   = fmt.Errorf("serialization error")
)

type OutputError struct {
	Message *string
	Err     error
	Cause   error
}

func (e *OutputError) Error() string {
	if e.Message != nil {
		return fmt.Sprintf("output error: %s: %v (cause: %v)", *e.Message, e.Err, e.Cause)
	}
	return fmt.Sprintf("output error: %v (cause: %v)", e.Err, e.Cause)
}

func (e *OutputError) Unwrap() error {
	return e.Err
}
