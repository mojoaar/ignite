package history

import (
	"os"
	"path/filepath"
	"testing"
)

func TestStore_ProjectCRUD(t *testing.T) {
	dir := filepath.Join(os.TempDir(), "ignite-test-history")
	os.MkdirAll(dir, 0700)
	defer os.RemoveAll(dir)

	store, err := OpenDB(filepath.Join(dir, "test.db"))
	if err != nil {
		t.Fatalf("OpenDB: %v", err)
	}
	defer store.Close()

	p := Project{ID: "test-1", Name: "Test", Tagline: "A test", Provider: "claude", Model: "opus"}
	if err := store.CreateProject(p); err != nil {
		t.Fatalf("CreateProject: %v", err)
	}

	projects, err := store.ListProjects()
	if err != nil {
		t.Fatalf("ListProjects: %v", err)
	}
	if len(projects) != 1 {
		t.Fatalf("expected 1, got %d", len(projects))
	}
	if projects[0].Name != "Test" {
		t.Errorf("name mismatch")
	}
}

func TestStore_MessageCRUD(t *testing.T) {
	dir := filepath.Join(os.TempDir(), "ignite-test-msg")
	os.MkdirAll(dir, 0700)
	defer os.RemoveAll(dir)

	store, err := OpenDB(filepath.Join(dir, "test.db"))
	if err != nil {
		t.Fatalf("OpenDB: %v", err)
	}
	defer store.Close()

	store.CreateProject(Project{ID: "proj-1", Name: "Msg Test"})
	store.AddMessage(Message{ID: "m1", ProjectID: "proj-1", Phase: "identity", Role: "user", Content: "hello"})
	store.AddMessage(Message{ID: "m2", ProjectID: "proj-1", Phase: "identity", Role: "assistant", Content: "hi"})

	msgs, err := store.GetMessages("proj-1")
	if err != nil {
		t.Fatalf("GetMessages: %v", err)
	}
	if len(msgs) != 2 {
		t.Fatalf("expected 2, got %d", len(msgs))
	}
}

func TestStore_UpdateProject(t *testing.T) {
	dir := filepath.Join(os.TempDir(), "ignite-test-update")
	os.MkdirAll(dir, 0700)
	defer os.RemoveAll(dir)

	store, err := OpenDB(filepath.Join(dir, "test.db"))
	if err != nil {
		t.Fatalf("OpenDB: %v", err)
	}
	defer store.Close()

	store.CreateProject(Project{ID: "p1", Name: "Original"})
	if err := store.UpdateProject(Project{ID: "p1", Name: "Updated", Tagline: "New tag", Path: "/tmp", Provider: "openai", Model: "gpt-4"}); err != nil {
		t.Fatalf("UpdateProject: %v", err)
	}

	proj, err := store.GetProject("p1")
	if err != nil {
		t.Fatalf("GetProject: %v", err)
	}
	if proj.Name != "Updated" {
		t.Errorf("expected Updated, got %s", proj.Name)
	}
	if proj.Tagline != "New tag" {
		t.Errorf("expected New tag, got %s", proj.Tagline)
	}
}

func TestStore_DeleteProject(t *testing.T) {
	dir := filepath.Join(os.TempDir(), "ignite-test-delete")
	os.MkdirAll(dir, 0700)
	defer os.RemoveAll(dir)

	store, err := OpenDB(filepath.Join(dir, "test.db"))
	if err != nil {
		t.Fatalf("OpenDB: %v", err)
	}
	defer store.Close()

	store.CreateProject(Project{ID: "p1", Name: "ToDelete"})
	store.AddMessage(Message{ID: "m1", ProjectID: "p1", Role: "user", Content: "msg"})

	if err := store.DeleteProject("p1"); err != nil {
		t.Fatalf("DeleteProject: %v", err)
	}

	projects, _ := store.ListProjects()
	if len(projects) != 0 {
		t.Error("expected 0 projects after delete")
	}

	msgs, _ := store.GetMessages("p1")
	if len(msgs) != 0 {
		t.Error("expected 0 messages after delete (cascade)")
	}
}

func TestStore_UpsertAndListModels(t *testing.T) {
	dir := filepath.Join(os.TempDir(), "ignite-test-models")
	os.MkdirAll(dir, 0700)
	defer os.RemoveAll(dir)

	store, err := OpenDB(filepath.Join(dir, "test.db"))
	if err != nil {
		t.Fatalf("OpenDB: %v", err)
	}
	defer store.Close()

	store.UpsertProviderModel("openai", "gpt-4", "GPT-4")
	store.UpsertProviderModel("openai", "gpt-4o", "GPT-4o")
	store.UpsertProviderModel("openai", "gpt-4", "GPT-4 Updated")

	models, err := store.ListCachedModels("openai")
	if err != nil {
		t.Fatalf("ListCachedModels: %v", err)
	}
	if len(models) != 2 {
		t.Fatalf("expected 2 models, got %d", len(models))
	}
}
