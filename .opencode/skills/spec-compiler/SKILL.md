---
name: spec-compiler
description: Turn spec/features.jsonl into parity proposals, a test plan, and an audit. Load after /research, before merging a parity row, or when spec-auditor needs a target list.
autoinvoke: false
---

# Spec compiler

You are the bridge between research and the build loop. You do not observe the product and
you do not write app code. You make research executable.

Read `docs/RESEARCH.md` first. It explains the six-layer pipeline and the OSS references.

## When to use

- `/research` has appended to `spec/features.jsonl`.
- A new feature needs a parity row.
- A `MUST` row is about to be marked `pass`.
- `spec-auditor` needs the list of observables and evidence.

## Commands

```bash
node scripts/spec.mjs check      # validate JSONL, ids, priorities, evidence paths
node scripts/spec.mjs propose    # spec/parity.proposed.csv + spec/parity.map.json
node scripts/spec.mjs plan       # spec/TEST-PLAN.md
node scripts/spec.mjs audit      # spec/AUDIT.md; exit 1 on a broken gate
node scripts/spec.mjs all        # check + propose + plan + audit
```

## Workflow

1. Run `check`. Fix every failure at the source in `features.jsonl`. Do not weaken the check.
2. Run `propose`. It assigns stable `p01`-style ids and never edits `parity.csv`.
3. Stop. A human reviews `parity.proposed.csv` and merges accepted rows into `parity.csv`
   with a `[contract]` commit. You do not edit the frozen queue.
4. Run `plan`. The build agent works from `spec/TEST-PLAN.md`.
5. Run `audit` before marking a row `pass`. Treat a `FAIL` as a blocker.
6. Spawn `spec-auditor` for the row. If it can make the test pass without the feature, the
   test is wrong. Fix the test or the feature; never the auditor.

## Feature line

`spec/features.jsonl` is one JSON object per line. The schema is
`spec/feature.schema.json`. Minimum shape:

```json
{"id":"edge.create","trigger":"dragging from an output port to an input port","observable":"an edge appears only when the drop lands on an input port","priority":"must","evidence":["traces/connect/network.jsonl#req-4","traces/connect/steps.jsonl#7"]}
```

Rules:

- `id`: lowercase, stable, `^[a-z0-9][a-z0-9._-]*$`.
- `trigger` and `observable`: behaviour, never implementation.
- `priority`: one of `must`, `should`, `could`, `wont`.
- `evidence`: at least one real `traces/…` path. Doc-only claims belong in `spec/INTEL.md`,
  not here.
- `wire`, `errors`, `preconditions`, `unsure`: optional string arrays.
- No function names, CSS classes, internal ids, or copied microcopy.

## Hard rules

- Never edit `parity.csv`, `docs/CONTRACT.md`, `scripts/check.sh`, `scripts/dod.sh`, or
  `GOAL.md`. They are frozen.
- Never generate a passing test stub. The compiler writes a plan; assertions are written by
  the build agent and attacked by `spec-auditor`.
- If a feature needs a new entity, event, or field, write `spec/CONTRACT-PROPOSAL.md` with the
  feature id, the proposed change, and the evidence. Stop and wait for a `[contract]` commit.
- Keep output deterministic. No timestamps, no random ids.
- If `features.jsonl` is empty, say so and stop. Do not invent features.

## Output contract

| File | Meaning |
|---|---|
| `spec/parity.proposed.csv` | proposed rows, same six columns as `parity.csv` |
| `spec/parity.map.json` | row id → feature id; stable once assigned |
| `spec/TEST-PLAN.md` | per-feature observable, evidence, and test checklist |
| `spec/AUDIT.md` | `PASS` / `FAIL` / `WARN` with the reason |
