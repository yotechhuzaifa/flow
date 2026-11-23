// Non-relative ambient declarations. Real modules come from the Wails build/generation step.
declare module "@wailsio/runtime" {
  export const Events: {
    On(name: string, cb: (e: { data: unknown }) => void): () => void
    Emit(name: string, data?: unknown): void
  }
}

declare module "@bindings/flow/api" {
  export function ListFlows(): Promise<unknown[]>
  export function GetFlow(id: string): Promise<unknown | null>
  export function SaveFlow(f: unknown): Promise<void>
  export function DeleteFlow(id: string): Promise<void>
  export function StartRun(flowID: string): Promise<string>
}
