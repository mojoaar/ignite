import { vi } from "vitest";
import "@testing-library/jest-dom/vitest";
import "./__mocks__/wails";

const store = new Map<string, string>();

Object.defineProperty(globalThis, "localStorage", {
  value: {
    getItem: vi.fn((key: string) => store.get(key) ?? null),
    setItem: vi.fn((key: string, value: string) => {
      store.set(key, value);
    }),
    removeItem: vi.fn((key: string) => {
      store.delete(key);
    }),
    clear: vi.fn(() => {
      store.clear();
    }),
  },
  writable: true,
});

if (!(globalThis as any).crypto?.randomUUID) {
  (globalThis as any).crypto = {
    ...(globalThis as any).crypto,
    randomUUID: () => "00000000-0000-0000-0000-000000000000",
  };
}
