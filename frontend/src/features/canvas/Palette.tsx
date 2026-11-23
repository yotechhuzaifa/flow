import type { NodeType } from "@shared/lib/types"

const TYPES: { type: NodeType; label: string }[] = [
  { type: "send", label: "Send" },
  { type: "delay", label: "Delay" },
  { type: "log", label: "Log" },
  { type: "if", label: "If" },
]

export default function Palette({ onAdd }: { onAdd: (t: NodeType) => void }) {
  return (
    <aside
      data-testid="palette"
      style={{ width: 160, borderRight: "1px solid #333", padding: 12, display: "grid", gap: 6 }}
    >
      {TYPES.map((t) => (
        <button key={t.type} data-testid={`palette-${t.type}`} onClick={() => onAdd(t.type)}>
          {t.label}
        </button>
      ))}
    </aside>
  )
}
