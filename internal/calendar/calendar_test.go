package calendar

import (
	"fmt"
	"testing"
	"time"

	"github.com/go-playground/validator/v10"

	"github.com/julianstephens/liturgical-time-index/internal/util"
)

func getConfiguredValidator() *validator.Validate {
	validate := validator.New()
	if err := validate.RegisterValidation("ISO8601", util.ValidateISO8601); err != nil {
		panic(fmt.Sprintf("Failed to register ISO8601 validation: %v", err))
	}
	return validate
}

func TestNewCalendarEngine(t *testing.T) {
	validate := getConfiguredValidator()
	ce := NewCalendarEngine(validate)

	if ce == nil {
		t.Error("NewCalendarEngine returned nil")
	}

	if ce.validate != validate {
		t.Error("CalendarEngine validator not set correctly")
	}
}

func TestGetEasterGregorian(t *testing.T) {
	validate := getConfiguredValidator()
	ce := NewCalendarEngine(validate)

	testYears := []int{2025, 2024, 2023, 2022, 2021, 1961}

	for _, year := range testYears {
		t.Run(fmt.Sprintf("Year%d", year), func(t *testing.T) {
			easter := ce.getEasterGregorian(year)

			// Verify that Easter is in April or late March
			month := easter.Month()
			if month != time.March && month != time.April {
				t.Errorf("Expected Easter in March or April, got %s", easter.Month())
			}

			// Verify that Easter year matches input year
			if easter.Year() != year {
				t.Errorf("Expected Easter year %d, got %d", year, easter.Year())
			}

			// Verify it's a valid date
			if easter.IsZero() {
				t.Error("Easter date should not be zero")
			}
		})
	}
}

func TestCalendarEngineIntegration(t *testing.T) {
	validate := getConfiguredValidator()
	ce := NewCalendarEngine(validate)

	// Test a date in 2025
	dayKey, err := ce.GetRomanDay("2025-03-09", RomanCalendar)
	if err != nil {
		t.Fatalf("GetRomanDay failed: %v", err)
	}

	if dayKey == nil {
		t.Fatal("GetRomanDay returned nil")
	}

	if dayKey.Date != "2025-03-09" {
		t.Errorf("Expected date 2025-03-09, got %s", dayKey.Date)
	}

	if dayKey.Tradition != RomanCalendar {
		t.Errorf("Expected tradition RomanCalendar, got %s", dayKey.Tradition)
	}

	if dayKey.Weekday != Sunday {
		t.Errorf("Expected weekday Sunday, got %s", dayKey.Weekday)
	}
}
