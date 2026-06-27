import { useEffect, useState, useCallback } from "react";
import { useTheme } from "@/hooks/useTheme";
import { useConversation } from "@/hooks/useConversation";
import { useChatStore } from "@/lib/store/chat";
import { Sidebar } from "@/components/sidebar/sidebar";
import { ChatPanel } from "@/components/chat/chat-panel";
import { StatusBar } from "@/components/status-bar/status-bar";
import { SettingsModal } from "@/components/settings/settings-modal";
import { ErrorBoundary } from "@/components/ErrorBoundary";
import { ExportChat, HasAPIKey, GetSettings, SelectDirectory, GetCachedModels, GenerateProjectFiles } from "@wails/go/main/App";
import { EventsOn, EventsOff } from "@wails/runtime";

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
  const [avatar, setAvatar] = useState("");
  const [userName, setUserName] = useState("");

  useEffect(() => {
    GetSettings().then((s) => {
      setProjectDir(s.default_project_dir);
      setAvatar(s.avatar || "");
      setUserName(s.name || "");
    }).catch(() => {});
  }, []);

  useEffect(() => {
    setProviderReady(false);
    HasAPIKey(provider).then(setProviderReady).catch(() => setProviderReady(false));
  }, [provider]);

  useEffect(() => {
    GetSettings().then((s) => {
      const saved = s.providers?.[provider]?.default_model;
      if (saved) {
        setModel(saved);
      } else {
        GetCachedModels(provider).then((m) => {
          if (m && m.length > 0) setModel(m[0].model_id);
        }).catch(() => {});
      }
    }).catch(() => {});
  }, [provider]);

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

  const handleGenerate = useCallback(async () => {
    if (!activeProjectId) return;
    const msgs = useChatStore.getState().messages;
    if (msgs.length === 0) return;
    const name = useChatStore.getState().projects.find((p) => p.id === activeProjectId)?.name || "project";
    try {
      await GenerateProjectFiles(provider, model, name, msgs.map((m) => ({ role: m.role, content: m.content })));
      alert("Files generated in ~/Development/" + name + "/");
    } catch (e) {
      alert("Generation failed: " + (e instanceof Error ? e.message : String(e)));
    }
  }, [activeProjectId, provider, model]);

  useEffect(() => {
    EventsOn("menu-export", () => { handleExport(); });
    EventsOn("menu-settings", () => { setSettingsOpen(true); });
    EventsOn("menu-open-folder", () => {
      SelectDirectory().then((dir: string) => setProjectDir(dir)).catch(() => {});
    });
    return () => {
      EventsOff("menu-export");
      EventsOff("menu-settings");
      EventsOff("menu-open-folder");
    };
  }, [handleExport]);

  return (
    <ErrorBoundary>
    <div className="flex h-screen flex-col bg-background text-text-primary">
      <div className="flex flex-1 overflow-hidden">
        <Sidebar />
        <ChatPanel onSend={handleSend} providerReady={providerReady} avatar={avatar} userName={userName} />
      </div>
      <StatusBar
        provider={provider}
        model={model}
        projectDir={activeProjectPath || projectDir}
        onProviderChange={setProvider}
        onModelChange={setModel}
        onOpenSettings={() => setSettingsOpen(true)}
        onExport={handleExport}
        onGenerate={activeProjectId ? handleGenerate : undefined}
      />
      <SettingsModal
        open={settingsOpen}
        onClose={() => {
          setSettingsOpen(false);
          HasAPIKey(provider).then(setProviderReady).catch(() => setProviderReady(false));
          GetSettings().then((s) => {
            const saved = s.providers?.[provider]?.default_model;
            if (saved) setModel(saved);
            setAvatar(s.avatar || "");
            setUserName(s.name || "");
          }).catch(() => {});
        }}
      />
    </div>
    </ErrorBoundary>
  );
}
export default App;
