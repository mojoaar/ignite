package templates

type ProjectContext struct {
	Name, Tagline, Description, License, GoModulePath        string
	Frontend, Backend, Database, Auth, Hosting, PackageManager string
	Features       []string
	Phases         []Phase
	TechStack      []TechItem
	Dependencies   []Dependency
	APIs           []APIEndpoint
	DBTables       []DBTable
	Performance    []PerfTarget
	Risks          []Risk
	BannedPackages []string
	EnvVars        []EnvVar
	Verification   []string
	Theme          ThemeConfig
	DevWorkflow    DevWorkflow
	Changelog      []ChangelogEntry
}

type Phase struct {
	Name, Description string
	Tasks             []string
}

type TechItem struct {
	Category, Choice, Version string
}

type Dependency struct {
	Name, Version, Why string
}

type APIEndpoint struct {
	Method, Path, Desc, Auth string
}

type DBTable struct {
	Name    string
	Columns []DBColumn
}

type DBColumn struct {
	Name, Type, Desc string
}

type PerfTarget struct {
	Metric, Target string
}

type Risk struct {
	Risk, Mitigation string
}

type EnvVar struct {
	Name, Desc, Default string
}

type ThemeConfig struct {
	Dark, Light ThemePalette
	Font        string
}

type ThemePalette struct {
	Background, Surface, SurfaceHover, Border, TextPrimary, TextSecondary string
	Accent, AccentSecondary, Success, Error                               string
}

type DevWorkflow struct {
	Setup, Dev, Build, Test, Lint, TypeCheck []string
}

type ChangelogEntry struct {
	Version, Date string
	Changes       []string
}

type ProjectFiles struct {
	ProjectMD string
	AgentsMD  string
	PlanMD    string
	ReadmeMD  string
}
