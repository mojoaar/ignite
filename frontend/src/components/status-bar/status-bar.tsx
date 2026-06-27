import { useState, useEffect } from "react";
import { Settings, Download, CheckCircle, XCircle } from "lucide-react";
import { GetCachedModels, HasAPIKey } from "@wails/go/main/App";

const PROVIDERS = [
  { id: "opencode-go", label: "OpenCode Go" },
  { id: "opencode-zen", label: "OpenCode Zen" },
  { id: "deepseek", label: "DeepSeek" },
];

interface StatusBarProps {
  provider: string;
  model: string;
  projectDir?: string;
  onProviderChange: (provider: string) => void;
  onModelChange: (model: string) => void;
  onOpenSettings: () => void;
  onExport: () => void;
}

export function StatusBar({
  provider,
  model,
  projectDir,
  onProviderChange,
  onModelChange,
  onOpenSettings,
  onExport,
}: StatusBarProps) {
  const [models, setModels] = useState<{ id: string; label: string }[]>([]);
  const [connected, setConnected] = useState<boolean | null>(null);

  useEffect(() => {
    HasAPIKey(provider)
      .then((hasKey) => setConnected(hasKey))
      .catch(() => setConnected(false));
  }, [provider]);

  useEffect(() => {
    GetCachedModels(provider)
      .then((m) => {
        if (m && m.length > 0) {
          const list = (m as { model_id: string; display_name: string }[]).map((x) => ({
            id: x.model_id,
            label: x.display_name || x.model_id,
          }));
          setModels(list);
          const found = list.find((x) => x.id === model);
          if (!found && list.length > 0) {
            onModelChange(list[0].id);
          }
        } else {
          setModels([]);
        }
      })
      .catch(() => setModels([]));
  }, [provider]);

  return (
    <div className="flex h-10 shrink-0 items-center gap-3 border-t border-border bg-surface px-4">
      {projectDir && (
        <span className="font-mono text-[11px] text-text-secondary truncate max-w-[200px]" title={projectDir}>
          {projectDir.replace(/^\/Users\/[^/]+\//, "~/")}
        </span>
      )}
      <select
        value={provider}
        onChange={(e) => onProviderChange(e.target.value)}
        className="h-7 rounded border border-border bg-background px-2 font-mono text-xs text-text-primary focus:border-accent focus:outline-none"
      >
        {PROVIDERS.map((p) => (
          <option key={p.id} value={p.id}>
            {p.label}
          </option>
        ))}
      </select>

      <select
        value={model}
        onChange={(e) => onModelChange(e.target.value)}
        className="h-7 rounded border border-border bg-background px-2 font-mono text-xs text-text-primary focus:border-accent focus:outline-none"
      >
        {models.map((m) => (
          <option key={m.id} value={m.id}>
            {m.label}
          </option>
        ))}
      </select>

      <span className="text-text-secondary" title={connected ? "Connected" : connected === false ? "Disconnected" : "Checking..."}>
        {connected === true ? (
          <CheckCircle className="h-4 w-4 text-success" />
        ) : connected === false ? (
          <XCircle className="h-4 w-4 text-error" />
        ) : (
          <span className="flex h-4 w-4 items-center justify-center text-xs">...</span>
        )}
      </span>

      <div className="flex-1" />

      <button
        onClick={onOpenSettings}
        className="rounded p-1 text-text-secondary hover:bg-surface-hover hover:text-text-primary"
        title="Settings"
      >
        <Settings className="h-4 w-4" />
      </button>

      <button
        onClick={onExport}
        className="rounded p-1 text-text-secondary hover:bg-surface-hover hover:text-text-primary"
        title="Export chat"
      >
        <Download className="h-4 w-4" />
      </button>
    </div>
  );
}
