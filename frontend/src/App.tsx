import { ReactFlowProvider } from "@xyflow/react"
import Canvas from "@features/canvas/Canvas"
import Palette from "@features/canvas/Palette"
import { useFlow } from "@shared/lib/useFlow"

export default function App() {
  const { flow, setFlow, addNode, moveNode, addEdge, removeNode, run, statuses } = useFlow()

  return (
    <div style={{ display: "flex", height: "100vh", fontFamily: "system-ui" }}>
      <Palette onAdd={addNode} />
      <div style={{ flex: 1, position: "relative" }} data-testid="canvas-root">
        <ReactFlowProvider>
          <Canvas
            flow={flow}
            onChange={setFlow}
            onMoveNode={moveNode}
            onAddEdge={addEdge}
            onRemoveNode={removeNode}
            statuses={statuses}
          />
        </ReactFlowProvider>
        <button
          data-testid="run-button"
          onClick={run}
          style={{ position: "absolute", top: 12, right: 12, padding: "6px 12px" }}
        >
          Run
        </button>
      </div>
    </div>
  )
}
