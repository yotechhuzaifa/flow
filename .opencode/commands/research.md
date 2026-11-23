---
description: Research the target flow builder and produce spec files (no app code)
agent: re-observer
---

Target: $ARGUMENTS

Load the `grabbing-flows` skill and follow it exactly.

Phase A — Layer 0 intel. Read `docs/RESEARCH.md`, then harvest in this order:

1. docs / Learning Center — `firecrawl_map`, then `firecrawl_scrape` one page per call.
2. changelog / release blog / "What's new" — `firecrawl_scrape`; record dates.
3. community feature requests — `firecrawl_search` with `site:` filters.
4. public API docs / collection schema — `firecrawl_scrape`.
5. demos / walkthroughs — `firecrawl_search`; record journeys, not frames.
6. open-source reference code — read-only; record repo + commit + path + license.

Write `spec/INTEL.md` and `spec/GAPS.md`. Cite a URL or repo path for every claim.
Never read the target's shipped JS or sourcemaps. Never copy target microcopy, assets, CSS,
icons, class names, or code.

Phase B — observe. Take the first three gaps and run three scenarios with the playwright MCP.
Save traces under `traces/<scenario>/`.

Phase C — extract. Write `spec/features.jsonl`, one JSON object per line, every line citing a
trace path. Then run `node scripts/spec.mjs check` and fix the source, not the checker.

Phase D — compile. Run `node scripts/spec.mjs propose` and `node scripts/spec.mjs plan`.
Write `spec/RESEARCH-DONE.md` when every gap is covered or marked `unobservable, reason: …`.

Stop after three scenarios and report: which gaps are now covered, which remain, and what you
could not observe through automation. Do not write application code. Do not edit `parity.csv`;
a human merges proposals with a `[contract]` commit.
