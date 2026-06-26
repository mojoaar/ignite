import { create } from "zustand";

type Mode = "dark" | "light";

interface ThemeState {
  mode: Mode;
  toggle: () => void;
  init: () => void;
}

const KEY = "ignite-mode";

export const useThemeStore = create<ThemeState>((set) => ({
  mode: "dark",
    toggle: () =>
    set((s) => {
      const next: Mode = s.mode === "dark" ? "light" : "dark";
      localStorage.setItem(KEY, next);
      const root = document.documentElement;
      root.setAttribute("data-mode", next);
      root.classList.toggle("dark", next === "dark");
      return { mode: next };
    }),
  init: () => {
    const stored = localStorage.getItem(KEY);
    const mode: Mode = stored === "light" ? "light" : "dark";
    const root = document.documentElement;
    root.setAttribute("data-mode", mode);
    root.classList.toggle("dark", mode === "dark");
    set({ mode });
  },
}));
