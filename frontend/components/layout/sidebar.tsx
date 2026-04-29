"use client"

import Link from "next/link"
import { usePathname } from "next/navigation"
import { CheckSquare, Settings } from "lucide-react"
import { cn } from "@/lib/utils"

interface SidebarProps {
  pendingCount?: number
}

export function Sidebar({ pendingCount = 0 }: SidebarProps) {
  const pathname = usePathname()

  const navItems = [
    {
      href: "/todos",
      label: "Todos",
      icon: CheckSquare,
      badge: pendingCount > 0 ? pendingCount : null,
    },
    {
      href: "/settings",
      label: "Settings",
      icon: Settings,
      badge: null,
    },
  ]

  return (
    <aside className="flex h-full w-60 flex-col border-r border-border bg-background">
      {/* App name */}
      <div className="flex h-14 items-center px-4">
        <span className="text-xl font-semibold text-foreground">Todo</span>
      </div>

      {/* Navigation */}
      <nav className="flex-1 space-y-1 px-3 py-2">
        {navItems.map(({ href, label, icon: Icon, badge }) => {
          const isActive = pathname.startsWith(href)
          return (
            <Link
              key={href}
              href={href}
              className={cn(
                "flex items-center gap-3 rounded-md px-3 py-2 text-sm transition-all duration-150",
                isActive
                  ? "bg-secondary font-medium text-foreground"
                  : "text-muted-foreground hover:bg-secondary hover:text-foreground"
              )}
            >
              <Icon className="h-4 w-4 shrink-0" />
              <span className="flex-1">{label}</span>
              {badge !== null && (
                <span className="flex h-5 min-w-5 items-center justify-center rounded-full bg-background px-1.5 text-xs font-medium text-muted-foreground">
                  {badge}
                </span>
              )}
            </Link>
          )
        })}
      </nav>

      {/* User area */}
      <div className="border-t border-border px-3 py-3">
        <div className="flex items-center gap-3 rounded-md px-3 py-2">
          <div className="flex h-7 w-7 shrink-0 items-center justify-center rounded-full bg-secondary text-xs font-medium text-foreground">
            U
          </div>
          <span className="truncate text-sm text-muted-foreground">ユーザー</span>
        </div>
      </div>
    </aside>
  )
}
