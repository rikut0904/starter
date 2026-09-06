package cli

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	starter "github.com/rikut0904/starter"
	"github.com/rikut0904/starter/internal/domain"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

type Config = domain.Config

type commandRunner func(string, string, ...string) error

type externalFailure struct {
	cwd  string
	name string
	args []string
	err  error
}

var runCommand commandRunner = func(dir, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Dir, cmd.Stdout, cmd.Stderr = dir, os.Stdout, os.Stderr
	return cmd.Run()
}

var profiles = []struct{ id, title string }{
	{"next-go", "Next.js + Go + Docker Compose + Makefile"},
	{"nextjs", "Next.jsのみ"},
	{"go", "Goのみ"},
	{"empty", "空環境"},
}

func RootCommand() *cobra.Command {
	var output, configPath string
	var profile string
	var force, nonInteractive bool
	cmd := &cobra.Command{Use: "starter <project-name>", Short: "開発用プロジェクトスターター", Args: cobra.MaximumNArgs(1), RunE: func(cmd *cobra.Command, args []string) error {
		cfg := Config{}
		if configPath != "" {
			if err := loadConfig(configPath, &cfg); err != nil {
				return err
			}
		}
		if len(args) == 1 {
			cfg.Name = args[0]
		}
		if profile != "" {
			cfg.Profile = profile
		}
		if cfg.Name == "" {
			cfg.Name = promptText("プロジェクト名", "my-project")
		}
		if err := validateProjectName(cfg.Name); err != nil {
			return err
		}
		if !nonInteractive {
			if err := interactiveConfig(&cfg); err != nil {
				return err
			}
		}
		if cfg.Profile == "" {
			cfg.Profile = "empty"
		}
		if output == "" {
			output = filepath.Join(".", cfg.Name)
		}
		abs, err := filepath.Abs(filepath.Clean(output))
		if err != nil {
			return err
		}
		fmt.Printf("生成先: %s\n", abs)
		if !force {
			fmt.Print("この場所に生成しますか？ [y/N] ")
			var answer string
			fmt.Scanln(&answer)
			if strings.ToLower(answer) != "y" {
				return errors.New("生成をキャンセルしました")
			}
		}
		return generate(abs, cfg, force, runCommand)
	}}
	cmd.Flags().StringVar(&configPath, "config", "", "YAMLまたはJSONの設定ファイル")
	cmd.Flags().StringVar(&output, "output", "", "生成先ディレクトリ")
	cmd.Flags().StringVar(&profile, "profile", "", "非対話時のプロファイル（next-go, nextjs, go, empty）")
	cmd.Flags().BoolVar(&force, "force", false, "既存ファイルを上書きする")
	cmd.Flags().BoolVar(&nonInteractive, "non-interactive", false, "対話を省略する")
	return cmd
}

func validateProjectName(name string) error {
	if name == "" || name == "." || name == ".." || filepath.Base(name) != name || strings.ContainsAny(name, `/\\`) {
		return fmt.Errorf("プロジェクト名は単一ディレクトリ名で指定してください: %q", name)
	}
	return nil
}

func loadConfig(path string, cfg *Config) error {
	b, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	if strings.HasSuffix(strings.ToLower(path), ".json") {
		return json.Unmarshal(b, cfg)
	}
	return yaml.Unmarshal(b, cfg)
}

type item struct{ title, desc string }

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }
func interactiveConfig(cfg *Config) error {
	choices := make([]list.Item, len(profiles))
	for i, p := range profiles {
		choices[i] = item{p.title, p.id}
	}
	l := list.New(choices, list.NewDefaultDelegate(), 60, 12)
	l.Title = "プロジェクト構成を選択してください"
	model := tea.NewProgram(listModel{list: l})
	result, err := model.Run()
	if err != nil {
		return err
	}
	selectedModel := result.(listModel)
	if selectedModel.cancelled || selectedModel.list.SelectedItem() == nil {
		return errors.New("構成選択をキャンセルしました")
	}
	selected := selectedModel.list.SelectedItem().(item)
	cfg.Profile = selected.desc
	if cfg.Profile == "next-go" || cfg.Profile == "nextjs" {
		if err := interactiveNextConfig(&cfg.Next); err != nil {
			return err
		}
	}
	return nil
}

type nextChoice struct {
	label string
	value *bool
}

type nextConfigModel struct {
	choices   []nextChoice
	cursor    int
	cancelled bool
}

func interactiveNextConfig(cfg *domain.NextConfig) error {
	// These are create-next-app's recommended defaults, exposed before generation.
	if !cfg.TypeScript && !cfg.ESLint && !cfg.Tailwind && !cfg.SrcDir && !cfg.AppRouter {
		cfg.TypeScript, cfg.ESLint, cfg.Tailwind, cfg.AppRouter = true, true, true, true
	}
	cfg.Configured = true
	m := nextConfigModel{choices: []nextChoice{
		{"TypeScriptを使用", &cfg.TypeScript},
		{"ESLintを使用", &cfg.ESLint},
		{"Tailwind CSSを使用", &cfg.Tailwind},
		{"src/ディレクトリを使用", &cfg.SrcDir},
		{"App Routerを使用", &cfg.AppRouter},
	}}
	result, err := tea.NewProgram(m).Run()
	if err != nil {
		return err
	}
	if result.(nextConfigModel).cancelled {
		return errors.New("Next.js設定をキャンセルしました")
	}
	return nil
}

func (m nextConfigModel) Init() tea.Cmd { return nil }
func (m nextConfigModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if key, ok := msg.(tea.KeyMsg); ok {
		switch key.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case " ":
			*m.choices[m.cursor].value = !*m.choices[m.cursor].value
		case "enter":
			return m, tea.Quit
		case "q", "ctrl+c", "esc":
			m.cancelled = true
			return m, tea.Quit
		}
	}
	return m, nil
}
func (m nextConfigModel) View() string {
	var b strings.Builder
	b.WriteString("\nNext.jsの設定を選択してください（Spaceで切替、Enterで決定）\n\n")
	for i, choice := range m.choices {
		cursor := "  "
		if i == m.cursor {
			cursor = "❯ "
		}
		mark := "[ ]"
		if *choice.value {
			mark = "[x]"
		}
		b.WriteString(fmt.Sprintf("%s%s %s\n", cursor, mark, choice.label))
	}
	return lipgloss.NewStyle().Margin(1, 2).Render(b.String())
}

type listModel struct {
	list      list.Model
	cancelled bool
}

func (m listModel) Init() tea.Cmd { return nil }
func (m listModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			return m, tea.Quit
		case "q", "ctrl+c", "esc":
			m.cancelled = true
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.list.SetSize(msg.Width, msg.Height)
	}
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}
func (m listModel) View() string { return lipgloss.NewStyle().Margin(1, 2).Render(m.list.View()) }
func promptText(label, fallback string) string {
	fmt.Printf("%s [%s]: ", label, fallback)
	r := bufio.NewReader(os.Stdin)
	s, _ := r.ReadString('\n')
	s = strings.TrimSpace(s)
	if s == "" {
		return fallback
	}
	return s
}

func generate(out string, cfg Config, force bool, runner commandRunner) error {
	if err := os.MkdirAll(out, 0755); err != nil {
		return err
	}
	if cfg.Profile != "next-go" && cfg.Profile != "nextjs" && cfg.Profile != "go" && cfg.Profile != "empty" {
		return fmt.Errorf("不明なプロファイル: %s", cfg.Profile)
	}
	if (cfg.Profile == "next-go" || cfg.Profile == "nextjs") && !cfg.Next.Configured && !cfg.Next.TypeScript && !cfg.Next.ESLint && !cfg.Next.Tailwind && !cfg.Next.SrcDir && !cfg.Next.AppRouter {
		cfg.Next.TypeScript, cfg.Next.ESLint, cfg.Next.Tailwind, cfg.Next.AppRouter = true, true, true, true
	}
	failures := []externalFailure{}
	if cfg.Profile == "next-go" {
		if failure := reportExternal(out, "npx", nextCommand("frontend", cfg.Next), runner); failure != nil {
			failures = append(failures, *failure)
		}
		if failure := reportExternal(filepath.Join(out, "backend"), "go", []string{"mod", "init", "backend"}, runner); failure != nil {
			failures = append(failures, *failure)
		}
	}
	if cfg.Profile == "nextjs" {
		if failure := reportExternal(filepath.Dir(out), "npx", nextCommand(filepath.Base(out), cfg.Next), runner); failure != nil {
			failures = append(failures, *failure)
		}
	}
	if cfg.Profile == "go" {
		if failure := reportExternal(out, "go", []string{"mod", "init", cfg.Name}, runner); failure != nil {
			failures = append(failures, *failure)
		}
	}
	files := commonFiles(cfg)
	for path, body := range files {
		if err := writeFile(out, path, body, force); err != nil {
			return err
		}
	}
	if len(failures) > 0 {
		fmt.Println("\n生成処理は完了しましたが、外部コマンドに失敗があります。")
		for _, failure := range failures {
			fmt.Printf("\n失敗: (cd %s && %s %s)\n原因: %v\n再実行: cd %s && %s %s\n", failure.cwd, failure.name, strings.Join(failure.args, " "), failure.err, failure.cwd, failure.name, strings.Join(failure.args, " "))
		}
	}
	return nil
}

func nextCommand(name string, cfg domain.NextConfig) []string {
	args := []string{"create-next-app@latest", name, "--yes"}
	if cfg.TypeScript {
		args = append(args, "--ts")
	} else {
		args = append(args, "--js")
	}
	if cfg.ESLint {
		args = append(args, "--eslint")
	} else {
		args = append(args, "--no-eslint")
	}
	if cfg.Tailwind {
		args = append(args, "--tailwind")
	} else {
		args = append(args, "--no-tailwind")
	}
	if cfg.SrcDir {
		args = append(args, "--src-dir")
	} else {
		args = append(args, "--no-src-dir")
	}
	if cfg.AppRouter {
		args = append(args, "--app")
	} else {
		args = append(args, "--no-app")
	}
	return args
}
func reportExternal(out, name string, args []string, runner commandRunner) *externalFailure {
	_ = os.MkdirAll(out, 0755)
	fmt.Printf("実行: %s %s\n", name, strings.Join(args, " "))
	if err := runner(out, name, args...); err != nil {
		return &externalFailure{cwd: out, name: name, args: append([]string(nil), args...), err: err}
	}
	return nil
}
func writeFile(root, rel, body string, force bool) error {
	path := filepath.Join(root, filepath.Clean(rel))
	if !strings.HasPrefix(path, filepath.Clean(root)+string(os.PathSeparator)) {
		return errors.New("不正な生成パス")
	}
	if _, err := os.Stat(path); err == nil && !force {
		return fmt.Errorf("既存ファイルのため停止しました: %s（--forceで上書き）", rel)
	}
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(body), 0644)
}

func commonFiles(cfg Config) map[string]string {
	result := map[string]string{}
	for _, name := range []string{"AGENTS.md", "CONTRIBUTING.md", "LICENSE", "README.md", "Makefile", ".gitignore", ".github/CODEOWNERS", ".github/dependabot.yml", ".github/workflows/ci.yml", ".github/workflows/dependency-review.yml"} {
		b, err := starter.Assets.ReadFile(filepath.Join("create", "common", name))
		if err == nil {
			result[name] = string(b)
		}
	}
	result["README.md"] = strings.ReplaceAll(result["README.md"], "{{PROJECT_NAME}}", cfg.Name)
	result["README.md"] = strings.ReplaceAll(result["README.md"], "{{PROFILE}}", cfg.Profile)
	result["AGENTS.md"] = strings.ReplaceAll(result["AGENTS.md"], "{{PROJECT_NAME}}", cfg.Name)
	result["AGENTS.md"] = strings.ReplaceAll(result["AGENTS.md"], "{{PROFILE}}", cfg.Profile)
	result["LICENSE"] = strings.ReplaceAll(strings.ReplaceAll(result["LICENSE"], "{{PROJECT_NAME}}", cfg.Name), "{{YEAR}}", fmt.Sprint(time.Now().Year()))
	result[".github/dependabot.yml"] = dependabot(cfg.Profile)
	result[".github/workflows/ci.yml"] = ciWorkflow(cfg.Profile)
	result["Makefile"] = profileMakefile(cfg.Profile)
	for _, name := range profileFiles(cfg.Profile) {
		b, err := starter.Assets.ReadFile(filepath.Join("create", cfg.Profile, name))
		if err == nil {
			body := strings.ReplaceAll(string(b), "{{MODULE}}", cfg.Name)
			body = strings.ReplaceAll(body, "github.com/example/starter-template-go", cfg.Name)
			body = strings.ReplaceAll(body, "github.com/rikut0904/starter/create/go", cfg.Name)
			if cfg.Profile == "next-go" {
				body = strings.ReplaceAll(body, "github.com/rikut0904/starter/create/next-go/backend", "backend")
				body = strings.ReplaceAll(body, "github.com/rikut0904/starter/create/next-go", "backend")
			} else {
				body = strings.ReplaceAll(body, "github.com/rikut0904/starter/create/next-go", cfg.Name)
			}
			result[name] = body
		}
	}
	return result
}

func dependabot(profile string) string {
	base := "version: 2\nupdates:\n"
	if profile == "next-go" {
		base += "  - package-ecosystem: npm\n    directory: /frontend\n    schedule: { interval: weekly }\n    open-pull-requests-limit: 5\n    labels: [dependencies]\n"
	}
	if profile == "nextjs" {
		base += "  - package-ecosystem: npm\n    directory: /\n    schedule: { interval: weekly }\n    open-pull-requests-limit: 5\n    labels: [dependencies]\n"
	}
	if profile == "next-go" || profile == "go" {
		base += "  - package-ecosystem: gomod\n    directory: /backend\n    schedule: { interval: weekly }\n    open-pull-requests-limit: 5\n    labels: [dependencies]\n"
	}
	base += "  - package-ecosystem: github-actions\n    directory: /\n    schedule: { interval: weekly }\n    open-pull-requests-limit: 5\n    labels: [dependencies]\n"
	return base
}

func ciWorkflow(profile string) string {
	check := "      - run: git diff --check\n"
	switch profile {
	case "next-go":
		check += "      - uses: actions/setup-node@v4\n        with: { node-version: 22, cache: npm, cache-dependency-path: frontend/package-lock.json }\n      - run: npm ci\n        working-directory: frontend\n      - run: npm run lint\n        working-directory: frontend\n      - uses: actions/setup-go@v5\n        with: { go-version: '1.26' }\n      - run: go test ./...\n        working-directory: backend\n      - run: docker compose config\n"
	case "nextjs":
		check += "      - uses: actions/setup-node@v4\n        with: { node-version: 22, cache: npm }\n      - run: npm ci\n      - run: npm run lint\n      - run: npm run build\n"
	case "go":
		check += "      - uses: actions/setup-go@v5\n        with: { go-version: '1.26' }\n      - run: go test ./...\n      - run: go vet ./...\n      - run: docker compose config\n"
	}
	return "name: CI\non:\n  pull_request:\n  push:\n    branches: [main]\npermissions:\n  contents: read\njobs:\n  basic:\n    runs-on: ubuntu-latest\n    steps:\n      - uses: actions/checkout@v4\n" + check
}

func profileMakefile(profile string) string {
	switch profile {
	case "next-go":
		return ".PHONY: help dev down logs test lint build format\nhelp:\n\t@echo \"dev down logs test lint build format\"\ndev:\n\tdocker compose up\ndown:\n\tdocker compose down\nlogs:\n\tdocker compose logs -f\ntest:\n\tcd backend && go test ./...\nlint:\n\tcd frontend && npm run lint\nbuild:\n\tdocker compose build\nformat:\n\tcd backend && gofmt -w .\n"
	case "nextjs":
		return ".PHONY: help dev down logs test lint build format\nhelp:\n\t@echo \"dev down logs test lint build format\"\ndev:\n\tnpm run dev\ndown:\n\tdocker compose down\nlogs:\n\tdocker compose logs -f\ntest:\n\tnpm test\nlint:\n\tnpm run lint\nbuild:\n\tnpm run build\nformat:\n\tnpx prettier --write .\n"
	case "go":
		return ".PHONY: help dev down logs test lint build format\nhelp:\n\t@echo \"dev down logs test lint build format\"\ndev:\n\tgo run .\ndown:\n\tdocker compose down\nlogs:\n\tdocker compose logs -f\ntest:\n\tgo test ./...\nlint:\n\tgo vet ./...\nbuild:\n\tgo build ./...\nformat:\n\tgofmt -w .\n"
	default:
		return ".PHONY: help dev down logs test lint build format\nhelp:\n\t@echo \"利用可能: dev down logs test lint build format\"\ndev:\n\t@echo \"開発環境を構築してください\"\ndown:\n\t@echo \"停止対象はありません\"\nlogs:\n\t@echo \"ログ対象はありません\"\ntest lint build format:\n\t@echo \"対象を追加してください\"\n"
	}
}
func profileFiles(profile string) []string {
	switch profile {
	case "next-go":
		return []string{"docker-compose.yml", "frontend/Dockerfile", "backend/Dockerfile", ".env.example", "backend/cmd/server/main.go", "backend/internal/domain/health.go", "backend/internal/usecase/health.go", "backend/internal/interface/http/handler.go", "backend/internal/infrastructure/config/config.go"}
	case "nextjs":
		return []string{"docker-compose.yml", "Dockerfile", ".env.example"}
	case "go":
		return []string{"docker-compose.yml", "Dockerfile", ".env.example", "cmd/server/main.go", "internal/domain/health.go", "internal/usecase/health.go", "internal/interface/http/handler.go", "internal/infrastructure/config/config.go"}
	default:
		return nil
	}
}
