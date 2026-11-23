---
description: Do exactly one parity iteration
agent: build
---

Read `GOAL.md`, `.agents/STATE.md`, `parity.csv`, `docs/CONTRACT.md`.

1. Pick the highest `MUST` row with status `todo`.
2. Overwrite the "Current token" block in `.agents/STATE.md`: row id, files, verification command.
3. Implement the smallest change that makes that row pass.
4. Run `./scripts/check.sh`. Paste the tail into `.agents/STATE.md`.
5. Green: set the row to `pass` in `parity.csv`, `git add -A && git commit -m "iter: <row> <what>"`.
   Red: fix twice. Still red: revert and split the row into two smaller rows.

Then run `./scripts/dod.sh`. If it exits 0, write `DONE.md` and stop.
