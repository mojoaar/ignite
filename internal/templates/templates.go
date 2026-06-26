package templates

import (
	"embed"
	"fmt"
)

//go:embed templates/*
var templateFS embed.FS

func EmbeddedEngine() (*Engine, error) {
	projectTmpl, err := templateFS.ReadFile("templates/project.md.tmpl")
	if err != nil {
		return nil, fmt.Errorf("template: read project.md.tmpl: %w", err)
	}
	agentsTmpl, err := templateFS.ReadFile("templates/AGENTS.md.tmpl")
	if err != nil {
		return nil, fmt.Errorf("template: read AGENTS.md.tmpl: %w", err)
	}
	planTmpl, err := templateFS.ReadFile("templates/PLAN.md.tmpl")
	if err != nil {
		return nil, fmt.Errorf("template: read PLAN.md.tmpl: %w", err)
	}
	readmeTmpl, err := templateFS.ReadFile("templates/README.md.tmpl")
	if err != nil {
		return nil, fmt.Errorf("template: read README.md.tmpl: %w", err)
	}
	return NewEngine(string(projectTmpl), string(agentsTmpl), string(planTmpl), string(readmeTmpl))
}
