package calendar

import (
	"errors"
	"testing"
)

func TestCalendarErrorError(t *testing.T) {
	testCases := []struct {
		name           string
		err            error
		cause          error
		message        *string
		expectedOutput string
	}{
		{
			name:           "error without message",
			err:            ErrParseDateFailed,
			cause:          errors.New("invalid format"),
			message:        nil,
			expectedOutput: "calendar error",
		},
		{
			name:           "error with message",
			err:            ErrParseDateFailed,
			cause:          errors.New("invalid format"),
			message:        stringPtr("custom message"),
			expectedOutput: "custom message",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			calendarErr := &CalendarError{
				Message: tc.message,
				Err:     tc.err,
				Cause:   tc.cause,
			}

			errMsg := calendarErr.Error()
			if errMsg == "" {
				t.Error("Error() returned empty string")
			}

			if tc.expectedOutput != "" && !containsString(errMsg, tc.expectedOutput) {
				t.Errorf("Error() = %q, expected to contain %q", errMsg, tc.expectedOutput)
			}
		})
	}
}

func TestCalendarErrorUnwrap(t *testing.T) {
	underlyingErr := ErrParseDateFailed
	calendarErr := &CalendarError{
		Err:   underlyingErr,
		Cause: errors.New("some cause"),
	}

	unwrapped := calendarErr.Unwrap()
	if unwrapped != underlyingErr {
		t.Errorf("Unwrap() = %v, expected %v", unwrapped, underlyingErr)
	}
}

func TestCalendarErrorIs(t *testing.T) {
	calendarErr := &CalendarError{
		Err:   ErrParseDateFailed,
		Cause: errors.New("some cause"),
	}

	// Test with errors.Is
	if !errors.Is(calendarErr, ErrParseDateFailed) {
		t.Error("errors.Is should return true for wrapped error")
	}

	// Test with different error
	if errors.Is(calendarErr, ErrUnsupportedCalendarTradition) {
		t.Error("errors.Is should return false for different error")
	}
}

func TestErrorConstants(t *testing.T) {
	testCases := []struct {
		name error
		want string
	}{
		{ErrParseDateFailed, "failed to parse date"},
		{ErrUnsupportedCalendarTradition, "unsupported calendar tradition"},
		{ErrValidationFailed, "validation failed"},
	}

	for _, tc := range testCases {
		t.Run(tc.want, func(t *testing.T) {
			if tc.name.Error() != tc.want {
				t.Errorf("Error() = %q, expected %q", tc.name.Error(), tc.want)
			}
		})
	}
}

func TestCalendarErrorWithNilCause(t *testing.T) {
	calendarErr := &CalendarError{
		Err:   ErrValidationFailed,
		Cause: nil,
	}

	errMsg := calendarErr.Error()
	if errMsg == "" {
		t.Error("Error() returned empty string for nil cause")
	}
}

// Helper functions
func stringPtr(s string) *string {
	return &s
}

func containsString(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
