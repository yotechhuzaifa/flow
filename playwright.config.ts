import { defineConfig, devices } from "@playwright/test"

export default defineConfig({
  testDir: "./tests/parity",
  timeout: 30_000,
  retries: process.env.CI ? 1 : 0,
  reporter: process.env.CI ? "line" : "list",
  use: {
    baseURL: process.env.BASE_URL ?? "http://localhost:5173",
    trace: "retain-on-failure", // the same traces the research skill consumes
  },
  projects: [{ name: "chromium", use: { ...devices["Desktop Chrome"] } }],
  // Locally and in CI the suite starts Vite itself. Point BASE_URL at another
  // server when one is already running.
  webServer: process.env.BASE_URL
    ? undefined
    : {
        command: "npm --prefix frontend run dev -- --port 5173 --strictPort",
        url: "http://localhost:5173",
        reuseExistingServer: !process.env.CI,
        timeout: 60_000,
      },
})
