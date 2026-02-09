package compile_test

import (
	"fmt"
	"testing"

	"github.com/julianstephens/liturgical-time-index/internal/calendar"
	"github.com/julianstephens/liturgical-time-index/internal/compile"
	"github.com/julianstephens/liturgical-time-index/internal/plan"
)

// TestBoundaryCondition_AshWednesdayToHolySaturday tests the Lent/Triduum transition boundary.
// This window covers the critical transition where the same weekdays (Wed-Sat) exist in two seasons.
func TestBoundaryCondition_AshWednesdayToHolySaturday(t *testing.T) {
	// For 2025: Ash Wednesday is 2025-03-05 (Wed), Holy Saturday is 2025-04-19 (Sat)
	testCases := []struct {
		date               string
		expectedSeason     calendar.LiturgicalSeason
		expectedWeekday    calendar.Weekday
		expectedCue        string
		expectedRbContains string
		description        string
	}{
		{
			date:               "2025-03-05",
			expectedSeason:     calendar.Lent,
			expectedWeekday:    calendar.Wednesday,
			expectedCue:        "Wednesday in Lent",
			expectedRbContains: "RB 8.3",
			description:        "Ash Wednesday - first day of Lent",
		},
		{
			date:               "2025-03-06",
			expectedSeason:     calendar.Lent,
			expectedWeekday:    calendar.Thursday,
			expectedCue:        "Thursday in Lent",
			expectedRbContains: "RB 8.4",
			description:        "Thursday in Week 1 of Lent",
		},
		{
			date:               "2025-03-09",
			expectedSeason:     calendar.Lent,
			expectedWeekday:    calendar.Sunday,
			expectedCue:        "First Sunday of Lent",
			expectedRbContains: "RB 9.1",
			description:        "First Sunday of Lent",
		},
		{
			date:               "2025-04-16",
			expectedSeason:     calendar.Lent,
			expectedWeekday:    calendar.Wednesday,
			expectedCue:        "Wednesday in Lent",
			expectedRbContains: "RB 8.3",
			description:        "Last Wednesday of Lent",
		},
		{
			date:               "2025-04-17",
			expectedSeason:     calendar.Triduum,
			expectedWeekday:    calendar.Thursday,
			expectedCue:        "Holy Thursday",
			expectedRbContains: "RB 10.4",
			description:        "Holy Thursday - transition from Lent to Triduum",
		},
		{
			date:               "2025-04-18",
			expectedSeason:     calendar.Triduum,
			expectedWeekday:    calendar.Friday,
			expectedCue:        "Good Friday",
			expectedRbContains: "RB 10.5",
			description:        "Good Friday",
		},
		{
			date:               "2025-04-19",
			expectedSeason:     calendar.Triduum,
			expectedWeekday:    calendar.Saturday,
			expectedCue:        "Holy Saturday",
			expectedRbContains: "RB 10.6",
			description:        "Holy Saturday - last day of Triduum",
		},
	}

	testPlan := createLentTriduumTransitionPlan()

	ce := calendar.NewCalendarEngine()
	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			dayKey, err := ce.GetRomanDay(tc.date, calendar.RomanCalendar)
			if err != nil {
				t.Fatalf("Failed to get Roman day for %s: %v", tc.date, err)
			}

			// Assert season is correct
			if dayKey.Season != tc.expectedSeason {
				t.Errorf("Expected season %s, got %s", tc.expectedSeason, dayKey.Season)
			}

			// Assert weekday is correct
			if dayKey.Weekday != tc.expectedWeekday {
				t.Errorf("Expected weekday %s, got %s", tc.expectedWeekday, dayKey.Weekday)
			}

			// Compile entry
			entry, err := compile.Compile(*dayKey, testPlan)
			if err != nil {
				t.Fatalf("Compile failed: %v", err)
			}

			// Assert output is stable and correct
			if entry.Cue != tc.expectedCue {
				t.Errorf("Expected cue %q, got %q", tc.expectedCue, entry.Cue)
			}

			// Assert RB refs contain expected value
			if len(entry.Rb) == 0 {
				t.Fatal("Expected RB refs, got none")
			}
			if entry.Rb[0].String() != tc.expectedRbContains {
				t.Errorf("Expected RB ref containing %q, got %q", tc.expectedRbContains, entry.Rb[0].String())
			}
		})
	}
}

// TestBoundaryCondition_EasterToPentecost tests the Easter season, validating consistent
// compilation through the 7-week window from Easter Sunday to Pentecost.
func TestBoundaryCondition_EasterToPentecost(t *testing.T) {
	// For 2025: Easter is 2025-04-20 (Sun), Pentecost is 2025-06-08 (Sun)
	testCases := []struct {
		date               string
		expectedSeason     calendar.LiturgicalSeason
		expectedWeekday    calendar.Weekday
		expectedCue        string
		expectedRbContains string
		description        string
	}{
		{
			date:               "2025-04-20",
			expectedSeason:     calendar.Eastertide,
			expectedWeekday:    calendar.Sunday,
			expectedCue:        "Pentecost Sunday",
			expectedRbContains: "RB 13.7",
			description:        "Easter Sunday - start of Easter season",
		},
		{
			date:               "2025-04-21",
			expectedSeason:     calendar.Eastertide,
			expectedWeekday:    calendar.Monday,
			expectedCue:        "Easter Monday",
			expectedRbContains: "RB 12.1",
			description:        "Easter Monday - week 1 of Easter",
		},
		{
			date:               "2025-04-27",
			expectedSeason:     calendar.Eastertide,
			expectedWeekday:    calendar.Sunday,
			expectedCue:        "Pentecost Sunday",
			expectedRbContains: "RB 13.7",
			description:        "Second Sunday of Easter",
		},
		{
			date:               "2025-05-04",
			expectedSeason:     calendar.Eastertide,
			expectedWeekday:    calendar.Sunday,
			expectedCue:        "Pentecost Sunday",
			expectedRbContains: "RB 13.7",
			description:        "Third Sunday of Easter",
		},
		{
			date:               "2025-05-25",
			expectedSeason:     calendar.Eastertide,
			expectedWeekday:    calendar.Sunday,
			expectedCue:        "Pentecost Sunday",
			expectedRbContains: "RB 13.7",
			description:        "Sixth Sunday of Easter - day before Ascension",
		},
		{
			date:               "2025-05-29",
			expectedSeason:     calendar.Eastertide,
			expectedWeekday:    calendar.Thursday,
			expectedCue:        "Ascension Thursday",
			expectedRbContains: "RB 12.4",
			description:        "Ascension Thursday (39 days after Easter)",
		},
		{
			date:               "2025-05-20",
			expectedSeason:     calendar.Eastertide,
			expectedWeekday:    calendar.Tuesday,
			expectedCue:        "Easter Tuesday",
			expectedRbContains: "RB 12.2",
			description:        "Mid-Easter season weekday",
		},
	}

	testPlan := createEasterSeasonPlan()

	ce := calendar.NewCalendarEngine()
	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			dayKey, err := ce.GetRomanDay(tc.date, calendar.RomanCalendar)
			if err != nil {
				t.Fatalf("Failed to get Roman day for %s: %v", tc.date, err)
			}

			if dayKey.Season != tc.expectedSeason {
				t.Errorf("Expected season %s, got %s", tc.expectedSeason, dayKey.Season)
			}

			if dayKey.Weekday != tc.expectedWeekday {
				t.Errorf("Expected weekday %s, got %s", tc.expectedWeekday, dayKey.Weekday)
			}

			entry, err := compile.Compile(*dayKey, testPlan)
			if err != nil {
				t.Fatalf("Compile failed: %v", err)
			}

			if entry.Cue != tc.expectedCue {
				t.Errorf("Expected cue %q, got %q", tc.expectedCue, entry.Cue)
			}

			if len(entry.Rb) == 0 {
				t.Fatal("Expected RB refs, got none")
			}
			if entry.Rb[0].String() != tc.expectedRbContains {
				t.Errorf("Expected RB ref containing %q, got %q", tc.expectedRbContains, entry.Rb[0].String())
			}
		})
	}
}

// TestBoundaryCondition_LateAdventToEpiphany tests the Advent/Christmas/Epiphany transition.
// This window spans three seasons and includes important fixed feasts.
func TestBoundaryCondition_LateAdventToEpiphany(t *testing.T) {
	// For 2024-2025: Advent starts Dec 1 2024, Christmas Dec 25, Epiphany Jan 6 2025
	testCases := []struct {
		date               string
		expectedSeason     calendar.LiturgicalSeason
		expectedWeekday    calendar.Weekday
		expectedCue        string
		expectedRbContains string
		description        string
	}{
		{
			date:               "2024-12-22",
			expectedSeason:     calendar.Advent,
			expectedWeekday:    calendar.Sunday,
			expectedCue:        "Fourth Sunday of Advent",
			expectedRbContains: "RB 3.4",
			description:        "Fourth Sunday of Advent",
		},
		{
			date:               "2024-12-23",
			expectedSeason:     calendar.Advent,
			expectedWeekday:    calendar.Monday,
			expectedCue:        "Monday before Christmas",
			expectedRbContains: "RB 2.1",
			description:        "Monday in final week of Advent",
		},
		{
			date:               "2024-12-24",
			expectedSeason:     calendar.Advent,
			expectedWeekday:    calendar.Tuesday,
			expectedCue:        "Tuesday before Christmas",
			expectedRbContains: "RB 2.2",
			description:        "Christmas Eve - last day of Advent",
		},
		{
			date:               "2024-12-25",
			expectedSeason:     calendar.Christmastide,
			expectedWeekday:    calendar.Wednesday,
			expectedCue:        "Christmas Day",
			expectedRbContains: "RB 4.3",
			description:        "Christmas Day - start of Christmastide",
		},
		{
			date:               "2024-12-26",
			expectedSeason:     calendar.Christmastide,
			expectedWeekday:    calendar.Thursday,
			expectedCue:        "St. Stephen",
			expectedRbContains: "RB 4.4",
			description:        "St. Stephen (day after Christmas)",
		},
		{
			date:               "2024-12-29",
			expectedSeason:     calendar.Christmastide,
			expectedWeekday:    calendar.Sunday,
			expectedCue:        "Holy Family Sunday",
			expectedRbContains: "RB 5.1",
			description:        "Holy Family Sunday after Christmas",
		},
		{
			date:               "2025-01-01",
			expectedSeason:     calendar.Christmastide,
			expectedWeekday:    calendar.Wednesday,
			expectedCue:        "Christmas Day",
			expectedRbContains: "RB 4.3",
			description:        "Solemnity of Mary, Mother of God (Wednesday plan)",
		},
		{
			date:               "2025-01-05",
			expectedSeason:     calendar.Christmastide,
			expectedWeekday:    calendar.Sunday,
			expectedCue:        "Holy Family Sunday",
			expectedRbContains: "RB 5.1",
			description:        "Last Sunday of Christmastide",
		},
		{
			date:               "2025-01-06",
			expectedSeason:     calendar.Epiphanytide,
			expectedWeekday:    calendar.Monday,
			expectedCue:        "Epiphany",
			expectedRbContains: "RB 6.1",
			description:        "Epiphany - transition from Christmastide",
		},
		{
			date:               "2025-01-07",
			expectedSeason:     calendar.Epiphanytide,
			expectedWeekday:    calendar.Tuesday,
			expectedCue:        "Tuesday after Epiphany",
			expectedRbContains: "RB 6.2",
			description:        "Tuesday after Epiphany",
		},
		{
			date:               "2025-01-12",
			expectedSeason:     calendar.Epiphanytide,
			expectedWeekday:    calendar.Sunday,
			expectedCue:        "First Sunday after Epiphany",
			expectedRbContains: "RB 7.1",
			description:        "First Sunday after Epiphany",
		},
	}

	testPlan := createAdventChristmasEpiphanyPlan()

	ce := calendar.NewCalendarEngine()
	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			dayKey, err := ce.GetRomanDay(tc.date, calendar.RomanCalendar)
			if err != nil {
				t.Fatalf("Failed to get Roman day for %s: %v", tc.date, err)
			}

			if dayKey.Season != tc.expectedSeason {
				t.Errorf("Expected season %s, got %s", tc.expectedSeason, dayKey.Season)
			}

			if dayKey.Weekday != tc.expectedWeekday {
				t.Errorf("Expected weekday %s, got %s", tc.expectedWeekday, dayKey.Weekday)
			}

			entry, err := compile.Compile(*dayKey, testPlan)
			if err != nil {
				t.Fatalf("Compile failed: %v", err)
			}

			if entry.Cue != tc.expectedCue {
				t.Errorf("Expected cue %q, got %q", tc.expectedCue, entry.Cue)
			}

			if len(entry.Rb) == 0 {
				t.Fatal("Expected RB refs, got none")
			}
			if entry.Rb[0].String() != tc.expectedRbContains {
				t.Errorf("Expected RB ref containing %q, got %q", tc.expectedRbContains, entry.Rb[0].String())
			}
		})
	}
}

// TestMatchingPrecedence_WeekdayOverride verifies that season + weekday overrides take precedence.
func TestMatchingPrecedence_WeekdayOverride(t *testing.T) {
	testPlan := createWeekdayOverridePlan()

	testCases := []struct {
		season      calendar.LiturgicalSeason
		weekday     calendar.Weekday
		expectedCue string
		expectedRb  string
		description string
	}{
		{
			season:      calendar.Advent,
			weekday:     calendar.Monday,
			expectedCue: "Advent Monday Override",
			expectedRb:  "RB 2.1",
			description: "Monday in Advent uses weekday override",
		},
		{
			season:      calendar.Advent,
			weekday:     calendar.Tuesday,
			expectedCue: "Advent Tuesday Override",
			expectedRb:  "RB 2.2",
			description: "Tuesday in Advent uses weekday override",
		},
		{
			season:      calendar.Advent,
			weekday:     calendar.Wednesday,
			expectedCue: "Advent Fallback",
			expectedRb:  "RB 2.3",
			description: "Wednesday in Advent uses fallback (no override)",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			dayKey := calendar.DayKey{
				Date:       "2025-12-01",
				Tradition:  calendar.RomanCalendar,
				Season:     tc.season,
				SeasonWeek: 1,
				Weekday:    tc.weekday,
			}

			entry, err := compile.Compile(dayKey, testPlan)
			if err != nil {
				t.Fatalf("Compile failed: %v", err)
			}

			if entry.Cue != tc.expectedCue {
				t.Errorf("Expected cue %q, got %q", tc.expectedCue, entry.Cue)
			}

			if len(entry.Rb) == 0 || entry.Rb[0].String() != tc.expectedRb {
				t.Errorf("Expected RB %q, got %q", tc.expectedRb, func() string {
					if len(entry.Rb) == 0 {
						return "none"
					}
					return entry.Rb[0].String()
				}())
			}
		})
	}
}

// TestMatchingPrecedence_SeasonFallback verifies that season fallback is used when weekday override is missing.
func TestMatchingPrecedence_SeasonFallback(t *testing.T) {
	testPlan := createFallbackOnlyPlan()

	testCases := []struct {
		season      calendar.LiturgicalSeason
		weekday     calendar.Weekday
		expectedCue string
		expectedRb  string
		description string
	}{
		{
			season:      calendar.Advent,
			weekday:     calendar.Monday,
			expectedCue: "Advent Fallback",
			expectedRb:  "RB 2.1",
			description: "Any weekday in Advent uses fallback",
		},
		{
			season:      calendar.Advent,
			weekday:     calendar.Friday,
			expectedCue: "Advent Fallback",
			expectedRb:  "RB 2.1",
			description: "Different weekday in Advent also uses fallback",
		},
		{
			season:      calendar.Lent,
			weekday:     calendar.Wednesday,
			expectedCue: "Lent Fallback",
			expectedRb:  "RB 8.1",
			description: "Lent fallback distinct from Advent fallback",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			dayKey := calendar.DayKey{
				Date:       "2025-03-05",
				Tradition:  calendar.RomanCalendar,
				Season:     tc.season,
				SeasonWeek: 1,
				Weekday:    tc.weekday,
			}

			entry, err := compile.Compile(dayKey, testPlan)
			if err != nil {
				t.Fatalf("Compile failed: %v", err)
			}

			if entry.Cue != tc.expectedCue {
				t.Errorf("Expected cue %q, got %q", tc.expectedCue, entry.Cue)
			}

			if len(entry.Rb) == 0 || entry.Rb[0].String() != tc.expectedRb {
				t.Errorf("Expected RB %q, got %q", tc.expectedRb, func() string {
					if len(entry.Rb) == 0 {
						return "none"
					}
					return entry.Rb[0].String()
				}())
			}
		})
	}
}

// TestMatchingPrecedence_DefaultFallback verifies that defaults are used when season is missing.
func TestMatchingPrecedence_DefaultFallback(t *testing.T) {
	testPlan := createDefaultFallbackPlan()

	testCases := []struct {
		season      calendar.LiturgicalSeason
		weekday     calendar.Weekday
		expectedCue string
		expectedRb  string
		description string
	}{
		{
			season:      calendar.Advent,
			weekday:     calendar.Monday,
			expectedCue: "Advent Monday",
			expectedRb:  "RB 2.1",
			description: "Advent has season plan, uses weekday entry",
		},
		{
			season:      calendar.Ordinary,
			weekday:     calendar.Monday,
			expectedCue: "Default Reading",
			expectedRb:  "RB 1",
			description: "Ordinary Time missing plan, uses default",
		},
		{
			season:      calendar.Eastertide,
			weekday:     calendar.Friday,
			expectedCue: "Default Reading",
			expectedRb:  "RB 1",
			description: "Easter missing plan, uses default",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			dayKey := calendar.DayKey{
				Date:       "2025-06-09",
				Tradition:  calendar.RomanCalendar,
				Season:     tc.season,
				SeasonWeek: 1,
				Weekday:    tc.weekday,
			}

			entry, err := compile.Compile(dayKey, testPlan)
			if err != nil {
				t.Fatalf("Compile failed: %v", err)
			}

			if entry.Cue != tc.expectedCue {
				t.Errorf("Expected cue %q, got %q", tc.expectedCue, entry.Cue)
			}

			if len(entry.Rb) == 0 || entry.Rb[0].String() != tc.expectedRb {
				t.Errorf("Expected RB %q, got %q", tc.expectedRb, func() string {
					if len(entry.Rb) == 0 {
						return "none"
					}
					return entry.Rb[0].String()
				}())
			}
		})
	}
}

// TestOutputStability verifies that compiled output is stable and consistent across multiple compilations.
func TestOutputStability(t *testing.T) {
	testPlan := createLentTriduumTransitionPlan()
	dayKey := calendar.DayKey{
		Date:       "2025-03-05",
		Tradition:  calendar.RomanCalendar,
		Season:     calendar.Lent,
		SeasonWeek: 1,
		Weekday:    calendar.Wednesday,
	}

	// Compile multiple times and verify output is consistent
	var firstEntry *plan.FormattedEntry
	for i := 0; i < 3; i++ {
		t.Run(fmt.Sprintf("Compilation%d", i+1), func(t *testing.T) {
			entry, err := compile.Compile(dayKey, testPlan)
			if err != nil {
				t.Fatalf("Compile failed: %v", err)
			}

			if i == 0 {
				firstEntry = entry
			} else {
				// Compare with first compilation
				if entry.Cue != firstEntry.Cue {
					t.Errorf("Cue mismatch: %q vs %q", entry.Cue, firstEntry.Cue)
				}
				if len(entry.Rb) != len(firstEntry.Rb) {
					t.Errorf("RB length mismatch: %d vs %d", len(entry.Rb), len(firstEntry.Rb))
				}
				for j, rb := range entry.Rb {
					if rb.String() != firstEntry.Rb[j].String() {
						t.Errorf("RB[%d] mismatch: %q vs %q", j, rb.String(), firstEntry.Rb[j].String())
					}
				}
				// Tags comparison
				if (entry.Tags == nil) != (firstEntry.Tags == nil) {
					t.Errorf("Tags nil state mismatch")
				}
				if entry.Tags != nil && firstEntry.Tags != nil {
					if len(*entry.Tags) != len(*firstEntry.Tags) {
						t.Errorf("Tags length mismatch: %d vs %d", len(*entry.Tags), len(*firstEntry.Tags))
					}
					for j, tag := range *entry.Tags {
						if tag != (*firstEntry.Tags)[j] {
							t.Errorf("Tag[%d] mismatch: %q vs %q", j, tag, (*firstEntry.Tags)[j])
						}
					}
				}
			}
		})
	}
}

// Helper functions to create test plans

func createLentTriduumTransitionPlan() plan.Plan {
	advTag := "advent"
	lentTag := "lent"
	triduumTag := "triduum"
	easterTag := "easter"

	return plan.Plan{
		Version: 1,
		Work:    "Test Plan",
		Witness: "test",
		Defaults: plan.PlanEntry{
			Cue: "Default Reading",
			Rb:  []string{"RB 1"},
		},
		Seasons: map[string]plan.SeasonPlan{
			string(calendar.Advent): {
				Weekdays: map[string]plan.PlanEntry{
					"mon": {Cue: "Monday in Advent", Rb: []string{"RB 2.1"}, Tags: &[]string{advTag}},
					"tue": {Cue: "Tuesday in Advent", Rb: []string{"RB 2.2"}, Tags: &[]string{advTag}},
					"wed": {Cue: "Wednesday in Advent", Rb: []string{"RB 2.3"}, Tags: &[]string{advTag}},
					"thu": {Cue: "Thursday in Advent", Rb: []string{"RB 2.4"}, Tags: &[]string{advTag}},
					"fri": {Cue: "Friday in Advent", Rb: []string{"RB 2.5"}, Tags: &[]string{advTag}},
					"sat": {Cue: "Saturday in Advent", Rb: []string{"RB 2.6"}, Tags: &[]string{advTag}},
					"sun": {Cue: "Sunday in Advent", Rb: []string{"RB 3.1"}, Tags: &[]string{advTag}},
				},
			},
			string(calendar.Lent): {
				Weekdays: map[string]plan.PlanEntry{
					"mon": {Cue: "Monday in Lent", Rb: []string{"RB 8.1"}, Tags: &[]string{lentTag}},
					"tue": {Cue: "Tuesday in Lent", Rb: []string{"RB 8.2"}, Tags: &[]string{lentTag}},
					"wed": {Cue: "Wednesday in Lent", Rb: []string{"RB 8.3"}, Tags: &[]string{lentTag}},
					"thu": {Cue: "Thursday in Lent", Rb: []string{"RB 8.4"}, Tags: &[]string{lentTag}},
					"fri": {Cue: "Friday in Lent", Rb: []string{"RB 8.5"}, Tags: &[]string{lentTag}},
					"sat": {Cue: "Saturday in Lent", Rb: []string{"RB 8.6"}, Tags: &[]string{lentTag}},
					"sun": {Cue: "First Sunday of Lent", Rb: []string{"RB 9.1"}, Tags: &[]string{lentTag}},
				},
			},
			string(calendar.Triduum): {
				Weekdays: map[string]plan.PlanEntry{
					"mon": {Cue: "Triduum Monday", Rb: []string{"RB 10.1"}, Tags: &[]string{triduumTag}},
					"tue": {Cue: "Triduum Tuesday", Rb: []string{"RB 10.2"}, Tags: &[]string{triduumTag}},
					"wed": {Cue: "Triduum Wednesday", Rb: []string{"RB 10.3"}, Tags: &[]string{triduumTag}},
					"thu": {Cue: "Holy Thursday", Rb: []string{"RB 10.4"}, Tags: &[]string{triduumTag}},
					"fri": {Cue: "Good Friday", Rb: []string{"RB 10.5"}, Tags: &[]string{triduumTag}},
					"sat": {Cue: "Holy Saturday", Rb: []string{"RB 10.6"}, Tags: &[]string{triduumTag}},
					"sun": {Cue: "Triduum Sunday", Rb: []string{"RB 11.1"}, Tags: &[]string{triduumTag}},
				},
			},
			string(calendar.Eastertide): {
				Weekdays: map[string]plan.PlanEntry{
					"mon": {Cue: "Easter Monday", Rb: []string{"RB 12.1"}, Tags: &[]string{easterTag}},
					"tue": {Cue: "Easter Tuesday", Rb: []string{"RB 12.2"}, Tags: &[]string{easterTag}},
					"wed": {Cue: "Easter Wednesday", Rb: []string{"RB 12.3"}, Tags: &[]string{easterTag}},
					"thu": {Cue: "Easter Thursday", Rb: []string{"RB 12.4"}, Tags: &[]string{easterTag}},
					"fri": {Cue: "Easter Friday", Rb: []string{"RB 12.5"}, Tags: &[]string{easterTag}},
					"sat": {Cue: "Easter Saturday", Rb: []string{"RB 12.6"}, Tags: &[]string{easterTag}},
					"sun": {Cue: "Easter Sunday", Rb: []string{"RB 13.1"}, Tags: &[]string{easterTag}},
				},
			},
		},
	}
}

func createEasterSeasonPlan() plan.Plan {
	easterTag := "easter"

	// Create a plan that distinguishes specific Sundays and special days
	return plan.Plan{
		Version: 1,
		Work:    "Easter Season Test Plan",
		Witness: "test",
		Defaults: plan.PlanEntry{
			Cue: "Default Reading",
			Rb:  []string{"RB 1"},
		},
		Seasons: map[string]plan.SeasonPlan{
			string(calendar.Eastertide): {
				Weekdays: map[string]plan.PlanEntry{
					"mon": {Cue: "Easter Monday", Rb: []string{"RB 12.1"}, Tags: &[]string{easterTag}},
					"tue": {Cue: "Easter Tuesday", Rb: []string{"RB 12.2"}, Tags: &[]string{easterTag}},
					"wed": {Cue: "Easter Wednesday", Rb: []string{"RB 12.3"}, Tags: &[]string{easterTag}},
					"thu": {Cue: "Ascension Thursday", Rb: []string{"RB 12.4"}, Tags: &[]string{easterTag}},
					"fri": {Cue: "Easter Friday", Rb: []string{"RB 12.5"}, Tags: &[]string{easterTag}},
					"sat": {Cue: "Easter Saturday", Rb: []string{"RB 12.6"}, Tags: &[]string{easterTag}},
					"sun": {Cue: "Pentecost Sunday", Rb: []string{"RB 13.7"}, Tags: &[]string{easterTag}},
				},
			},
		},
	}
}

func createAdventChristmasEpiphanyPlan() plan.Plan {
	advTag := "advent"
	christmasTag := "christmas"
	epiphanyTag := "epiphany"

	return plan.Plan{
		Version: 1,
		Work:    "Advent/Christmas/Epiphany Test Plan",
		Witness: "test",
		Defaults: plan.PlanEntry{
			Cue: "Default Reading",
			Rb:  []string{"RB 1"},
		},
		Seasons: map[string]plan.SeasonPlan{
			string(calendar.Advent): {
				Weekdays: map[string]plan.PlanEntry{
					"mon": {Cue: "Monday before Christmas", Rb: []string{"RB 2.1"}, Tags: &[]string{advTag}},
					"tue": {Cue: "Tuesday before Christmas", Rb: []string{"RB 2.2"}, Tags: &[]string{advTag}},
					"wed": {Cue: "Christmas Eve", Rb: []string{"RB 2.2"}, Tags: &[]string{advTag}},
					"thu": {Cue: "Thursday in Advent", Rb: []string{"RB 2.4"}, Tags: &[]string{advTag}},
					"fri": {Cue: "Friday in Advent", Rb: []string{"RB 2.5"}, Tags: &[]string{advTag}},
					"sat": {Cue: "Saturday in Advent", Rb: []string{"RB 2.6"}, Tags: &[]string{advTag}},
					"sun": {Cue: "Fourth Sunday of Advent", Rb: []string{"RB 3.4"}, Tags: &[]string{advTag}},
				},
			},
			string(calendar.Christmastide): {
				Weekdays: map[string]plan.PlanEntry{
					"mon": {Cue: "Monday after Christmas", Rb: []string{"RB 4.1"}, Tags: &[]string{christmasTag}},
					"tue": {Cue: "Tuesday after Christmas", Rb: []string{"RB 4.2"}, Tags: &[]string{christmasTag}},
					"wed": {Cue: "Christmas Day", Rb: []string{"RB 4.3"}, Tags: &[]string{christmasTag}},
					"thu": {Cue: "St. Stephen", Rb: []string{"RB 4.4"}, Tags: &[]string{christmasTag}},
					"fri": {Cue: "Holy Innocents", Rb: []string{"RB 4.5"}, Tags: &[]string{christmasTag}},
					"sat": {Cue: "Saturday after Christmas", Rb: []string{"RB 4.6"}, Tags: &[]string{christmasTag}},
					"sun": {Cue: "Holy Family Sunday", Rb: []string{"RB 5.1"}, Tags: &[]string{christmasTag}},
				},
				Fallback: &plan.PlanEntry{
					Cue: "Mary Mother of God",
					Rb:  []string{"RB 4.3"},
				},
			},
			string(calendar.Epiphanytide): {
				Weekdays: map[string]plan.PlanEntry{
					"mon": {Cue: "Epiphany", Rb: []string{"RB 6.1"}, Tags: &[]string{epiphanyTag}},
					"tue": {Cue: "Tuesday after Epiphany", Rb: []string{"RB 6.2"}, Tags: &[]string{epiphanyTag}},
					"wed": {Cue: "Wednesday after Epiphany", Rb: []string{"RB 6.3"}, Tags: &[]string{epiphanyTag}},
					"thu": {Cue: "Thursday after Epiphany", Rb: []string{"RB 6.4"}, Tags: &[]string{epiphanyTag}},
					"fri": {Cue: "Friday after Epiphany", Rb: []string{"RB 6.5"}, Tags: &[]string{epiphanyTag}},
					"sat": {Cue: "Saturday after Epiphany", Rb: []string{"RB 6.6"}, Tags: &[]string{epiphanyTag}},
					"sun": {Cue: "First Sunday after Epiphany", Rb: []string{"RB 7.1"}, Tags: &[]string{epiphanyTag}},
				},
			},
		},
	}
}

func createWeekdayOverridePlan() plan.Plan {
	advTag := "advent"

	return plan.Plan{
		Version: 1,
		Work:    "Weekday Override Test Plan",
		Witness: "test",
		Defaults: plan.PlanEntry{
			Cue: "Default Reading",
			Rb:  []string{"RB 1"},
		},
		Seasons: map[string]plan.SeasonPlan{
			string(calendar.Advent): {
				Weekdays: map[string]plan.PlanEntry{
					"mon": {Cue: "Advent Monday Override", Rb: []string{"RB 2.1"}, Tags: &[]string{advTag}},
					"tue": {Cue: "Advent Tuesday Override", Rb: []string{"RB 2.2"}, Tags: &[]string{advTag}},
					"sun": {Cue: "Sunday in Advent", Rb: []string{"RB 3.1"}, Tags: &[]string{advTag}},
				},
				Fallback: &plan.PlanEntry{
					Cue: "Advent Fallback",
					Rb:  []string{"RB 2.3"},
				},
			},
		},
	}
}

func createFallbackOnlyPlan() plan.Plan {
	advTag := "advent"
	lentTag := "lent"

	return plan.Plan{
		Version: 1,
		Work:    "Fallback Only Test Plan",
		Witness: "test",
		Defaults: plan.PlanEntry{
			Cue: "Default Reading",
			Rb:  []string{"RB 1"},
		},
		Seasons: map[string]plan.SeasonPlan{
			string(calendar.Advent): {
				Fallback: &plan.PlanEntry{
					Cue:  "Advent Fallback",
					Rb:   []string{"RB 2.1"},
					Tags: &[]string{advTag},
				},
			},
			string(calendar.Lent): {
				Fallback: &plan.PlanEntry{
					Cue:  "Lent Fallback",
					Rb:   []string{"RB 8.1"},
					Tags: &[]string{lentTag},
				},
			},
		},
	}
}

func createDefaultFallbackPlan() plan.Plan {
	advTag := "advent"

	return plan.Plan{
		Version: 1,
		Work:    "Default Fallback Test Plan",
		Witness: "test",
		Defaults: plan.PlanEntry{
			Cue: "Default Reading",
			Rb:  []string{"RB 1"},
		},
		Seasons: map[string]plan.SeasonPlan{
			string(calendar.Advent): {
				Weekdays: map[string]plan.PlanEntry{
					"mon": {Cue: "Advent Monday", Rb: []string{"RB 2.1"}, Tags: &[]string{advTag}},
					"tue": {Cue: "Advent Tuesday", Rb: []string{"RB 2.2"}, Tags: &[]string{advTag}},
					"wed": {Cue: "Advent Wednesday", Rb: []string{"RB 2.3"}, Tags: &[]string{advTag}},
					"thu": {Cue: "Advent Thursday", Rb: []string{"RB 2.4"}, Tags: &[]string{advTag}},
					"fri": {Cue: "Advent Friday", Rb: []string{"RB 2.5"}, Tags: &[]string{advTag}},
					"sat": {Cue: "Advent Saturday", Rb: []string{"RB 2.6"}, Tags: &[]string{advTag}},
					"sun": {Cue: "Sunday in Advent", Rb: []string{"RB 3.1"}, Tags: &[]string{advTag}},
				},
			},
		},
	}
}
