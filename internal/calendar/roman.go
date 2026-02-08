package calendar

import (
	"math"
	"strconv"
	"time"

	"github.com/julianstephens/liturgical-time-index/internal/util"
)

// GenerateRomanCalendar generates a list of DayKey entries for each day in the specified year and tradition.
// It iterates through each day of the year, checks if it's a valid date, and then generates a DayKey for that date.
func (ce *CalendarEngine) GenerateRomanCalendar(year string, tradition CalendarTradition) ([]DayKey, error) {
	result := []DayKey{}

	for _, month := range []string{"01", "02", "03", "04", "05", "06", "07", "08", "09", "10", "11", "12"} {
		for day := 1; day <= 31; day++ {
			date := year + "-" + month + "-" + padZero(day)
			_, err := time.Parse(util.DateFormat, date)
			if err == nil {
				// valid date, generate DayKey
				dayKey, err := ce.GetRomanDay(date, tradition)
				if err != nil {
					return nil, err
				}

				if err := ce.validate.Struct(dayKey); err != nil {
					return nil, &CalendarError{
						Err:   ErrValidationFailed,
						Cause: err,
					}
				}

				result = append(result, *dayKey)
			}
		}
	}
	return result, nil
}

// GetRomanDay generates a DayKey for a given date and tradition by determining the season, season week, and weekday.
func (ce *CalendarEngine) GetRomanDay(date string, tradition CalendarTradition) (*DayKey, error) {
	if tradition != RomanCalendar {
		return nil, &CalendarError{
			Err: ErrUnsupportedCalendarTradition,
		}
	}

	season, err := ce.GetRomanSeason(date, tradition)
	if err != nil {
		return nil, err
	}

	seasonWeek, err := ce.GetRomanSeasonWeek(date, season, tradition)
	if err != nil {
		return nil, err
	}

	weekday, err := ce.GetRomanWeekday(date)
	if err != nil {
		return nil, err
	}

	return &DayKey{
		Date:       date,
		Tradition:  tradition,
		Season:     season,
		SeasonWeek: seasonWeek,
		Weekday:    weekday,
	}, nil
}

// GetRomanSeason determines the liturgical season for a given date and tradition by calculating key feast dates and comparing them to the input date.
func (ce *CalendarEngine) GetRomanSeason(date string, tradition CalendarTradition) (LiturgicalSeason, error) {
	if tradition == RomanCalendar {
		return ce.getRomanSeason(date)
	}
	return "", &CalendarError{
		Err: ErrUnsupportedCalendarTradition,
	}
}

// getRomanSeason determines the liturgical season for a given date in the Roman calendar tradition.
// It calculates the dates of key movable feasts like Easter and Ash Wednesday to determine the season.
// The logic is based on the general rules for the Roman liturgical calendar, with specific date ranges for each season.
func (ce *CalendarEngine) getRomanSeason(date string) (LiturgicalSeason, error) {
	parsed, err := time.Parse(util.DateFormat, date)
	if err != nil {
		return "", &CalendarError{
			Err:   ErrParseDateFailed,
			Cause: err,
		}
	}
	parsed = parsed.Truncate(24 * time.Hour)

	month := parsed.Month()
	day := parsed.Day()
	easterDay := ce.getEasterGregorian(parsed.Year())
	easterDay = easterDay.Truncate(24 * time.Hour)
	ashWednesday := easterDay.AddDate(0, 0, -46)
	ashWednesday = ashWednesday.Truncate(24 * time.Hour)
	holyThursday := easterDay.AddDate(0, 0, -3)
	holyThursday = holyThursday.Truncate(24 * time.Hour)
	pentecost := easterDay.AddDate(0, 0, 49)
	pentecost = pentecost.Truncate(24 * time.Hour)
	nov27 := time.Date(parsed.Year(), time.November, 27, 0, 0, 0, 0, time.Local)
	sundayAfterNov27 := nov27
	if nov27.Weekday() != time.Sunday {
		sundayAfterNov27 = nov27.AddDate(0, 0, int(time.Sunday-nov27.Weekday()+7)%7)
	}
	sundayAfterNov27 = sundayAfterNov27.Truncate(24 * time.Hour)

	switch month {
	case time.November:
		if parsed.Equal(sundayAfterNov27) || parsed.After(sundayAfterNov27) {
			return Advent, nil
		}
		return Ordinary, nil
	case time.December:
		if day >= 25 {
			return Christmastide, nil
		}
		return Advent, nil
	case time.January:
		if day < 6 {
			return Christmastide, nil
		}
		return Epiphanytide, nil
	case time.February:
		if ashWednesday.Month() == time.February && (parsed.Equal(ashWednesday) || parsed.After(ashWednesday)) {
			return Lent, nil
		}
		return Epiphanytide, nil
	case time.March, time.April, time.May:
		if parsed.Before(ashWednesday) {
			return Epiphanytide, nil
		}
		if (parsed.Equal(ashWednesday) || parsed.After(ashWednesday)) && parsed.Before(holyThursday) {
			return Lent, nil
		}
		if (parsed.Equal(holyThursday) || parsed.After(holyThursday)) && parsed.Before(easterDay) {
			return Triduum, nil
		}
		if (parsed.Equal(easterDay) || parsed.After(easterDay)) &&
			(parsed.Equal(pentecost) || parsed.Before(pentecost)) {
			return Easter, nil
		}
		return Ordinary, nil
	default:
		return Ordinary, nil
	}
}

// GetRomanWeekday determines the weekday for a given date string in ISO8601 format.
func (ce *CalendarEngine) GetRomanWeekday(date string) (Weekday, error) {
	parsed, err := time.Parse(util.DateFormat, date)
	if err != nil {
		return "", &CalendarError{
			Err:   ErrParseDateFailed,
			Cause: err,
		}
	}

	weekday := parsed.Weekday()
	switch weekday {
	case time.Sunday:
		return Sunday, nil
	case time.Monday:
		return Monday, nil
	case time.Tuesday:
		return Tuesday, nil
	case time.Wednesday:
		return Wednesday, nil
	case time.Thursday:
		return Thursday, nil
	case time.Friday:
		return Friday, nil
	case time.Saturday:
		return Saturday, nil
	default:
		return "", nil
	}
}

// GetRomanSeasonWeek calculates the week number within the liturgical season for a given date, season, and tradition.
// It determines the start date of the season and calculates how many weeks have passed since that date.
func (ce *CalendarEngine) GetRomanSeasonWeek(
	date string,
	season LiturgicalSeason,
	tradition CalendarTradition,
) (int, error) {
	seasonStartDate, err := ce.getRomanSeasonStartDate(date, season, tradition)
	if err != nil {
		return 0, err
	}
	seasonStartDate = seasonStartDate.Truncate(24 * time.Hour)

	parsed, err := time.Parse(util.DateFormat, date)
	if err != nil {
		return 0, &CalendarError{
			Err:   ErrParseDateFailed,
			Cause: err,
		}
	}
	parsed = parsed.Truncate(24 * time.Hour)

	daysSinceSeasonStart := int(parsed.Sub(seasonStartDate).Hours() / 24)
	if daysSinceSeasonStart < 0 {
		return 0, &CalendarError{
			Message: ptr("date is before the start of the season"),
			Err:     ErrValidationFailed,
		}
	}

	if daysSinceSeasonStart == 0 {
		return 1, nil
	}
	calculatedWeek := math.Floor(float64(daysSinceSeasonStart) / 7)

	return 1 + int(calculatedWeek), nil
}

// getRomanSeasonStartDate returns the start date of a given liturgical season for a specific date and tradition.
// It calculates the start date based on the rules of the Roman calendar tradition, using key feast dates like Easter and Ash Wednesday.
func (ce *CalendarEngine) getRomanSeasonStartDate(
	date string,
	season LiturgicalSeason,
	tradition CalendarTradition,
) (time.Time, error) {
	if tradition != RomanCalendar {
		return time.Time{}, &CalendarError{
			Err: ErrUnsupportedCalendarTradition,
		}
	}

	parsed, err := time.Parse(util.DateFormat, date)
	if err != nil {
		return time.Time{}, &CalendarError{
			Err:   ErrParseDateFailed,
			Cause: err,
		}
	}

	switch season {
	case Advent:
		nov27 := time.Date(parsed.Year(), time.November, 27, 0, 0, 0, 0, time.Local)
		sundayAfterNov27 := nov27
		if nov27.Weekday() != time.Sunday {
			sundayAfterNov27 = nov27.AddDate(0, 0, int(time.Sunday-nov27.Weekday()+7)%7)
		}
		return sundayAfterNov27, nil
	case Christmastide:
		if parsed.Month() == time.December {
			return time.Date(parsed.Year(), time.December, 25, 0, 0, 0, 0, time.Local), nil
		}
		return time.Date(parsed.Year()-1, time.December, 25, 0, 0, 0, 0, time.Local), nil
	case Epiphanytide:
		return time.Date(parsed.Year(), time.January, 6, 0, 0, 0, 0, time.Local), nil
	case Lent:
		easterDay := ce.getEasterGregorian(parsed.Year())
		return easterDay.AddDate(0, 0, -46), nil
	case Triduum:
		easterDay := ce.getEasterGregorian(parsed.Year())
		return easterDay.AddDate(0, 0, -3), nil
	case Easter:
		return ce.getEasterGregorian(parsed.Year()), nil
	case Ordinary:
		easterDay := ce.getEasterGregorian(parsed.Year())
		pentacost := easterDay.AddDate(0, 0, 49)
		return pentacost.AddDate(0, 0, 1), nil
	default:
		return time.Time{}, nil
	}
}

// Holidays generates a map of key holidays for a given year and tradition, including their dates, seasons, season weeks, and weekdays.
func (ce *CalendarEngine) Holidays(year int, tradition CalendarTradition) (map[string]DayKey, error) {
	easterDay := ce.getEasterGregorian(year)
	ashWednesday := easterDay.AddDate(0, 0, -46)
	holyThursday := easterDay.AddDate(0, 0, -3)
	goodFriday := easterDay.AddDate(0, 0, -2)
	easterMonday := easterDay.AddDate(0, 0, 1)
	pentecost := easterDay.AddDate(0, 0, 49)

	holidayDates := []time.Time{ashWednesday, holyThursday, goodFriday, easterDay, easterMonday, pentecost}
	holidayNames := []string{
		"Ash Wednesday",
		"Holy Thursday",
		"Good Friday",
		"Easter Sunday",
		"Easter Monday",
		"Pentecost",
	}

	holidays := make(map[string]DayKey)
	for i, holiday := range holidayNames {
		dateStr := holidayDates[i].Format(util.DateFormat)
		dayKey, err := ce.GetRomanDay(dateStr, tradition)
		if err != nil {
			return nil, err
		}
		holidays[holiday] = *dayKey
	}

	return holidays, nil
}

func padZero(num int) string {
	if num < 10 {
		return "0" + strconv.Itoa(num)
	}
	return strconv.Itoa(num)
}

func ptr[T any](s T) *T {
	return &s
}
