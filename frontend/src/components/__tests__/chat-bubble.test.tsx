import { describe, it, expect } from "vitest";
import { render, screen } from "@testing-library/react";
import { ChatBubble } from "@/components/chat/chat-bubble";

const userMsg = { id: "1", project_id: "p1", phase: "identity", role: "user" as const, content: "hello", created_at: "" };
const aiMsg = { id: "2", project_id: "p1", phase: "identity", role: "assistant" as const, content: "hi there", created_at: "" };

describe("ChatBubble", () => {
  it("renders user message", () => {
    render(<ChatBubble message={userMsg} />);
    expect(screen.getByText("hello")).toBeInTheDocument();
  });

  it("renders AI message with markdown", () => {
    render(<ChatBubble message={aiMsg} />);
    expect(screen.getByText("hi there")).toBeInTheDocument();
  });

  it("shows streaming cursor when isStreaming is true", () => {
    render(<ChatBubble message={aiMsg} isStreaming />);
    const cursor = document.querySelector(".animate-blink");
    expect(cursor).toBeTruthy();
  });

  it("shows user avatar letter when no avatar image", () => {
    render(<ChatBubble message={userMsg} userName="Morten" />);
    expect(screen.queryByText("M")).toBeTruthy();
  });
});
