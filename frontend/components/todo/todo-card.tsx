"use client"

import { useTransition } from "react"
import { cn } from "@/lib/utils"
import type { Todo } from "@/types/todo"

interface TodoCardProps {
  todo: Todo
  onToggle: (id: string, done: boolean) => void
  onClick: (todo: Todo) => void
}

export function TodoCard({ todo, onToggle, onClick }: TodoCardProps) {
  const [isPending, startTransition] = useTransition()
  const isDone = todo.status === "completed"

  function handleCheck(e: React.ChangeEvent<HTMLInputElement>) {
    e.stopPropagation()
    startTransition(() => {
      onToggle(todo.id, e.target.checked)
    })
  }

  return (
    <div
      className={cn(
        "group flex cursor-pointer gap-3 rounded-lg border border-border bg-card p-4 shadow-sm transition-all duration-150",
        "hover:border-slate-300 hover:shadow-md",
        isPending && "opacity-60"
      )}
      onClick={() => onClick(todo)}
      role="button"
      tabIndex={0}
      onKeyDown={(e) => e.key === "Enter" && onClick(todo)}
    >
      {/* チェックボックス */}
      <div
        className="mt-0.5 shrink-0"
        onClick={(e) => e.stopPropagation()}
      >
        <label className="relative flex h-4 w-4 cursor-pointer items-center">
          <input
            type="checkbox"
            checked={isDone}
            onChange={handleCheck}
            className="peer sr-only"
            aria-label={`${todo.title}を${isDone ? "未完了" : "完了"}にする`}
          />
          <span
            className={cn(
              "flex h-4 w-4 items-center justify-center rounded-sm border transition-colors duration-150",
              isDone
                ? "border-green-500 bg-green-500"
                : "border-border bg-background peer-focus-visible:ring-2 peer-focus-visible:ring-ring"
            )}
          >
            {isDone && (
              <svg className="h-3 w-3 text-white" viewBox="0 0 12 12" fill="none">
                <path
                  d="M2 6l3 3 5-5"
                  stroke="currentColor"
                  strokeWidth="2"
                  strokeLinecap="round"
                  strokeLinejoin="round"
                />
              </svg>
            )}
          </span>
        </label>
      </div>

      {/* テキスト */}
      <div className="min-w-0 flex-1">
        <p
          className={cn(
            "text-sm leading-snug transition-colors duration-150",
            isDone
              ? "text-muted-foreground line-through"
              : "text-foreground"
          )}
        >
          {todo.title}
        </p>
        {todo.description && (
          <p className="mt-1 line-clamp-2 text-xs text-muted-foreground">
            {todo.description}
          </p>
        )}
      </div>
    </div>
  )
}
