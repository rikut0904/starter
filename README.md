# 開発用プロジェクトスターター

新規開発プロジェクトの基本構成、開発ルール、CI、Dependabot設定などをCLIで生成するためのスターターです。

## 基本方針

- CLI内のラジオボタン形式の画面で構成を選択する
- 生成先はコマンド実行ディレクトリ配下とする
- LICENSE、`AGENTS.md`、README、CI、Dependabotなどを生成する
- リポジトリ作成後の外部Webhookや自動コミットに依存しない

## ローカルセットアップ

```bash
make init/mac   # macOS
make init/linux # Linux
make init/win   # Windows
```

CLIを削除する場合は `make uninstall` を使用します。

セットアップ後は、任意の作業ディレクトリで利用できます。

```bash
starter my-project
```

## 選択できる構成

1. Next.js + Go + Docker Compose + Makefile
2. Next.jsのみ
3. Goのみ
4. 空環境

構成や追加機能は、`starter` 実行後のCLI画面で選択します。

認証機能で「共通認証（common-id）を使用」を選択した場合、`common-id` がPATHにあるか確認し、導入済みなら `common-id install` を実行します。見つからない場合は、CLIが次の導入手順を表示します。

```bash
git clone https://github.com/rikut0904/common-id.git
cd common-id
make init/commond
```

## 生成先の例

```text
/work/projects/
└── my-project/
    ├── AGENTS.md
    ├── LICENSE
    ├── Makefile
    └── ...
```

## リポジトリ構成

```text
.
├── create/       # CLIが利用するテンプレート・生成用素材の配置領域
├── docs/         # 詳細仕様・運用ドキュメント
├── Makefile      # starter CLIのローカルセットアップ
├── LICENSE
├── LICENSE_JA
└── README.md
```

`create/` はstarter CLIが参照する生成用素材を管理するディレクトリです。プロファイル別のテンプレート、共通ファイル、生成時に利用する設定・定義を配置します。CLIで作成されるプロジェクト自体は `create/` 配下ではなく、CLIを実行したカレントディレクトリ配下に生成します。

## ドキュメント

- [全体仕様](docs/specification.md)
- [CLI仕様とMakefile](docs/cli.md)
- [生成プロファイル](docs/profiles.md)
- [運用・Dependabot・CI](docs/operations.md)

## ライセンス

MIT License。詳細は [LICENSE](LICENSE) および [LICENSE_JA](LICENSE_JA) を参照してください。
