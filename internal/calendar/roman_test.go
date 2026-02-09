package calendar

import (
	"strings"
	"testing"
	"time"
)

func TestGetRomanWeekday(t *testing.T) {
	ce := NewCalendarEngine()

	testCases := []struct {
		date     string
		expected Weekday
	}{
		{"2025-03-09", Sunday},    // Sunday
		{"2025-03-10", Monday},    // Monday
		{"2025-03-11", Tuesday},   // Tuesday
		{"2025-03-12", Wednesday}, // Wednesday
		{"2025-03-13", Thursday},  // Thursday
		{"2025-03-14", Friday},    // Friday
		{"2025-03-15", Saturday},  // Saturday
	}

	for _, tc := range testCases {
		t.Run(tc.date, func(t *testing.T) {
			weekday, err := ce.GetRomanWeekday(tc.date)
			if err != nil {
				t.Fatalf("GetRomanWeekday failed: %v", err)
			}

			if weekday != tc.expected {
				t.Errorf("Expected weekday %s, got %s", tc.expected, weekday)
			}
		})
	}
}

func TestGetRomanWeekdayInvalidDate(t *testing.T) {
	ce := NewCalendarEngine()

	_, err := ce.GetRomanWeekday("invalid-date")
	if err == nil {
		t.Error("Expected error for invalid date")
	}

	if !strings.Contains(err.Error(), "failed to parse date") {
		t.Errorf("Expected parse error, got: %v", err)
	}
}

func TestGetRomanSeason(t *testing.T) {
	ce := NewCalendarEngine()

	testCases := []struct {
		date string
	}{
		{"2025-01-10"},
		{"2025-03-30"},
		{"2025-04-20"},
		{"2025-06-09"},
		{"2025-12-01"},
		{"2025-12-25"},
	}

	validSeasons := map[LiturgicalSeason]bool{
		Advent:        true,
		Christmastide: true,
		Epiphanytide:  true,
		Lent:          true,
		Triduum:       true,
		Easter:        true,
		Ordinary:      true,
	}

	for _, tc := range testCases {
		t.Run(tc.date, func(t *testing.T) {
			season, err := ce.GetRomanSeason(tc.date, RomanCalendar)
			if err != nil {
				t.Fatalf("GetRomanSeason failed: %v", err)
			}

			if season == "" {
				t.Errorf("Expected non-empty season for date %s", tc.date)
			}

			if !validSeasons[season] {
				t.Errorf("Got invalid season: %s", season)
			}
		})
	}
}

func TestGetRomanSeasonInvalidTradition(t *testing.T) {
	ce := NewCalendarEngine()

	_, err := ce.GetRomanSeason("2025-03-09", CalendarTradition("Invalid"))
	if err == nil {
		t.Error("Expected error for invalid tradition")
	}

	if !strings.Contains(err.Error(), "unsupported calendar tradition") {
		t.Errorf("Expected unsupported tradition error, got: %v", err)
	}
}

func TestGetRomanSeasonWeek(t *testing.T) {
	ce := NewCalendarEngine()

	testCases := []struct {
		date   string
		season LiturgicalSeason
	}{
		{"2025-01-06", Epiphanytide},
		{"2025-01-13", Epiphanytide},
		{"2025-03-05", Lent},
		{"2025-04-20", Easter},
		{"2025-06-09", Ordinary},
	}

	for _, tc := range testCases {
		t.Run(tc.date, func(t *testing.T) {
			week, err := ce.GetRomanSeasonWeek(tc.date, tc.season, RomanCalendar)
			if err != nil {
				t.Fatalf("GetRomanSeasonWeek failed: %v", err)
			}

			if week < 1 {
				t.Errorf("Expected positive week number, got %d for date %s", week, tc.date)
			}
		})
	}
}

func TestGetRomanSeasonWeekBeforeSeasonStart(t *testing.T) {
	ce := NewCalendarEngine()

	// Try to get season week for a date before the season starts
	_, err := ce.GetRomanSeasonWeek("2025-01-01", Lent, RomanCalendar)
	if err == nil {
		t.Error("Expected error for date before season start")
	}

	if !strings.Contains(err.Error(), "date is before the start of the season") {
		t.Errorf("Expected 'before season start' error, got: %v", err)
	}
}

func TestGetRomanDay(t *testing.T) {
	ce := NewCalendarEngine()

	testCases := []struct {
		date            string
		expectedSeason  LiturgicalSeason
		expectedWeekday Weekday
	}{
		{"2025-03-09", Lent, Sunday},
		{"2025-04-20", Easter, Sunday},
		{"2025-12-25", Christmastide, Thursday},
	}

	for _, tc := range testCases {
		t.Run(tc.date, func(t *testing.T) {
			dayKey, err := ce.GetRomanDay(tc.date, RomanCalendar)
			if err != nil {
				t.Fatalf("GetRomanDay failed: %v", err)
			}

			if dayKey.Date != tc.date {
				t.Errorf("Expected date %s, got %s", tc.date, dayKey.Date)
			}

			if dayKey.Season != tc.expectedSeason {
				t.Errorf("Expected season %s, got %s", tc.expectedSeason, dayKey.Season)
			}

			if dayKey.Weekday != tc.expectedWeekday {
				t.Errorf("Expected weekday %s, got %s", tc.expectedWeekday, dayKey.Weekday)
			}

			if dayKey.Tradition != RomanCalendar {
				t.Errorf("Expected tradition RomanCalendar, got %s", dayKey.Tradition)
			}

			if dayKey.SeasonWeek < 1 {
				t.Errorf("Expected positive season week, got %d", dayKey.SeasonWeek)
			}
		})
	}
}

func TestGetRomanDayInvalidTradition(t *testing.T) {
	ce := NewCalendarEngine()

	_, err := ce.GetRomanDay("2025-03-09", CalendarTradition("Invalid"))
	if err == nil {
		t.Error("Expected error for invalid tradition")
	}

	if !strings.Contains(err.Error(), "unsupported calendar tradition") {
		t.Errorf("Expected unsupported tradition error, got: %v", err)
	}
}

func TestGetRomanDayInvalidDate(t *testing.T) {
	ce := NewCalendarEngine()

	_, err := ce.GetRomanDay("invalid-date", RomanCalendar)
	if err == nil {
		t.Error("Expected error for invalid date")
	}
}

func TestGenerateRomanCalendar(t *testing.T) {
	ce := NewCalendarEngine()

	days, err := ce.GenerateRomanCalendar("2024", RomanCalendar)
	if err != nil {
		t.Fatalf("GenerateRomanCalendar failed: %v", err)
	}

	// 2024 is a leap year with 366 days
	if len(days) != 366 {
		t.Errorf("Expected 366 days for 2024 (leap year), got %d", len(days))
	}

	// Check that all days have required fields
	for i, day := range days {
		if day.Date == "" {
			t.Errorf("Day %d has empty date", i)
		}

		if day.Tradition != RomanCalendar {
			t.Errorf("Day %d has wrong tradition: %s", i, day.Tradition)
		}

		if day.Season == "" {
			t.Errorf("Day %d has empty season", i)
		}

		if day.SeasonWeek < 1 {
			t.Errorf("Day %d has invalid season week: %d", i, day.SeasonWeek)
		}

		if day.Weekday == "" {
			t.Errorf("Day %d has empty weekday", i)
		}
	}

	// Check that the first day is January 1
	if days[0].Date != "2024-01-01" {
		t.Errorf("Expected first day to be 2024-01-01, got %s", days[0].Date)
	}

	// Check that the last day is December 31
	if days[len(days)-1].Date != "2024-12-31" {
		t.Errorf("Expected last day to be 2024-12-31, got %s", days[len(days)-1].Date)
	}
}

func TestGenerateRomanCalendarLeapYear(t *testing.T) {
	ce := NewCalendarEngine()

	days, err := ce.GenerateRomanCalendar("2024", RomanCalendar)
	if err != nil {
		t.Fatalf("GenerateRomanCalendar failed: %v", err)
	}

	// 2024 is a leap year with 366 days
	if len(days) != 366 {
		t.Errorf("Expected 366 days for 2024 (leap year), got %d", len(days))
	}
}

func TestGenerateRomanCalendarInvalidTradition(t *testing.T) {
	ce := NewCalendarEngine()

	_, err := ce.GenerateRomanCalendar("2025", CalendarTradition("Invalid"))
	if err == nil {
		t.Error("Expected error for invalid tradition")
	}
}

func TestHolidays(t *testing.T) {
	ce := NewCalendarEngine()

	holidays, err := ce.Holidays(2024, RomanCalendar)
	if err != nil {
		t.Fatalf("Holidays failed: %v", err)
	}

	// Check that all expected holidays are present
	expectedHolidays := []string{
		"Ash Wednesday",
		"Holy Thursday",
		"Good Friday",
		"Easter Sunday",
		"Easter Monday",
		"Pentecost",
	}

	for _, holiday := range expectedHolidays {
		if _, ok := holidays[holiday]; !ok {
			t.Errorf("Missing expected holiday: %s", holiday)
		}
	}

	// Check that holidays have the required fields
	for name, day := range holidays {
		if day.Date == "" {
			t.Errorf("Holiday %s has empty date", name)
		}

		if day.Tradition != RomanCalendar {
			t.Errorf("Holiday %s has wrong tradition: %s", name, day.Tradition)
		}

		if day.Season == "" {
			t.Errorf("Holiday %s has empty season", name)
		}

		if day.Weekday == "" {
			t.Errorf("Holiday %s has empty weekday", name)
		}
	}
}

func TestGetRomanSeasonStartDate(t *testing.T) {
	ce := NewCalendarEngine()

	testCases := []struct {
		date   string
		season LiturgicalSeason
	}{
		{"2025-01-10", Epiphanytide},
		{"2025-03-05", Lent},
		{"2025-04-20", Easter},
	}

	for _, tc := range testCases {
		t.Run(tc.date+":"+string(tc.season), func(t *testing.T) {
			startDate, err := ce.getRomanSeasonStartDate(tc.date, tc.season, RomanCalendar)
			if err != nil {
				t.Fatalf("getRomanSeasonStartDate failed: %v", err)
			}

			if startDate.IsZero() {
				t.Errorf("Expected non-zero start date for season %s", tc.season)
			}

			// Verify that the start date is before or equal to the given date
			parsed, _ := time.Parse("2006-01-02", tc.date)
			// Truncate to day level to avoid timezone issues
			startDateTrunc := startDate.Truncate(24 * time.Hour)
			parsedTrunc := parsed.Truncate(24 * time.Hour)
			if startDateTrunc.After(parsedTrunc) {
				t.Errorf("Season start date %v should not be after %v", startDateTrunc, parsedTrunc)
			}
		})
	}
}

func TestGetRomanSeasonStartDateInvalidTradition(t *testing.T) {
	ce := NewCalendarEngine()

	_, err := ce.getRomanSeasonStartDate("2025-03-09", Lent, CalendarTradition("Invalid"))
	if err == nil {
		t.Error("Expected error for invalid tradition")
	}

	if !strings.Contains(err.Error(), "unsupported calendar tradition") {
		t.Errorf("Expected unsupported tradition error, got: %v", err)
	}
}

func TestDayKeyValidation(t *testing.T) {
	ce := NewCalendarEngine()

	// Get a valid DayKey
	dayKey, err := ce.GetRomanDay("2025-03-09", RomanCalendar)
	if err != nil {
		t.Fatalf("GetRomanDay failed: %v", err)
	}

	// Verify all required fields are populated
	if dayKey.Date == "" {
		t.Error("DayKey Date is empty")
	}
	if dayKey.Season == "" {
		t.Error("DayKey Season is empty")
	}
	if dayKey.Weekday == "" {
		t.Error("DayKey Weekday is empty")
	}
	if dayKey.Tradition == "" {
		t.Error("DayKey Tradition is empty")
	}
	if dayKey.SeasonWeek < 1 {
		t.Error("DayKey SeasonWeek is less than 1")
	}
}

func TestAdventBoundary(t *testing.T) {
	ce := NewCalendarEngine()

	// Test the boundary for Advent (Sunday after Nov 27)
	// In 2025, Nov 27 is a Thursday, so the Sunday after is Nov 30
	dayBefore, _ := ce.GetRomanDay("2025-11-29", RomanCalendar)
	dayOn, _ := ce.GetRomanDay("2025-11-30", RomanCalendar)

	if dayBefore.Season != Ordinary {
		t.Errorf("Expected day before Advent to be Ordinary Time, got %s", dayBefore.Season)
	}

	if dayOn.Season != Advent {
		t.Errorf("Expected day on Advent start to be Advent, got %s", dayOn.Season)
	}
}

func TestChristmasBoundary(t *testing.T) {
	ce := NewCalendarEngine()

	// Test Christmas boundaries
	dayBefore, _ := ce.GetRomanDay("2025-12-24", RomanCalendar)
	dayOn, _ := ce.GetRomanDay("2025-12-25", RomanCalendar)
	dayAfter, _ := ce.GetRomanDay("2025-01-06", RomanCalendar)

	if dayBefore.Season != Advent {
		t.Errorf("Expected Christmas Eve to be Advent, got %s", dayBefore.Season)
	}

	if dayOn.Season != Christmastide {
		t.Errorf("Expected Christmas Day to be Christmastide, got %s", dayOn.Season)
	}

	if dayAfter.Season != Epiphanytide {
		t.Errorf("Expected Jan 6 to be Epiphanytide, got %s", dayAfter.Season)
	}
}

func TestEpiphanyBoundary(t *testing.T) {
	ce := NewCalendarEngine()

	// Test Epiphany boundaries
	dayAfterChristmas, _ := ce.GetRomanDay("2025-01-05", RomanCalendar)
	dayOfEpiphany, _ := ce.GetRomanDay("2025-01-06", RomanCalendar)

	if dayAfterChristmas.Season != Christmastide {
		t.Errorf("Expected Jan 5 to be Christmastide, got %s", dayAfterChristmas.Season)
	}

	if dayOfEpiphany.Season != Epiphanytide {
		t.Errorf("Expected Jan 6 to be Epiphanytide, got %s", dayOfEpiphany.Season)
	}
}

func TestEasterSeasonBoundary(t *testing.T) {
	ce := NewCalendarEngine()

	easterDate := "2025-04-20"

	easterDay, err := ce.GetRomanDay(easterDate, RomanCalendar)
	if err != nil || easterDay == nil {
		t.Fatalf("Failed to get Easter day: %v", err)
	}

	afterPentecostDay, err := ce.GetRomanDay("2025-06-09", RomanCalendar)
	if err != nil || afterPentecostDay == nil {
		t.Fatalf("Failed to get day after Pentecost: %v", err)
	}

	if easterDay.Season != Easter {
		t.Errorf("Expected Easter Sunday to be Easter season, got %s", easterDay.Season)
	}

	if afterPentecostDay.Season != Ordinary {
		t.Errorf("Expected day after Pentecost to be Ordinary Time, got %s", afterPentecostDay.Season)
	}
}
