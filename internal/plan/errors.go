package plan

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidRbRefFormat = errors.New("invalid RB reference format")
	ErrInvalidStartVerse  = errors.New("invalid start verse")
	ErrInvalidEndVerse    = errors.New("invalid end verse")
	ErrPlanFileRead       = errors.New("failed to read plan file")
	ErrParsePlanFailed    = errors.New("failed to parse plan file")
	ErrInvalidPlanEntry   = errors.New("invalid plan entry")
)

type PlanError struct {
	Message *string `json:"message,omitempty"`
	Err     error   `json:"error"`
	Cause   error   `json:"cause,omitempty"`
}

func (e *PlanError) Error() string {
	if e.Message != nil {
		return fmt.Sprintf("plan error: %s, cause: %v", *e.Message, e.Cause)
	}
	return fmt.Sprintf("plan error: %v, cause: %v", e.Err, e.Cause)
}

func (e *PlanError) Unwrap() error {
	return e.Err
}

var (
	ErrRbRefParseFailed      = errors.New("failed to parse RB reference")
	ErrRbRefValidationFailed = errors.New("RB reference validation failed")
)

type RbRefError struct {
	Message *string `json:"message,omitempty"`
	Err     error   `json:"error"`
	Cause   error   `json:"cause,omitempty"`
}

func (e *RbRefError) Error() string {
	if e.Message != nil {
		return fmt.Sprintf("RB reference error: %s, cause: %v", *e.Message, e.Cause)
	}
	return fmt.Sprintf("RB reference error: %v, cause: %v", e.Err, e.Cause)
}

func (e *RbRefError) Unwrap() error {
	return e.Err
}
