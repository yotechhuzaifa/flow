#!/usr/bin/env bash
# Definition of done. The loop stops only when this exits 0. Never let an agent edit it.
set -uo pipefail
cd "$(dirname "$0")/.."

fail=0
say() { printf '  %s %s\n' "$1" "$2"; }

bash ./scripts/check.sh >/tmp/dod.log 2>&1 \
  && say "ok  " "check.sh" \
  || { say "FAIL" "check.sh"; tail -n 20 /tmp/dod.log; fail=1; }

if [ -f parity.csv ]; then
  missing=$(awk -F',' 'NR>1 && $3=="MUST" && $6!="pass" {print $1}' parity.csv | tr -d ' \r')
  if [ -z "$missing" ]; then say "ok  " "all MUST rows pass"
  else say "FAIL" "MUST rows open: $(echo $missing | tr '\n' ' ')"; fail=1; fi
fi

if git rev-parse --git-dir >/dev/null 2>&1; then
  for f in GOAL.md docs/CONTRACT.md scripts/check.sh scripts/dod.sh; do
    git diff --quiet HEAD -- "$f" 2>/dev/null || { say "FAIL" "frozen file changed: $f"; fail=1; }
  done
  git grep -qE 'TODO\(agent\)|FIXME\(agent\)' -- . 2>/dev/null && { say "FAIL" "agent TODO markers left"; fail=1; }
fi

[ "$fail" -eq 0 ] && { echo "DOD: DONE"; exit 0; } || { echo "DOD: NOT DONE"; exit 1; }
