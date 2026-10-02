# Next.js 推奨技術スタック

この構成は、フロントエンドから認証・サーバー処理までNext.jsで完結させ、個人開発やMVPを最短で構築するための標準スタックです。

## 推奨構成

| 分類 | 技術 | 用途 |
|---|---|---|
| Framework | Next.js / React / TypeScript | UI、SSR、Server Components、Server Functions |
| Styling | Tailwind CSS | レイアウトとデザイン |
| UI | shadcn/ui + Base UI | 再利用可能なUIコンポーネント |
| Form | React Hook Form | 複雑なフォーム状態管理 |
| Validation | Zod | 入力・環境変数・外部データ検証 |
| Auth | Auth0 または Auth.js | 認証・OAuth/OIDC |
| Database | PostgreSQL / Supabase | 永続化 |
| DB Access | Drizzle ORM または Supabase SDK | DB操作 |
| Storage | Supabase Storage / S3互換 | 画像・添付ファイル |
| Server State | TanStack Query（必要時） | クライアントキャッシュ・再取得 |
| Payment | Stripe（任意） | 決済・サブスクリプション |
| Error Monitoring | Sentry | エラー監視 |
| CI | GitHub Actions | lint / test / build |
| Deploy | Vercel | Next.jsデプロイ |

## 認証の選び方

Next.jsだけで完結するWebアプリではAuth.jsも有力です。既存のAuth0運用や、将来別APIを追加する可能性が高い場合はAuth0を選びます。

## React Hook Formの扱い

RHFはスターターに入れて構いませんが、すべてのフォームで使う必要はありません。複雑なフォームではRHF + Zod、単純なフォームでは標準form + Server Functionsを使い分けます。

## 推奨アーキテクチャ

```text
Browser
  ↓
Next.js App Router
├─ Server Components
├─ Server Functions
├─ Route Handlers
└─ Client Components
       ↓
PostgreSQL / Supabase
```

## 採用ルール

- UIはTailwind CSS + shadcn/uiを標準とする。
- Zodは入力境界で標準利用する。
- React Hook Formは複雑なフォームで利用する。
- TanStack Queryはクライアント側キャッシュが必要な場合だけ追加する。
- DBがSupabaseなら小規模案件ではORMを省略してもよい。
- Stripe、Redis等は必要になるまで追加しない。
- Vercel固有機能へ過度に依存せず、移植可能性を残す。

## 公式URL

- Next.js: https://nextjs.org/
- React: https://react.dev/
- TypeScript: https://www.typescriptlang.org/
- Tailwind CSS: https://tailwindcss.com/
- shadcn/ui: https://ui.shadcn.com/
- Base UI: https://base-ui.com/
- React Hook Form: https://react-hook-form.com/
- Zod: https://zod.dev/
- Auth0: https://auth0.com/
- Auth.js: https://authjs.dev/
- Supabase: https://supabase.com/
- Drizzle ORM: https://orm.drizzle.team/
- TanStack Query: https://tanstack.com/query/latest
- Stripe: https://stripe.com/
- Sentry: https://sentry.io/
- Vercel: https://vercel.com/
