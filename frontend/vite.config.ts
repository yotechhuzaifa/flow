import { defineConfig } from "vite"
import react from "@vitejs/plugin-react"
import { fileURLToPath, URL } from "node:url"

export default defineConfig({
  plugins: [react()],
  esbuild: { target: "esnext" },
  build: { target: "esnext", outDir: "dist", rollupOptions: { external: ["@wailsio/runtime", "@bindings/flow/api"] } },
  server: { port: 5173 },
  optimizeDeps: { exclude: ["@wailsio/runtime", "@bindings/flow/api"] },
  resolve: {
    alias: {
      "@bindings": fileURLToPath(new URL("./bindings", import.meta.url)),
      "@shared": fileURLToPath(new URL("./src/shared", import.meta.url)),
      "@features": fileURLToPath(new URL("./src/features", import.meta.url)),
      // inside Wails this resolves to the generated runtime; in a browser it is a no-op stub
      "@wailsio/runtime": fileURLToPath(new URL("./src/shared/lib/wails-stub.ts", import.meta.url)),
    },
  },
})
