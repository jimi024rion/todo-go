"use client"

import { Sidebar } from "./sidebar"
import { MobileNav } from "./mobile-nav"

interface AppShellProps {
  children: React.ReactNode
  pendingCount?: number
  title?: string
}

export function AppShell({ children, pendingCount = 0, title }: AppShellProps) {
  return (
    <div className="flex h-screen">
      {/* デスクトップ: 左サイドバー固定表示 */}
      <div className="hidden md:flex md:shrink-0">
        <Sidebar pendingCount={pendingCount} />
      </div>

      {/* メインコンテンツ */}
      <div className="flex flex-1 flex-col overflow-hidden">
        {/* モバイル: ハンバーガー + ページタイトル（標準的なモバイル UX） */}
        <header className="flex h-14 items-center border-b border-border px-4 md:hidden">
          <MobileNav pendingCount={pendingCount} />
          {title && (
            <span className="ml-3 text-base font-semibold text-foreground">{title}</span>
          )}
        </header>

        <main className="flex-1 overflow-auto">{children}</main>
      </div>
    </div>
  )
}
