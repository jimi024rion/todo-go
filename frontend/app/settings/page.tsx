import { AppShell } from "@/components/layout/app-shell"

export default function SettingsPage() {
  return (
    <AppShell>
      <div className="mx-auto max-w-2xl px-4 py-8">
        <h1 className="mb-6 text-xl font-semibold text-foreground">設定</h1>

        {/* プロフィールカード */}
        <div className="rounded-lg border border-border bg-card p-6 shadow-sm">
          <div className="flex items-center gap-4">
            <div className="flex h-12 w-12 items-center justify-center rounded-full bg-secondary text-lg font-semibold text-foreground">
              U
            </div>
            <div>
              <p className="text-sm font-medium text-foreground">ユーザー</p>
              <p className="text-xs text-muted-foreground">user@example.com</p>
            </div>
          </div>

          <div className="mt-6 border-t border-border pt-4">
            <p className="mb-3 text-xs text-muted-foreground">
              Firebase Auth が設定されると、ここにログイン情報が表示されます。
            </p>
            <button
              disabled
              className="rounded-md border border-border px-4 py-2 text-sm text-muted-foreground opacity-50 cursor-not-allowed"
            >
              ログアウト（Firebase Auth 設定後に有効）
            </button>
          </div>
        </div>
      </div>
    </AppShell>
  )
}
