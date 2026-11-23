# embedded frontend

Generated. Do not edit.

`npm --prefix frontend run build` writes `frontend/dist` and syncs it here through
`scripts/sync-webui.mjs`. `go:embed all:dist` needs at least one file at compile time, so this
README keeps `go build`/`go test` working on a fresh clone before the first frontend build.

Everything else in this directory is ignored by git.
