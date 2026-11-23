# STATE

Append-only. Newest at the bottom.

## Current token
- Task: (none)
- Files:
- Verify with:

## Log
| # | Row | Check | Result | Note |
|---|-----|-------|--------|------|
| 0 | scaffold | ./scripts/check.sh | — | initialised |
| 1 | setup audit | ./scripts/check.sh | fail → pass | missing `internal/webui/dist` broke `go build`; vite alias pointed at a non-existent stub; `@playwright/test` was not a dependency. Dependency added: `@playwright/test` (the parity suite in `tests/parity` imports it; no runtime cost). |
