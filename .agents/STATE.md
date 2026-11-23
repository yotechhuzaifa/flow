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
| 2 | setup audit | `check.sh` · `npx playwright test` · `opencode debug config` | pass | check.sh PASS; f01 passes against Vite; config parses with all three MCP servers. Initial commit `2c7eeca` lets the `dod.sh` frozen-file guard compare against HEAD. `dod.sh` now fails only on the 12 open MUST rows. |

## Known gaps (setup audit, not fixed here)
- Desktop build needs the Wails v3 `build/` assets, which are not vendored. `wails3` cannot install on a
  host without `pkg-config`/GTK4. Headless (`cmd/server`, Vite, Playwright) is fully usable.
- `frontend/src/shared/api/client.ts` talks only to generated Wails bindings; there is no HTTP client for
  server mode yet. Later rows that need persistence/live runs in the browser must add it.
- The generated binding path `@bindings/flow/api` is a declared placeholder; verify it after
  `wails3 generate bindings` on a machine with the WebView toolchain.
