import { describe, it, expect, beforeEach, vi } from "vitest";

vi.mock("@wails/go/main/App", () => ({
  ListProjects: vi.fn(() => Promise.resolve([])),
  CreateProject: vi.fn(() => Promise.resolve()),
  GetProject: vi.fn(() => Promise.resolve(null)),
  GetMessages: vi.fn(() => Promise.resolve([])),
  DeleteProject: vi.fn(() => Promise.resolve()),
  UpdateProject: vi.fn(() => Promise.resolve()),
}));
vi.mock("@wails/runtime", () => ({ EventsOn: vi.fn(() => () => {}), EventsOff: vi.fn() }));

import { render, screen } from "@testing-library/react";
import { Sidebar } from "@/components/sidebar/sidebar";
import { useChatStore } from "@/lib/store/chat";

describe("Sidebar", () => {
  beforeEach(() => {
    useChatStore.setState({
      projects: [],
      activeProjectId: null,
      messages: [],
    });
  });

  it("shows empty state when no projects", () => {
    render(<Sidebar />);
    expect(screen.getByText("No projects yet")).toBeInTheDocument();
  });

  it("renders projects from store", () => {
    useChatStore.setState({
      projects: [
        { id: "p1", name: "Test Project", tagline: "", path: "", provider: "", model: "", created_at: "", updated_at: "" },
      ],
    });
    render(<Sidebar />);
    expect(screen.getByText("Test Project")).toBeInTheDocument();
  });

  it("shows brand name", () => {
    render(<Sidebar />);
    expect(screen.getByText("Ignite")).toBeInTheDocument();
  });

  it("shows New Project button", () => {
    render(<Sidebar />);
    expect(screen.getByText("New Project")).toBeInTheDocument();
  });
});
