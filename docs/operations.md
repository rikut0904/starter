# 運用・Dependabot・CI

## 1. Dependabot

プロファイルに応じて`.github/dependabot.yml`を生成する。

| プロファイル | 監視対象 |
|---|---|
| `next-go` | frontendのnpm、backendのGo Modules、GitHub Actions |
| `nextjs` | npm、GitHub Actions |
| `go` | backendのGo Modules、GitHub Actions |
| `empty` | GitHub Actions |

標準で週次更新、同時オープンPR数の上限、`dependencies`等のラベル、GitHub Actions自身の監視を設定する。

## 2. CI

`.github/workflows/ci.yml`をPRとmainへのpushで実行する。

| プロファイル | 必須チェック |
|---|---|
| `next-go` | frontend lint/test/build、backend fmt/vet/test/build、Compose検証 |
| `nextjs` | format/lint/test/build |
| `go` | backend fmt/vet/test/build、Compose検証 |
| `empty` | YAML、Markdown、設定ファイルの基本検証 |

依存関係変更を含むPRでは`dependency-review.yml`で脆弱性・ライセンスを確認する。存在しないコマンドをCIへ記載しない。

## 3. GitHub Projects

Dependabot PRのProject登録・ステータス同期は任意機能とする。Project ID等はSecrets/Variablesから取得し、同じPRを重複登録しない。PRコード実行とProject書き込みはジョブを分離し、権限を最小化する。

## 4. セキュリティと検証

- 秘密情報をテンプレートへ埋め込まない
- `.env`、秘密鍵、トークンを`.gitignore`に含める
- Actionsの権限をWorkflow単位で明示する
- PR本文やタイトルをシェルコマンドとして評価しない
- ログへ認証情報を出力しない
- `pull_request_target`でPRブランチのコードを実行しない

生成テストの成功だけで、ブラウザ表示、Docker起動、GitHub認証、Project連携、本番環境の正常性を証明したことにしない。
