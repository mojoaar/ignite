import { describe, it, expect, beforeEach } from "vitest";
import { useThemeStore } from "../theme";

describe("theme store", () => {
  beforeEach(() => {
    localStorage.clear();
    useThemeStore.setState({ mode: "dark" });
  });

  it("defaults to dark", () => {
    expect(useThemeStore.getState().mode).toBe("dark");
  });

  it("toggle switches dark <-> light", () => {
    useThemeStore.getState().toggle();
    expect(useThemeStore.getState().mode).toBe("light");
    useThemeStore.getState().toggle();
    expect(useThemeStore.getState().mode).toBe("dark");
  });

  it("persists to localStorage", () => {
    useThemeStore.getState().toggle();
    expect(localStorage.getItem("ignite-mode")).toBe("light");
  });

  it("init reads from localStorage", () => {
    localStorage.setItem("ignite-mode", "light");
    useThemeStore.getState().init();
    expect(useThemeStore.getState().mode).toBe("light");
  });

  it("init and toggle update data-mode attribute", () => {
    localStorage.setItem("ignite-mode", "light");
    useThemeStore.getState().init();
    expect(document.documentElement.getAttribute("data-mode")).toBe("light");
    useThemeStore.getState().toggle();
    expect(document.documentElement.getAttribute("data-mode")).toBe("dark");
  });
});
