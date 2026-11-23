// Mirrors docs/CONTRACT.md. If this file changes, the contract changed.
export type NodeType = "send" | "delay" | "log" | "if"
export type NodeStatus = "idle" | "running" | "ok" | "failed" | "skipped"

export interface FlowNode {
  id: string
  type: NodeType
  pos: { x: number; y: number }
  config: Record<string, unknown>
}

export interface FlowEdge {
  id: string
  from: { node: string; port: string }
  to: { node: string; port: string }
}

export interface Flow {
  id: string
  name: string
  updatedAt: string
  nodes: FlowNode[]
  edges: FlowEdge[]
  vars: Record<string, string>
}

export const OUT_PORTS: Record<NodeType, string[]> = {
  send: ["out"],
  delay: ["out"],
  log: ["out"],
  if: ["true", "false"],
}
