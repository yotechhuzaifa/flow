---
description: Attack the current parity tests with the spec-auditor. No app code.
agent: spec-auditor
---

Read `spec/TEST-PLAN.md`, `parity.csv`, and `tests/parity/`.

Target: $ARGUMENTS

For every open `MUST` row, try to make its test pass without the feature working:

1. Stub the UI, hard-code the expected value, skip an assertion, or add a test-only path.
2. If the test still passes, the test is weak. Add a failing test under `tests/parity/` that
   proves the gap.
3. Do not edit app code. Do not weaken an existing test. Do not edit frozen files.
4. Report the row id, the trick that passed, and the failing test you added.

If `spec/parity.proposed.csv` has rows that a human merged, audit those too. Stop after the
first three confirmed gaps and report; do not rewrite the suite.
