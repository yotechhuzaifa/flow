#!/usr/bin/env bash
# Outer loop: one fresh opencode turn per iteration. Survives crashes, enforces hard stops.
set -uo pipefail
cd "$(dirname "$0")/.."

MAX_ITERS="${MAX_ITERS:-25}"
MAX_STALL="${MAX_STALL:-2}"
mkdir -p .agents/logs

stall=0; prev=""
for i in $(seq 1 "$MAX_ITERS"); do
  echo "── iteration $i/$MAX_ITERS ──"

  # stall detection works with or without git
  if git rev-parse --git-dir >/dev/null 2>&1; then
    # include HEAD: a committed iteration cleans the tree, but the commit still counts as progress
    cur=$( { git rev-parse --verify HEAD 2>/dev/null || echo no-head; git status --porcelain 2>/dev/null; } | sha1sum | cut -d' ' -f1 )
  else
    cur=$(find . -type f -newermt '-1 minute' -not -path './node_modules/*' 2>/dev/null | sort | sha1sum | cut -d' ' -f1)
  fi
  [ "$cur" = "$prev" ] && stall=$((stall+1)) || stall=0
  prev="$cur"
  if [ "$stall" -ge "$MAX_STALL" ]; then
    printf 'BLOCKED: no file changes for %s iterations. See .agents/STATE.md\n' "$MAX_STALL" > BLOCKED.md
    echo "stalled"; exit 2
  fi

  if [ "$i" -eq 1 ]; then prompt="/start"; else prompt="/next"; fi
  timeout "${TURN_TIMEOUT:-1800}" opencode run "$prompt" 2>&1 | tee ".agents/logs/iter-$i.log"

  if bash ./scripts/dod.sh; then echo "DONE at iteration $i"; exit 0; fi
done

echo "MAX_ITERS reached"; exit 1
