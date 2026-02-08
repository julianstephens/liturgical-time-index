package calendar

import (
	"time"

	"github.com/go-playground/validator/v10"
)

type CalendarEngine struct {
	validate *validator.Validate
}

type CalendarTradition string

const (
	RomanCalendar CalendarTradition = "Roman"
)

type LiturgicalSeason string

const (
	Advent        LiturgicalSeason = "Advent"
	Christmastide LiturgicalSeason = "Christmastide"
	Epiphanytide  LiturgicalSeason = "Epiphanytide"
	Lent          LiturgicalSeason = "Lent"
	Triduum       LiturgicalSeason = "Paschal Triduum"
	Easter        LiturgicalSeason = "Easter"
	Ordinary      LiturgicalSeason = "Ordinary Time"
)

type Weekday string

const (
	Sunday    Weekday = "Sunday"
	Monday    Weekday = "Monday"
	Tuesday   Weekday = "Tuesday"
	Wednesday Weekday = "Wednesday"
	Thursday  Weekday = "Thursday"
	Friday    Weekday = "Friday"
	Saturday  Weekday = "Saturday"
)

type DayKey struct {
	Date       string            `validate:"required,ISO8601" json:"date"`
	Tradition  CalendarTradition `validate:"required"         json:"tradition"`
	Season     LiturgicalSeason  `validate:"required"         json:"season"`
	SeasonWeek int               `validate:"required,gte=1"   json:"season_week"`
	Weekday    Weekday           `validate:"required"         json:"weekday"`
}

func NewCalendarEngine(validate *validator.Validate) *CalendarEngine {
	return &CalendarEngine{
		validate: validate,
	}
}

// getEasterGregorian computes the date of Easter for a given year using Butcher's algorithm for the Gregorian calendar.
// For years before 1583, it uses a simpler algorithm based on the Julian calendar.
func (ce *CalendarEngine) getEasterGregorian(year int) time.Time {
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
