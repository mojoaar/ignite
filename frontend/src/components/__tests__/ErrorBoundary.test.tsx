import { describe, it, expect } from "vitest";
import { render, screen } from "@testing-library/react";
import { ErrorBoundary } from "@/components/ErrorBoundary";

function Crash(): React.ReactElement {
  throw new Error("test crash");
}

function Okay(): React.ReactElement {
  return <div>all good</div>;
}

describe("ErrorBoundary", () => {
  it("renders children normally", () => {
    render(<ErrorBoundary><Okay /></ErrorBoundary>);
    expect(screen.getByText("all good")).toBeInTheDocument();
  });

  it("shows error UI on crash", () => {
    render(<ErrorBoundary><Crash /></ErrorBoundary>);
    expect(screen.getByText("Something went wrong")).toBeInTheDocument();
    expect(screen.getByText("test crash")).toBeInTheDocument();
    expect(screen.getByText("Reload")).toBeInTheDocument();
  });
});
