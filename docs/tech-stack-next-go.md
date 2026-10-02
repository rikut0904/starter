# Next.js + Go 推奨技術スタック

この構成は、WebフロントとAPIを分離しながら、個人開発でも運用しやすくするための標準スタックです。UIはNext.js、ビジネスAPIはGo、認証はAuth0、DBはPostgreSQLを基本とします。

## 推奨構成

| 分類 | 技術 | 用途 |
|---|---|---|
| Frontend | Next.js / React / TypeScript | UI、SSR、Server Components、Web側の処理 |
| Styling | Tailwind CSS | レイアウトとデザイン |
| UI | shadcn/ui + Base UI | 再利用可能なUIコンポーネント |
| Form | React Hook Form | 複雑なフォーム状態管理 |
| Validation | Zod | Next.js側の入力・スキーマ検証 |
| Auth | Auth0 | OAuth 2.0 / OIDC、JWT、Google等のログイン |
| Backend | Go + Echo | REST API、ビジネスロジック |
| API Contract | OpenAPI + oapi-codegen | API仕様、Go型・サーバーコード生成 |
| Database | PostgreSQL | 永続化 |
| ORM | GORM | GoからのDBアクセス |
| Migration | goose / Atlas | DBスキーマ変更管理 |
| Object Storage | S3互換ストレージ | 画像・添付ファイル |
| Payment | Stripe（任意） | 課金・サブスクリプション |
| Observability | Sentry / OpenTelemetry（段階導入） | エラー・トレース |
| CI | GitHub Actions | lint / test / build / OpenAPI差分検査 |
| Deploy | Vercel + Railway | Frontend / Go API / PostgreSQL |

## 認証フロー

```text
Browser
  ↓
Auth0
  ↓ Access Token (JWT)
Next.js
  ↓
Go API
  ↓ JWT検証 + scope/permission
PostgreSQL
```

Go APIを独立したResource Serverとして扱えるため、今後Flutterや別Webクライアントを追加する場合にも拡張しやすい構成です。

## API設計

`openapi/openapi.yaml` をAPI契約の中心とします。

```text
openapi/openapi.yaml
       ↓
 oapi-codegen
       ↓
Go models / server interface
       ↓
Go implementation
```

フロント側もOpenAPIからTypeScript Clientを生成する構成へ拡張できます。Next.jsとGoで同じDTOを別々に手書きすることを減らすのが目的です。

## React Hook Formの扱い

依存関係として標準搭載して構いません。ただし、すべての`<form>`をRHFへ置き換えるルールにはしません。動的フィールド、複数ステップ、即時エラー表示などがあるフォームで使い、単純な検索やServer Action送信ではHTML formを優先します。

## 推奨ディレクトリ例

```text
project/
├── frontend/
│   ├── app/
│   ├── components/
│   │   └── ui/
│   ├── lib/
│   │   ├── auth/
│   │   ├── api/
│   │   └── validation/
│   └── package.json
├── backend/
│   ├── cmd/server/
│   ├── internal/domain/
│   ├── internal/usecase/
│   ├── internal/interface/http/
│   └── internal/infrastructure/
├── openapi/
│   └── openapi.yaml
├── migrations/
├── docker-compose.yml
└── Makefile
```

## 採用ルール

- Auth0はWeb + 独立Go APIを基本とする新規案件で第一候補。
- Firebase中心のモバイル案件はFirebase Authを利用してよい。
- Redisはキャッシュ・rate limit・queueなど明確な用途がある場合だけ追加。
- Stripeは課金がある場合だけ追加。
- OpenTelemetryは初期MVPでは必須にせず、本番運用・分散構成で追加。
- DB Migrationは本番利用する案件では必須とする。

## 公式URL

- Next.js: https://nextjs.org/
- React: https://react.dev/
- Tailwind CSS: https://tailwindcss.com/
- shadcn/ui: https://ui.shadcn.com/
- React Hook Form: https://react-hook-form.com/
- Zod: https://zod.dev/
- Auth0: https://auth0.com/
- Go: https://go.dev/
- Echo: https://echo.labstack.com/
- OpenAPI: https://www.openapis.org/
- oapi-codegen: https://github.com/oapi-codegen/oapi-codegen
- PostgreSQL: https://www.postgresql.org/
- GORM: https://gorm.io/
- Stripe: https://stripe.com/
- Sentry: https://sentry.io/
- OpenTelemetry: https://opentelemetry.io/
- Vercel: https://vercel.com/
- Railway: https://railway.com/
