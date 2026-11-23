/**
 * f01 — the app renders the canvas.
 *
 * Runs in server mode: `npm run dev` (vite) talks to `cmd/server` over /api and /events.
 * Locally: BASE_URL=http://localhost:5173 npx playwright test
 */
import { expect, test } from "@playwright/test"

const BASE = process.env.BASE_URL ?? "http://localhost:5173"

test("f01: canvas and palette render", async ({ page }) => {
  await page.goto(BASE)
  await expect(page.getByTestId("palette")).toBeVisible()
  await expect(page.getByTestId("canvas-root")).toBeVisible()
  await expect(page.getByTestId("run-button")).toBeVisible()
})
