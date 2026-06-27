import { useEffect, useState, useCallback } from "react";
import { useTheme } from "@/hooks/useTheme";
import { useConversation } from "@/hooks/useConversation";
import { useChatStore } from "@/lib/store/chat";
import { Sidebar } from "@/components/sidebar/sidebar";
import { ChatPanel } from "@/components/chat/chat-panel";
import { StatusBar } from "@/components/status-bar/status-bar";
import { SettingsModal } from "@/components/settings/settings-modal";
import { ExportChat, HasAPIKey, GetSettings, SelectDirectory } from "@wails/go/main/App";
import { EventsOn } from "@wails/runtime";

function App() {
  useTheme();
  const { sendMessage, subscribeToStream } = useConversation();
  const activeProjectId = useChatStore((s) => s.activeProjectId);
  const activeProjectPath = useChatStore((s) => s.activeProjectPath);
  const messages = useChatStore((s) => s.messages);

  const [provider, setProvider] = useState("opencode-go");
  const [model, setModel] = useState("");
  const [settingsOpen, setSettingsOpen] = useState(false);
  const [providerReady, setProviderReady] = useState(false);
  const [projectDir, setProjectDir] = useState("");

  useEffect(() => {
    GetSettings().then((s) => setProjectDir(s.default_project_dir)).catch(() => {});
  }, []);

  useEffect(() => {
    setProviderReady(false);
    HasAPIKey(provider).then(setProviderReady).catch(() => setProviderReady(false));
  }, [provider]);

  useEffect(() => {
    GetSettings().then((s) => {
      const saved = s.providers?.[provider]?.default_model;
      if (saved) setModel(saved);
    }).catch(() => {});
  }, [provider]);

  useEffect(() => {
    const unsubExport = EventsOn("menu-export", () => { handleExport(); });
    const unsubSettings = EventsOn("menu-settings", () => { setSettingsOpen(true); });
    const unsubFolder = EventsOn("menu-open-folder", () => {
      SelectDirectory().then((dir: string) => setProjectDir(dir)).catch((e: unknown) => console.error("SelectDirectory failed:", e));
    });
    return () => {
      typeof unsubExport === "function" && unsubExport();
      typeof unsubSettings === "function" && unsubSettings();
      typeof unsubFolder === "function" && unsubFolder();
    };
  }, [messages]);

  useEffect(() => {
    const cleanup = subscribeToStream();
    return cleanup;
  }, [subscribeToStream]);

  const handleSend = useCallback(
    (content: string) => {
      if (!activeProjectId) return;
      sendMessage(provider, model, content);
    },
    [provider, model, activeProjectId, sendMessage]
  );

  const handleExport = useCallback(async () => {
    const md = await ExportChat(messages.map((m) => ({
      ...m,
      role: m.role,
    })));
    const blob = new Blob([md], { type: "text/markdown" });
    const url = URL.createObjectURL(blob);
    const a = document.createElement("a");
    a.href = url;
    a.download = "conversation.md";
    a.click();
    URL.revokeObjectURL(url);
  }, [messages]);

  return (
    <div className="flex h-screen flex-col bg-background text-text-primary">
      <div className="flex flex-1 overflow-hidden">
        <Sidebar />
        <ChatPanel onSend={handleSend} providerReady={providerReady} />
      </div>
      <StatusBar
        provider={provider}
        model={model}
        projectDir={activeProjectPath || projectDir}
        onProviderChange={setProvider}
        onModelChange={setModel}
        onOpenSettings={() => setSettingsOpen(true)}
        onExport={handleExport}
      />
      <SettingsModal
        open={settingsOpen}
        onClose={() => {
          setSettingsOpen(false);
          HasAPIKey(provider).then(setProviderReady).catch(() => setProviderReady(false));
          GetSettings().then((s) => {
            const saved = s.providers?.[provider]?.default_model;
            if (saved) setModel(saved);
          }).catch(() => {});
        }}
      />
    </div>
  );
}
export default App;
