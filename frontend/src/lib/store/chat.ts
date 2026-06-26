import { create } from "zustand";

export type Role = string;

export interface Message {
  id: string;
  project_id: string;
  phase: string;
  role: Role;
  content: string;
  created_at: string;
}

export interface Project {
  id: string;
  name: string;
  tagline: string;
  path: string;
  provider: string;
  model: string;
  created_at: string;
  updated_at: string;
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
        project_id: s.activeProjectId ?? "",
        phase: "",
        role: "assistant",
        content,
        created_at: new Date().toISOString(),
      }],
      streamingContent: "",
      isStreaming: false,
    })),
}));
