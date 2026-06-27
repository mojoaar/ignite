import { useEffect, useState } from "react";
import { Plus, Flame, Trash2, Pencil } from "lucide-react";
import { Button } from "@/components/ui/button";
import { useChatStore } from "@/lib/store/chat";
import { cn } from "@/lib/utils";
import { ListProjects, CreateProject, GetProject, GetMessages, DeleteProject, UpdateProject } from "@wails/go/main/App";
import { EventsOn, EventsOff } from "@wails/runtime";

export function Sidebar() {
  const projects = useChatStore((s) => s.projects);
  const activeProjectId = useChatStore((s) => s.activeProjectId);
  const setProjects = useChatStore((s) => s.setProjects);
  const setActiveProject = useChatStore((s) => s.setActiveProject);
  const setActiveProjectPath = useChatStore((s) => s.setActiveProjectPath);
  const setProjectName = useChatStore((s) => s.setProjectName);
  const removeProject = useChatStore((s) => s.removeProject);
  const setMessages = useChatStore((s) => s.setMessages);

  const [editingId, setEditingId] = useState<string | null>(null);
  const [editValue, setEditValue] = useState("");

  const startRename = (id: string, currentName: string) => {
    setEditingId(id);
    setEditValue(currentName);
  };

  const commitRename = async () => {
    if (!editingId) return;
    const val = editValue.trim();
    try {
      await UpdateProject({
        id: editingId,
        name: val,
        tagline: "",
        path: "",
        provider: "",
        model: "",
        created_at: "",
        updated_at: "",
      });
      setProjectName(editingId, val, "");
      useChatStore.getState().projects.map((p) =>
        p.id === editingId ? { ...p, name: val } : p
      );
    } catch {}
    setEditingId(null);
  };

  useEffect(() => {
    ListProjects()
      .then((p) => setProjects(p ?? []))
      .catch(() => {});
  }, [setProjects]);

  useEffect(() => {
    EventsOn("menu-new-project", () => { handleNewProject(); });
    return () => { EventsOff("menu-new-project"); };
  }, []);

  const handleNewProject = async () => {
    const id = crypto.randomUUID();
    const now = new Date().toISOString();
    const project = {
      id,
      name: "",
      tagline: "",
      path: "",
      provider: "",
      model: "",
      created_at: now,
      updated_at: now,
    };
    try {
      await CreateProject(project);
      const updated = await ListProjects();
      setProjects(updated ?? []);
      handleSelectProject(id);
    } catch {}
  };

  const handleSelectProject = async (id: string) => {
    setActiveProject(id);
    try {
      const msgs = await GetMessages(id);
      setMessages(msgs ?? []);
      const proj = await GetProject(id);
      if (proj) {
        setActiveProjectPath(proj.path || "");
        if (proj.name) {
          setProjects(
            useChatStore.getState().projects.map((p) =>
              p.id === id ? { ...p, name: proj.name, tagline: proj.tagline } : p
            )
          );
        }
      }
    } catch {}
  };

  const handleDeleteProject = async (id: string, e: React.MouseEvent) => {
    e.stopPropagation();
    try {
      await DeleteProject(id);
      removeProject(id);
      const updated = await ListProjects();
      setProjects(updated ?? []);
    } catch {}
  };

  return (
    <aside className="flex h-full w-[260px] shrink-0 flex-col border-r border-border bg-surface">
      <div className="flex items-center gap-2 px-4 py-4">
        <Flame className="h-6 w-6 text-accent" />
        <span className="font-mono text-lg font-semibold text-text-primary">
          Ignite
        </span>
      </div>

      <div className="px-3 py-1">
        <p className="font-mono text-[10px] font-medium uppercase tracking-wider text-text-secondary">
          Projects
        </p>
      </div>

      <div className="flex-1 overflow-y-auto px-2">
        {projects.length === 0 && (
          <p className="px-2 py-8 text-center text-sm text-text-secondary">
            No projects yet
          </p>
        )}
        {projects.map((project) => (
          <div
            key={project.id}
            onClick={() => handleSelectProject(project.id)}
            className={cn(
              "group flex w-full cursor-pointer items-center gap-2 rounded-md px-3 py-2 text-left text-sm transition-colors",
              project.id === activeProjectId
                ? "bg-accent/15 text-text-primary"
                : "text-text-secondary hover:bg-surface-hover hover:text-text-primary"
            )}
          >
            <span className="flex h-2 w-2 shrink-0 rounded-full bg-accent" />
            {editingId === project.id ? (
              <input
                autoFocus
                value={editValue}
                onChange={(e) => setEditValue(e.target.value)}
                onBlur={commitRename}
                onKeyDown={(e) => {
                  if (e.key === "Enter") commitRename();
                  if (e.key === "Escape") setEditingId(null);
                }}
                onClick={(e) => e.stopPropagation()}
                className="w-full truncate rounded border border-accent bg-background px-1 py-0 font-mono text-xs text-text-primary outline-none"
              />
            ) : (
              <span
                className="truncate font-mono text-xs"
              >
                {project.name || "Untitled"}
              </span>
            )}
            <button
              onClick={(e) => { e.stopPropagation(); startRename(project.id, project.name || "Untitled"); }}
              className="shrink-0 rounded p-0.5 opacity-0 transition-opacity group-hover:opacity-100 hover:bg-surface-hover hover:text-accent"
              title="Rename project"
            >
              <Pencil className="h-3 w-3" />
            </button>
            <button
              onClick={(e) => handleDeleteProject(project.id, e)}
              className="ml-auto shrink-0 rounded p-0.5 opacity-0 transition-opacity group-hover:opacity-100 hover:bg-surface-hover hover:text-error"
              title="Delete project"
            >
              <Trash2 className="h-3 w-3" />
            </button>
          </div>
        ))}
      </div>

      <div className="border-t border-border p-3">
        <Button
          onClick={handleNewProject}
          className="w-full"
          variant="outline"
          size="sm"
        >
          <Plus className="mr-2 h-4 w-4" />
          New Project
        </Button>
      </div>
    </aside>
  );
}
