package plan

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"

	"github.com/julianstephens/go-utils/generic"
)

type Plan struct {
	Version  int                   `yaml:"version"`
	Work     string                `yaml:"work"`
	Witness  string                `yaml:"witness"`
	Defaults PlanEntry             `yaml:"defaults"`
	Seasons  map[string]SeasonPlan `yaml:"seasons"`
}

type SeasonPlan struct {
	Weekdays map[string]PlanEntry `yaml:"weekdays"`
	Fallback *PlanEntry           `yaml:"fallback"`
}

// LoadPlan reads a YAML file from the given path and unmarshals it into a Plan struct.
func LoadPlan(planPath string) (*Plan, error) {
	var plan Plan

	bytes, err := os.ReadFile(filepath.Clean(planPath))
	if err != nil {
		return nil, &PlanError{
			Err:   ErrPlanFileRead,
			Cause: err,
		}
	}
	if err = yaml.Unmarshal(bytes, &plan); err != nil {
		return nil, &PlanError{
			Err:   ErrParsePlanFailed,
			Cause: err,
		}
	}

	return &plan, nil
}

// ValidatePlan loads a plan from the specified path and validates its contents.
// It checks that the plan can be loaded successfully and that all entries conform to expected formats and rules.
func ValidatePlan(planPath string) error {
	plan, err := LoadPlan(planPath)
	if err != nil {
		return err
	}

	if err := plan.validate(); err != nil {
		return err
	}

	return nil
}

// validate checks the structure and content of the Plan to ensure it meets the required criteria.
// It verifies that each season has valid weekday entries, that there are no duplicate weekdays,
// and that all RB references are properly formatted.
func (p *Plan) validate() error {
	if _, err := p.Defaults.validate(); err != nil {
		return &PlanError{
			Message: generic.Ptr("invalid default plan entry"),
			Err:     ErrInvalidPlanEntry,
		}
	}

	for seasonName, seasonPlan := range p.Seasons {
		weekdaysCovered := make(map[string]bool)
		for weekday, entry := range seasonPlan.Weekdays {
			if weekdaysCovered[weekday] {
				return &PlanError{
					Message: generic.Ptr("duplicate weekday " + weekday + " in season " + seasonName),
					Err:     ErrInvalidPlanEntry,
				}
			}
			weekdaysCovered[weekday] = true
			if _, err := entry.validate(); err != nil {
				return &PlanError{
					Message: generic.Ptr("invalid plan entry in season " + seasonName + " for weekday " + weekday),
					Err:     ErrInvalidPlanEntry,
				}
			}
		}
		if seasonPlan.Fallback != nil {
			if _, err := seasonPlan.Fallback.validate(); err != nil {
				return &PlanError{
					Message: generic.Ptr("invalid fallback plan entry in season " + seasonName),
					Err:     ErrInvalidPlanEntry,
				}
			}
		}
		if len(weekdaysCovered) == 0 && seasonPlan.Fallback == nil {
			return &PlanError{
				Message: generic.Ptr("season " + seasonName + " must have at least one weekday entry or a fallback"),
				Err:     ErrInvalidPlanEntry,
			}
		}
		if len(weekdaysCovered) > 7 {
			return &PlanError{
				Message: generic.Ptr("season " + seasonName + " cannot have more than 7 weekday entries"),
				Err:     ErrInvalidPlanEntry,
			}
		}
		if len(weekdaysCovered) < 7 && seasonPlan.Fallback == nil {
			return &PlanError{
				Message: generic.Ptr(
					"season " + seasonName + " must have a fallback if it does not cover all 7 weekdays",
				),
				Err: ErrInvalidPlanEntry,
			}
		}
		if len(weekdaysCovered) == 7 && seasonPlan.Fallback == nil {
			if err := validateWeekdays(generic.Keys(seasonPlan.Weekdays)); err != nil {
				return &PlanError{
					Message: generic.Ptr("invalid weekday in season " + seasonName + " and no fallback provided"),
					Err:     ErrInvalidPlanEntry,
				}
			}
		}
	}
	return nil
}

func validateWeekdays(weekdays []string) error {
	if !generic.Contains(weekdays, "mon") && !generic.Contains(weekdays, "monday") {
		return &PlanError{
			Message: generic.Ptr("missing required weekday 'mon' or 'monday'"),
			Err:     ErrInvalidPlanEntry,
		}
	}
	if !generic.Contains(weekdays, "tue") && !generic.Contains(weekdays, "tuesday") {
		return &PlanError{
			Message: generic.Ptr("missing required weekday 'tue' or 'tuesday'"),
			Err:     ErrInvalidPlanEntry,
		}
	}
	if !generic.Contains(weekdays, "wed") && !generic.Contains(weekdays, "wednesday") {
		return &PlanError{
			Message: generic.Ptr("missing required weekday 'wed' or 'wednesday'"),
			Err:     ErrInvalidPlanEntry,
		}
	}
	if !generic.Contains(weekdays, "thu") && !generic.Contains(weekdays, "thursday") {
		return &PlanError{
			Message: generic.Ptr("missing required weekday 'thu' or 'thursday'"),
			Err:     ErrInvalidPlanEntry,
		}
	}
	if !generic.Contains(weekdays, "fri") && !generic.Contains(weekdays, "friday") {
		return &PlanError{
			Message: generic.Ptr("missing required weekday 'fri' or 'friday'"),
			Err:     ErrInvalidPlanEntry,
		}
	}
	if !generic.Contains(weekdays, "sat") && !generic.Contains(weekdays, "saturday") {
		return &PlanError{
			Message: generic.Ptr("missing required weekday 'sat' or 'saturday'"),
			Err:     ErrInvalidPlanEntry,
		}
	}
	if !generic.Contains(weekdays, "sun") && !generic.Contains(weekdays, "sunday") {
		return &PlanError{
			Message: generic.Ptr("missing required weekday 'sun' or 'sunday'"),
			Err:     ErrInvalidPlanEntry,
		}
	}
	return nil
}
