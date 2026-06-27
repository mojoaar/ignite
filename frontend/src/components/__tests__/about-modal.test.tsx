import { describe, it, expect, vi } from "vitest";

vi.mock("@wails/go/main/App", () => ({
  GetVersion: vi.fn(() => Promise.resolve("0.1.6")),
}));

import { render, screen } from "@testing-library/react";
import { AboutModal } from "@/components/settings/about-modal";

describe("AboutModal", () => {
  it("renders title and version", () => {
    render(<AboutModal open={true} onClose={() => {}} />);
    expect(screen.getByText("Ignite")).toBeInTheDocument();
  });

  it("shows author info", () => {
    render(<AboutModal open={true} onClose={() => {}} />);
    expect(screen.getByText("Morten Johansen")).toBeInTheDocument();
  });
});
