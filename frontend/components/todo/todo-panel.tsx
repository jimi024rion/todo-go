"use client"

import { useRef, useState } from "react"
import { Trash2 } from "lucide-react"
import { Sheet, SheetContent, SheetHeader, SheetTitle } from "@/components/ui/sheet"
import { cn } from "@/lib/utils"
import type { Todo, UpdateTodoInput } from "@/types/todo"

interface TodoPanelProps {
  todo: Todo | null
  onClose: () => void
  onUpdate: (id: string, input: UpdateTodoInput) => Promise<void>
  onDelete: (id: string) => Promise<void>
}

export function TodoPanel({ todo, onClose, onUpdate, onDelete }: TodoPanelProps) {
  return (
    <Sheet open={!!todo} onOpenChange={(open) => !open && onClose()}>
      <SheetContent side="right" className="w-full max-w-md overflow-y-auto">
        {todo && (
          <PanelContent
            todo={todo}
            onClose={onClose}
            onUpdate={onUpdate}
            onDelete={onDelete}
          />
        )}
      </SheetContent>
    </Sheet>
  )
}

function PanelContent({
  todo,
  onClose,
  onUpdate,
  onDelete,
}: {
  todo: Todo
  onClose: () => void
  onUpdate: (id: string, input: UpdateTodoInput) => Promise<void>
  onDelete: (id: string) => Promise<void>
}) {
  const isDone = todo.status === "completed"
  const [isDeleting, setIsDeleting] = useState(false)

  async function handleDelete() {
    if (!confirm("このタスクを削除しますか？")) return
    setIsDeleting(true)
    await onDelete(todo.id)
    onClose()
  }

  async function handleToggleDone() {
    await onUpdate(todo.id, { status: isDone ? "pending" : "completed" })
  }

  return (
    <>
      <SheetHeader className="mb-6">
        <SheetTitle className="text-left text-base font-semibold text-foreground">
          タスクの詳細
        </SheetTitle>
      </SheetHeader>

      <div className="space-y-6">
        {/* タイトル */}
        <EditableField
          label="タイトル"
          value={todo.title}
          multiline={false}
          onSave={(val) => onUpdate(todo.id, { title: val })}
          strikethrough={isDone}
        />

        {/* 説明文 */}
        <EditableField
          label="説明"
          value={todo.description}
          multiline={true}
          placeholder="説明を追加..."
          onSave={(val) => onUpdate(todo.id, { description: val })}
          strikethrough={false}
        />

        {/* メタ情報 */}
        <div className="text-xs text-muted-foreground">
          <p>作成日: {new Date(todo.created_at).toLocaleString("ja-JP")}</p>
          <p>更新日: {new Date(todo.updated_at).toLocaleString("ja-JP")}</p>
        </div>
      </div>

      {/* フッター */}
      <div className="mt-8 flex gap-2 border-t border-border pt-4">
        <button
          onClick={handleToggleDone}
          className={cn(
            "flex-1 rounded-md px-4 py-2 text-sm font-medium transition-colors duration-150",
            "focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring",
            isDone
              ? "bg-secondary text-secondary-foreground hover:bg-secondary/80"
              : "bg-primary text-primary-foreground hover:bg-primary/90"
          )}
        >
          {isDone ? "未完了に戻す" : "完了にする"}
        </button>
        <button
          onClick={handleDelete}
          disabled={isDeleting}
          className={cn(
            "rounded-md px-4 py-2 text-sm font-medium text-destructive",
            "border border-destructive/30 hover:bg-destructive hover:text-destructive-foreground",
            "transition-colors duration-150",
            "focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring",
            "disabled:opacity-50"
          )}
          aria-label="タスクを削除"
        >
          <Trash2 className="h-4 w-4" />
        </button>
      </div>
    </>
  )
}

interface EditableFieldProps {
  label: string
  value: string
  multiline: boolean
  placeholder?: string
  onSave: (value: string) => Promise<void>
  strikethrough: boolean
}

function EditableField({
  label,
  value,
  multiline,
  placeholder,
  onSave,
  strikethrough,
}: EditableFieldProps) {
  const [editing, setEditing] = useState(false)
  const [draft, setDraft] = useState(value)
  const ref = useRef<HTMLInputElement & HTMLTextAreaElement>(null)

  function startEdit() {
    setDraft(value)
    setEditing(true)
    // 次フレームでフォーカス
    setTimeout(() => ref.current?.focus(), 0)
  }

  async function handleBlur() {
    setEditing(false)
    const trimmed = draft.trim()
    if (trimmed !== value) {
      await onSave(trimmed)
    }
  }

  const displayValue = value || placeholder

  return (
    <div>
      <p className="mb-1 text-xs font-medium text-muted-foreground">{label}</p>
      {editing ? (
        multiline ? (
          <textarea
            ref={ref as React.RefObject<HTMLTextAreaElement>}
            value={draft}
            onChange={(e) => setDraft(e.target.value)}
            onBlur={handleBlur}
            rows={4}
            className={cn(
              "w-full resize-none rounded-md border border-border bg-background px-3 py-2 text-sm",
              "focus:outline-none focus:ring-2 focus:ring-ring"
            )}
          />
        ) : (
          <input
            ref={ref as React.RefObject<HTMLInputElement>}
            type="text"
            value={draft}
            onChange={(e) => setDraft(e.target.value)}
            onBlur={handleBlur}
            className={cn(
              "w-full rounded-md border border-border bg-background px-3 py-2 text-sm",
              "focus:outline-none focus:ring-2 focus:ring-ring"
            )}
          />
        )
      ) : (
        <p
          onClick={startEdit}
          className={cn(
            "cursor-text rounded-md px-3 py-2 text-sm hover:bg-secondary/50 transition-colors duration-150",
            value ? "text-foreground" : "text-muted-foreground italic",
            strikethrough && "line-through text-muted-foreground"
          )}
        >
          {displayValue || placeholder}
        </p>
      )}
    </div>
  )
}
