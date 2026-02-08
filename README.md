# LTI — Liturgical Time Index (v0)

Local-first CLI that compiles a Roman-season liturgical calendar into daily practice entries
(cue line + RB citation references). Outputs ICS and Markdown.

## What this is

- A compiler: (year + tradition + plan.yaml) -> daily entries -> artifacts (ICS/MD)
- Strict about references and plan coverage

## What this is not (v0)

- Saints calendar, lectionary texts, prayers, tracking, accounts, UI

## Quick start

### Build artifacts

```bash
go run ./cmd/lti build --year 2026 --tradition roman \
  --plan data/rb_plan.yaml \
  --out-ics out/2026.ics \
  --out-md out/2026.md
```

### Show a date

```sh
go run ./cmd/lti today --date 2026-02-08 --tradition roman --plan data/rb_plan.yaml
```

### Validate plan

```sh
go run ./cmd/lti validate --plan data/rb_plan.yaml
```

### Plan editing

Edit `data/rb_plan.yaml`. You can set per-season weekday overrides and fallbacks.
RB references are validated (Prologue and chapter/verse forms).

## Notes

Roman season boundaries are computed with:

- Easter (Gregorian)

- Ash Wednesday (Easter - 46 days)

- Triduum (Holy Thu -> Holy Sat)

- Pentecost (Easter + 49)

- Advent start approximated as the Sunday on or after Nov 27

- Week numbering is a simple 7-day index from each season start.
