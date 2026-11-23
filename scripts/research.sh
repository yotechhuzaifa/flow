#!/usr/bin/env bash
# One research cycle: intel + observation + extraction + compile.
# Human input is the target URL in spec/TARGET.md (or the first argument).
set -uo pipefail
cd "$(dirname "$0")/.."

mkdir -p .agents/logs

TARGET_FILE="spec/TARGET.md"
if [ ! -f "$TARGET_FILE" ]; then
  echo "research: missing $TARGET_FILE" >&2
  exit 1
fi

TARGET="${1:-$(sed -n 's/^[[:space:]]*url:[[:space:]]*//p' "$TARGET_FILE" | head -n1 | tr -d '\r')}"
if [ -z "$TARGET" ] || [ "$TARGET" = "unset" ]; then
  echo "research: set the url: line in $TARGET_FILE, or pass a URL" >&2
  exit 1
fi

if ! command -v opencode >/dev/null 2>&1; then
  echo "research: opencode is not on PATH" >&2
  exit 1
fi

echo "research: target $TARGET"
opencode run "/research $TARGET" 2>&1 | tee .agents/logs/research.log

if ! node scripts/spec.mjs check; then
  echo "research: features.jsonl is not valid; inspect .agents/logs/research.log" >&2
  exit 1
fi

node scripts/spec.mjs propose
node scripts/spec.mjs plan
node scripts/spec.mjs audit || true

echo
echo "next: review spec/parity.proposed.csv; a human merges selected rows into parity.csv"
echo "      with a [contract] commit, then run ./scripts/loop.sh"
