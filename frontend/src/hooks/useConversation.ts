import { useCallback } from "react";
import { useChatStore } from "@/lib/store/chat";
import { SendMessageStream, AddMessage } from "@wails/go/main/App";
import { EventsOn, EventsOff } from "@wails/runtime";

export function useConversation() {
  const addMessage = useChatStore((s) => s.addMessage);
  const appendStreamChunk = useChatStore((s) => s.appendStreamChunk);
  const startStreaming = useChatStore((s) => s.startStreaming);
  const finishStreaming = useChatStore((s) => s.finishStreaming);
  const activeProjectId = useChatStore((s) => s.activeProjectId);
  const messages = useChatStore((s) => s.messages);

  const sendMessage = useCallback(
    async (provider: string, model: string, content: string) => {
      const userMsg = {
        id: crypto.randomUUID(),
        project_id: activeProjectId ?? "",
        phase: "",
        role: "user" as const,
        content,
        created_at: new Date().toISOString(),
      };

      addMessage(userMsg);

      if (activeProjectId) {
        AddMessage(userMsg).catch(() => {});
      }

      startStreaming();

      const apiMessages = [
        ...messages.map((m) => ({ role: m.role, content: m.content })),
        { role: "user", content },
      ];

      try {
        await SendMessageStream(provider, model, apiMessages);
        const streamedContent = useChatStore.getState().streamingContent;
        finishStreaming(streamedContent);

        if (activeProjectId) {
          const finalText = streamedContent || "[empty response]";
          const assistantMsg = {
            id: crypto.randomUUID(),
            project_id: activeProjectId,
            phase: "",
            role: "assistant" as const,
            content: finalText,
            created_at: new Date().toISOString(),
          };
          AddMessage(assistantMsg).catch(() => {});
        }
      } catch (err) {
        useChatStore.setState({ isStreaming: false });
      }
    },
    [activeProjectId, messages, addMessage, startStreaming, finishStreaming]
  );

  const subscribeToStream = useCallback(() => {
    EventsOn("stream-chunk", (chunk: string) => {
      appendStreamChunk(chunk);
    });

    return () => {
      EventsOff("stream-chunk");
    };
  }, [appendStreamChunk]);

  return { sendMessage, subscribeToStream };
}
