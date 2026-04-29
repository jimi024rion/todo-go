"use client"

import { useCallback, useRef, useState } from "react"
import { Trash2 } from "lucide-react"
import { Sheet, SheetContent, SheetTitle } from "@/components/ui/sheet"
import { cn } from "@/lib/utils"
import { useIsMobile } from "@/lib/use-media-query"
import type { Todo, UpdateTodoInput } from "@/types/todo"

const MIN_WIDTH = 280
const MAX_WIDTH = 800
const DEFAULT_WIDTH = 400

interface TodoPanelProps {
  todo: Todo | null
  onClose: () => void
  onUpdate: (id: string, input: UpdateTodoInput) => Promise<void>
  onDelete: (id: string) => Promise<void>
}

export function TodoPanel({ todo, onClose, onUpdate, onDelete }: TodoPanelProps) {
  const isMobile = useIsMobile()
  const [panelWidth, setPanelWidth] = useState(DEFAULT_WIDTH)
  const isDragging = useRef(false)

  const handleDragStart = useCallback((e: React.MouseEvent) => {
    e.preventDefault()
    isDragging.current = true
    document.body.style.cursor = "ew-resize"
    document.body.style.userSelect = "none"

    function onMouseMove(ev: MouseEvent) {
      if (!isDragging.current) return
      const newWidth = window.innerWidth - ev.clientX
      setPanelWidth(Math.max(MIN_WIDTH, Math.min(MAX_WIDTH, newWidth)))
    }

    function onMouseUp() {
      isDragging.current = false
      document.body.style.cursor = ""
      document.body.style.userSelect = ""
      document.removeEventListener("mousemove", onMouseMove)
      document.removeEventListener("mouseup", onMouseUp)
    }

    document.addEventListener("mousemove", onMouseMove)
    document.addEventListener("mouseup", onMouseUp)
  }, [])

  return (
    <Sheet open={!!todo} onOpenChange={(open) => !open && onClose()}>
      {isMobile ? (
        <SheetContent
          side="bottom"
          className="p-0 overflow-hidden flex flex-col rounded-t-2xl h-[85dvh]"
        >
          {todo && (
            <PanelContent
              todo={todo}
              onClose={onClose}
              onUpdate={onUpdate}
              onDelete={onDelete}
              isMobile
            />
          )}
        </SheetContent>
      ) : (
        <SheetContent
          side="right"
          style={{ width: panelWidth, maxWidth: panelWidth }}
          className="p-0 overflow-hidden flex flex-col"
        >
          <div
            onMouseDown={handleDragStart}
            title="ドラッグして幅を変更"
            className="absolute left-0 top-0 h-full w-1 z-10 cursor-ew-resize hover:bg-primary/40 active:bg-primary/60 transition-colors duration-100"
          />
          {todo && (
            <PanelContent
              todo={todo}
              onClose={onClose}
              onUpdate={onUpdate}
              onDelete={onDelete}
              isMobile={false}
            />
          )}
        </SheetContent>
      )}
    </Sheet>
  )
}

function PanelContent({
  todo,
  onClose,
  onUpdate,
  onDelete,
  isMobile,
}: {
  todo: Todo
  onClose: () => void
  onUpdate: (id: string, input: UpdateTodoInput) => Promise<void>
  onDelete: (id: string) => Promise<void>
  isMobile: boolean
}) {
  const isDone = todo.status === "completed"
  const [isDeleting, setIsDeleting] = useState(false)
  const scrollRef = useRef<HTMLDivElement>(null)
  const touchStartY = useRef(0)
  const touchStartTime = useRef(0)

  // シート全体でスワイプ閉じを検知（スクロール中は干渉しない）
  function handleTouchStart(e: React.TouchEvent) {
    touchStartY.current = e.touches[0].clientY
    touchStartTime.current = Date.now()
  }

  function handleTouchEnd(e: React.TouchEvent) {
    const deltaY = e.changedTouches[0].clientY - touchStartY.current
    const deltaTime = Date.now() - touchStartTime.current
    const velocity = deltaY / deltaTime // px/ms
    const scrollTop = scrollRef.current?.scrollTop ?? 0

    // スクロールエリアが先頭にある場合のみ判定
    // 距離が大きい(60px) OR 素早いフリック(0.3px/ms かつ 20px以上)
    if (scrollTop === 0 && (deltaY > 60 || (velocity > 0.3 && deltaY > 20))) {
      onClose()
    }
  }

  async function handleDelete() {
    if (!confirm("このタスクを削除しますか？")) return
    setIsDeleting(true)
    await onDelete(todo.id)
    onClose()
  }

  async function fullUpdate(patch: UpdateTodoInput) {
    await onUpdate(todo.id, {
      title: todo.title,
      description: todo.description,
      status: todo.status,
      ...patch,
    })
  }

  async function handleToggleDone() {
    await fullUpdate({ status: isDone ? "pending" : "completed" })
  }

  return (
    <div
      className="flex flex-col h-full w-full"
      onTouchStart={isMobile ? handleTouchStart : undefined}
      onTouchEnd={isMobile ? handleTouchEnd : undefined}
    >
      {/* ヘッダー */}
      <div className="relative flex items-center justify-between px-5 py-4 border-b border-border shrink-0">
        {isMobile && (
          <div className="absolute top-2 left-1/2 -translate-x-1/2 w-10 h-1 rounded-full bg-border" />
        )}
        <SheetTitle className="text-base font-semibold text-foreground">
          タスクの詳細
        </SheetTitle>
      </div>

      {/* スクロール可能なコンテンツエリア */}
      <div ref={scrollRef} className="flex-1 overflow-y-auto px-5 py-5 space-y-6">
        <EditableField
          label="タイトル"
          value={todo.title}
          multiline={false}
          onSave={(val) => fullUpdate({ title: val })}
          strikethrough={isDone}
        />
        <EditableField
          label="説明"
          value={todo.description}
          multiline={true}
          placeholder="説明を追加..."
          onSave={(val) => fullUpdate({ description: val })}
          strikethrough={false}
        />
      </div>

      {/* フッター */}
      <div className="flex items-center gap-2 px-5 py-4 border-t border-border shrink-0">
        <button
          onClick={handleToggleDone}
          className={cn(
            "flex-1 rounded-md px-4 py-2.5 text-sm font-medium transition-colors duration-150",
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
            "shrink-0 rounded-md px-3 py-2.5 text-sm font-medium text-destructive",
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
    </div>
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

  function startEdit() {
    setDraft(value)
    setEditing(true)
    // autoFocus が iOS のユーザージェスチャー内でキーボードを起動する
  }

  async function handleBlur() {
    setEditing(false)
    const trimmed = draft.trim()
    if (trimmed !== value) {
      await onSave(trimmed)
    }
  }

  return (
    <div className="w-full">
      <p className="mb-1.5 text-xs font-medium text-muted-foreground">{label}</p>
      {editing ? (
        multiline ? (
          <textarea
            // eslint-disable-next-line jsx-a11y/no-autofocus
            autoFocus
            value={draft}
            onChange={(e) => setDraft(e.target.value)}
            onBlur={handleBlur}
            rows={4}
            className={cn(
              "w-full resize-none rounded-md border border-border bg-background px-3 py-2 text-base",
              "focus:outline-none focus:ring-2 focus:ring-ring"
            )}
          />
        ) : (
          <input
            // eslint-disable-next-line jsx-a11y/no-autofocus
            autoFocus
            type="text"
            value={draft}
            onChange={(e) => setDraft(e.target.value)}
            onBlur={handleBlur}
            className={cn(
              "w-full rounded-md border border-border bg-background px-3 py-2 text-base",
              "focus:outline-none focus:ring-2 focus:ring-ring"
            )}
          />
        )
      ) : (
        <p
          onPointerDown={startEdit}
          className={cn(
            "w-full cursor-text rounded-md px-3 py-2 text-sm min-h-[2.5rem]",
            "active:bg-secondary/70 transition-colors duration-100",
            value ? "text-foreground" : "text-muted-foreground italic",
            strikethrough && "line-through text-muted-foreground"
          )}
        >
          {value || placeholder}
        </p>
      )}
    </div>
  )
}
