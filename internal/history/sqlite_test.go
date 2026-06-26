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
