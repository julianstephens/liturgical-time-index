package compile

import "fmt"

var (
	ErrSeasonNotFound = fmt.Errorf("season not found in plan")
)

type CompileError struct {
	Message *string
	Err     error
	Cause   error
}

func (e *CompileError) Error() string {
	if e.Message != nil {
		return fmt.Sprintf("compile error: %s, err: %v (cause: %v)", *e.Message, e.Err, e.Cause)
	}
	return fmt.Sprintf("compile error: %v (cause: %v)", e.Err, e.Cause)
}

func (e *CompileError) Unwrap() error {
	return e.Err
}
