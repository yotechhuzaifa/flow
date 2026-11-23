#!/usr/bin/env bash
# One-time setup. Free tools only. Nothing large is written into the workspace.
set -uo pipefail
cd "$(dirname "$0")/.."

say() { printf '\n== %s ==\n' "$*"; }
ok()  { printf '  ok   %s\n' "$*"; }
bad() { printf '  FAIL %s\n' "$*"; }

say "base tools"
for t in git node npx python3; do
  command -v "$t" >/dev/null 2>&1 && ok "$t" || bad "$t missing"
done

say "go toolchain (kept outside the workspace)"
if command -v go >/dev/null 2>&1; then ok "go $(go version | awk '{print $3}')"
else echo "  run: bash scripts/toolchain.sh   (installs to ~/.local/share/gotoolchain)"; fi

say "wails3 cli"
if command -v wails3 >/dev/null 2>&1; then ok "wails3 $(wails3 version 2>/dev/null | head -1)"
else
  echo "  go install github.com/wailsapp/wails/v3/cmd/wails3@latest"
  echo "  needs pkg-config + GTK4/WebKitGTK headers; the headless path does not"
fi

say "opencode"
if command -v opencode >/dev/null 2>&1; then
  opencode auth list 2>/dev/null | grep -qi opencode && ok "console/go connected" \
    || echo "  run: opencode auth login   (OpenCode Console, Go plan)"
else bad "opencode missing — https://opencode.ai/docs"; fi

say "mcp servers (npx, nothing installed globally)"
echo "  playwright: fetched on demand by npx @playwright/mcp"
echo "  context7:   remote, no key needed"
opencode mcp list 2>/dev/null || true

say "frontend deps"
if [ -d frontend/node_modules ]; then ok "frontend/node_modules present"
elif [ -f frontend/package-lock.json ]; then echo "  npm --prefix frontend ci   (gitignored, ~68 MB — your machine, not a snapshot)"
else echo "  npm --prefix frontend install"; fi

say "parity suite (root package)"
if [ -d node_modules/@playwright ]; then ok "playwright runner installed"
else echo "  npm install   (root, for @playwright/test; gitignored)"; fi
if ls "$HOME/.cache/ms-playwright"/chromium-* >/dev/null 2>&1; then ok "chromium browser present"
else echo "  npx playwright install chromium"; fi

say "workspace size"
if git rev-parse --git-dir >/dev/null 2>&1; then
  tracked=$(git ls-files 2>/dev/null | wc -l | tr -d ' ')
  echo "  tracked files: ${tracked:-0} (keep node_modules, dist and caches out)"
fi

say "git"
if git rev-parse --git-dir >/dev/null 2>&1; then ok "repo"
else git init -q && ok "git init (needed: dod.sh guards frozen files, loop.sh detects stalls)"; fi
git config user.email >/dev/null 2>&1 || echo "  set git user.email and user.name before the loop commits"

say "gates"
bash scripts/check.sh >/dev/null 2>&1 && ok "check.sh green" || echo "  check.sh red (expected before deps/build)"
bash scripts/dod.sh   >/dev/null 2>&1 && ok "dod.sh green"   || echo "  dod.sh red (expected: MUST rows open)"

cat <<'NEXT'

Next:
  bash scripts/toolchain.sh     # if go is missing
  npm install && npm --prefix frontend ci
  make check                    # headless gate
  make e2e                      # parity suite (starts Vite)
  make dev                      # native window — needs wails3, build/ assets and GTK4/WebKitGTK
  opencode  →  /start           # or ./scripts/loop.sh for unattended
NEXT
