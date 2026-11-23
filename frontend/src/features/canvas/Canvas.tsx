import {
  Background,
  Controls,
  ReactFlow,
  type Edge,
  type Node,
  type OnConnect,
  type NodeChange,
  applyNodeChanges,
} from "@xyflow/react"
import { useCallback, useMemo } from "react"
import type { Flow, FlowEdge, NodeStatus } from "@shared/lib/types"
import StatusNode from "../nodes/StatusNode"

const nodeTypes = { status: StatusNode }

function toRfNodes(flow: Flow, statuses: Record<string, NodeStatus>): Node[] {
  return flow.nodes.map((n) => ({
    id: n.id,
    type: "status",
    position: n.pos,
    data: { label: n.type, status: statuses[n.id] ?? "idle", config: n.config },
  }))
}

function toRfEdges(flow: Flow): Edge[] {
  return flow.edges.map((e) => ({
    id: e.id,
    source: e.from.node,
    sourceHandle: e.from.port,
    target: e.to.node,
    targetHandle: e.to.port ?? "in",
  }))
}

export default function Canvas({
  flow,
  onChange,
  onMoveNode,
  onAddEdge,
  onRemoveNode,
  statuses,
}: {
  flow: Flow
  onChange: (fn: (f: Flow) => Flow) => void
  onMoveNode: (id: string, x: number, y: number) => void
  onAddEdge: (e: Omit<FlowEdge, "id">) => void
  onRemoveNode: (id: string) => void
  statuses: Record<string, NodeStatus>
}) {
  const nodes = useMemo(() => toRfNodes(flow, statuses), [flow, statuses])
  const edges = useMemo(() => toRfEdges(flow), [flow])

  const onNodesChange = useCallback(
    (changes: NodeChange[]) => {
      onChange((f) => {
        const next = applyNodeChanges(changes, toRfNodes(f, statuses))
        return {
          ...f,
          nodes: f.nodes.map((n) => {
            const rf = next.find((x) => x.id === n.id)
            return rf ? { ...n, pos: { x: rf.position.x, y: rf.position.y } } : n
          }),
        }
      })
    },
    [onChange, statuses],
  )

  const onConnect: OnConnect = useCallback(
    (c) => {
      if (!c.source || !c.target) return
      onAddEdge({
        from: { node: c.source, port: c.sourceHandle ?? "out" },
        to: { node: c.target, port: c.targetHandle ?? "in" },
      })
    },
    [onAddEdge],
  )

  return (
    <ReactFlow
      data-testid="canvas"
      nodes={nodes}
      edges={edges}
      nodeTypes={nodeTypes}
      onNodesChange={onNodesChange}
      onConnect={onConnect}
      onNodeDragStop={(_, n) => onMoveNode(n.id, n.position.x, n.position.y)}
      onNodesDelete={(ds) => ds.forEach((d) => onRemoveNode(d.id))}
      fitView
    >
      <Background />
      <Controls />
    </ReactFlow>
  )
}
