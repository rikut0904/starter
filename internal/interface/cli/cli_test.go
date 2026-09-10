package cli

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/rikut0904/starter/internal/domain"
)

func TestGenerateEmptyCreatesCommonFilesAndProtectsExisting(t *testing.T) {
	dir := t.TempDir()
	if err := generate(filepath.Join(dir, "demo"), Config{Name: "demo", Profile: "empty"}, false, func(string, string, ...string) error { return nil }); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"AGENTS.md", "README.md", "Makefile", ".gitignore", ".github/workflows/ci.yml"} {
		if _, err := os.Stat(filepath.Join(dir, "demo", name)); err != nil {
			t.Errorf("missing %s: %v", name, err)
		}
	}
	if err := generate(filepath.Join(dir, "demo"), Config{Name: "demo", Profile: "empty"}, false, func(string, string, ...string) error { return nil }); err == nil {
		t.Fatal("expected collision error")
	}
}

func TestLicenseYearIsReplacedAndHolderIsFixed(t *testing.T) {
	files := commonFiles(Config{Name: "sample-project", Profile: "empty"})
	license := files["LICENSE"]
	if strings.Contains(license, "{YEAR}") {
		t.Fatal("license placeholders were not replaced")
	}
	if !strings.Contains(license, "Copyright (c) "+fmt.Sprint(time.Now().Year())+" rikut0904") {
		t.Fatalf("unexpected license copyright line: %s", license)
	}
}

func TestGenerateRunsRequiredCommandsOnce(t *testing.T) {
	dir := t.TempDir()
	calls := []string{}
	runner := func(cwd, name string, args ...string) error {
		calls = append(calls, name+" "+args[0])
		return errors.New("simulated")
	}
	if err := generate(filepath.Join(dir, "demo"), Config{Name: "demo", Profile: "next-go"}, true, runner); err != nil {
		t.Fatal(err)
	}
	if len(calls) != 2 {
		t.Fatalf("got %d external calls, want 2", len(calls))
	}
	if _, err := os.Stat(filepath.Join(dir, "demo", "README.md")); err != nil {
		t.Fatalf("common files must be generated after external failures: %v", err)
	}
}

func TestLoadConfigJSON(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.json")
	if err := os.WriteFile(path, []byte(`{"name":"sample","profile":"go"}`), 0600); err != nil {
		t.Fatal(err)
	}
	var cfg Config
	if err := loadConfig(path, &cfg); err != nil {
		t.Fatal(err)
	}
	if cfg.Name != "sample" || cfg.Profile != "go" {
		t.Fatalf("unexpected config: %+v", cfg)
	}
}

func TestGoProfileIncludesCleanArchitecture(t *testing.T) {
	files := commonFiles(Config{Name: "example.com/demo", Profile: "go"})
	for _, name := range []string{"cmd/server/main.go", "internal/domain/health.go", "internal/usecase/health.go"} {
		if _, ok := files[name]; !ok {
			t.Fatalf("missing generated file %s", name)
		}
	}
}

func TestNextGoBackendUsesBackendModule(t *testing.T) {
	files := commonFiles(Config{Name: "sample", Profile: "next-go"})
	if !strings.Contains(files["backend/cmd/server/main.go"], `"backend/internal/usecase"`) {
		t.Fatal("next-go main.go must import the backend module")
	}
	if strings.Contains(files["backend/cmd/server/main.go"], "starter/create/next-go") {
		t.Fatal("template import path leaked into generated main.go")
	}
}

func TestListItemDisplaysTitleAndDescription(t *testing.T) {
	choice := item{title: "Goのみ", desc: "go"}
	if choice.Title() != "Goのみ" || choice.Description() != "go" {
		t.Fatalf("unexpected list item: title=%q description=%q", choice.Title(), choice.Description())
	}
}

func TestListModelEnterQuits(t *testing.T) {
	model := listModel{list: list.New([]list.Item{item{title: "Goのみ", desc: "go"}}, list.NewDefaultDelegate(), 40, 10)}
	_, command := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if command == nil {
		t.Fatal("Enter must return a quit command")
	}
}

func TestNextCommandIsNonInteractive(t *testing.T) {
	args := nextCommand("frontend", domain.NextConfig{TypeScript: true, ESLint: true, Tailwind: true, AppRouter: true})
	joined := strings.Join(args, " ")
	if !strings.Contains(joined, "--yes") || strings.Contains(joined, " --y ") || strings.Contains(joined, "--no-lint") {
		t.Fatalf("unexpected create-next-app arguments: %s", joined)
	}
}

func TestGoDependabotTargetsRootModule(t *testing.T) {
	config := dependabot("go")
	if !strings.Contains(config, "package-ecosystem: gomod\n    directory: /\n") {
		t.Fatalf("go profile must target the root Go module: %s", config)
	}
	if strings.Contains(config, "directory: /backend") {
		t.Fatal("go profile must not target /backend")
	}
}

func TestNextGoDependabotTargetsBackendModule(t *testing.T) {
	config := dependabot("next-go")
	if !strings.Contains(config, "package-ecosystem: gomod\n    directory: /backend\n") {
		t.Fatalf("next-go profile must target the backend module: %s", config)
	}
}

func TestNextGoMakefileUsesUpTarget(t *testing.T) {
	makefile := profileMakefile("next-go")
	if !strings.Contains(makefile, ".PHONY: help up down") || !strings.Contains(makefile, "\nup:\n") {
		t.Fatalf("next-go Makefile must provide up target: %s", makefile)
	}
	if strings.Contains(makefile, "\ndev:\n") {
		t.Fatal("next-go Makefile must not provide dev target")
	}
}

func TestGoMakefileUsesCorrectEntrypoint(t *testing.T) {
	makefile := profileMakefile("go")
	if !strings.Contains(makefile, "\nup:\n\tdocker compose up\n") {
		t.Fatalf("go Makefile must provide docker up target: %s", makefile)
	}
	if !strings.Contains(makefile, "\nrun:\n\tgo run ./cmd/server\n") {
		t.Fatalf("go Makefile must run cmd/server: %s", makefile)
	}
	if strings.Contains(makefile, "\tgo run .\n") {
		t.Fatal("go Makefile must not run the repository root")
	}
}
