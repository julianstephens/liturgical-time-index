package plan

import "github.com/julianstephens/liturgical-time-index/internal/calendar"

type PlanEntry struct {
	Cue  string    `yaml:"cue"`
	Rb   []string  `yaml:"rb"`
	Tags *[]string `yaml:"tags,omitempty"`
}

type FormattedEntry struct {
	Key  calendar.DayKey `yaml:"key"`
	Cue  string          `yaml:"cue"`
	Rb   []RbRef         `yaml:"rb"`
	Tags *[]string       `yaml:"tags,omitempty"`
}

func (e *PlanEntry) Validate() (*FormattedEntry, error) {
	refs := make([]RbRef, len(e.Rb))
	for i, rbRef := range e.Rb {
		ref, err := NewRbRef(rbRef)
		if err != nil {
			return nil, err
		}
		refs[i] = *ref
	}
	return &FormattedEntry{
		Cue:  e.Cue,
		Rb:   refs,
		Tags: e.Tags,
	}, nil
}
