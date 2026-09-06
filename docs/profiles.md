# 生成プロファイル

## 1. Next.js + Go (`next-go`)

```text
project/
├── frontend/              # Next.js
├── backend/               # Go API
├── docker-compose.yml
├── Makefile
└── .env.example
```

frontend/backend用のDockerfile、Compose設定、Goモジュール、Next.js最小構成を生成する。生成時には次のコマンドを各1回だけ試行する。

```bash
npx create-next-app@latest frontend -y
go mod init backend
```

コマンドが利用できない、ネットワークや依存関係の問題で失敗するなどの場合は、CLIが原因と再実行コマンドを表示する。CLIは自動で繰り返さず、共通ファイルや生成可能なファイルの処理を継続するか、失敗状況を示して安全に終了する。生成先のMakefileには`dev`、`down`、`logs`、`test`、`lint`、`build`、`format`、`help`を用意する。コンテナ内のサーバーは`0.0.0.0`にbindし、コンテナ間通信にはサービス名を使用する。
なお、frontendは`npx create-next-app@latest frontend -y`にて基本プロジェクトを実装するようにする。また、backendは`go mod init backend`にて初期化するようにする。

## 2. Next.js (`nextjs`)

Next.js App Router、TypeScript、lint、format、test、build設定を生成する。生成時に`npx create-next-app@latest <project-name> -y`を1回だけ試行する。失敗時は利用者が再実行するための実際のコマンドとエラー概要を表示する。Dockerfile、Compose設定、MakefileはCLIで追加選択できる。

## 3. Go (`go`)

`go mod init <project-name>`を1回だけ試行してGoモジュールを初期化する。失敗時は利用者が再実行するための実際のコマンドを表示する。また、docker-compose.yml、Dockerfile、Makefileを生成する。

## 4. 空環境 (`empty`)

言語やフレームワークを固定せず、共通運用ファイルとプレースホルダーのMakefileだけを生成する。未導入のツールを暗黙にインストールしない。

## 5. プロファイル別ファイル

| ファイル | `next-go` | `nextjs` | `go` | `empty` |
|---|:---:|:---:|:---:|:---:|
| `frontend/` | ○ | - | - | - |
| `backend/` | ○ | - | - | - |
| `docker-compose.yml` | ○ | ○ | ○ | - |
| `Dockerfile` | ○ | ○ | ○ | - |
| `package.json` | ○ | ○ | - | - |
| `go.mod` | ○ | - | ○ | - |
| `.env.example` | ○ | ○ | ○ | - |
