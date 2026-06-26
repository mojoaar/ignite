import { useEffect, useState, useCallback } from "react";
import { useTheme } from "@/hooks/useTheme";
import { useConversation } from "@/hooks/useConversation";
import { useChatStore } from "@/lib/store/chat";
import { Sidebar } from "@/components/sidebar/sidebar";
import { ChatPanel } from "@/components/chat/chat-panel";
import { StatusBar } from "@/components/status-bar/status-bar";
import { SettingsModal } from "@/components/settings/settings-modal";
import { ExportChat } from "@wails/go/main/App";

function App() {
  useTheme();
  const { sendMessage, subscribeToStream } = useConversation();
  const activeProjectId = useChatStore((s) => s.activeProjectId);
  const messages = useChatStore((s) => s.messages);

  const [provider, setProvider] = useState("opencode-go");
  const [model, setModel] = useState("gpt-4o");
  const [settingsOpen, setSettingsOpen] = useState(false);

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
        <ChatPanel onSend={handleSend} />
      </div>
      <StatusBar
        provider={provider}
        model={model}
        onProviderChange={setProvider}
        onModelChange={setModel}
        onOpenSettings={() => setSettingsOpen(true)}
        onExport={handleExport}
      />
      <SettingsModal
        open={settingsOpen}
        onClose={() => setSettingsOpen(false)}
      />
    </div>
  );
}
export default App;
