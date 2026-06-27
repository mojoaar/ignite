package scanner

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func AnalyzePath(path string) string {
	expanded := expandPath(path)

	info, err := os.Stat(expanded)
	if err != nil {
		return fmt.Sprintf("Path not found: %s", path)
	}

	if !info.IsDir() {
		data, err := os.ReadFile(expanded)
		if err != nil {
			return fmt.Sprintf("Could not read file: %s", path)
		}
		return fmt.Sprintf("File: %s\n```\n%s\n```", filepath.Base(path), truncate(string(data), 8000))
	}

	return scanDirectory(expanded)
}

func expandPath(path string) string {
	if strings.HasPrefix(path, "~/") {
		home, err := os.UserHomeDir()
		if err == nil {
			return filepath.Join(home, path[2:])
		}
	}
	return path
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max] + "\n... (truncated)"
}

type detector struct {
	file     string
	label    string
	expander func(string) string
}

var detectors = []detector{
	{file: "package.json", label: "Node.js package", expander: expandPackageJSON},
	{file: "go.mod", label: "Go module", expander: expandGoMod},
	{file: "Cargo.toml", label: "Rust crate", expander: readFirstLines},
	{file: "pyproject.toml", label: "Python project", expander: readFirstLines},
	{file: "requirements.txt", label: "Python requirements", expander: readFirstLines},
	{file: "Gemfile", label: "Ruby Gemfile", expander: readFirstLines},
	{file: "composer.json", label: "PHP Composer", expander: readFirstLines},
	{file: "build.gradle", label: "Gradle build", expander: readFirstLines},
	{file: "build.gradle.kts", label: "Gradle build (KTS)", expander: readFirstLines},
	{file: "pom.xml", label: "Maven POM", expander: readFirstLines},
	{file: "wails.json", label: "Wails config", expander: readFirstLines},
	{file: "Dockerfile", label: "Dockerfile", expander: readFirstLines},
	{file: "docker-compose.yml", label: "Docker Compose", expander: readFirstLines},
	{file: "docker-compose.yaml", label: "Docker Compose", expander: readFirstLines},
	{file: "README.md", label: "README", expander: readFirstLines},
}

var globDetectors = []struct {
	pattern string
	label   string
}{
	{"next.config.*", "Next.js config"},
	{"tsconfig.json", "TypeScript config"},
	{"vite.config.*", "Vite config"},
	{"tailwind.config.*", "Tailwind config (v3)"},
	{"postcss.config.*", "PostCSS config"},
	{"tailwind.css", "Tailwind CSS"},
	{"svelte.config.*", "SvelteKit config"},
	{"astro.config.*", "Astro config"},
	{"eslint.config.*", "ESLint config"},
	{".eslintrc*", "ESLint config"},
}

func scanDirectory(dir string) string {
	var out strings.Builder
	out.WriteString(fmt.Sprintf("Directory: %s\n\n## Detected\n\n", dir))

	count := 0

	for _, d := range detectors {
		target := filepath.Join(dir, d.file)
		if _, err := os.Stat(target); err == nil {
			content := d.expander(target)
			out.WriteString(fmt.Sprintf("### %s (%s)\n%s\n\n", d.label, d.file, content))
			count++
		}
	}

	for _, gd := range globDetectors {
		matches, _ := filepath.Glob(filepath.Join(dir, gd.pattern))
		for _, m := range matches {
			out.WriteString(fmt.Sprintf("### %s (%s)\n%s\n\n", gd.label, filepath.Base(m), readFirstLines(m)))
			count++
		}
	}

	if count == 0 {
		out.WriteString("No project files found.\n")
	}

	return out.String()
}

func readFirstLines(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return "(could not read)"
	}
	return "```\n" + truncate(string(data), 2000) + "\n```"
}

func expandPackageJSON(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return "(could not read)"
	}
	var sb strings.Builder
	sb.WriteString("```json\n")
	content := string(data)
	sb.WriteString(truncate(content, 2000))
	sb.WriteString("\n```")

	deps := extractJSONSection(content, "dependencies")
	devDeps := extractJSONSection(content, "devDependencies")
	scripts := extractJSONSection(content, "scripts")

	if deps != "" {
		sb.WriteString("\n\n**Dependencies:**\n```json\n")
		sb.WriteString(deps)
		sb.WriteString("\n```")
	}
	if devDeps != "" {
		sb.WriteString("\n\n**Dev Dependencies:**\n```json\n")
		sb.WriteString(devDeps)
		sb.WriteString("\n```")
	}
	if scripts != "" {
		sb.WriteString("\n\n**Scripts:**\n```json\n")
		sb.WriteString(scripts)
		sb.WriteString("\n```")
	}

	return sb.String()
}

func expandGoMod(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return "(could not read)"
	}
	content := string(data)
	var sb strings.Builder
	sb.WriteString("```\n")
	for _, line := range strings.Split(content, "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "module ") || strings.HasPrefix(trimmed, "go ") || strings.HasPrefix(trimmed, "require ") {
			sb.WriteString(line + "\n")
		}
	}
	sb.WriteString("\n```")

	reqBlock := extractRequireBlock(content)
	if reqBlock != "" {
		sb.WriteString("\n\n**Require:**\n```\n")
		sb.WriteString(reqBlock)
		sb.WriteString("\n```")
	}

	return sb.String()
}

func extractJSONSection(content, key string) string {
	searchKey := `"` + key + `"`
	idx := strings.Index(content, searchKey)
	if idx < 0 {
		return ""
	}
	rest := content[idx+len(searchKey):]
	colonIdx := strings.Index(rest, ":")
	if colonIdx < 0 {
		return ""
	}
	rest = rest[colonIdx+1:]
	depth := 0
	start := -1
	end := -1
	for i, c := range rest {
		if c == '{' || c == '[' {
			if depth == 0 {
				start = i
			}
			depth++
		} else if c == '}' || c == ']' {
			depth--
			if depth == 0 && start >= 0 {
				end = i + 1
				break
			}
		}
	}
	if start >= 0 && end > start {
		return strings.TrimSpace(rest[start:end])
	}
	return ""
}

func extractRequireBlock(content string) string {
	lines := strings.Split(content, "\n")
	inRequire := false
	var rb strings.Builder
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "require (") || strings.HasPrefix(trimmed, "require(") {
			inRequire = true
			continue
		}
		if inRequire {
			if trimmed == ")" {
				break
			}
			rb.WriteString(line + "\n")
		}
	}
	return rb.String()
}

func AnalyzePathContent(path string) string {
	expanded := expandPath(path)
	info, err := os.Stat(expanded)
	if err != nil {
		return ""
	}

	if !info.IsDir() {
		data, err := os.ReadFile(expanded)
		if err != nil {
			return ""
		}
		return string(data)
	}

	keyFiles := []string{
		"README.md", "readme.md", "README", "readme",
		"package.json", "go.mod", "Cargo.toml", "pyproject.toml",
		"wails.json", "agents.md", "plan.md", "AGENTS.md", "PLAN.md",
		"docker-compose.yml", "Dockerfile", "Makefile",
	}

	var sb strings.Builder
	for _, f := range keyFiles {
		fp := filepath.Join(expanded, f)
		data, err := os.ReadFile(fp)
		if err != nil {
			continue
		}
		maxLen := 3000
		content := string(data)
		if len(content) > maxLen {
			for maxLen > 0 && content[maxLen]&0xC0 == 0x80 {
				maxLen--
			}
			content = content[:maxLen] + "\n... (truncated)"
		}
		sb.WriteString(fmt.Sprintf("=== %s ===\n%s\n\n", f, content))
	}
	if sb.Len() == 0 {
		return ""
	}
	return sb.String()
}
