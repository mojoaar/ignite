package templates

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/Masterminds/sprig/v3"
)

type Engine struct {
	projectTmpl, agentsTmpl, planTmpl, readmeTmpl *template.Template
}

func NewEngine(projectTemplate, agentsTemplate, planTemplate, readmeTemplate string) (*Engine, error) {
	funcMap := sprig.TxtFuncMap()

	projectTmpl, err := template.New("project.md").Funcs(funcMap).Parse(projectTemplate)
	if err != nil {
		return nil, fmt.Errorf("template engine: parse project.md: %w", err)
	}

	agentsTmpl, err := template.New("agents.md").Funcs(funcMap).Parse(agentsTemplate)
	if err != nil {
		return nil, fmt.Errorf("template engine: parse agents.md: %w", err)
	}

	planTmpl, err := template.New("plan.md").Funcs(funcMap).Parse(planTemplate)
	if err != nil {
		return nil, fmt.Errorf("template engine: parse plan.md: %w", err)
	}

	readmeTmpl, err := template.New("README.md").Funcs(funcMap).Parse(readmeTemplate)
	if err != nil {
		return nil, fmt.Errorf("template engine: parse README.md: %w", err)
	}

	return &Engine{projectTmpl, agentsTmpl, planTmpl, readmeTmpl}, nil
}

func (e *Engine) Generate(ctx ProjectContext) (*ProjectFiles, error) {
	var pb, ab, plb, rb bytes.Buffer

	if err := e.projectTmpl.Execute(&pb, ctx); err != nil {
		return nil, fmt.Errorf("template engine: project.md: %w", err)
	}
	if err := e.agentsTmpl.Execute(&ab, ctx); err != nil {
		return nil, fmt.Errorf("template engine: agents.md: %w", err)
	}
	if err := e.planTmpl.Execute(&plb, ctx); err != nil {
		return nil, fmt.Errorf("template engine: plan.md: %w", err)
	}
	if err := e.readmeTmpl.Execute(&rb, ctx); err != nil {
		return nil, fmt.Errorf("template engine: README.md: %w", err)
	}

	return &ProjectFiles{pb.String(), ab.String(), plb.String(), rb.String()}, nil
}
