package scanner

import (
	"os"
	"path/filepath"
	"testing"
)

func createTempDir(t *testing.T, files map[string]string) string {
	t.Helper()
	dir := t.TempDir()
	for name, content := range files {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0644); err != nil {
			t.Fatalf("write %s: %v", name, err)
		}
	}
	return dir
}

func TestAnalyzePath_GoModule(t *testing.T) {
	dir := createTempDir(t, map[string]string{
		"go.mod": "module github.com/test/app\n\ngo 1.21\n\nrequire (\n\tgithub.com/gorilla/mux v1.8.0\n)\n",
		"main.go": "package main\n\nfunc main() {}",
	})
	result := AnalyzePath(dir)
	if result == "" {
		t.Error("expected non-empty result for Go module")
	}
	if !contains(result, "Go") {
		t.Error("expected Go detection in", result)
	}
}

func TestAnalyzePath_NodeJS(t *testing.T) {
	dir := createTempDir(t, map[string]string{
		"package.json": `{"name":"test","dependencies":{"react":"^19.0.0","zustand":"^5.0.0"}}`,
	})
	result := AnalyzePath(dir)
	if result == "" {
		t.Error("expected non-empty result for Node.js project")
	}
	if !contains(result, "Node") {
		t.Error("expected Node.js detection in", result)
	}
	if !contains(result, "react") {
		t.Error("expected react dependency in", result)
	}
}

func TestAnalyzePath_Wails(t *testing.T) {
	dir := createTempDir(t, map[string]string{
		"wails.json": `{"name":"test","outputfilename":"test"}`,
	})
	result := AnalyzePath(dir)
	if result == "" {
		t.Error("expected non-empty result for Wails project")
	}
	if !contains(result, "Wails") {
		t.Error("expected Wails detection in", result)
	}
}

func TestAnalyzePath_EmptyDir(t *testing.T) {
	dir := createTempDir(t, map[string]string{})
	result := AnalyzePath(dir)
	if !contains(result, "No project files found") {
		t.Error("expected 'No project files found' message")
	}
}

func TestAnalyzePath_FileContent(t *testing.T) {
	dir := createTempDir(t, map[string]string{
		"README.md": "# Test Project\n\nA test project for validation.",
		"package.json": `{"name":"test"}`,
	})
	result := AnalyzePathContent(dir)
	if !contains(result, "# Test Project") {
		t.Error("expected README content in", result)
	}
}

func TestAnalyzePath_FileDirectRead(t *testing.T) {
	dir := createTempDir(t, map[string]string{
		"hello.txt": "Hello, World!",
	})
	result := AnalyzePath(filepath.Join(dir, "hello.txt"))
	if !contains(result, "Hello") {
		t.Error("expected file content in", result)
	}
}

func contains(s, substr string) bool {
	return len(s) > 0 && len(substr) > 0 && len(s) >= len(substr) && (s == substr || containsInner(s, substr))
}

func containsInner(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
