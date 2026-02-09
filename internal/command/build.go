package command

import (
	"fmt"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/julianstephens/go-utils/cliutil"
	"github.com/julianstephens/go-utils/helpers"
	"github.com/julianstephens/liturgical-time-index/internal/calendar"
	"github.com/julianstephens/liturgical-time-index/internal/compile"
	"github.com/julianstephens/liturgical-time-index/internal/output"
	"github.com/julianstephens/liturgical-time-index/internal/plan"
)

type BuildCmd struct {
	Year         string  `name:"year"      help:"The year to build the index for."`
	Plan         string  `name:"plan"      help:"The path to the plan file to build the index from."             default:"./plan.yaml"`
	Tradition    string  `name:"tradition" help:"The liturgical tradition to build the index for."               default:"roman"       enum:"roman"`
	ICSPath      *string `name:"out"       help:"The path to output the ICalendar file to (e.g. ./calendar.ics)"                                    required:"" xor:"md,out"`
	MarkdownPath *string `name:"md"        help:"The path to output the Markdown file to (e.g. ./calendar.md)"                                      required:"" xor:"md,out"`
	MarkdownType string  `name:"type" help:"Whether to output the full calendar or a specific season." default:"annual" enum:"annual,advent,christmastide,epiphanytide,lent,triduum,easter,ordinary"`
	Verbose      bool    `name:"verbose"   help:"Enable verbose logging."`
}

func (c *BuildCmd) Run() error {
	ce := calendar.NewCalendarEngine()

	p, err := plan.LoadAndValidatePlan(c.Plan)
	if err != nil {
		cliutil.PrintError("Unable to load and validate plan file")
		return err
	}

	tradition := calendar.CalendarTradition(c.Tradition)
	if tradition == "" {
		cliutil.PrintError(fmt.Sprintf("Unsupported tradition: %s", c.Tradition))
		return fmt.Errorf("unsupported tradition: %s", c.Tradition)
	}

	calendar, err := ce.GenerateRomanCalendar(c.Year, tradition)
	if err != nil {
		cliutil.PrintError("Unable to generate calendar")
		return err
	}

	entries := make([]plan.FormattedEntry, len(calendar))
	caser := cases.Title(language.English)
	for i, day := range calendar {
		if c.MarkdownType != "annual" && string(day.Season) != caser.String(c.MarkdownType) {
			continue
		}
		entry, err := compile.Compile(day, *p)
		if err != nil {
			cliutil.PrintError("Unable to compile calendar and plan into entries")
			return err
		}

		if c.Verbose {
			cliutil.PrintInfo(
				fmt.Sprintf(
					"Generated calendar entry for %s: %s, season week %d, weekday %s",
					day.Date,
					day.Season,
					day.SeasonWeek,
					day.Weekday,
				),
			)
		}
		entries[i] = *entry
	}

	if c.ICSPath != nil {
		if helpers.Exists(*c.ICSPath) {
			cliutil.PrintError(fmt.Sprintf("Output file already exists: %s", *c.ICSPath))
			return fmt.Errorf("output file already exists: %s", *c.ICSPath)
		}

		if err := output.ICS(entries, *c.ICSPath); err != nil {
			cliutil.PrintError("Unable to output entries to ICS")
			return err
		}
	}

	if c.MarkdownPath != nil {
		if helpers.Exists(*c.MarkdownPath) {
			cliutil.PrintError(fmt.Sprintf("Output file already exists: %s", *c.MarkdownPath))
			return fmt.Errorf("output file already exists: %s", *c.MarkdownPath)
		}

		if err := output.Markdown(entries, *c.MarkdownPath); err != nil {
			cliutil.PrintError("Unable to output entries to Markdown")
			return err
		}
	}

	return nil
}
