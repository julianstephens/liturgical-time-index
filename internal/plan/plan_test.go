package plan_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/julianstephens/liturgical-time-index/internal/plan"
)

var testDataDir = "testdata"

func TestLoadPlan_ValidPlan(t *testing.T) {
	planPath := filepath.Join(testDataDir, "valid_plan.yml")

	p, err := plan.LoadPlan(planPath)
	if err != nil {
		t.Fatalf("LoadPlan failed for valid plan: %v", err)
	}

	if p == nil {
		t.Fatal("LoadPlan returned nil for valid plan")
	}

	if p.Version == 0 {
		t.Error("Plan version should be set")
	}

	if p.Work == "" {
		t.Error("Plan work should be set")
	}
}

func TestLoadPlan_FileNotFound(t *testing.T) {
	planPath := filepath.Join(testDataDir, "nonexistent_plan.yml")

	p, err := plan.LoadPlan(planPath)
	if err == nil {
		t.Errorf("LoadPlan should error for non-existent file, got nil. Result: %v", p)
	}

	if p != nil {
		t.Errorf("LoadPlan should return nil on error, got: %v", p)
	}
}

func TestLoadPlan_MalformedYAML(t *testing.T) {
	planPath := filepath.Join(testDataDir, "malformed_plan.yml")

	p, err := plan.LoadPlan(planPath)
	if err == nil {
		t.Errorf("LoadPlan should error for malformed YAML, got nil. Result: %v", p)
	}

	if p != nil {
		t.Errorf("LoadPlan should return nil on error, got: %v", p)
	}
}

func TestValidatePlan_ValidPlan(t *testing.T) {
	planPath := filepath.Join(testDataDir, "valid_plan.yml")

	err := plan.LoadAndValidatePlan(planPath)
	if err != nil {
		t.Fatalf("LoadAndValidatePlan failed for valid plan: %v", err)
	}
}

func TestValidatePlan_IncompleteSeason(t *testing.T) {
	planPath := filepath.Join(testDataDir, "incomplete_season_plan.yml")

	err := plan.LoadAndValidatePlan(planPath)
	if err == nil {
		t.Error("LoadAndValidatePlan should error for season without all weekdays and no fallback")
	}
}

func TestValidatePlan_DuplicateWeekday(t *testing.T) {
	planPath := filepath.Join(testDataDir, "duplicate_weekday_plan.yml")

	err := plan.LoadAndValidatePlan(planPath)
	if err == nil {
		t.Error("LoadAndValidatePlan should error for duplicate weekday in season")
	}
}

func TestValidatePlan_TooManyWeekdays(t *testing.T) {
	planPath := filepath.Join(testDataDir, "too_many_weekdays_plan.yml")

	err := plan.LoadAndValidatePlan(planPath)
	if err == nil {
		t.Error("LoadAndValidatePlan should error for more than 7 weekdays in season")
	}
}

func TestValidatePlan_NoWeekdaysNoFallback(t *testing.T) {
	planPath := filepath.Join(testDataDir, "no_weekdays_no_fallback_plan.yml")

	err := plan.LoadAndValidatePlan(planPath)
	if err == nil {
		t.Error("LoadAndValidatePlan should error for season with no weekdays and no fallback")
	}
}

func TestValidatePlan_PartialWeekdaysWithFallback(t *testing.T) {
	planPath := filepath.Join(testDataDir, "partial_weekdays_with_fallback_plan.yml")

	err := plan.LoadAndValidatePlan(planPath)
	if err != nil {
		t.Fatalf("LoadAndValidatePlan should succeed for partial weekdays with fallback: %v", err)
	}
}

func TestValidatePlan_AllWeekdaysNoFallback(t *testing.T) {
	planPath := filepath.Join(testDataDir, "all_weekdays_no_fallback_plan.yml")

	err := plan.LoadAndValidatePlan(planPath)
	if err != nil {
		t.Fatalf("LoadAndValidatePlan should succeed for all weekdays without fallback: %v", err)
	}
}

func init() {
	// Create valid plan
	createTestFileIfNotExists("valid_plan.yml", `
version: 1
work: "Rule of Saint Benedict"
witness: "latin"
defaults:
  responsory: "RB 1.1"
seasons:
  Advent:
    weekdays:
      mon: { responsory: "RB 2.1" }
      tue: { responsory: "RB 2.2" }
      wed: { responsory: "RB 2.3" }
      thu: { responsory: "RB 2.4" }
      fri: { responsory: "RB 2.5" }
      sat: { responsory: "RB 2.6" }
      sun: { responsory: "RB 3.1" }
  Christmastide:
    weekdays:
      mon: { responsory: "RB 4.1" }
      tue: { responsory: "RB 4.2" }
      wed: { responsory: "RB 4.3" }
      thu: { responsory: "RB 4.4" }
      fri: { responsory: "RB 4.5" }
      sat: { responsory: "RB 4.6" }
      sun: { responsory: "RB 5.1" }
  Epiphanytide:
    weekdays:
      mon: { responsory: "RB 6.1" }
      tue: { responsory: "RB 6.2" }
      wed: { responsory: "RB 6.3" }
      thu: { responsory: "RB 6.4" }
      fri: { responsory: "RB 6.5" }
      sat: { responsory: "RB 6.6" }
      sun: { responsory: "RB 7.1" }
  Lent:
    weekdays:
      mon: { responsory: "RB 8.1" }
      tue: { responsory: "RB 8.2" }
      wed: { responsory: "RB 8.3" }
      thu: { responsory: "RB 8.4" }
      fri: { responsory: "RB 8.5" }
      sat: { responsory: "RB 8.6" }
      sun: { responsory: "RB 9.1" }
  "Paschal Triduum":
    weekdays:
      mon: { responsory: "RB 10.1" }
      tue: { responsory: "RB 10.2" }
      wed: { responsory: "RB 10.3" }
      thu: { responsory: "RB 10.4" }
      fri: { responsory: "RB 10.5" }
      sat: { responsory: "RB 10.6" }
      sun: { responsory: "RB 11.1" }
  Easter:
    weekdays:
      mon: { responsory: "RB 12.1" }
      tue: { responsory: "RB 12.2" }
      wed: { responsory: "RB 12.3" }
      thu: { responsory: "RB 12.4" }
      fri: { responsory: "RB 12.5" }
      sat: { responsory: "RB 12.6" }
      sun: { responsory: "RB 13.1" }
  "Ordinary Time":
    weekdays:
      mon: { responsory: "RB 14.1" }
      tue: { responsory: "RB 14.2" }
      wed: { responsory: "RB 14.3" }
      thu: { responsory: "RB 14.4" }
      fri: { responsory: "RB 14.5" }
      sat: { responsory: "RB 14.6" }
      sun: { responsory: "RB 15.1" }
`)

	// Create plan with missing seasons
	createTestFileIfNotExists("missing_seasons_plan.yml", `
version: 1
work: "Rule of Saint Benedict"
witness: "latin"
defaults:
  responsory: "RB 1.1"
seasons: {}
`)

	// Create plan with incomplete season
	createTestFileIfNotExists("incomplete_season_plan.yml", `
version: 1
work: "Rule of Saint Benedict"
witness: "latin"
defaults:
  responsory: "RB 1.1"
seasons:
  Advent:
    weekdays:
      mon: { responsory: "RB 2.1" }
`)

	// Create plan with duplicate weekday
	createTestFileIfNotExists("duplicate_weekday_plan.yml", `
version: 1
work: "Rule of Saint Benedict"
witness: "latin"
defaults:
  responsory: "RB 1.1"
seasons:
  Advent:
    weekdays:
      mon: { responsory: "RB 2.1" }
      mon: { responsory: "RB 2.2" }
`)

	// Create plan with too many weekdays
	createTestFileIfNotExists("too_many_weekdays_plan.yml", `
version: 1
work: "Rule of Saint Benedict"
witness: "latin"
defaults:
  responsory: "RB 1.1"
seasons:
  Advent:
    weekdays:
      mon: { responsory: "RB 2.1" }
      tue: { responsory: "RB 2.2" }
      wed: { responsory: "RB 2.3" }
      thu: { responsory: "RB 2.4" }
      fri: { responsory: "RB 2.5" }
      sat: { responsory: "RB 2.6" }
      sun: { responsory: "RB 3.1" }
      extra: { responsory: "RB 3.2" }
`)

	// Create plan with invalid RB reference
	createTestFileIfNotExists("invalid_rb_reference_plan.yml", `
version: 1
work: "Rule of Saint Benedict"
witness: "latin"
defaults:
  responsory: "RB 1.1"
seasons:
  Advent:
    weekdays:
      mon: { responsory: "RB 74.1" }
      tue: { responsory: "RB 2.2" }
      wed: { responsory: "RB 2.3" }
      thu: { responsory: "RB 2.4" }
      fri: { responsory: "RB 2.5" }
      sat: { responsory: "RB 2.6" }
      sun: { responsory: "RB 3.1" }
`)

	// Create plan with no weekdays and no fallback
	createTestFileIfNotExists("no_weekdays_no_fallback_plan.yml", `
version: 1
work: "Rule of Saint Benedict"
witness: "latin"
defaults:
  responsory: "RB 1.1"
seasons:
  Advent:
    weekdays: {}
`)

	// Create plan with partial weekdays and fallback
	createTestFileIfNotExists("partial_weekdays_with_fallback_plan.yml", `
version: 1
work: "Rule of Saint Benedict"
witness: "latin"
defaults:
  responsory: "RB 1.1"
seasons:
  Advent:
    weekdays:
      mon: { responsory: "RB 2.1" }
      tue: { responsory: "RB 2.2" }
      wed: { responsory: "RB 2.3" }
    fallback: { responsory: "RB 1.1" }
`)

	// Create plan with all weekdays and no fallback
	createTestFileIfNotExists("all_weekdays_no_fallback_plan.yml", `
version: 1
work: "Rule of Saint Benedict"
witness: "latin"
defaults:
  responsory: "RB 1.1"
seasons:
  Advent:
    weekdays:
      mon: { responsory: "RB 2.1" }
      tue: { responsory: "RB 2.2" }
      wed: { responsory: "RB 2.3" }
      thu: { responsory: "RB 2.4" }
      fri: { responsory: "RB 2.5" }
      sat: { responsory: "RB 2.6" }
      sun: { responsory: "RB 3.1" }
`)

	// Create malformed YAML
	createTestFileIfNotExists("malformed_plan.yml", `
version: 1
work: "Rule of Saint Benedict
witness: "latin"
`)
}

func createTestFileIfNotExists(filename, content string) {
	path := filepath.Join(testDataDir, filename)
	if _, err := os.Stat(path); err == nil {
		return // file already exists
	}
	if err := os.MkdirAll(testDataDir, 0750); err != nil {
		// ignore errors - directory may already exist
	}
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		// ignore errors - file may already exist or be created elsewhere
	}
}
