import { describe, it, expect, vi } from "vitest";
import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { ChatInput } from "@/components/chat/chat-input";

describe("ChatInput", () => {
  it("sends on Enter", async () => {
    const onSend = vi.fn();
    render(<ChatInput onSend={onSend} />);
    await userEvent.type(screen.getByPlaceholderText(/Type your response/), "hello{Enter}");
    expect(onSend).toHaveBeenCalledWith("hello");
  });

  it("does not send on Shift+Enter", async () => {
    const onSend = vi.fn();
    render(<ChatInput onSend={onSend} />);
    await userEvent.type(screen.getByPlaceholderText(/Type your response/), "hello{Shift>}{Enter}{/Shift}");
    expect(onSend).not.toHaveBeenCalled();
  });

  it("does not send empty text", async () => {
    const onSend = vi.fn();
    render(<ChatInput onSend={onSend} />);
    await userEvent.type(screen.getByPlaceholderText(/Type your response/), "   {Enter}");
    expect(onSend).not.toHaveBeenCalled();
  });

  it("disables input and button when disabled prop is true", () => {
    render(<ChatInput onSend={vi.fn()} disabled />);
    const ta = screen.getByPlaceholderText(/Type your response/) as HTMLTextAreaElement;
    expect(ta.disabled).toBe(true);
  });
});
