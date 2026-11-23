import { Handle, Position, type NodeProps } from "@xyflow/react"
import { OUT_PORTS, type NodeStatus, type NodeType } from "@shared/lib/types"

const COLOR: Record<NodeStatus, string> = {
  idle: "#555",
  running: "#e0a800",
  ok: "#2e7d32",
  failed: "#c62828",
  skipped: "#777",
}

export default function StatusNode({ id, data }: NodeProps) {
  const label = String(data.label ?? "")
  const status = (data.status as NodeStatus) ?? "idle"
  const outs = OUT_PORTS[label as NodeType] ?? ["out"]

  return (
    <div
      data-testid={`node-${label}`}
      data-status={status}
      style={{
        padding: "8px 12px",
        minWidth: 110,
        background: "#1e1e1e",
        color: "#eee",
        border: `2px solid ${COLOR[status]}`,
        borderRadius: 6,
        fontSize: 13,
      }}
    >
      <Handle type="target" position={Position.Left} id="in" data-testid={`in-${id}`} />
      <div style={{ fontWeight: 600 }}>{label}</div>
      <div style={{ fontSize: 10, opacity: 0.6 }}>{status}</div>
      {outs.map((p, i) => (
        <Handle
          key={p}
          type="source"
          position={Position.Right}
          id={p}
          data-testid={`out-${id}-${p}`}
          style={{ top: `${((i + 1) / (outs.length + 1)) * 100}%` }}
        />
      ))}
    </div>
  )
}
