---
description: Research the target flow builder and produce spec files (no app code)
agent: re-observer
---

Target: $ARGUMENTS

Load the `grabbing-flows` skill and follow it exactly.

Phase A — intel. Fetch their public docs and changelog with webfetch, one page per call.
Write `spec/INTEL.md` and `spec/GAPS.md`. Cite a URL for every claim.

Phase B — observe. Take the first three gaps and run three scenarios with the playwright MCP.
Save traces under `traces/<scenario>/`.

Phase C — extract. Write `spec/features.jsonl`, one JSON object per line, every line citing a
trace path. Validate against `spec/feature.schema.json`.

Stop after three scenarios and report: which gaps are now covered, which remain, and what you
could not observe through automation. Do not write application code.
