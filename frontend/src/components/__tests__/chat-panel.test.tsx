import { describe, it, expect, beforeEach } from "vitest";
import { render, screen } from "@testing-library/react";
import { ChatPanel } from "@/components/chat/chat-panel";
import { useChatStore } from "@/lib/store/chat";

describe("ChatPanel", () => {
  beforeEach(() => {
    useChatStore.setState({
      projects: [],
      activeProjectId: null,
      activeProjectPath: "",
      messages: [],
      streamingContent: "",
      isStreaming: false,
    });
  });

  it("shows welcome banner when no project is selected", () => {
    render(<ChatPanel onSend={() => {}} />);
    expect(screen.getByText("Welcome to Ignite")).toBeInTheDocument();
  });

  it("shows start conversation prompt when project selected but no messages", () => {
    useChatStore.setState({ activeProjectId: "p1" });
    render(<ChatPanel onSend={() => {}} />);
    expect(screen.getByText(/Start a conversation/)).toBeInTheDocument();
  });

  it("renders messages from store", () => {
    useChatStore.setState({
      activeProjectId: "p1",
      messages: [
        { id: "1", project_id: "p1", phase: "identity", role: "user", content: "hello", created_at: "" },
        { id: "2", project_id: "p1", phase: "identity", role: "assistant", content: "hi", created_at: "" },
      ],
    });
    render(<ChatPanel onSend={() => {}} />);
    expect(screen.getByText("hello")).toBeInTheDocument();
    expect(screen.getByText("hi")).toBeInTheDocument();
  });
});
