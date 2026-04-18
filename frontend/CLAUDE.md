@AGENTS.md

# Frontend Development Guide

## Design System

**必ず `DESIGN.md` に従って実装すること。** 色・スペーシング・コンポーネントスタイルはすべて DESIGN.md が正。

Key: Orange accent (`#F97316`), warm off-white background (`#F7F6F2`), stone text palette.

## Page Structure

```
/           → Todo一覧（APIキー未設定時は /register にリダイレクト）
/register   → 初回ユーザー登録画面（APIキー設定済みなら / にリダイレクト）
```

## Components

| ファイル | 種別 | 役割 |
|---|---|---|
| `app/page.tsx` | Server | メインページ・SSR・認証チェック |
| `app/register/page.tsx` | Server | 登録ページ・認証チェック |
| `components/RegisterForm.tsx` | Client | 名前・メール入力フォーム → `/api/register` |
| `components/TodoList.tsx` | Client | フィルタータブ・アクティブ一覧・完了セクション |
| `components/TodoCard.tsx` | Client | カード表示・ステータス変更・編集・削除 |
| `components/EditTodoModal.tsx` | Client | タイトル・説明の編集モーダル |
| `components/StatusBadge.tsx` | Server | ステータスバッジ（pill形式） |
| `components/TodoForm.tsx` | Client | 新規Todo作成フォーム |
| `components/ResetKeyButton.tsx` | Client | APIキーリセット |

## API Routes

| Route | 処理 |
|---|---|
| `POST /api/register` | ユーザー作成 + APIキー発行 + Cookie設定 |
| `POST /api/auth` | APIキーをCookieに保存 |
| `DELETE /api/auth` | Cookie削除 |
| `GET/POST /api/todos` | Todo一覧取得・作成（バックエンドへのプロキシ） |
| `GET/PUT/DELETE /api/todos/[id]` | Todo詳細・更新・削除 |

## Key Logic

### 登録フロー
1. `GET /` → Cookie `api_key` なし → `redirect('/register')`
2. `/register` → RegisterForm送信 → `POST /api/register`
3. API Route: `POST /v1/users` → `POST /v1/api-keys` → Cookie設定
4. `router.push('/')` でメイン画面へ

### フィルタリング
- クライアントサイドでフィルタリング（APIコールなし）
- `activeFilter: "all" | "pending" | "in_progress"`
- 完了Todoは常に下の別セクションに分離

### Todo編集
- `EditTodoModal` に `todo`, `onClose`, `onSaved` を渡す
- `PUT /api/todos/:id` で更新
- 保存後 `onUpdated(updatedTodo)` で親の状態を更新

## Development Commands

```bash
npm run dev    # 開発サーバー起動 (localhost:3000)
npm run build  # プロダクションビルド
```

バックエンドは別途 `make air` で起動が必要（localhost:8080）
