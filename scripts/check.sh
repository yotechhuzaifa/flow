#!/usr/bin/env bash
# Fast gate for the loop. Target: under 60 seconds.
set -uo pipefail
cd "$(dirname "$0")/.."

fail=0
run() {
  local name="$1"; shift
  printf '\n== %s ==\n' "$name"
  if "$@" 2>&1 | tail -n 25; then printf '  ok   %s\n' "$name"
  else printf '  FAIL %s\n' "$name"; fail=1; fi
}

# research output must stay machine-readable and cited
if [ -f spec/features.jsonl ]; then
  run "spec jsonl" python3 -c '
import json
n=0
for i,l in enumerate(open("spec/features.jsonl"),1):
    l=l.strip()
    if not l: continue
    d=json.loads(l)
    assert "id" in d and "observable" in d, f"line {i}: missing id/observable"
    assert d.get("evidence"), f"line {i}: no evidence"
    n+=1
print(f"{n} features")'
fi

if command -v go >/dev/null 2>&1; then
  run "gofmt" bash -c 'test -z "$(gofmt -l internal cmd)"'
  run "go vet"  go vet -tags server ./...
  run "go test" go test -tags server ./... -count=1 -timeout 120s
  if go build -tags server ./... >/dev/null 2>&1; then
    echo "  ok   go build (server, no webview)"
  else
    echo "  FAIL go build (server)"; fail=1
  fi
  # desktop build needs the native webview; reported, does not block
  if go build ./cmd/app/ >/tmp/gocmd.log 2>&1; then
    echo "  ok   go build (desktop)"
  else
    echo "  note desktop build needs gtk4/webkitgtk — see docs/platform-setup.md"
  fi
else
  echo "== go =="; echo "  skip: no go toolchain on PATH"
fi

if [ -f frontend/package.json ]; then
  if [ ! -d frontend/node_modules ]; then
    if [ -f frontend/package-lock.json ]; then
      run "npm ci" npm --prefix frontend ci --no-audit --no-fund
    else
      run "npm install" npm --prefix frontend install --no-audit --no-fund
    fi
  fi
  run "typecheck" npm --prefix frontend run -s typecheck
  run "frontend build" npm --prefix frontend run -s build
  # embedded assets must not go stale
  if command -v rsync >/dev/null 2>&1; then
    if ! rsync -rc --delete --dry-run frontend/dist/ internal/webui/dist/ >/dev/null 2>&1; then
      echo "  note internal/webui/dist is stale: cp -r frontend/dist internal/webui/"
    fi
  fi
fi

echo
[ "$fail" -eq 0 ] && echo "CHECK: PASS" || echo "CHECK: FAIL"
exit "$fail"
