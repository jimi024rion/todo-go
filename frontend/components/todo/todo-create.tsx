"use client"

import { useRef, useState } from "react"
import { Plus } from "lucide-react"
import { cn } from "@/lib/utils"

interface TodoCreateProps {
  onSubmit: (title: string) => Promise<void>
}

export function TodoCreate({ onSubmit }: TodoCreateProps) {
  const [value, setValue] = useState("")
  const [isSubmitting, setIsSubmitting] = useState(false)
  const inputRef = useRef<HTMLInputElement>(null)

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    const title = value.trim()
    if (!title || isSubmitting) return

    setIsSubmitting(true)
    try {
      await onSubmit(title)
      setValue("")
    } finally {
      setIsSubmitting(false)
    }
  }

  function handleKeyDown(e: React.KeyboardEvent) {
    if (e.key === "Escape") {
      setValue("")
      inputRef.current?.blur()
    }
  }

  return (
    <form onSubmit={handleSubmit} className="flex items-center gap-2">
      <div className="relative flex-1">
        <Plus className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
        <input
          ref={inputRef}
          type="text"
          value={value}
          onChange={(e) => setValue(e.target.value)}
          onKeyDown={handleKeyDown}
          placeholder="タスクを追加..."
          disabled={isSubmitting}
          className={cn(
            "w-full rounded-md border border-border bg-card py-2 pl-9 pr-3 text-sm",
            "placeholder:text-muted-foreground",
            "focus:outline-none focus:ring-2 focus:ring-ring focus:border-primary",
            "transition-all duration-150",
            "disabled:opacity-50"
          )}
        />
      </div>
      <button
        type="submit"
        disabled={!value.trim() || isSubmitting}
        className={cn(
          "rounded-md bg-primary px-3 py-2 text-sm font-medium text-primary-foreground",
          "hover:bg-primary/90 transition-colors duration-150",
          "disabled:opacity-40 disabled:cursor-not-allowed",
          "focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
        )}
      >
        追加
      </button>
    </form>
  )
}
