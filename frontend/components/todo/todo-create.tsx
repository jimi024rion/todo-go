"use client"

import { useRef, useState } from "react"
import { Plus } from "lucide-react"
import { cn } from "@/lib/utils"
import type { CreateTodoInput } from "@/types/todo"

interface TodoCreateProps {
  onSubmit: (input: CreateTodoInput) => Promise<void>
}

export function TodoCreate({ onSubmit }: TodoCreateProps) {
  const [title, setTitle] = useState("")
  const [description, setDescription] = useState("")
  const [expanded, setExpanded] = useState(false)
  const [isSubmitting, setIsSubmitting] = useState(false)
  const titleRef = useRef<HTMLInputElement>(null)
  const formRef = useRef<HTMLDivElement>(null)

  function handleFocus() {
    setExpanded(true)
  }

  // フォームの外側がクリックされたときに折りたたむ
  function handleBlur(e: React.FocusEvent) {
    // 次のフォーカス先がフォーム内であれば折りたたまない
    if (formRef.current?.contains(e.relatedTarget as Node)) return
    if (!title.trim() && !description.trim()) {
      setExpanded(false)
    }
  }

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    const trimmedTitle = title.trim()
    if (!trimmedTitle || isSubmitting) return

    setIsSubmitting(true)
    try {
      await onSubmit({ title: trimmedTitle, description: description.trim() || undefined })
      setTitle("")
      setDescription("")
      setExpanded(false)
    } finally {
      setIsSubmitting(false)
    }
  }

  function handleCancel() {
    setTitle("")
    setDescription("")
    setExpanded(false)
    titleRef.current?.blur()
  }

  function handleTitleKeyDown(e: React.KeyboardEvent) {
    if (e.key === "Escape") handleCancel()
    // Shift+Enter で説明欄にフォーカス移動
    if (e.key === "Enter" && !e.shiftKey) {
      e.preventDefault()
      if (title.trim()) handleSubmit(e as unknown as React.FormEvent)
    }
  }

  function handleDescriptionKeyDown(e: React.KeyboardEvent) {
    if (e.key === "Escape") handleCancel()
  }

  return (
    <div ref={formRef} onBlur={handleBlur}>
      <form onSubmit={handleSubmit}>
        <div
          className={cn(
            "rounded-lg border bg-card shadow-sm transition-all duration-150",
            expanded ? "border-primary/50 shadow-md" : "border-border"
          )}
        >
          {/* タイトル入力 */}
          <div className="flex items-center gap-2 px-3 py-2.5">
            <Plus
              className={cn(
                "h-4 w-4 shrink-0 transition-colors duration-150",
                expanded ? "text-primary" : "text-muted-foreground"
              )}
            />
            <input
              ref={titleRef}
              type="text"
              value={title}
              onChange={(e) => setTitle(e.target.value)}
              onFocus={handleFocus}
              onKeyDown={handleTitleKeyDown}
              placeholder="タスクを追加..."
              disabled={isSubmitting}
              className={cn(
                "flex-1 bg-transparent text-sm text-foreground",
                "placeholder:text-muted-foreground",
                "focus:outline-none",
                "disabled:opacity-50"
              )}
            />
          </div>

          {/* 展開エリア: 説明欄 + ボタン */}
          {expanded && (
            <>
              <div className="border-t border-border/60 px-3 py-2">
                <textarea
                  value={description}
                  onChange={(e) => setDescription(e.target.value)}
                  onKeyDown={handleDescriptionKeyDown}
                  placeholder="説明を追加（任意）"
                  rows={3}
                  disabled={isSubmitting}
                  className={cn(
                    "w-full resize-none bg-transparent text-sm text-foreground",
                    "placeholder:text-muted-foreground",
                    "focus:outline-none",
                    "disabled:opacity-50"
                  )}
                />
              </div>

              <div className="flex items-center justify-end gap-2 border-t border-border/60 px-3 py-2">
                <button
                  type="button"
                  onClick={handleCancel}
                  disabled={isSubmitting}
                  className={cn(
                    "rounded-md px-3 py-1.5 text-sm text-muted-foreground",
                    "hover:bg-secondary hover:text-foreground transition-colors duration-150",
                    "focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring",
                    "disabled:opacity-50"
                  )}
                >
                  キャンセル
                </button>
                <button
                  type="submit"
                  disabled={!title.trim() || isSubmitting}
                  className={cn(
                    "rounded-md bg-primary px-3 py-1.5 text-sm font-medium text-primary-foreground",
                    "hover:bg-primary/90 transition-colors duration-150",
                    "focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring",
                    "disabled:opacity-40 disabled:cursor-not-allowed"
                  )}
                >
                  {isSubmitting ? "追加中..." : "追加"}
                </button>
              </div>
            </>
          )}
        </div>
      </form>
    </div>
  )
}
