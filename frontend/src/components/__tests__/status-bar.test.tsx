import { describe, it, expect, vi } from "vitest";

vi.mock("@wails/go/main/App", () => ({
  GetCachedModels: vi.fn(() => Promise.resolve([])),
  HasAPIKey: vi.fn(() => Promise.resolve(false)),
}));

import { render, screen } from "@testing-library/react";
import { StatusBar } from "@/components/status-bar/status-bar";
import { useChatStore } from "@/lib/store/chat";

describe("StatusBar", () => {
  it("renders provider options", () => {
    render(<StatusBar provider="opencode-go" model="" onProviderChange={() => {}} onModelChange={() => {}} onOpenSettings={() => {}} onExport={() => {}} />);
    expect(screen.getByText("OpenCode Go")).toBeInTheDocument();
  });

  it("shows connection indicator", () => {
    render(<StatusBar provider="opencode-go" model="" onProviderChange={() => {}} onModelChange={() => {}} onOpenSettings={() => {}} onExport={() => {}} />);
    const checkCircles = document.querySelectorAll(".lucide");
    expect(checkCircles.length).toBeGreaterThan(0);
  });

  it("shows generate button when onGenerate is provided", () => {
    render(<StatusBar provider="opencode-go" model="" onProviderChange={() => {}} onModelChange={() => {}} onOpenSettings={() => {}} onExport={() => {}} onGenerate={() => {}} />);
    const btn = screen.getByTitle("Generate project files");
    expect(btn).toBeInTheDocument();
  });

  it("hides phase badge when no messages exist", () => {
    useChatStore.setState({ messages: [] });
    render(<StatusBar provider="opencode-go" model="" onProviderChange={() => {}} onModelChange={() => {}} onOpenSettings={() => {}} onExport={() => {}} />);
    expect(screen.queryByText(/Phase/)).toBeNull();
  });
});
