.DEFAULT_GOAL := help
BIN_DIR ?= $(shell go env GOPATH)/bin

.PHONY: help init init/mac init/linux init/win uninstall test fmt
help:
	@echo "make init/mac|init/linux|init/win  CLIをビルドしてインストール"
	@echo "make uninstall                    CLIを削除（確認あり）"
init:
	@echo "OS別に make init/mac, make init/linux, make init/win を実行してください"
init/mac:
	@test "$(shell uname -s)" = "Darwin" || (echo "macOSで実行してください"; exit 1)
	@command -v go >/dev/null || (echo "Goをインストールしてください"; exit 1)
	@mkdir -p "$(BIN_DIR)" && go build -o "$(BIN_DIR)/starter" ./cmd/starter
	@echo "インストール完了: $(BIN_DIR)/starter"
	@case "$${SHELL##*/}" in \
		zsh) printf '%s\n' 'zsh: echo '\''export PATH="$$PATH:$(BIN_DIR)"'\'' >> ~/.zshrc && source ~/.zshrc' ;; \
		bash) printf '%s\n' 'bash: echo '\''export PATH="$$PATH:$(BIN_DIR)"'\'' >> ~/.bashrc && source ~/.bashrc' ;; \
		fish) echo "fish: fish_add_path $(BIN_DIR)" ;; \
		*) echo "ログインシェルを判定できませんでした（SHELL=$${SHELL:-未設定}）。$(BIN_DIR)をPATHへ追加してください" ;; \
	esac
init/linux:
	@test "$(shell uname -s)" = "Linux" || (echo "Linuxで実行してください"; exit 1)
	@command -v go >/dev/null || (echo "Goをインストールしてください"; exit 1)
	@mkdir -p "$(BIN_DIR)" && go build -o "$(BIN_DIR)/starter" ./cmd/starter
	@echo "インストール完了: $(BIN_DIR)/starter"
	@case "$${SHELL##*/}" in \
		zsh) printf '%s\n' 'zsh: echo '\''export PATH="$$PATH:$(BIN_DIR)"'\'' >> ~/.zshrc && source ~/.zshrc' ;; \
		bash) printf '%s\n' 'bash: echo '\''export PATH="$$PATH:$(BIN_DIR)"'\'' >> ~/.bashrc && source ~/.bashrc' ;; \
		fish) echo "fish: fish_add_path $(BIN_DIR)" ;; \
		*) echo "ログインシェルを判定できませんでした（SHELL=$${SHELL:-未設定}）。$(BIN_DIR)をPATHへ追加してください" ;; \
	esac
init/win:
	@command -v go >/dev/null || (echo "Goをインストールしてください（PowerShellまたはGit Bash）"; exit 1)
	@mkdir -p "$(BIN_DIR)" && go build -o "$(BIN_DIR)/starter.exe" ./cmd/starter
	@echo "インストール完了: $(BIN_DIR)/starter.exe"
	@echo "PowerShell: [Environment]::SetEnvironmentVariable('Path', \"$$env:Path;$(BIN_DIR)\", 'User')"
	@echo "設定後、PowerShellを再起動してください"
uninstall:
	@printf "$(BIN_DIR)/starter を削除しますか？ [y/N] "; read answer; test "$$answer" = "y" || (echo "キャンセルしました"; exit 1); rm -f "$(BIN_DIR)/starter" "$(BIN_DIR)/starter.exe"
test:
	go test ./...
fmt:
	gofmt -w $$(find . -type f -name '*.go' -not -path './.git/*' -not -path './.gocache/*' -not -path './.gomodcache/*')
