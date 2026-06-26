import { create } from "zustand";

export type Role = "user" | "assistant" | "system";

export interface Message {
  id: string;
  projectId: string;
  phase: string;
  role: Role;
  content: string;
  createdAt: string;
}

interface Project {
  id: string;
  name: string;
  tagline: string;
  path: string;
  provider: string;
  model: string;
  createdAt: string;
  updatedAt: string;
}

interface ChatState {
  projects: Project[];
  activeProjectId: string | null;
  messages: Message[];
  streamingContent: string;
  isStreaming: boolean;

  setProjects: (projects: Project[]) => void;
  setActiveProject: (id: string) => void;
  addMessage: (msg: Message) => void;
  setMessages: (msgs: Message[]) => void;
  appendStreamChunk: (chunk: string) => void;
  startStreaming: () => void;
  finishStreaming: (content: string) => void;
}

export const useChatStore = create<ChatState>((set) => ({
  projects: [],
  activeProjectId: null,
  messages: [],
  streamingContent: "",
  isStreaming: false,

  setProjects: (projects) => set({ projects }),
  setActiveProject: (id) => set({ activeProjectId: id, messages: [], streamingContent: "" }),
  addMessage: (msg) =>
    set((s) => ({ messages: [...s.messages, msg] })),
  setMessages: (msgs) => set({ messages: msgs }),
  appendStreamChunk: (chunk) =>
    set((s) => ({ streamingContent: s.streamingContent + chunk })),
  startStreaming: () => set({ streamingContent: "", isStreaming: true }),
  finishStreaming: (content) =>
    set((s) => ({
      messages: [...s.messages, {
        id: crypto.randomUUID(),
        projectId: s.activeProjectId ?? "",
        phase: "",
        role: "assistant",
        content,
        createdAt: new Date().toISOString(),
      }],
      streamingContent: "",
      isStreaming: false,
    })),
}));
