package templates

import (
	"strings"
	"testing"
)

func TestEngine_Generate(t *testing.T) {
	projectTmpl := `# {{.Name}}
**Tagline:** {{.Tagline}}
## Tech Stack
{{range .TechStack}}| {{.Category}} | {{.Choice}} |
{{end}}`
	agentsTmpl := `# Agents for {{.Name}}
- Build: {{index .DevWorkflow.Build 0}}`
	planTmpl := `# Plan for {{.Name}}
{{range .Phases}}### {{.Name}}
{{range .Tasks}}- [ ] {{.}}
{{end}}{{end}}`
	readmeTmpl := `# {{.Name}}
> {{.Tagline}}
## Features
{{range .Features}}- {{.}}
{{end}}`

	engine, err := NewEngine(projectTmpl, agentsTmpl, planTmpl, readmeTmpl)
	if err != nil {
		t.Fatalf("NewEngine: %v", err)
	}

	ctx := ProjectContext{
		Name:    "test",
		Tagline: "Just a test",
		TechStack: []TechItem{
			{Category: "Frontend", Choice: "React"},
			{Category: "Backend", Choice: "Go"},
		},
		Features: []string{"Fast", "Secure"},
		Phases:   []Phase{{Name: "Phase 0", Tasks: []string{"Setup"}}},
		DevWorkflow: DevWorkflow{
			Build: []string{"go build ./..."},
			Test:  []string{"go test ./..."},
		},
	}

	files, err := engine.Generate(ctx)
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}

	if !strings.Contains(files.ProjectMD, "# test") {
		t.Error("project.md missing name")
	}
	if !strings.Contains(files.AgentsMD, "go build ./...") {
		t.Error("agents.md missing build cmd")
	}
	if !strings.Contains(files.PlanMD, "Phase 0") {
		t.Error("plan.md missing phase")
	}
	if !strings.Contains(files.ReadmeMD, "Fast") {
		t.Error("README.md missing feature")
	}
}

func TestEmbeddedEngine(t *testing.T) {
	engine, err := EmbeddedEngine()
	if err != nil {
		t.Fatalf("EmbeddedEngine: %v", err)
	}
	ctx := ProjectContext{
		Name:    "embedded-test",
		Tagline: "Testing embed",
		License: "MIT",
		Features: []string{"Fast"},
		Phases: []Phase{
			{Name: "Phase 0", Tasks: []string{"Setup"}},
		},
		TechStack: []TechItem{
			{Category: "Frontend", Choice: "React"},
		},
	}
	files, err := engine.Generate(ctx)
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}
	if !strings.Contains(files.ProjectMD, "# embedded-test") {
		t.Error("project.md missing project name")
	}
	if !strings.Contains(files.AgentsMD, "React") {
		t.Error("agents.md missing tech stack")
	}
	if !strings.Contains(files.PlanMD, "Phase 0") {
		t.Error("plan.md missing phase")
	}
	if !strings.Contains(files.ReadmeMD, "Fast") {
		t.Error("README.md missing feature")
	}
}
