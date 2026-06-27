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
  activeProjectPath: string;
  messages: Message[];
  streamingContent: string;
  isStreaming: boolean;

  setProjects: (projects: Project[]) => void;
  setProjectName: (id: string, name: string, tagline: string) => void;
  removeProject: (id: string) => void;
  setActiveProject: (id: string) => void;
  setActiveProjectPath: (path: string) => void;
  addMessage: (msg: Message) => void;
  setMessages: (msgs: Message[]) => void;
  appendStreamChunk: (chunk: string) => void;
  startStreaming: () => void;
  finishStreaming: (content: string) => void;
}

export const useChatStore = create<ChatState>((set) => ({
  projects: [],
  activeProjectId: null,
  activeProjectPath: "",
  messages: [],
  streamingContent: "",
  isStreaming: false,

  setProjects: (projects) => set({ projects }),
  setProjectName: (id, name, tagline) =>
    set((s) => ({
      projects: s.projects.map((p) =>
        p.id === id ? { ...p, name, tagline } : p
      ),
    })),
  removeProject: (id) =>
    set((s) => ({
      projects: s.projects.filter((p) => p.id !== id),
      activeProjectId: s.activeProjectId === id ? null : s.activeProjectId,
      messages: s.activeProjectId === id ? [] : s.messages,
    })),
  setActiveProject: (id) => set({ activeProjectId: id, messages: [], streamingContent: "" }),
  setActiveProjectPath: (path) => set({ activeProjectPath: path }),
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
