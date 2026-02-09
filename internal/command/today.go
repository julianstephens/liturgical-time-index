package command

import (
	"fmt"
	"strconv"
	"time"

	"github.com/julianstephens/go-utils/cliutil"
	"github.com/julianstephens/liturgical-time-index/internal"
	"github.com/julianstephens/liturgical-time-index/internal/calendar"
	"github.com/julianstephens/liturgical-time-index/internal/compile"
	"github.com/julianstephens/liturgical-time-index/internal/plan"
)

type TodayCmd struct {
	Date      *string `name:"date"      help:"The date to get the entry for (e.g. 2024-12-25). If not provided, defaults to today's date."`
	Tradition string  `name:"tradition" help:"The liturgical tradition to get the entry for."                                              default:"roman"       enum:"roman"`
	Plan      string  `name:"plan"      help:"The path to the plan file to use for looking up the entry."                                  default:"./plan.yaml"`
}

func (c *TodayCmd) Run() error {
	p, err := plan.LoadAndValidatePlan(c.Plan)
	if err != nil {
		cliutil.PrintError("Unable to load and validate plan file")
		return err
	}

	if c.Date == nil {
		today := time.Now().Format(internal.DateFormat)
		c.Date = &today
	}

	formattedDate, err := time.Parse(internal.DateFormat, *c.Date)
	if err != nil {
		cliutil.PrintError(fmt.Sprintf("Invalid date format: %s. Expected format: %s", *c.Date, internal.DateFormat))
		return fmt.Errorf("invalid date format: %s. expected format: %s", *c.Date, internal.DateFormat)
	}

	ce := calendar.NewCalendarEngine()
	calendar, err := ce.GenerateRomanCalendar(
		strconv.Itoa(formattedDate.Year()),
		calendar.CalendarTradition(c.Tradition),
	)
	if err != nil {
		cliutil.PrintError("Unable to generate calendar")
		return err
	}

	var entry *plan.FormattedEntry
	for _, day := range calendar {
		if day.Date == *c.Date {
			e, err := compile.Compile(day, *p)
			if err != nil {
				cliutil.PrintError("Unable to compile calendar and plan into entry")
				return err
			}
			entry = e
			break
		}
	}

	if entry == nil {
		cliutil.PrintError(fmt.Sprintf("No entry found for date: %s", *c.Date))
		return fmt.Errorf("no entry found for date: %s", *c.Date)
	}

	fmt.Println()
	cliutil.PrintColored(*c.Date, cliutil.ColorBlue)
	cliutil.PrintColored(
		fmt.Sprintf("Season: %s, Week: %d, Weekday: %s", entry.Key.Season, entry.Key.SeasonWeek, entry.Key.Weekday),
		cliutil.ColorBold,
	)
	fmt.Println("---")
	fmt.Println()
	cliutil.PrintColored(entry.Cue, cliutil.ColorMagenta)
	fmt.Println()
	for _, rb := range entry.Rb {
		fmt.Println("- " + rb.String())
	}
	fmt.Println()

	return nil
}
