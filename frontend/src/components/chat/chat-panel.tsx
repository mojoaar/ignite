import { useEffect, useRef, useCallback, useState } from "react";
import { ArrowDown } from "lucide-react";
import { useChatStore } from "@/lib/store/chat";
import { ChatBubble } from "./chat-bubble";
import { ChatInput } from "./chat-input";

interface ChatPanelProps {
  onSend: (content: string) => void;
}

export function ChatPanel({ onSend }: ChatPanelProps) {
  const messages = useChatStore((s) => s.messages);
  const isStreaming = useChatStore((s) => s.isStreaming);
  const streamingContent = useChatStore((s) => s.streamingContent);
  const scrollRef = useRef<HTMLDivElement>(null);
  const [showScrollButton, setShowScrollButton] = useState(false);

  useEffect(() => {
    const el = scrollRef.current;
    if (!el) return;
    const onScroll = () => {
      const atBottom = el.scrollHeight - el.scrollTop - el.clientHeight < 80;
      setShowScrollButton(!atBottom);
    };
    el.addEventListener("scroll", onScroll);
    return () => el.removeEventListener("scroll", onScroll);
  }, []);

  const scrollToBottom = useCallback(() => {
    if (scrollRef.current) {
      scrollRef.current.scrollTop = scrollRef.current.scrollHeight;
    }
  }, []);

  useEffect(() => {
    if (!showScrollButton) {
      scrollToBottom();
    }
  }, [messages, streamingContent, showScrollButton, scrollToBottom]);

  return (
    <div className="flex flex-1 flex-col bg-background">
      <div ref={scrollRef} className="flex-1 overflow-y-auto">
        {messages.length === 0 && !isStreaming && (
          <div className="flex h-full items-center justify-center">
            <p className="font-mono text-sm text-text-secondary">
              Start a conversation to provision your project.
            </p>
          </div>
        )}
        {messages.map((msg) => (
          <ChatBubble key={msg.id} message={msg} />
        ))}
        {isStreaming && streamingContent && (
          <ChatBubble
            message={{
              id: "stream",
              projectId: "",
              phase: "",
              role: "assistant",
              content: streamingContent,
              createdAt: new Date().toISOString(),
            }}
            isStreaming
          />
        )}
      </div>

      {showScrollButton && (
        <button
          onClick={scrollToBottom}
          className="absolute bottom-20 right-6 z-10 flex h-8 w-8 items-center justify-center rounded-full border border-border bg-surface shadow-lg hover:bg-surface-hover"
        >
          <ArrowDown className="h-4 w-4 text-text-secondary" />
        </button>
      )}

      <ChatInput onSend={onSend} disabled={isStreaming} />
    </div>
  );
}
