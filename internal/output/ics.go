package output

import (
	"fmt"
	"strings"
	"time"

	ics "github.com/arran4/golang-ical"

	"github.com/julianstephens/go-utils/generic"
	"github.com/julianstephens/go-utils/helpers"
	"github.com/julianstephens/liturgical-time-index/internal/plan"
)

// ICS takes a slice of FormattedEntry and an output path, and writes the entries to an ICalendar file.
func ICS(entries []plan.FormattedEntry, outputPath string) error {
	cal := ics.NewCalendar()
	cal.SetMethod(ics.MethodPublish)
	now := time.Now().UTC()

	for _, entry := range entries {
		formattedDate := strings.ReplaceAll(entry.Key.Date, "-", "")
		event := cal.AddEvent(fmt.Sprintf("%s-%d-%s", entry.Key.Season, entry.Key.SeasonWeek, formattedDate))
		event.SetSummary(entry.Cue)
		event.SetDescription(fmt.Sprintf("%s\n\nRb references:\n%s", entry.Cue, formatRbRefs(entry.Rb)))

		event.SetDtStampTime(now)
		event.SetProperty(ics.ComponentPropertyDtStart, formattedDate)
		event.SetProperty(ics.ComponentPropertyDtEnd, formattedDate)
	}

	serialized := cal.Serialize(ics.WithNewLineWindows)
	if err := helpers.AtomicFileWrite(outputPath, []byte(serialized)); err != nil {
		return &OutputError{
			Message: generic.Ptr("failed to write ICS file"),
			Err:     ErrSerializationFailed,
			Cause:   err,
		}
	}

	return nil
}

func formatRbRefs(rbRefs []plan.RbRef) string {
	var formatted string

	for _, ref := range rbRefs {
		formatted += fmt.Sprintf("- %s\n", ref.String())
	}

	return formatted
}
