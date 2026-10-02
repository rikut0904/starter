# Go 推奨技術スタック

この構成は、REST API、バックエンドサービス、CLI、バッチなどをGo単体で構築するための標準スタックです。

## 推奨構成

| 分類 | 技術 | 用途 |
|---|---|---|
| Language | Go | API、CLI、バッチ |
| HTTP | Echo / net/http | HTTPルーティング |
| API Contract | OpenAPI + oapi-codegen | API仕様・型・サーバー生成 |
| Validation | go-playground/validator 等 | 入力検証 |
| Auth | Auth0 JWT / Firebase Admin SDK | API認証 |
| Database | PostgreSQL | 永続化 |
| ORM | GORM | DBアクセス |
| Migration | goose / Atlas | DBスキーマ変更管理 |
| Storage | S3互換Object Storage | ファイル保存 |
| Cache | Redis（必要時） | cache / rate limit / queue |
| Logging | slog | 構造化ログ |
| Observability | OpenTelemetry / Sentry | trace・error monitoring |
| Test | testing + httptest | Unit / HTTP test |
| CI | GitHub Actions | test / vet / static analysis / build |
| Deploy | Railway / container platform | APIデプロイ |

## 基本アーキテクチャ

```text
cmd/server
   ↓
internal/interface/http
   ↓
internal/usecase
   ↓
internal/domain
   ↓
internal/infrastructure
   ├─ PostgreSQL
   ├─ Object Storage
   └─ External APIs
```

starterの既存クリーンアーキテクチャを維持します。ただし、小規模なAPIで層を増やすこと自体が目的にならないよう、必要に応じて構造を簡略化して構いません。

## API設計

OpenAPIを追加し、API仕様を実装前後の共通契約として利用します。oapi-codegenによってGoの型やサーバーインターフェースを生成し、手書きDTOの重複を減らします。

## 認証

### Auth0
Web APIを独立したResource Serverとして設計する場合の第一候補です。Access Tokenの署名、issuer、audience、scope/permissionをGo側で検証します。

### Firebase Authentication
Flutter/Android/Firebase中心のサービスではFirebase Admin SDKを利用してID Tokenを検証します。

## DB Migration

GORM AutoMigrateだけを本番スキーマ管理の唯一の手段にせず、gooseやAtlasでMigrationを履歴管理します。

## Logging / Observability

Go標準のlog/slogを基本ログとして利用します。サービスが増えた場合や本番運用ではOpenTelemetryを追加し、trace idをHTTPリクエスト間で引き継げる設計にします。

## 採用ルール

- Echoはstarter標準として維持してよい。
- 新規APIにはOpenAPI導入を推奨する。
- PostgreSQL + GORMを標準とする。
- 本番利用ではMigrationツールを利用する。
- Redisは明確な用途がある場合だけ追加する。
- 認証はAuth0を基本、Firebase系案件ではFirebase Authを使う。
- ログはslogを基本とし、Observabilityは段階導入する。

## 公式URL

- Go: https://go.dev/
- Echo: https://echo.labstack.com/
- OpenAPI: https://www.openapis.org/
- oapi-codegen: https://github.com/oapi-codegen/oapi-codegen
- PostgreSQL: https://www.postgresql.org/
- GORM: https://gorm.io/
- goose: https://github.com/pressly/goose
- Atlas: https://atlasgo.io/
- Auth0: https://auth0.com/
- Firebase Authentication: https://firebase.google.com/docs/auth
- Redis: https://redis.io/
- OpenTelemetry: https://opentelemetry.io/
- Sentry: https://sentry.io/
- Railway: https://railway.com/
