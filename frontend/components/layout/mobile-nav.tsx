"use client"

import { X } from "lucide-react"
import { Menu } from "lucide-react"
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

      {/* showCloseButton=false でデフォルトの×ボタン（はみ出す）を無効化 */}
      <SheetContent side="left" className="p-0 w-60" showCloseButton={false}>
        {/* ×ボタンをサイドバー内の右上に配置 */}
        <SheetClose
          className="absolute top-3 right-3 z-10 inline-flex h-8 w-8 items-center justify-center rounded-md text-muted-foreground hover:bg-secondary hover:text-foreground transition-colors duration-150"
          aria-label="メニューを閉じる"
        >
          <X className="h-4 w-4" />
        </SheetClose>

        <Sidebar pendingCount={pendingCount} />
      </SheetContent>
    </Sheet>
  )
}
