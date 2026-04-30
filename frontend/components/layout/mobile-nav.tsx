"use client"

import { ArrowLeft, Menu } from "lucide-react"
import { Sheet, SheetClose, SheetContent, SheetTrigger } from "@/components/ui/sheet"
import { Sidebar } from "./sidebar"

interface MobileNavProps {
  pendingCount?: number
}

export function MobileNav({ pendingCount = 0 }: MobileNavProps) {
  return (
    <Sheet>
      <SheetTrigger
        className="inline-flex h-9 w-9 items-center justify-center rounded-md hover:bg-secondary transition-colors duration-150 md:hidden"
        aria-label="メニューを開く"
      >
        <Menu className="h-5 w-5 text-foreground" />
      </SheetTrigger>

      {/* フルスクリーンドロワー: バックドロップなし、幅 100% */}
      <SheetContent side="left" className="p-0 w-full" showCloseButton={false}>
        {/* ← 戻るボタン + アプリ名をヘッダーに配置 */}
        <div className="flex h-14 items-center gap-2 border-b border-border px-4">
          <SheetClose
            className="inline-flex h-9 w-9 items-center justify-center rounded-md text-muted-foreground hover:bg-secondary hover:text-foreground transition-colors duration-150"
            aria-label="メニューを閉じる"
          >
            <ArrowLeft className="h-5 w-5" />
          </SheetClose>
          <span className="text-xl font-semibold text-foreground">todos</span>
        </div>

        {/* ナビゲーション（Sidebar からヘッダーを除いて表示） */}
        <Sidebar pendingCount={pendingCount} hideHeader className="w-full border-r-0" />
      </SheetContent>
    </Sheet>
  )
}
