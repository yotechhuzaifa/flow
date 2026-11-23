// Copy the Vite build into the Go embed directory.
// Keeps `internal/webui/dist` and `frontend/dist` from drifting apart.
import { cpSync, existsSync, mkdirSync } from "node:fs"
import { dirname, join } from "node:path"
import { fileURLToPath } from "node:url"

const root = join(dirname(fileURLToPath(import.meta.url)), "..")
const src = join(root, "frontend", "dist")
const dst = join(root, "internal", "webui", "dist")

if (!existsSync(join(src, "index.html"))) {
  console.error("sync-webui: frontend/dist/index.html missing; run the frontend build first")
  process.exit(1)
}

mkdirSync(dst, { recursive: true })
cpSync(src, dst, { recursive: true, force: true })
console.log(`sync-webui: ${src} -> ${dst}`)
