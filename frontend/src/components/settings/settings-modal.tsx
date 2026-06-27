import { useState, useEffect, useCallback } from "react";
import { Eye, EyeOff, Check, Loader2, XCircle } from "lucide-react";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import {
  GetSettings,
  SaveSettings,
  SetAPIKey,
  HasAPIKey,
  ValidateProviderKey,
  ListProviderModels,
} from "@wails/go/main/App";
import { useThemeStore } from "@/lib/store/theme";
import { cn } from "@/lib/utils";

type Tab = "providers" | "appearance";

interface ProviderFields {
  endpoint: string;
  defaultModel: string;
  apiKey: string;
  showKey: boolean;
  hasKey: boolean;
  validating: boolean;
  valid: boolean | null;
}

interface SettingsData {
  providers: Record<string, { endpoint: string; default_model: string }>;
  default_provider: string;
  appearance: string;
  default_license: string;
  default_project_dir: string;
  font: string;
}

const PROVIDER_IDS = ["opencode-go", "opencode-zen", "claude", "deepseek"];

const PROVIDER_LABELS: Record<string, string> = {
  "opencode-go": "OpenCode Go",
  "opencode-zen": "OpenCode Zen",
  claude: "Claude",
  deepseek: "DeepSeek",
};

const PROVIDER_ENDPOINTS: Record<string, string> = {
  "opencode-go": "https://opencode.ai/zen/go/v1",
  "opencode-zen": "https://opencode.ai/zen/v1",
  claude: "https://api.anthropic.com/v1/messages",
  deepseek: "https://api.deepseek.com/v1",
};

const PROVIDER_DEFAULT_MODELS: Record<string, string> = {
  "opencode-go": "",
  "opencode-zen": "",
  claude: "",
  deepseek: "",
};

interface SettingsModalProps {
  open: boolean;
  onClose: () => void;
}

export function SettingsModal({ open, onClose }: SettingsModalProps) {
  const [tab, setTab] = useState<Tab>("providers");
  const [settings, setSettings] = useState<SettingsData | null>(null);
  const [providerFields, setProviderFields] = useState<Record<string, ProviderFields>>({});
  const [saving, setSaving] = useState(false);
  const [selectedProvider, setSelectedProvider] = useState("opencode-go");
  const [providerModels, setProviderModels] = useState<Record<string, string[]>>({});
  const [loadingModels, setLoadingModels] = useState(false);

  useEffect(() => {
    const onKeyDown = (e: KeyboardEvent) => { if (e.key === "Escape") onClose(); };
    if (open) {
      document.addEventListener("keydown", onKeyDown);
      return () => document.removeEventListener("keydown", onKeyDown);
    }
  }, [open, onClose]);

  const loadProviderModels = async (provider: string) => {
    if (providerModels[provider]) return;
    setLoadingModels(true);
    try {
      const models = await ListProviderModels(provider);
      setProviderModels((prev) => ({
        ...prev,
        [provider]: (models ?? []).map((m: { id: string }) => m.id),
      }));
    } catch {}
    setLoadingModels(false);
  };

  useEffect(() => {
    loadProviderModels(selectedProvider);
  }, [selectedProvider, open]);

  useEffect(() => {
    if (!open) return;
    GetSettings()
      .then((s) => {
        setSettings(s);
        setSelectedProvider(s.default_provider || "opencode-go");
        const fields: Record<string, ProviderFields> = {};
        for (const pid of PROVIDER_IDS) {
          fields[pid] = {
            endpoint: s.providers?.[pid]?.endpoint ?? PROVIDER_ENDPOINTS[pid] ?? "",
            defaultModel: s.providers?.[pid]?.default_model ?? PROVIDER_DEFAULT_MODELS[pid] ?? "",
            apiKey: "",
            showKey: false,
            hasKey: false,
            validating: false,
            valid: null,
          };
        }
        setProviderFields(fields);
        for (const pid of PROVIDER_IDS) {
          HasAPIKey(pid).then((has) => {
            setProviderFields((prev) => ({
              ...prev,
              [pid]: { ...prev[pid], hasKey: has },
            }));
          }).catch(() => {});
        }
      })
      .catch(() => {});
  }, [open]);

  const setField = useCallback(
    (provider: string, field: keyof ProviderFields, value: unknown) => {
      setProviderFields((prev) => ({
        ...prev,
        [provider]: { ...prev[provider], [field]: value },
      }));
    },
    []
  );

  const handleSave = async () => {
    if (!settings) return;
    setSaving(true);
    try {
      for (const pid of PROVIDER_IDS) {
        const f = providerFields[pid];
        if (f?.apiKey) {
          await SetAPIKey(pid, f.apiKey);
        }
        await SaveSettings({
          ...settings,
          default_provider: selectedProvider,
          providers: {
            ...settings.providers,
            [pid]: { endpoint: f.endpoint, default_model: f.defaultModel },
          },
        } as Parameters<typeof SaveSettings>[0]);
      }
      onClose();
    } catch {}
    setSaving(false);
  };

  const handleValidate = async (provider: string) => {
    const key = providerFields[provider]?.apiKey;
    if (!key) return;
    setField(provider, "validating", true);
    setField(provider, "valid", null);
    try {
      await ValidateProviderKey(provider, key);
      setField(provider, "valid", true);
    } catch {
      setField(provider, "valid", false);
    }
    setField(provider, "validating", false);
  };

  if (!settings) return null;

  return (
    <Dialog open={open} onOpenChange={(v) => { if (!v) onClose(); }}>
      <DialogContent className="max-h-[80vh] max-w-xl overflow-y-auto border-border bg-surface text-text-primary">
        <DialogHeader>
          <DialogTitle className="font-mono text-lg">Settings</DialogTitle>
        </DialogHeader>

        <div className="flex gap-1 rounded-md bg-background p-1">
          {(["providers", "appearance"] as Tab[]).map((t) => (
            <button
              key={t}
              onClick={() => setTab(t)}
              className={cn(
                "flex-1 rounded px-3 py-1.5 text-sm font-medium capitalize transition-colors",
                tab === t
                  ? "bg-accent text-white"
                  : "text-text-secondary hover:text-text-primary"
              )}
            >
              {t}
            </button>
          ))}
        </div>

        {tab === "providers" && (
          <div className="space-y-4">
            <div className="space-y-2">
              <Label className="text-text-secondary">Provider</Label>
              <Select value={selectedProvider} onValueChange={(v) => v && setSelectedProvider(v)}>
                <SelectTrigger className="bg-background">
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  {PROVIDER_IDS.map((pid) => (
                    <SelectItem key={pid} value={pid}>
                      {PROVIDER_LABELS[pid]}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>

            <div className="space-y-3 rounded-md border border-border bg-background p-4">
              <div className="space-y-2">
                <Label className="text-text-secondary" htmlFor="ep">Endpoint</Label>
                <Input
                  id="ep"
                  value={providerFields[selectedProvider]?.endpoint ?? ""}
                  onChange={(e) => setField(selectedProvider, "endpoint", e.target.value)}
                  className="bg-surface font-mono text-sm"
                />
              </div>

              <div className="space-y-2">
                <Label className="text-text-secondary" htmlFor="dm">Default Model</Label>
                <Select
                  value={providerFields[selectedProvider]?.defaultModel ?? ""}
                  onValueChange={(v) => {
                    if (!v) return;
                    setField(selectedProvider, "defaultModel", v);
                  }}
                >
                  <SelectTrigger className="bg-background">
                    <SelectValue placeholder={loadingModels ? "Loading models..." : "Select a model"} />
                  </SelectTrigger>
                  <SelectContent>
                    {(providerModels[selectedProvider] ?? []).map((m) => (
                      <SelectItem key={m} value={m}>
                        {m}
                      </SelectItem>
                    ))}
                    {!providerModels[selectedProvider] && (
                      <SelectItem value={providerFields[selectedProvider]?.defaultModel ?? ""}>
                        {providerFields[selectedProvider]?.defaultModel ?? ""}
                      </SelectItem>
                    )}
                  </SelectContent>
                </Select>
              </div>

              <div className="space-y-2">
                <Label className="text-text-secondary" htmlFor="key">API Key</Label>
                <div className="relative">
                  <Input
                    id="key"
                    type={providerFields[selectedProvider]?.showKey ? "text" : "password"}
                    placeholder={providerFields[selectedProvider]?.hasKey ? "••••••••" : "Enter API key"}
                    value={providerFields[selectedProvider]?.apiKey ?? ""}
                    onChange={(e) => setField(selectedProvider, "apiKey", e.target.value)}
                    className="bg-surface font-mono text-sm pr-16"
                  />
                  <div className="absolute right-1 top-1/2 -translate-y-1/2 flex gap-0.5">
                    <button
                      type="button"
                      onClick={() =>
                        setField(selectedProvider, "showKey", !providerFields[selectedProvider]?.showKey)
                      }
                      className="rounded p-1 text-text-secondary hover:text-text-primary"
                    >
                      {providerFields[selectedProvider]?.showKey ? (
                        <EyeOff className="h-4 w-4" />
                      ) : (
                        <Eye className="h-4 w-4" />
                      )}
                    </button>
                    <button
                      type="button"
                      onClick={() => handleValidate(selectedProvider)}
                      disabled={providerFields[selectedProvider]?.validating || !providerFields[selectedProvider]?.apiKey}
                      className="rounded p-1 text-text-secondary hover:text-text-primary disabled:opacity-50"
                      title="Validate key"
                    >
                      {providerFields[selectedProvider]?.validating ? (
                        <Loader2 className="h-4 w-4 animate-spin" />
                      ) : providerFields[selectedProvider]?.valid === true ? (
                        <Check className="h-4 w-4 text-success" />
                      ) : providerFields[selectedProvider]?.valid === false ? (
                        <XCircle className="h-4 w-4 text-error" />
                      ) : (
                        <Check className="h-4 w-4" />
                      )}
                    </button>
                  </div>
                </div>
              </div>
            </div>
          </div>
        )}

        {tab === "appearance" && (
          <div className="space-y-4">
            <div className="space-y-2">
              <Label className="text-text-secondary">Theme</Label>
              <Select
                value={settings.appearance}
                onValueChange={(v) => {
                  if (!v) return;
                  setSettings((s) => s && { ...s, appearance: v });
                  useThemeStore.getState().setMode(v as "dark" | "light");
                }}
              >
                <SelectTrigger className="bg-background">
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="dark">dark</SelectItem>
                  <SelectItem value="light">light</SelectItem>
                </SelectContent>
              </Select>
            </div>

            <div className="space-y-2">
              <Label className="text-text-secondary">Default License</Label>
              <Select
                value={settings.default_license}
                onValueChange={(v) => {
                  if (!v) return;
                  setSettings((s) => s && { ...s, default_license: v });
                }}
              >
                <SelectTrigger className="bg-background">
                  <SelectValue placeholder="Select a license" />
                </SelectTrigger>
                <SelectContent>
                  {[
                    "MIT",
                    "Apache-2.0",
                    "GPL-3.0",
                    "AGPL-3.0",
                    "BSD-3-Clause",
                    "BSD-2-Clause",
                    "MPL-2.0",
                    "Unlicense",
                    "Proprietary",
                  ].map((lic) => (
                    <SelectItem key={lic} value={lic}>
                      {lic}
                    </SelectItem>
                  ))}
                  {![
                    "MIT",
                    "Apache-2.0",
                    "GPL-3.0",
                    "AGPL-3.0",
                    "BSD-3-Clause",
                    "BSD-2-Clause",
                    "MPL-2.0",
                    "Unlicense",
                    "Proprietary",
                  ].includes(settings.default_license) && (
                    <SelectItem value={settings.default_license}>
                      {settings.default_license}
                    </SelectItem>
                  )}
                </SelectContent>
              </Select>
            </div>

            <div className="space-y-2">
              <Label className="text-text-secondary">Default Project Directory</Label>
              <Input
                value={settings.default_project_dir}
                onChange={(e) =>
                  setSettings((s) => s && { ...s, default_project_dir: e.target.value })
                }
                className="bg-surface font-mono text-sm"
              />
            </div>

            <div className="space-y-2">
              <Label className="text-text-secondary">Font</Label>
              <Select
                value={settings.font}
                  onValueChange={(v) => {
                    if (!v) return;
                    setSettings((s) => s && { ...s, font: v });
                    const families: Record<string, string> = {
                      "JetBrains Mono": "\"JetBrains Mono\", monospace",
                      "Fira Code": "\"Fira Code\", monospace",
                      "Cascadia Code": "\"Cascadia Code\", monospace",
                      "IBM Plex Mono": "\"IBM Plex Mono\", monospace",
                      "Source Code Pro": "\"Source Code Pro\", monospace",
                      "Inconsolata": "\"Inconsolata\", monospace",
                      "Ubuntu Mono": "\"Ubuntu Mono\", monospace",
                      "DejaVu Sans Mono": "\"DejaVu Sans Mono\", monospace",
                      "Roboto Mono": "\"Roboto Mono\", monospace",
                      "Monoid": "Monoid, monospace",
                    };
                    document.documentElement.style.setProperty(
                      "--font-mono",
                      families[v] || `\"${v}\", monospace`
                    );
                  }}
              >
                <SelectTrigger className="bg-background">
                  <SelectValue placeholder="Select a font" />
                </SelectTrigger>
                <SelectContent>
                  {[
                    "JetBrains Mono",
                    "Fira Code",
                    "Cascadia Code",
                    "IBM Plex Mono",
                    "Source Code Pro",
                    "Inconsolata",
                    "Ubuntu Mono",
                    "DejaVu Sans Mono",
                    "Roboto Mono",
                    "Monoid",
                  ].map((f) => (
                    <SelectItem key={f} value={f}>
                      {f}
                    </SelectItem>
                  ))}
                  {![
                    "JetBrains Mono",
                    "Fira Code",
                    "Cascadia Code",
                    "IBM Plex Mono",
                    "Source Code Pro",
                    "Inconsolata",
                    "Ubuntu Mono",
                    "DejaVu Sans Mono",
                    "Roboto Mono",
                    "Monoid",
                  ].includes(settings.font) && (
                    <SelectItem value={settings.font}>
                      {settings.font}
                    </SelectItem>
                  )}
                </SelectContent>
              </Select>
            </div>
          </div>
        )}

        <div className="flex justify-end gap-2 pt-2">
          <Button variant="outline" onClick={onClose}>
            Cancel
          </Button>
          <Button onClick={handleSave} disabled={saving}>
            {saving ? "Saving..." : "Save"}
          </Button>
        </div>
      </DialogContent>
    </Dialog>
  );
}
