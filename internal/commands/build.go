package commands

import (
	"github.com/go-playground/validator/v10"

	"github.com/julianstephens/liturgical-time-index/internal/calendar"
)

type BuildCmd struct {
	Year string `arg:"" name:"year" help:"The year to build the index for."`
}

func (c *BuildCmd) Run(validate *validator.Validate) error {
	ce := calendar.NewCalendarEngine(validate)
	holidays, err := ce.Holidays(2026, calendar.RomanCalendar)
	if err != nil {
		return err
	}

	for name, day := range holidays {
		println(name, day.Date, day.Season, day.SeasonWeek, day.Weekday)
	}

	return nil
}
