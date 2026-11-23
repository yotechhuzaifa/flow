# GOAL

A desktop app where a user can build and run a flow on a canvas.

- Select node types from a palette and place them on the canvas
- Connect node ports with edges
- Save a flow and reopen it after restart
- Run a flow and watch each node's status change live

## Out of scope
Multi-user, cloud sync, accounts, sharing, imports from other products, plugin marketplace, mobile,
themes, i18n, keyboard-only accessibility work beyond basics.

## Definition of Done
`./scripts/dod.sh` exits 0.

It requires:
- `go build ./...`, `go vet ./...`, `go test ./...` clean
- `npm run typecheck` clean
- every `MUST` row in `parity.csv` is `pass`
- the frozen files are unmodified

## Guardrails
- Max 25 iterations, 2 stalls (no tree change) → stop
- One row per iteration
- Commit every green iteration
