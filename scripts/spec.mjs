#!/usr/bin/env node
/**
 * Spec compiler for the research -> parity -> test pipeline.
 * No dependencies. Node 18+.
 *
 *   node scripts/spec.mjs check    validate spec/features.jsonl and evidence paths
 *   node scripts/spec.mjs propose  write spec/parity.proposed.csv + spec/parity.map.json
 *   node scripts/spec.mjs plan     write spec/TEST-PLAN.md
 *   node scripts/spec.mjs audit    write spec/AUDIT.md and fail on broken gates
 *   node scripts/spec.mjs all      check + propose + plan + audit
 *
 * parity.csv and docs/CONTRACT.md are frozen. This script never edits them.
 * It writes proposals only; a human [contract] commit is still the merge gate.
 */
import { existsSync, readFileSync, readdirSync, writeFileSync } from "node:fs"
import { dirname, join, relative } from "node:path"
import { fileURLToPath } from "node:url"

const root = join(dirname(fileURLToPath(import.meta.url)), "..")
const at = (...p) => join(root, ...p)
const FEATURES = at("spec", "features.jsonl")
const PARITY = at("parity.csv")
const MAP = at("spec", "parity.map.json")
const PROPOSED = at("spec", "parity.proposed.csv")
const PLAN = at("spec", "TEST-PLAN.md")
const AUDIT = at("spec", "AUDIT.md")
const RESEARCH_DONE = at("spec", "RESEARCH-DONE.md")
const TESTS = at("tests", "parity")
const PRIORITIES = new Set(["must", "should", "could", "wont"])
const ID_RE = /^[a-z0-9][a-z0-9._-]*$/

const cmd = process.argv[2] || "check"

function isUrl(value) {
  return /^https?:\/\//i.test(value)
}

function localPath(value) {
  return String(value).split("#")[0].trim()
}

function csvLine(line) {
  const out = []
  let cur = ""
  let quoted = false
  for (let i = 0; i < line.length; i++) {
    const c = line[i]
    if (quoted) {
      if (c === '"' && line[i + 1] === '"') { cur += '"'; i++ }
      else if (c === '"') quoted = false
      else cur += c
    } else if (c === '"') quoted = true
    else if (c === ",") { out.push(cur); cur = "" }
    else cur += c
  }
  out.push(cur)
  return out
}

function csvEscape(value) {
  const s = String(value ?? "")
  return /[",\n]/.test(s) ? `"${s.replace(/"/g, '""')}"` : s
}

function readParity() {
  if (!existsSync(PARITY)) return []
  const lines = readFileSync(PARITY, "utf8").split(/\r?\n/).filter((l) => l.trim())
  if (!lines.length) return []
  const header = csvLine(lines[0]).map((h) => h.trim())
  return lines.slice(1).map((line, i) => {
    const cols = csvLine(line)
    const row = { _line: i + 2 }
    header.forEach((h, j) => { row[h] = (cols[j] ?? "").trim() })
    return row
  })
}

function readMap() {
  if (!existsSync(MAP)) return {}
  return JSON.parse(readFileSync(MAP, "utf8"))
}

function readFeatures() {
  if (!existsSync(FEATURES)) return { features: [], problems: [] }
  const features = []
  const problems = []
  const seen = new Set()
  const lines = readFileSync(FEATURES, "utf8").split(/\r?\n/)

  lines.forEach((raw, index) => {
    const n = index + 1
    const line = raw.trim()
    if (!line) return
    let f
    try { f = JSON.parse(line) } catch (e) {
      problems.push(`line ${n}: invalid JSON (${e.message})`)
      return
    }
    const bad = (msg) => problems.push(`line ${n} (${f.id || "no-id"}): ${msg}`)

    if (typeof f.id !== "string" || !ID_RE.test(f.id)) bad("id must match ^[a-z0-9][a-z0-9._-]*$")
    if (typeof f.id === "string") {
      if (seen.has(f.id)) bad(`duplicate id "${f.id}"`)
      seen.add(f.id)
    }
    if (typeof f.trigger !== "string" || !f.trigger.trim()) bad("trigger is required")
    if (typeof f.observable !== "string" || !f.observable.trim()) bad("observable is required")
    if (!PRIORITIES.has(f.priority)) bad(`priority must be one of ${[...PRIORITIES].join(", ")}`)

    if (!Array.isArray(f.evidence) || f.evidence.length === 0) {
      bad("evidence is required")
    } else {
      let trace = 0
      for (const e of f.evidence) {
        if (typeof e !== "string" || !e.trim()) { bad("evidence entries must be non-empty strings"); continue }
        if (isUrl(e)) continue
        const p = localPath(e)
        if (p.startsWith("traces/")) trace++
        else if (!p.startsWith("spec/")) bad(`evidence path must be under traces/ or spec/: ${e}`)
        if (!existsSync(at(p))) bad(`evidence file missing: ${p}`)
      }
      if (!trace) bad("evidence needs at least one traces/ path")
    }

    for (const key of ["preconditions", "errors", "wire", "unsure"]) {
      if (f[key] === undefined) continue
      if (!Array.isArray(f[key]) || f[key].some((x) => typeof x !== "string" || !x.trim())) {
        bad(`${key} must be an array of non-empty strings`)
      }
    }
    features.push(f)
  })

  return { features, problems }
}

function cmdCheck() {
  if (!existsSync(FEATURES)) {
    console.log("spec: no spec/features.jsonl yet — run /research first")
    return 0
  }
  const { features, problems } = readFeatures()
  for (const p of problems) console.error(`FAIL ${p}`)

  if (existsSync(MAP)) {
    try {
      const ids = new Set(features.map((f) => f.id))
      for (const [row, feature] of Object.entries(readMap())) {
        if (!ids.has(feature)) {
          console.error(`FAIL ${relative(root, MAP)}: ${row} -> unknown feature "${feature}"`)
          problems.push("map")
        }
      }
    } catch (e) {
      console.error(`FAIL ${relative(root, MAP)}: ${e.message}`)
      problems.push("map")
    }
  }

  if (problems.length) return 1
  console.log(`spec: ${features.length} feature(s) valid`)
  return 0
}

function cmdPropose() {
  const check = cmdCheck()
  if (check) {
    console.error("spec: fix check failures before proposing")
    return check
  }
  const { features } = readFeatures()
  const used = new Set(readParity().map((r) => r.id))
  const map = readMap()
  const mapped = new Set(Object.values(map))
  let n = 0

  for (const f of features) {
    if (mapped.has(f.id)) continue
    let id
    do { n++; id = `p${String(n).padStart(2, "0")}` } while (used.has(id) || Object.hasOwn(map, id))
    used.add(id)
    map[id] = f.id
    mapped.add(f.id)
  }

  const header = "id,feature,priority,depends_on,effort,status"
  const rows = Object.entries(map).map(([id, featureId]) => {
    const f = features.find((x) => x.id === featureId)
    // dod.sh parses parity.csv with awk -F','. Keep the feature column comma-free so every
    // MUST row is counted, including rows whose prose would otherwise need CSV quoting.
    const feature = String(f.observable).replace(/,/g, ";").replace(/\s+/g, " ").trim()
    return [id, feature, f.priority.toUpperCase(), "", "", "todo"].map(csvEscape).join(",")
  })
  writeFileSync(PROPOSED, [header, ...rows].join("\n") + "\n")
  writeFileSync(MAP, JSON.stringify(map, null, 2) + "\n")
  console.log(`spec: wrote ${relative(root, PROPOSED)} (${rows.length} row(s))`)
  console.log(`spec: wrote ${relative(root, MAP)}`)
  console.log("spec: parity.csv is frozen; merge selected rows with a human [contract] commit")
  return 0
}

function cmdPlan() {
  const check = cmdCheck()
  if (check) return check
  if (!existsSync(FEATURES)) {
    writeFileSync(PLAN, "# Test plan\n\nNo features yet. Run `/research <target>` first.\n")
    console.log(`spec: wrote ${relative(root, PLAN)} (placeholder)`)
    return 0
  }

  const { features } = readFeatures()
  const map = readMap()
  const rowFor = Object.fromEntries(Object.entries(map).map(([row, feature]) => [feature, row]))
  const out = [
    "# Test plan",
    "",
    "Generated from `spec/features.jsonl` by `scripts/spec.mjs`. Do not edit by hand.",
    "",
  ]

  for (const f of features) {
    out.push(`## ${f.id} — ${f.priority}`)
    out.push("")
    out.push(`- Parity row: ${rowFor[f.id] ? `\`${rowFor[f.id]}\` (proposed)` : "unmapped — run `node scripts/spec.mjs propose`"}`)
    out.push(`- Observable: ${f.observable}`)
    out.push(`- Trigger: ${f.trigger}`)
    if (f.preconditions?.length) out.push(`- Preconditions: ${f.preconditions.join("; ")}`)
    if (f.errors?.length) out.push(`- Expected errors: ${f.errors.join("; ")}`)
    if (f.wire?.length) out.push(`- Wire: ${f.wire.join(", ")}`)
    out.push(`- Evidence: ${f.evidence.join(", ")}`)
    out.push("")
    out.push("Checklist:")
    out.push("- [ ] Playwright test under `tests/parity/`")
    out.push("- [ ] One assertion per observable outcome")
    out.push("- [ ] Preconditions set up explicitly")
    out.push("- [ ] Every expected error exercised")
    out.push("- [ ] Every interactive element has a `data-testid`")
    out.push("- [ ] `spec-auditor` cannot make it pass without the feature")
    out.push("")
  }

  writeFileSync(PLAN, out.join("\n"))
  console.log(`spec: wrote ${relative(root, PLAN)} (${features.length} feature(s))`)
  return 0
}

function cmdAudit() {
  const parity = readParity()
  const problems = []
  const warnings = []

  for (const row of parity) {
    if (row.priority !== "MUST") continue
    const file = join(TESTS, `${row.id}.spec.ts`)
    if (!existsSync(file)) {
      problems.push(`${row.id}: MUST row has no tests/parity/${row.id}.spec.ts`)
      continue
    }
    const text = readFileSync(file, "utf8")
    if (!text.includes(row.id)) warnings.push(`${row.id}: test file does not mention ${row.id}`)
    if (/test\.(skip|fixme)\s*\(/.test(text)) problems.push(`${row.id}: MUST test is skipped or fixme`)
  }

  const open = parity.filter((r) => r.priority === "MUST" && r.status !== "pass").map((r) => r.id)
  if (open.length) warnings.push(`MUST rows still open: ${open.join(" ")}`)

  if (existsSync(FEATURES)) {
    const { features, problems: featureProblems } = readFeatures()
    featureProblems.forEach((p) => problems.push(p))
    let map = {}
    if (existsSync(MAP)) {
      try { map = readMap() } catch (e) { problems.push(`spec/parity.map.json: ${e.message}`) }
    }
    const mapped = new Set(Object.values(map))
    const unmapped = features.filter((f) => !mapped.has(f.id)).map((f) => f.id)
    if (unmapped.length) warnings.push(`features not mapped to parity: ${unmapped.join(" ")}`)
    if (!existsSync(RESEARCH_DONE)) warnings.push("spec/RESEARCH-DONE.md is missing")
  }

  const lines = ["# Spec audit", ""]
  for (const p of problems) lines.push(`- FAIL ${p}`)
  for (const w of warnings) lines.push(`- WARN ${w}`)
  if (!problems.length && !warnings.length) lines.push("- PASS no issues")
  lines.push("", `Result: ${problems.length ? "FAIL" : "PASS"}`, "")
  writeFileSync(AUDIT, lines.join("\n"))
  console.log(`spec: wrote ${relative(root, AUDIT)} (${problems.length} fail, ${warnings.length} warn)`)
  return problems.length ? 1 : 0
}

function main() {
  switch (cmd) {
    case "check": return cmdCheck()
    case "propose": return cmdPropose()
    case "plan": return cmdPlan()
    case "audit": return cmdAudit()
    case "all": {
      const check = cmdCheck()
      if (check) return check
      const propose = cmdPropose()
      if (propose) return propose
      const plan = cmdPlan()
      if (plan) return plan
      return cmdAudit()
    }
    default:
      console.error(`unknown command "${cmd}"`)
      console.error("usage: node scripts/spec.mjs [check|propose|plan|audit|all]")
      return 2
  }
}

process.exit(main())
