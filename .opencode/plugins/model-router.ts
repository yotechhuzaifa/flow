/**
 * model-router.ts — OpenCode V2 plugin. Auto model choice, no extra subscription.
 *
 * stage 1: prompt hook classifies the request, remembers the choice per session.
 * stage 2: experimental.model.resolve applies it at request time.
 *
 * Ladder (all included in Go, so quota is the budget, not money):
 *   grunt  glm-5.3-flash     6,320 req/5h   renames, tests, small edits
 *   daily  qwen4-max           463 req/5h   normal feature work
 *   hard   kimi-k3.1-pro       147 req/5h   architecture, Wails, executor
 *   bulk   hy4               2,710 req/5h   sweeps across many files
 */
import { Plugin } from "@opencode/plugin"

type Ref = { providerID: string; id: string; variant?: string }

const LADDER = {
  grunt: { providerID: "opencode-go", id: "glm-5.3-flash", variant: "low" },
  daily: { providerID: "opencode-go", id: "qwen4-max", variant: "high" },
  hard: { providerID: "opencode-go", id: "kimi-k3.1-pro", variant: "high" },
  bulk: { providerID: "opencode-go", id: "hy4", variant: "high" },
} satisfies Record<string, Ref>

type Kind = keyof typeof LADDER

const HARD = /\b(architect|design|schema|contract|contract|executor|topolog|wails|binding|event|race|why|root cause|refactor)/i
const GRUNT = /\b(rename|typo|docstring|comment|format|lint|fixture|test case|boilerplate|readme)\b/i
const BULK = /\b(every|all rows|sweep|for each|batch)\b/i

function classify(text: string): Kind {
  if (BULK.test(text)) return "bulk"
  if (HARD.test(text)) return "hard"
  if (GRUNT.test(text)) return "grunt"
  return text.length < 240 ? "grunt" : "daily"
}

export default Plugin.define({
  id: "model-router",
  async setup(ctx) {
    const decisions = new Map<string, Ref>()
    const pinned = new Set<string>()

    await ctx.session.hook("prompt", (event) => {
      const kind = classify(event.prompt.text ?? "")
      decisions.set(event.sessionID, LADDER[kind])
    })

    await ctx.session.hook("experimental.model.resolve", (event) => {
      if (event.kind !== "primary") return
      if (pinned.has(event.sessionID)) return
      const decided = decisions.get(event.sessionID)
      if (!decided) return
      // respect a model the human picked in /models
      const onDefault = event.model.providerID === "opencode-go" && event.model.id === "deepseek-v4.1-flash"
      if (!onDefault) { pinned.add(event.sessionID); return }
      event.model = { ...decided }
    })
  },
})
