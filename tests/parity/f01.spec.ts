/**
 * f01 — the app renders the canvas.
 *
 * Runs against the Vite dev server (started by playwright.config.ts). The canvas and palette
 * render without a backend; persistence and live runs are added by later rows.
 */
import { expect, test } from "@playwright/test"

const BASE = process.env.BASE_URL ?? "http://localhost:5173"

test("f01: canvas and palette render", async ({ page }) => {
  await page.goto(BASE)
  await expect(page.getByTestId("palette")).toBeVisible()
  await expect(page.getByTestId("canvas-root")).toBeVisible()
  await expect(page.getByTestId("run-button")).toBeVisible()
})
