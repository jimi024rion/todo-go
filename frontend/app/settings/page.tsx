"use client"

import { useEffect } from "react"
import { useRouter } from "next/navigation"
import { AppShell } from "@/components/layout/app-shell"
import { useAuth } from "@/lib/auth-context"

export default function SettingsPage() {
  const { user, loading, signOut } = useAuth()
  const router = useRouter()

  useEffect(() => {
    if (!loading && !user) {
      router.replace("/login")
    }
  }, [user, loading, router])

  const handleSignOut = async () => {
    await signOut()
    router.replace("/login")
  }

  if (loading) {
    return (
      <div className="flex min-h-screen items-center justify-center bg-background">
        <div className="h-6 w-6 animate-spin rounded-full border-2 border-border border-t-foreground" />
      </div>
    )
  }

  if (!user) return null

  const initial = user.displayName?.[0]?.toUpperCase() ?? user.email?.[0]?.toUpperCase() ?? "U"

  return (
    <AppShell>
      <div className="mx-auto max-w-2xl px-4 py-8">
        <h1 className="mb-6 text-xl font-semibold text-foreground">設定</h1>

        <div className="rounded-lg border border-border bg-card p-6 shadow-sm">
          <div className="flex items-center gap-4">
            {user.photoURL ? (
              <img
                src={user.photoURL}
                alt={user.displayName ?? "User"}
                className="h-12 w-12 rounded-full object-cover"
              />
            ) : (
              <div className="flex h-12 w-12 items-center justify-center rounded-full bg-secondary text-lg font-semibold text-foreground">
                {initial}
              </div>
            )}
            <div>
              <p className="text-sm font-medium text-foreground">
                {user.displayName ?? "ユーザー"}
              </p>
              <p className="text-xs text-muted-foreground">{user.email}</p>
            </div>
          </div>

          <div className="mt-6 border-t border-border pt-4">
            <button
              onClick={handleSignOut}
              className="rounded-md border border-border px-4 py-2 text-sm text-foreground transition-colors hover:bg-secondary"
            >
              ログアウト
            </button>
          </div>
        </div>
      </div>
    </AppShell>
  )
}
