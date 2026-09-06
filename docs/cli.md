# CLI仕様とMakefile

## 技術スタック

CLI本体はGoで実装する。単一バイナリとしてmacOS、Linux、Windowsへ配布し、OSごとのセットアップ処理とプロジェクト生成処理を同じCLIで管理する。

| 領域 | 採用技術 | 用途 |
|---|---|---|
| 言語 | Go | CLI本体、生成処理、環境判定 |
| コマンド構成 | Cobra | コマンド、ヘルプ、補助オプション |
| ターミナルUI | Bubble Tea | ラジオボタン・チェックボックス画面 |
| UI部品 | Bubbles | リスト、選択、入力部品 |
| UIスタイル | Lip Gloss | 色、余白、レイアウト |
| テンプレート | `embed` + `text/template` | `create/`配下の素材を埋め込み・展開 |
| 設定 | YAMLまたはJSON | 非対話実行時の選択内容 |
| テスト | Go標準testing | CLI、生成結果、環境判定 |

Cobraはコマンド構造とヘルプを担当し、Bubble Tea/Bubblesは対話画面を担当する。生成ロジックはTUIから分離し、非対話実行と単体テストからも呼び出せる構造にする。

Next.jsやGoのプロジェクト生成では、CLIが外部コマンドをそれぞれ1回だけ試行する。失敗した場合はエラー内容と利用者が再実行するためのコマンドを表示し、利用者の手動対応へ切り替える。自動リトライ、失敗の隠蔽、無断での別コマンドへの置き換えは行わない。

公式資料: [Cobra](https://github.com/spf13/cobra)、[Bubble Tea](https://github.com/charmbracelet/bubbletea)

## 1. ローカルセットアップ

`make init/*`は、starter CLIをそのPCで利用可能にするためのOS別セットアップ用ターゲットとする。OS別ターゲット自身が前提条件を確認し、CLIのビルドとインストールまで行う。プロジェクト構成の選択は、セットアップ後に`starter`で行う。

```bash
make init/mac
make init/linux
make init/win
```

| ターゲット | 内容 |
|---|---|
| `make init` | OS別ターゲットの使い方を表示 |
| `make init/mac` | macOSの確認、設定、CLIインストール |
| `make init/linux` | Linuxの確認、設定、CLIインストール |
| `make init/win` | Windowsの確認、設定、CLIインストール |
| `make help` | Makeターゲット一覧を表示 |
| `make uninstall` | CLIを削除。確認必須 |

前提条件の確認に失敗した場合、OS別ターゲットは設定・インストールを実行せず、不足しているツールや対応方法を表示する。

インストール先は`BIN_DIR`で変更できる。未指定時はOSごとのGo標準のbinディレクトリを利用する。

```bash
make init/mac BIN_DIR="$HOME/.local/bin"
```

シェル設定ファイルやWindowsのユーザー環境変数は自動編集せず、`SHELL`から利用環境を判定して、該当するbash・zsh・fishの登録コマンドだけを表示する。WindowsではPowerShellの設定方法を表示する。

## OS別仕様

### macOS

- Apple SiliconとIntelの両方に対応する
- zsh/bash向けにPATH登録方法を表示する
- Goのbinディレクトリへユーザー権限でインストールする

### Linux

- CPUアーキテクチャとGoの実行条件を確認する
- root権限を要求せず、ユーザー領域へインストールする
- bash/zsh/fish向けにPATH登録方法を表示する

### Windows

- PowerShellを標準案内とする
- Windows向けバイナリとユーザー用binディレクトリを利用する
- `make`がない場合はGit Bash、WSL等の導入が必要であることを表示する
- PATH登録方法をPowerShell用に表示する

## 2. プロジェクト生成

```bash
starter my-project
```

正式なコマンド形式は `starter <project-name>` とする。引数なしの時にはproject-nameをはじめに聞くようにする。また `--help` 指定時は使い方を表示する。

CLIはstarterリポジトリ内の `create/` をテンプレート・生成用素材の参照元として利用する。生成結果は、コマンドを実行したカレントディレクトリ配下の `./<project-name>/` に保存する。

構成は単一選択のラジオボタン、Docker・Makefile・データベース等は複数選択のチェックボックス、プロジェクト名や著作権者は入力画面で指定する。

```text
プロジェクト構成を選択してください

  ◉ Next.js + Go + Docker Compose + Makefile
  ○ Next.jsのみ
  ○ Goのみ
  ○ 空環境

↑↓ 移動  Enter 決定  q 終了
```

## 3. 補助オプション

構成選択はCLI画面を基本とし、以下は自動実行や再現性が必要な場合の補助機能とする。

| オプション | 内容 | 初期値 |
|---|---|---|
| `--config` | 選択内容を記載した設定ファイル | 未指定 |
| `--output` | 生成先 | `./<name>` |
| `--force` | 既存ファイルを上書き | `false` |
| `--non-interactive` | 対話を省略 | `false` |
| `--profile` | 非対話時の構成（`next-go`、`nextjs`、`go`、`empty`） | 未指定 |
