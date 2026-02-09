package output

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"

	md "github.com/nao1215/markdown"

	"github.com/julianstephens/go-utils/generic"
	"github.com/julianstephens/liturgical-time-index/internal/plan"
)

// Markdown takes a slice of FormattedEntry and an output path, and writes the entries to a Markdown file in a tabular format.
func Markdown(entries []plan.FormattedEntry, outputPath string) (retErr error) {
	f, err := os.OpenFile(filepath.Clean(outputPath), os.O_CREATE|os.O_WRONLY|os.O_EXCL, 0600)
	if err != nil {
		return &OutputError{
			Message: generic.Ptr("failed to open Markdown file"),
			Err:     ErrOpenOutputFileFailed,
			Cause:   err,
		}
	}
	defer func() {
		if closeErr := f.Close(); closeErr != nil && retErr == nil {
			retErr = &OutputError{
				Message: generic.Ptr("failed to close Markdown file"),
				Err:     ErrCloseOutputFileFailed,
				Cause:   closeErr,
			}
		}
	}()

	if err := md.NewMarkdown(f).Table(md.TableSet{
		Header: []string{"Date", "Season", "Season Week", "Weekday", "Cue", "RB References"},
		Rows: generic.Map(entries, func(entry plan.FormattedEntry) []string {
			return []string{
				entry.Key.Date,
				entry.Key.Season.String(),
				strconv.Itoa(entry.Key.SeasonWeek),
				entry.Key.Weekday.String(),
				entry.Cue,
				strings.Join(generic.Map(entry.Rb, func(ref plan.RbRef) string { return ref.String() }), "; "),
			}
		}),
	}).Build(); err != nil {
		return &OutputError{
			Message: generic.Ptr("failed to write Markdown file"),
			Err:     ErrSerializationFailed,
			Cause:   err,
		}
	}

	return nil
}
