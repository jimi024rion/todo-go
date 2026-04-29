"use client"

import { useCallback, useEffect, useRef, useState } from "react"
import { Trash2 } from "lucide-react"
import { Sheet, SheetContent, SheetTitle } from "@/components/ui/sheet"
import { cn } from "@/lib/utils"
import { useIsMobile } from "@/lib/use-media-query"
import type { Todo, UpdateTodoInput } from "@/types/todo"

const MIN_WIDTH = 280
const MAX_WIDTH = 800
const DEFAULT_WIDTH = 400

type HeightMode = "partial" | "full"

interface TodoPanelProps {
  todo: Todo | null
  onClose: () => void
  onUpdate: (id: string, input: UpdateTodoInput) => Promise<void>
  onDelete: (id: string) => Promise<void>
}

export function TodoPanel({ todo, onClose, onUpdate, onDelete }: TodoPanelProps) {
  const isMobile = useIsMobile()
  const [panelWidth, setPanelWidth] = useState(DEFAULT_WIDTH)
  const [heightMode, setHeightMode] = useState<HeightMode>("partial")
  const isDragging = useRef(false)

  // シートを閉じたときに高さをリセット
  useEffect(() => {
    if (!todo) setHeightMode("partial")
  }, [todo])

  const handleDragStart = useCallback((e: React.MouseEvent) => {
    e.preventDefault()
    isDragging.current = true
    document.body.style.cursor = "ew-resize"
    document.body.style.userSelect = "none"

    function onMouseMove(ev: MouseEvent) {
      if (!isDragging.current) return
      setPanelWidth(Math.max(MIN_WIDTH, Math.min(MAX_WIDTH, window.innerWidth - ev.clientX)))
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
          className={cn(
            "p-0 overflow-hidden flex flex-col transition-[height,border-radius] duration-300",
            heightMode === "full"
              ? "h-[100dvh] rounded-t-none"
              : "h-[85dvh] rounded-t-2xl"
          )}
        >
          {todo && (
            <PanelContent
              todo={todo}
              onClose={onClose}
              onUpdate={onUpdate}
              onDelete={onDelete}
              heightMode={heightMode}
              onExpand={() => setHeightMode("full")}
              onCollapse={() => setHeightMode("partial")}
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
            className="absolute left-0 top-0 h-full w-1 z-10 cursor-ew-resize hover:bg-primary/40 active:bg-primary/60 transition-colors duration-100"
          />
          {todo && (
            <PanelContent
              todo={todo}
              onClose={onClose}
              onUpdate={onUpdate}
              onDelete={onDelete}
              heightMode="full"
              onExpand={() => {}}
              onCollapse={() => {}}
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
  heightMode,
  onExpand,
  onCollapse,
  isMobile,
}: {
  todo: Todo
  onClose: () => void
  onUpdate: (id: string, input: UpdateTodoInput) => Promise<void>
  onDelete: (id: string) => Promise<void>
  heightMode: HeightMode
  onExpand: () => void
  onCollapse: () => void
  isMobile: boolean
}) {
  const isDone = todo.status === "completed"
  const [isDeleting, setIsDeleting] = useState(false)
  const containerRef = useRef<HTMLDivElement>(null)
  const scrollRef = useRef<HTMLDivElement>(null)

  // ネイティブタッチイベントでドラッグ検知（passive: false が必要なため useEffect で登録）
  useEffect(() => {
    if (!isMobile) return
    const container = containerRef.current
    if (!container) return
    // 以降 container は non-null
    const el = container

    let startY = 0
    let startTime = 0
    let latestDY = 0
    let isGesture = false

    function onTouchStart(e: TouchEvent) {
      startY = e.touches[0].clientY
      startTime = Date.now()
      latestDY = 0
      isGesture = false
      el.style.transition = "none"
    }

    function onTouchMove(e: TouchEvent) {
      const dy = e.touches[0].clientY - startY
      const scrollTop = scrollRef.current?.scrollTop ?? 0

      const isCloseGesture = dy > 5 && scrollTop === 0
      const isExpandGesture = dy < -5

      if (isCloseGesture || isExpandGesture) {
        isGesture = true
      }

      if (isGesture) {
        e.preventDefault()
        latestDY = dy
        const visual = dy < 0 ? dy * 0.35 : dy
        el.style.transform = `translateY(${visual}px)`
      }
    }

    function onTouchEnd() {
      if (!isGesture) return

      const velocity = latestDY / (Date.now() - startTime)
      el.style.transition = "transform 0.25s ease"
      el.style.transform = ""

      const isDownSwipe = latestDY > 80 || (velocity > 0.4 && latestDY > 20)
      const isUpSwipe = latestDY < -60 || (velocity < -0.4 && latestDY < -20)

      if (isDownSwipe) {
        if (heightMode === "full") {
          // full → partial
          onCollapse()
        } else {
          // partial → close
          onClose()
        }
      } else if (isUpSwipe) {
        // partial → full
        if (heightMode === "partial") onExpand()
      }
    }

    el.addEventListener("touchstart", onTouchStart, { passive: true })
    el.addEventListener("touchmove", onTouchMove, { passive: false })
    el.addEventListener("touchend", onTouchEnd, { passive: true })

    return () => {
      el.removeEventListener("touchstart", onTouchStart)
      el.removeEventListener("touchmove", onTouchMove)
      el.removeEventListener("touchend", onTouchEnd)
    }
  }, [isMobile, heightMode, onClose, onExpand, onCollapse])

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

  return (
    <div ref={containerRef} className="flex flex-col h-full w-full">
      {/* ヘッダー */}
      <div className="relative flex items-center justify-between px-5 py-4 border-b border-border shrink-0">
        {isMobile && (
          <div className="absolute top-2 left-1/2 -translate-x-1/2 w-10 h-1 rounded-full bg-border" />
        )}
        <SheetTitle className="text-base font-semibold text-foreground">
          タスクの詳細
        </SheetTitle>
      </div>

      {/* スクロール可能なコンテンツ */}
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
          onClick={() => fullUpdate({ status: isDone ? "pending" : "completed" })}
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

function EditableField({ label, value, multiline, placeholder, onSave, strikethrough }: EditableFieldProps) {
  const [editing, setEditing] = useState(false)
  const [draft, setDraft] = useState(value)

  function startEdit() {
    setDraft(value)
    setEditing(true)
  }

  async function handleBlur() {
    setEditing(false)
    const trimmed = draft.trim()
    if (trimmed !== value) await onSave(trimmed)
  }

  return (
    <div className="w-full">
      <p className="mb-1.5 text-xs font-medium text-muted-foreground">{label}</p>
      {editing ? (
        multiline ? (
          <textarea
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
