"use client"

import { Sidebar } from "./sidebar"
import { MobileNav } from "./mobile-nav"

interface AppShellProps {
  children: React.ReactNode
  pendingCount?: number
}

export function AppShell({ children, pendingCount = 0 }: AppShellProps) {
  return (
    <div className="flex h-screen">
      {/* デスクトップ: 左サイドバー固定表示 */}
      <div className="hidden md:flex md:shrink-0">
        <Sidebar pendingCount={pendingCount} />
      </div>

      {/* メインコンテンツ */}
      <div className="flex flex-1 flex-col overflow-hidden">
        {/* モバイル: トップバーにハンバーガーメニュー */}
        <header className="flex h-14 items-center border-b border-border px-4 md:hidden">
          <MobileNav pendingCount={pendingCount} />
          <span className="ml-2 text-lg font-semibold">Todo</span>
        </header>

        <main className="flex-1 overflow-auto">{children}</main>
      </div>
    </div>
  )
}
