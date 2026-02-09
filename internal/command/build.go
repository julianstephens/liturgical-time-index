package command

import (
	"github.com/julianstephens/liturgical-time-index/internal/calendar"
)

type BuildCmd struct {
	Year string `arg:"" name:"year" help:"The year to build the index for."`
}

func (c *BuildCmd) Run() error {
	ce := calendar.NewCalendarEngine()
	holidays, err := ce.Holidays(2026, calendar.RomanCalendar)
	if err != nil {
		return err
	}

	for name, day := range holidays {
		println(name, day.Date, day.Season, day.SeasonWeek, day.Weekday)
	}

	return nil
}
