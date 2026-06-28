import { useCallback } from "react";
import { useChatStore } from "@/lib/store/chat";
import { SendMessage, SendMessageStream, AddMessage, AnalyzePath, AnalyzePathContent, FetchURL, SetProjectMeta } from "@wails/go/main/App";
import { EventsOn, EventsOff } from "@wails/runtime";

const PATH_PATTERN = /(~\/\S+|(?:\/Users\/|\/home\/|\/tmp\/)\S+)/g;

const SYSTEM_PROMPT = `You are Ignite's provisioning assistant. Follow this 5-phase interview.

Phase 1: Identity & Vision — project name, tagline, description, target audience.
  When confirmed, output: {"project":{"name":"...","tagline":"..."}}

Phase 2: Tech Stack — frontend, backend, database, auth, hosting. Offer guidance.

Phase 3: Features & Architecture — core features, API design, DB schema.

Phase 4: Roadmap & Quality — phases 0-4, performance targets, risks.

Phase 5: Generation — output the full project spec ready for template rendering.

Based on the conversation progress, determine the current phase and continue.
Ask one question at a time. Be concise. Use the JSON format to report decisions.`;

const PHASE_LABELS = ["identity", "tech-stack", "features", "roadmap", "generation"];

const PROJ_JSON = /\{"project"\s*:\s*\{"name"\s*:\s*"([^"]+)"\s*,\s*"tagline"\s*:\s*"([^"]*)"\}\}/;

function extractPaths(content: string): string[] {
  const matches = content.match(PATH_PATTERN);
  return matches ? [...new Set(matches.map((m) => m.replace(/[?!.,;:]+$/, "")))] : [];
}

export function useConversation() {
  const addMessage = useChatStore((s) => s.addMessage);
  const appendStreamChunk = useChatStore((s) => s.appendStreamChunk);
  const startStreaming = useChatStore((s) => s.startStreaming);
  const finishStreaming = useChatStore((s) => s.finishStreaming);
  const activeProjectId = useChatStore((s) => s.activeProjectId);
  const messages = useChatStore((s) => s.messages);
  const setProjectName = useChatStore((s) => s.setProjectName);

  const sendMessage = useCallback(
    async (provider: string, model: string, content: string) => {
      const userMsgCount = messages.filter((m) => m.role === "user").length;
      const phaseIdx = Math.min(Math.floor(userMsgCount / 3), 4);
      const phase = PHASE_LABELS[phaseIdx] || "generation";
      const userMsg = {
        id: crypto.randomUUID(),
        project_id: activeProjectId ?? "",
        phase,
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
        { role: "system", content: SYSTEM_PROMPT },
        ...messages.map((m) => ({ role: m.role, content: m.content })),
        { role: "user", content },
      ];

      if (messages.length === 0) {
        const allProjects = useChatStore.getState().projects;
        const others = allProjects.filter((p) => p.id !== activeProjectId && p.name);
        if (others.length > 0) {
          const recent = others.slice(0, 5).map((p) =>
            `- ${p.name}${p.tagline ? ": " + p.tagline : ""}`
          ).join("\n");
          apiMessages.splice(1, 0, {
            role: "system",
            content: `User's previous projects:\n${recent}\n\nReference these when relevant to the user's current project.`,
          });
        }
      }

      const paths = extractPaths(content);
      let scans = "";
      const urlMatches = content.match(/https?:\/\/\S+/g);
      let urlContent = "";
      if (urlMatches) {
        const urls = [...new Set(urlMatches)].map((u) => u.replace(/[.,;:!?]+$/, ""));
        const urlResults = await Promise.allSettled(
          urls.map((u) => FetchURL(u))
        );
        urlContent = urlResults
          .map((r) => (r.status === "fulfilled" ? r.value : ""))
          .filter(Boolean)
          .join("\n\n");
      }
      if (paths.length > 0) {
        const results = await Promise.allSettled(
          paths.map((p) => AnalyzePath(p))
        );
        const contentResults = await Promise.allSettled(
          paths.map((p) => AnalyzePathContent(p))
        );
        scans = (results
          .map((r) => (r.status === "fulfilled" ? r.value : ""))
          .filter(Boolean)
          .join("\n\n---\n\n"));
        const contents = contentResults
          .map((r) => (r.status === "fulfilled" ? r.value : ""))
          .filter(Boolean)
          .join("\n\n");
        if (scans || urlContent) {
          apiMessages.splice(-1, 0, {
            role: "system",
            content: `Additional context from user message:\n\n${scans}${contents ? "\n\nFile contents:\n\n" + contents : ""}${urlContent ? "\n\nURL content:\n\n" + urlContent : ""}`,
          });
        }
      }

      try {
        await SendMessageStream(provider, model, apiMessages);
        const streamedContent = useChatStore.getState().streamingContent;
        finishStreaming(streamedContent);

        const match = streamedContent.match(PROJ_JSON);
        if (match && activeProjectId) {
          const name = match[1];
          const tagline = match[2] || "";
          setProjectName(activeProjectId, name, tagline);
          SetProjectMeta(activeProjectId, name, tagline).catch(() => {});
        }

        if (activeProjectId) {
          const finalText = streamedContent || "[empty response]";
          const assistantMsg = {
            id: crypto.randomUUID(),
            project_id: activeProjectId,
            phase,
            role: "assistant" as const,
            content: finalText,
            created_at: new Date().toISOString(),
          };
          AddMessage(assistantMsg).catch(() => {});
        }
      } catch (err) {
        useChatStore.setState({ isStreaming: false });
        try {
          const resp = await SendMessage(provider, model, apiMessages);
          finishStreaming(resp.content || "[empty response]");
          if (activeProjectId) {
            AddMessage({
              id: crypto.randomUUID(),
              project_id: activeProjectId,
              phase,
              role: "assistant",
              content: resp.content || "[empty response]",
              created_at: new Date().toISOString(),
            }).catch(() => {});
          }
        } catch {}
      }
    },
    [activeProjectId, messages, addMessage, startStreaming, finishStreaming, setProjectName]
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
