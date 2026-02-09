package compile

import (
	"github.com/julianstephens/liturgical-time-index/internal/calendar"
	"github.com/julianstephens/liturgical-time-index/internal/plan"
)

// Compile compiles a plan for a given day key, applying defaults and fallbacks as necessary.
func Compile(key calendar.DayKey, p plan.Plan) (*plan.FormattedEntry, error) {
	defaults := p.Defaults
	formattedDefaults, err := defaults.Validate()
	if err != nil {
		return nil, err
	}
	defaultEntry := &plan.FormattedEntry{
		Key:  key,
		Cue:  formattedDefaults.Cue,
		Rb:   formattedDefaults.Rb,
		Tags: formattedDefaults.Tags,
	}

	seasonPlan, ok := p.Seasons[string(key.Season)]
	if !ok {
		return defaultEntry, nil
	}

	weekday, ok := seasonPlan.Weekdays[string(key.Weekday)]
	if !ok {
		if seasonPlan.Fallback != nil {
			fallback, err := seasonPlan.Fallback.Validate()
			if err != nil {
				return nil, err
			}

			return &plan.FormattedEntry{
				Key:  key,
				Cue:  fallback.Cue,
				Rb:   fallback.Rb,
				Tags: fallback.Tags,
			}, nil
		}

		return defaultEntry, nil
	}

	formattedEntry, err := weekday.Validate()
	if err != nil {
		return nil, err
	}
	formattedEntry.Key = key

	return formattedEntry, nil
}
