# 開発用プロジェクトスターター 全体仕様

## 1. 目的

新規プロジェクトのディレクトリ構成、Docker、Makefile、CI、Dependabot、Issue/PRテンプレートなどをCLIで標準化する。リポジトリ作成後の外部Webhookや自動コミットには依存せず、CLI実行後に利用者が内容を確認して初回コミットする。

CLI本体はGoで実装し、単一バイナリとしてmacOS、Linux、Windowsで利用できるようにする。TUIはBubble Tea系ライブラリを採用し、生成ロジックは画面から分離する。

## 2. 基本フロー

```text
make init/<os>
       ↓
starter <project-name>
       ↓
CLI内で構成をラジオボタン選択
       ↓
追加機能をチェックボックス選択
       ↓
カレントディレクトリ配下へ生成
       ↓
利用者が確認してGit管理・初回コミット
```

## 3. 全プロファイル共通ファイル

### 3.1 `create/` の位置づけ

starterリポジトリの `create/` は、CLIが参照する生成用素材を管理する領域である。プロファイル別テンプレート、共通ファイル、生成時の設定・定義を配置する。

`create/` は生成先ではない。CLIで作成するプロジェクトは、CLIを実行したカレントディレクトリ配下の `./<project-name>/` に生成する。`--output`指定時は指定先を事前表示し、既存ファイルを無断で上書きしない。

```text
starter/
└── create/
    ├── common/
    ├── next-go/
    ├── nextjs/
    ├── go/
    └── empty/
```

生成先の `<project-name>/` に以下を生成する。

```text
AGENTS.md
CONTRIBUTING.md
LICENSE
LICENSE_JA
README.md
Makefile
.gitignore
.github/
├── dependabot.yml
└── workflows/
    ├── ci.yml
    └── dependency-review.yml
```

`AGENTS.md`には、目的、構成、開発・テスト・ビルドコマンド、規約、シークレットの扱い、検証範囲を記載する。サブディレクトリ固有の指示は下位の`AGENTS.md`へ分割できる。

## 4. 生成先の安全性

- デフォルトはCLIを実行したディレクトリ配下の `./<project-name>/`
- `--output`指定時は対象パスを事前表示する
- 既存ファイルはデフォルトで上書きしない
- `--force`指定時だけ上書きする
- パスを正規化し、意図しない場所へ書き込まない
- 秘密情報、絶対パス、認証情報を生成物へ含めない
- CLIは自動でcommit、push、外部リポジトリ操作を行わない

外部生成コマンドを利用するプロファイルでは、`npx create-next-app`や`go mod init`を各1回だけ試行する。失敗時はエラー内容と利用者向けの再実行コマンドを表示し、自動リトライや失敗の隠蔽は行わない。

## 5. ローカルセットアップ

`make init/*`はstarter CLIをそのPCで利用可能にするためのOS別セットアップを提供する。前提条件の確認、設定、CLIのビルド・インストールは各OS別ターゲットの内部処理とし、独立したcheck/path/installターゲットは設けない。

| ターゲット | 内容 |
|---|---|
| `make init/mac` | macOSの確認、設定、CLIインストール |
| `make init/linux` | Linuxの確認、設定、CLIインストール |
| `make init/win` | Windowsの確認、設定、CLIインストール |
| `make uninstall` | CLIを削除。確認必須 |

前提条件の確認に失敗した場合、OS別ターゲットは設定・インストールを実行してはならない。シェル設定ファイルやWindowsのユーザー環境変数を自動編集せず、必要な操作を案内する。

## 6. 受け入れ条件

- CLI内で構成をラジオボタン形式に選択できる
- `next-go`、`nextjs`、`go`、`empty`を生成できる
- 生成物がコマンド実行ディレクトリ配下に作成される
- 共通ファイルが全プロファイルで生成される
- 既存ファイルが無断で上書きされない
- macOS、Linux、WindowsでOS別の初期化ターゲットを利用できる
- アンインストールを`make uninstall`で実行できる
- 必要なプロファイルで外部生成コマンドを各1回試行できる
- 外部コマンド失敗時にエラー内容と手動実行コマンドが表示される
