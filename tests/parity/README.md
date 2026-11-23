# Parity tests

One file per `MUST` row in `parity.csv`. Named `f02.spec.ts` etc.

Each test drives the real UI and asserts the observable behaviour of the row. Tests are added by
the agent, never weakened, and `spec-auditor` specifically tries to prove they can pass without the
feature working.

Run all: `npx playwright test`
