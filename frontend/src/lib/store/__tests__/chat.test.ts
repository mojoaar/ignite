import { describe, it, expect, beforeEach } from "vitest";
import { useChatStore } from "../chat";
import type { Message, Project } from "../chat";

function makeMsg(overrides: Partial<Message> = {}): Message {
  return {
    id: "m1",
    project_id: "p1",
    phase: "identity",
    role: "user",
    content: "hello",
    created_at: new Date().toISOString(),
    ...overrides,
  };
}

function makeProj(overrides: Partial<Project> = {}): Project {
  return {
    id: "p1",
    name: "Test Project",
    tagline: "A test",
    path: "",
    provider: "",
    model: "",
    created_at: new Date().toISOString(),
    updated_at: new Date().toISOString(),
    ...overrides,
  };
}

describe("chat store", () => {
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

  describe("setProjects", () => {
    it("sets the project list", () => {
      useChatStore.getState().setProjects([makeProj()]);
      expect(useChatStore.getState().projects).toHaveLength(1);
    });
  });

  describe("setActiveProject", () => {
    it("sets active project and clears messages and streaming", () => {
      useChatStore.setState({ messages: [makeMsg()], streamingContent: "stream..." });
      useChatStore.getState().setActiveProject("p1");
      expect(useChatStore.getState().activeProjectId).toBe("p1");
      expect(useChatStore.getState().messages).toEqual([]);
      expect(useChatStore.getState().streamingContent).toBe("");
    });
  });

  describe("setActiveProjectPath", () => {
    it("sets the active project path", () => {
      useChatStore.getState().setActiveProjectPath("/tmp/test");
      expect(useChatStore.getState().activeProjectPath).toBe("/tmp/test");
    });
  });

  describe("setProjectName", () => {
    it("updates name and tagline of matching project", () => {
      useChatStore.setState({ projects: [makeProj({ id: "p1", name: "Old" }), makeProj({ id: "p2", name: "Other" })] });
      useChatStore.getState().setProjectName("p1", "New Name", "New Tagline");
      const proj = useChatStore.getState().projects.find((p) => p.id === "p1")!;
      expect(proj.name).toBe("New Name");
      expect(proj.tagline).toBe("New Tagline");
      expect(useChatStore.getState().projects.find((p) => p.id === "p2")!.name).toBe("Other");
    });
  });

  describe("removeProject", () => {
    it("removes the project from the list", () => {
      useChatStore.setState({ projects: [makeProj({ id: "p1" }), makeProj({ id: "p2" })] });
      useChatStore.getState().removeProject("p1");
      expect(useChatStore.getState().projects).toHaveLength(1);
      expect(useChatStore.getState().projects[0].id).toBe("p2");
    });

    it("clears activeProjectId if the removed project was active", () => {
      useChatStore.setState({ projects: [makeProj({ id: "p1" })], activeProjectId: "p1", messages: [makeMsg()] });
      useChatStore.getState().removeProject("p1");
      expect(useChatStore.getState().activeProjectId).toBeNull();
      expect(useChatStore.getState().messages).toEqual([]);
    });

    it("keeps activeProjectId if a different project was removed", () => {
      useChatStore.setState({ projects: [makeProj({ id: "p1" }), makeProj({ id: "p2" })], activeProjectId: "p1" });
      useChatStore.getState().removeProject("p2");
      expect(useChatStore.getState().activeProjectId).toBe("p1");
    });
  });

  describe("addMessage", () => {
    it("appends a message", () => {
      useChatStore.getState().addMessage(makeMsg());
      expect(useChatStore.getState().messages).toHaveLength(1);
    });
  });

  describe("setMessages", () => {
    it("replaces all messages", () => {
      useChatStore.setState({ messages: [makeMsg()] });
      useChatStore.getState().setMessages([makeMsg({ id: "m2" })]);
      expect(useChatStore.getState().messages).toHaveLength(1);
      expect(useChatStore.getState().messages[0].id).toBe("m2");
    });
  });

  describe("streaming", () => {
    it("startStreaming clears content and sets streaming flag", () => {
      useChatStore.setState({ streamingContent: "old" });
      useChatStore.getState().startStreaming();
      expect(useChatStore.getState().streamingContent).toBe("");
      expect(useChatStore.getState().isStreaming).toBe(true);
    });

    it("appendStreamChunk accumulates content", () => {
      useChatStore.setState({ streamingContent: "Hello" });
      useChatStore.getState().appendStreamChunk(" world");
      expect(useChatStore.getState().streamingContent).toBe("Hello world");
    });

    it("finishStreaming appends assistant message and resets streaming", () => {
      useChatStore.setState({ activeProjectId: "p1", streamingContent: "Hello world", isStreaming: true, messages: [] });
      useChatStore.getState().finishStreaming("Hello world");
      expect(useChatStore.getState().isStreaming).toBe(false);
      expect(useChatStore.getState().streamingContent).toBe("");
      expect(useChatStore.getState().messages).toHaveLength(1);
      const msg = useChatStore.getState().messages[0];
      expect(msg.role).toBe("assistant");
      expect(msg.content).toBe("Hello world");
      expect(msg.project_id).toBe("p1");
    });
  });
});
