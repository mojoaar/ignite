import { vi } from "vitest";

vi.mock("@wails/runtime", () => ({
  EventsOn: vi.fn(() => () => {}),
  EventsOff: vi.fn(),
}));

vi.mock("@wails/go/main/App", () => ({
  GetVersion: vi.fn(() => Promise.resolve("0.1.6")),
  GetSettings: vi.fn(() => Promise.resolve({ appearance: "dark", default_project_dir: "~/Development", font: "JetBrains Mono", name: "", avatar: "", window_width: 1024, window_height: 768, default_provider: "opencode-go", default_license: "AGPL-3.0", providers: {} })),
  HasAPIKey: vi.fn(() => Promise.resolve(false)),
  GetCachedModels: vi.fn(() => Promise.resolve([])),
  SaveSettings: vi.fn(() => Promise.resolve()),
  SetAPIKey: vi.fn(() => Promise.resolve()),
  ValidateProviderKey: vi.fn(() => Promise.resolve()),
  ListProjects: vi.fn(() => Promise.resolve([])),
  CreateProject: vi.fn(() => Promise.resolve()),
  GetProject: vi.fn(() => Promise.resolve({ id: "", name: "", tagline: "", path: "", provider: "", model: "", created_at: "", updated_at: "" })),
  GetMessages: vi.fn(() => Promise.resolve([])),
  AddMessage: vi.fn(() => Promise.resolve()),
  UpdateProject: vi.fn(() => Promise.resolve()),
  DeleteProject: vi.fn(() => Promise.resolve()),
  SelectDirectory: vi.fn(() => Promise.resolve("~/Development")),
  SendMessageStream: vi.fn(() => Promise.resolve()),
  ExportChat: vi.fn(() => Promise.resolve("")),
  Greet: vi.fn(() => Promise.resolve("")),
  GetCachedModels: vi.fn(() => Promise.resolve([])),
  GetVersion: vi.fn(() => Promise.resolve("0.1.6")),
  AnalyzePath: vi.fn(() => Promise.resolve("")),
  AnalyzePathContent: vi.fn(() => Promise.resolve("")),
  FetchURL: vi.fn(() => Promise.resolve("")),
  SetProjectMeta: vi.fn(() => Promise.resolve()),
  ResizeWindow: vi.fn(() => Promise.resolve()),
  GenerateProjectFiles: vi.fn(() => Promise.resolve()),
}));
