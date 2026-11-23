// Browser stand-in for @wailsio/runtime. The Wails build resolves the real module instead.
export const Events = {
  On(_name: string, _cb: (e: { data: unknown }) => void): () => void {
    return () => {}
  },
  Emit(_name: string, _data?: unknown): void {},
}
