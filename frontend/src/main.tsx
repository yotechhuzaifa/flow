import "@xyflow/react/dist/style.css"
import React from "react"
import { createRoot } from "react-dom/client"
import App from "./App"

// Wails runtime must be imported for side effects, or events never arrive.
// It is only present inside the desktop shell; the browser dev server ignores it.
void import("@wailsio/runtime").catch(() => {
  // running in a plain browser: no runtime, no events
})

createRoot(document.getElementById("root")!).render(<App />)
