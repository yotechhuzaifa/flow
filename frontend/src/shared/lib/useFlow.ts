import { useCallback, useEffect, useRef, useState } from "react"
import { api } from "@shared/api/client"
import type { Flow, FlowEdge, FlowNode, NodeStatus, NodeType } from "./types"

const empty = (): Flow => ({
  id: "local",
  name: "Untitled",
  updatedAt: "",
  nodes: [],
  edges: [],
  vars: {},
})

let seq = 0
const nid = (p: string) => `${p}_${Date.now().toString(36)}_${seq++}`

export function useFlow() {
  const [flow, setFlow] = useState<Flow>(empty)
  const [statuses, setStatuses] = useState<Record<string, NodeStatus>>({})
  const saveTimer = useRef<number | undefined>(undefined)

  // load once
  useEffect(() => {
    api.getFlow("local").then((f) => f && setFlow(f)).catch(() => {})
  }, [])

  // autosave, debounced
  const persist = useCallback((next: Flow) => {
    window.clearTimeout(saveTimer.current)
    saveTimer.current = window.setTimeout(() => api.saveFlow(next).catch(() => {}), 400)
  }, [])

  const update = useCallback(
    (fn: (f: Flow) => Flow) => {
      setFlow((cur) => {
        const next = fn(cur)
        persist(next)
        return next
      })
    },
    [persist],
  )

  const addNode = useCallback(
    (type: NodeType) => {
      update((f) => ({
        ...f,
        nodes: [...f.nodes, { id: nid(type), type, pos: { x: 80 + f.nodes.length * 40, y: 80 }, config: {} }],
      }))
    },
    [update],
  )

  const moveNode = useCallback(
    (id: string, x: number, y: number) => {
      update((f) => ({ ...f, nodes: f.nodes.map((n) => (n.id === id ? { ...n, pos: { x, y } } : n)) }))
    },
    [update],
  )

  const addEdge = useCallback(
    (e: Omit<FlowEdge, "id">) => {
      update((f) => {
        const dup = f.edges.some(
          (x) => x.from.node === e.from.node && x.from.port === e.from.port && x.to.node === e.to.node,
        )
        if (dup || e.from.node === e.to.node) return f
        return { ...f, edges: [...f.edges, { ...e, id: nid("e") }] }
      })
    },
    [update],
  )

  const removeNode = useCallback(
    (id: string) => {
      update((f) => ({
        ...f,
        nodes: f.nodes.filter((n) => n.id !== id),
        edges: f.edges.filter((e) => e.from.node !== id && e.to.node !== id),
      }))
    },
    [update],
  )

  const run = useCallback(async () => {
    setStatuses({})
    try {
      const id = await api.runFlow(flow.id)
      if (!id) simulate(flow, setStatuses) // browser dev: no backend
    } catch {
      simulate(flow, setStatuses)
    }
  }, [flow.id])

  // live status from the backend, plus a local fallback for browser dev
  useEffect(() => {
    const off = api.onRunEvent((evt) => {
      if (evt.type === "node.started") setStatuses((s) => ({ ...s, [evt.nodeId]: "running" }))
      if (evt.type === "node.finished")
        setStatuses((s) => ({ ...s, [evt.nodeId]: evt.status as NodeStatus }))
    })
    return off
  }, [])

  return { flow, setFlow: update, addNode, moveNode, addEdge, removeNode, run, statuses }
}

export type { FlowNode }

/** Browser-only stand-in for the Go runner, so the UI is testable without the desktop shell. */
function simulate(flow: Flow, set: (fn: (s: Record<string, NodeStatus>) => Record<string, NodeStatus>) => void) {
  const order = topo(flow)
  let failed = false
  order.forEach((id, i) => {
    window.setTimeout(() => {
      if (failed) {
        set((s) => ({ ...s, [id]: "skipped" }))
        return
      }
      set((s) => ({ ...s, [id]: "running" }))
      window.setTimeout(() => {
        const ok = flow.nodes.find((n) => n.id === id)?.type !== "if" // 'if' always fails in the stub
        if (!ok) failed = true
        set((s) => ({ ...s, [id]: ok ? "ok" : "failed" }))
      }, 250)
    }, i * 400)
  })
}

function topo(flow: Flow): string[] {
  const indeg: Record<string, number> = {}
  const next: Record<string, string[]> = {}
  flow.nodes.forEach((n) => (indeg[n.id] = 0))
  flow.edges.forEach((e) => {
    next[e.from.node] = [...(next[e.from.node] ?? []), e.to.node]
    indeg[e.to.node] = (indeg[e.to.node] ?? 0) + 1
  })
  const q = Object.keys(indeg).filter((id) => indeg[id] === 0)
  const out: string[] = []
  while (q.length) {
    const id = q.shift()!
    out.push(id)
    ;(next[id] ?? []).forEach((d) => {
      if (--indeg[d] === 0) q.push(d)
    })
  }
  return out
}
