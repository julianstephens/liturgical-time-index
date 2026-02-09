package calendar_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/julianstephens/liturgical-time-index/internal/calendar"
)

func TestNewCalendarEngine(t *testing.T) {
	ce := calendar.NewCalendarEngine()

	if ce == nil {
		t.Error("NewCalendarEngine returned nil")
	}
}

func TestGetEasterGregorian(t *testing.T) {
	ce := calendar.NewCalendarEngine()

	testYears := []int{2025, 2024, 2023, 2022, 2021, 1961}

	for _, year := range testYears {
		t.Run(fmt.Sprintf("Year%d", year), func(t *testing.T) {
			easter := ce.GetEasterGregorian(year)

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
	ce := calendar.NewCalendarEngine()

	// Test a date in 2025
	dayKey, err := ce.GetRomanDay("2025-03-09", calendar.RomanCalendar)
	if err != nil {
		t.Fatalf("GetRomanDay failed: %v", err)
	}

	if dayKey == nil {
		t.Fatal("GetRomanDay returned nil")
	}

	if dayKey.Date != "2025-03-09" {
		t.Errorf("Expected date 2025-03-09, got %s", dayKey.Date)
	}

	if dayKey.Tradition != calendar.RomanCalendar {
		t.Errorf("Expected tradition RomanCalendar, got %s", dayKey.Tradition)
	}

	if dayKey.Weekday != calendar.Sunday {
		t.Errorf("Expected weekday Sunday, got %s", dayKey.Weekday)
	}
}
