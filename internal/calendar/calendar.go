package calendar

import (
	"time"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/julianstephens/go-utils/generic"
	"github.com/julianstephens/liturgical-time-index/internal"
)

type CalendarEngine struct{}

type CalendarTradition string

const (
	RomanCalendar CalendarTradition = "roman"
)

type LiturgicalSeason string

const (
	Advent        LiturgicalSeason = "advent"
	Christmastide LiturgicalSeason = "christmastide"
	Epiphanytide  LiturgicalSeason = "epiphanytide"
	Lent          LiturgicalSeason = "lent"
	Triduum       LiturgicalSeason = "triduum"
	Eastertide    LiturgicalSeason = "eastertide"
	Ordinary      LiturgicalSeason = "ordinary"
)

func (s LiturgicalSeason) String() string {
	switch s {
	case Advent:
		return "Advent"
	case Christmastide:
		return "Christmas"
	case Epiphanytide:
		return "Epiphany"
	case Lent:
		return "Lent"
	case Triduum:
		return "Paschal Triduum"
	case Eastertide:
		return "Eastertide"
	case Ordinary:
		return "Ordinary Time"
	default:
		caser := cases.Title(language.English)
		return caser.String(string(s))
	}
}

type Weekday string

const (
	Sunday    Weekday = "sun"
	Monday    Weekday = "mon"
	Tuesday   Weekday = "tue"
	Wednesday Weekday = "wed"
	Thursday  Weekday = "thu"
	Friday    Weekday = "fri"
	Saturday  Weekday = "sat"
)

func (w Weekday) String() string {
	switch w {
	case Sunday:
		return "Sunday"
	case Monday:
		return "Monday"
	case Tuesday:
		return "Tuesday"
	case Wednesday:
		return "Wednesday"
	case Thursday:
		return "Thursday"
	case Friday:
		return "Friday"
	case Saturday:
		return "Saturday"
	default:
		caser := cases.Title(language.English)
		return caser.String(string(w))
	}
}

type DayKey struct {
	Date       string            `json:"date"`
	Tradition  CalendarTradition `json:"tradition"`
	Season     LiturgicalSeason  `json:"season"      validate:"required"`
	SeasonWeek int               `json:"season_week" validate:"required,gte=1"`
	Weekday    Weekday           `json:"weekday"     validate:"required"`
}

func NewCalendarEngine() *CalendarEngine {
	return &CalendarEngine{}
}

// GetEasterGregorian computes the date of Easter for a given year using Butcher's algorithm for the Gregorian calendar.
// For years before 1583, it uses a simpler algorithm based on the Julian calendar.
func (ce *CalendarEngine) GetEasterGregorian(year int) time.Time {
	var a, b, c, d, e, r int

	a = year % 19
	if year >= 1583 {
		var f, g, h, i, k, l, m int
		b = year / 100
		c = year % 100
		d = b / 4
		e = b % 4
		f = (b + 8) / 25
		g = (b - f + 1) / 3
		h = (19*a + b - d - g + 15) % 30
		i = c / 4
		k = c % 4
		l = (32 + 2*e + 2*i - h - k) % 7
		m = (a + 11*h + 22*l) / 451
		r = 22 + h + l - 7*m
	} else {
		b = year % 7
		c = year % 4
		d = (19*a + 15) % 30
		e = (2*c + 4*b - d + 34) % 7
		r = 22 + d + e
	}

	return time.Date(year, time.March, r, 0, 0, 0, 0, time.Local)
}

// validate checks that the provided DayKey has valid values for its fields,
// including correct date format, valid season and tradition, and appropriate
// season week and weekday values.
func (ce *CalendarEngine) validate(dayKey *DayKey) error {
	if dayKey.Date == "" {
		return &CalendarError{
			Message: generic.Ptr("date is required"),
			Err:     ErrValidationFailed,
		}
	}
	_, err := time.Parse(internal.DateFormat, dayKey.Date)
	if err != nil {
		return &CalendarError{
			Message: generic.Ptr("invalid date format, expected YYYY-MM-DD"),
			Err:     ErrValidationFailed,
			Cause:   err,
		}
	}

	if dayKey.Season == "" {
		return &CalendarError{
			Message: generic.Ptr("season is required"),
			Err:     ErrValidationFailed,
		}
	}
	parsedSeason := LiturgicalSeason(dayKey.Season)
	switch parsedSeason {
	case Advent, Christmastide, Epiphanytide, Lent, Triduum, Eastertide, Ordinary:
		// valid season
	default:
		return &CalendarError{
			Message: generic.Ptr("invalid season"),
			Err:     ErrValidationFailed,
		}
	}

	if dayKey.Tradition == "" {
		return &CalendarError{
			Message: generic.Ptr("tradition is required"),
			Err:     ErrValidationFailed,
		}
	}
	parsedTradition := CalendarTradition(dayKey.Tradition)
	if parsedTradition != RomanCalendar {
		return &CalendarError{
			Message: generic.Ptr("unsupported calendar tradition"),
			Err:     ErrUnsupportedCalendarTradition,
		}
	}

	if dayKey.SeasonWeek < 1 {
		return &CalendarError{
			Message: generic.Ptr("season week must be at least 1"),
			Err:     ErrValidationFailed,
		}
	}
	if dayKey.SeasonWeek > 53 {
		return &CalendarError{
			Message: generic.Ptr("season week cannot be greater than 53"),
			Err:     ErrValidationFailed,
		}
	}

	if dayKey.Weekday == "" {
		return &CalendarError{
			Message: generic.Ptr("weekday is required"),
			Err:     ErrValidationFailed,
		}
	}
	parsedWeekday := Weekday(dayKey.Weekday)
	switch parsedWeekday {
	case Sunday, Monday, Tuesday, Wednesday, Thursday, Friday, Saturday:
		// valid weekday
	default:
		return &CalendarError{
			Message: generic.Ptr("invalid weekday"),
			Err:     ErrValidationFailed,
		}
	}

	return nil
}
