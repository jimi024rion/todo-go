"use client"

import { Menu } from "lucide-react"
import { Sheet, SheetContent, SheetTrigger } from "@/components/ui/sheet"
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
      <SheetContent side="left" className="p-0 w-60">
        <Sidebar pendingCount={pendingCount} />
      </SheetContent>
    </Sheet>
  )
}
