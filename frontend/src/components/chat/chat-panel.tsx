import { useEffect, useRef, useCallback, useState } from "react";
import { ArrowDown } from "lucide-react";
import { useChatStore } from "@/lib/store/chat";
import { ChatBubble } from "./chat-bubble";
import { ChatInput } from "./chat-input";

interface ChatPanelProps {
  onSend: (content: string) => void;
  providerReady?: boolean;
  avatar?: string;
  userName?: string;
}

export function ChatPanel({ onSend, providerReady = false, avatar, userName }: ChatPanelProps) {
  const messages = useChatStore((s) => s.messages);
  const isStreaming = useChatStore((s) => s.isStreaming);
  const streamingContent = useChatStore((s) => s.streamingContent);
  const activeProjectId = useChatStore((s) => s.activeProjectId);
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
    <div className="relative flex flex-1 flex-col bg-background">
      <div ref={scrollRef} className="flex-1 overflow-y-auto">
        {messages.length === 0 && !isStreaming && (
          <div className="flex h-full items-center justify-center">
            {!activeProjectId ? (
              <div className="max-w-sm text-center space-y-3">
                <h2 className="font-mono text-lg text-text-primary">Welcome to Ignite</h2>
                <p className="font-mono text-sm text-text-secondary leading-relaxed">
                  Provisioning with a heartbeat. Create a new project to begin
                  a guided AI conversation that will produce a full spec,
                  agent guide, implementation plan, and README.
                </p>
                {!providerReady && (
                  <p className="font-mono text-xs text-text-secondary/70">
                    Configure an AI provider in Settings to get started.
                  </p>
                )}
              </div>
            ) : (
              <p className="font-mono text-sm text-text-secondary">
                Start a conversation to provision your project.
              </p>
            )}
          </div>
        )}
        {messages.map((msg) => (
          <ChatBubble key={msg.id} message={msg} avatar={avatar} userName={userName} />
        ))}
        {isStreaming && (
          <ChatBubble
            message={{
              id: "stream",
              project_id: "",
              phase: "",
              role: "assistant",
              content: streamingContent,
              created_at: new Date().toISOString(),
            }}
            isStreaming
            avatar={avatar}
            userName={userName}
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

      <ChatInput onSend={onSend} disabled={isStreaming || !activeProjectId || !providerReady} />
    </div>
  );
}
