"use client"

import { createPortal } from "react-dom"
import { useCallback, useEffect, useRef, useState } from "react"
import { Trash2, Plus, X } from "lucide-react"
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query"
import { Sheet, SheetContent } from "@/components/ui/sheet"
import { cn } from "@/lib/utils"
import { useIsMobile } from "@/lib/use-media-query"
import { tagApi, todoApi } from "@/lib/api"
import type { Tag, Todo, UpdateTodoInput } from "@/types/todo"

const MIN_DESKTOP_WIDTH = 280
const MAX_DESKTOP_WIDTH = 800
const DEFAULT_DESKTOP_WIDTH = 400

type HeightMode = "partial" | "full"

interface TodoPanelProps {
  todo: Todo | null
  onClose: () => void
  onUpdate: (id: string, input: UpdateTodoInput) => Promise<void>
  onDelete: (id: string) => Promise<void>
}

export function TodoPanel({ todo, onClose, onUpdate, onDelete }: TodoPanelProps) {
  const isMobile = useIsMobile()
  const [panelWidth, setPanelWidth] = useState(DEFAULT_DESKTOP_WIDTH)
  const [heightMode, setHeightMode] = useState<HeightMode>("partial")
  // スクロールコンテナの ref を親で持ち、BottomSheet のスワイプ判定に使う
  const scrollRef = useRef<HTMLDivElement>(null)
  const isDraggingDesktop = useRef(false)

  useEffect(() => {
    if (!todo) setHeightMode("partial")
  }, [todo])

  const handleDesktopDragStart = useCallback((e: React.MouseEvent) => {
    e.preventDefault()
    isDraggingDesktop.current = true
    document.body.style.cursor = "ew-resize"
    document.body.style.userSelect = "none"
    function onMove(ev: MouseEvent) {
      if (!isDraggingDesktop.current) return
      setPanelWidth(Math.max(MIN_DESKTOP_WIDTH, Math.min(MAX_DESKTOP_WIDTH, window.innerWidth - ev.clientX)))
    }
    function onUp() {
      isDraggingDesktop.current = false
      document.body.style.cursor = ""
      document.body.style.userSelect = ""
      document.removeEventListener("mousemove", onMove)
      document.removeEventListener("mouseup", onUp)
    }
    document.addEventListener("mousemove", onMove)
    document.addEventListener("mouseup", onUp)
  }, [])

  const content = todo ? (
    <PanelContent
      todo={todo}
      onClose={onClose}
      onUpdate={onUpdate}
      onDelete={onDelete}
      scrollRef={scrollRef}
    />
  ) : null

  if (isMobile) {
    return (
      <MobileBottomSheet
        open={!!todo}
        onClose={onClose}
        heightMode={heightMode}
        onExpand={() => setHeightMode("full")}
        onCollapse={() => setHeightMode("partial")}
        scrollRef={scrollRef}
      >
        {content}
      </MobileBottomSheet>
    )
  }

  return (
    <Sheet open={!!todo} onOpenChange={(open) => !open && onClose()}>
      <SheetContent
        side="right"
        style={{ width: panelWidth, maxWidth: panelWidth }}
        className="p-0 overflow-hidden flex flex-col"
      >
        <div
          onMouseDown={handleDesktopDragStart}
          className="absolute left-0 top-0 h-full w-1 z-10 cursor-ew-resize hover:bg-primary/40 transition-colors duration-100"
        />
        {content}
      </SheetContent>
    </Sheet>
  )
}

// ─── カスタム ボトムシート ──────────────────────────────────────────

function MobileBottomSheet({
  open,
  onClose,
  heightMode,
  onExpand,
  onCollapse,
  scrollRef,
  children,
}: {
  open: boolean
  onClose: () => void
  heightMode: HeightMode
  onExpand: () => void
  onCollapse: () => void
  scrollRef: React.RefObject<HTMLDivElement | null>
  children: React.ReactNode
}) {
  const [mounted, setMounted] = useState(false)
  const [visible, setVisible] = useState(false)
  const sheetRef = useRef<HTMLDivElement>(null)

  // open/close アニメーション
  useEffect(() => {
    if (open) {
      setMounted(true)
      // 2フレーム待って transition が効くようにする
      requestAnimationFrame(() => requestAnimationFrame(() => setVisible(true)))
    } else {
      setVisible(false)
      const t = setTimeout(() => setMounted(false), 300)
      return () => clearTimeout(t)
    }
  }, [open])

  // body スクロールロック
  useEffect(() => {
    if (!open) return
    document.body.style.overflow = "hidden"
    return () => { document.body.style.overflow = "" }
  }, [open])

  // スワイプ操作（passive:false が必要なので useEffect で登録）
  useEffect(() => {
    if (!open || !mounted) return
    const sheet = sheetRef.current
    if (!sheet) return

    let startY = 0
    let startTime = 0
    let latestDY = 0
    let isGesture = false

    function onTouchStart(e: TouchEvent) {
      startY = e.touches[0].clientY
      startTime = Date.now()
      latestDY = 0
      isGesture = false
      sheet!.style.transition = "none"
    }

    function onTouchMove(e: TouchEvent) {
      const dy = e.touches[0].clientY - startY
      const scrollTop = scrollRef.current?.scrollTop ?? 0

      if (!isGesture) {
        if (dy > 8 && scrollTop === 0) isGesture = true        // 下スワイプ（閉じる or 縮小）
        else if (dy < -8 && heightMode === "partial") isGesture = true  // 上スワイプ（全画面）
      }

      if (isGesture) {
        e.preventDefault()
        latestDY = dy
        // 上方向は 0.3 の抵抗を付与、下方向はそのまま追従
        const visual = dy < 0 ? dy * 0.3 : dy
        sheet!.style.transform = `translateY(${visual}px)`
      }
    }

    function onTouchEnd() {
      if (!isGesture) return

      const velocity = latestDY / Math.max(1, Date.now() - startTime) // px/ms
      sheet!.style.transition = "transform 0.25s cubic-bezier(0.32,0.72,0,1)"
      sheet!.style.transform = ""

      const isDown = latestDY > 80 || (velocity > 0.4 && latestDY > 20)
      const isUp   = latestDY < -60 || (velocity < -0.4 && latestDY < -20)

      if (isDown) {
        heightMode === "full" ? onCollapse() : onClose()
      } else if (isUp) {
        onExpand()
      }
    }

    sheet.addEventListener("touchstart", onTouchStart, { passive: true })
    sheet.addEventListener("touchmove",  onTouchMove,  { passive: false })
    sheet.addEventListener("touchend",   onTouchEnd,   { passive: true })
    return () => {
      sheet.removeEventListener("touchstart", onTouchStart)
      sheet.removeEventListener("touchmove",  onTouchMove)
      sheet.removeEventListener("touchend",   onTouchEnd)
    }
  }, [open, mounted, heightMode, scrollRef, onClose, onExpand, onCollapse])

  if (!mounted) return null

  return createPortal(
    <div className="fixed inset-0 z-50 touch-none">
      {/* バックドロップ */}
      <div
        className={cn(
          "absolute inset-0 bg-black/30 transition-opacity duration-300",
          visible ? "opacity-100" : "opacity-0"
        )}
        onClick={onClose}
      />

      {/* シート本体（このdivにtransformを直接適用） */}
      <div
        ref={sheetRef}
        className={cn(
          "absolute bottom-0 left-0 right-0 bg-background flex flex-col overflow-hidden",
          "transition-[height,border-radius,transform] duration-300",
          heightMode === "full" ? "h-[100dvh] rounded-t-none" : "h-[60dvh] rounded-t-2xl",
          visible ? "translate-y-0" : "translate-y-full"
        )}
      >
        {children}
      </div>
    </div>,
    document.body
  )
}

// ─── パネルコンテンツ ────────────────────────────────────────────────

function PanelContent({
  todo,
  onClose,
  onUpdate,
  onDelete,
  scrollRef,
}: {
  todo: Todo
  onClose: () => void
  onUpdate: (id: string, input: UpdateTodoInput) => Promise<void>
  onDelete: (id: string) => Promise<void>
  scrollRef: React.RefObject<HTMLDivElement | null>
}) {
  const isDone = todo.status === "completed"
  const [isDeleting, setIsDeleting] = useState(false)

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
    <div className="flex flex-col h-full w-full touch-auto">
      {/* ヘッダー */}
      <div className="relative flex items-center px-5 py-4 border-b border-border shrink-0">
        <div className="absolute top-2 left-1/2 -translate-x-1/2 w-10 h-1 rounded-full bg-border" />
        <h2 className="text-base font-semibold text-foreground">
          タスクの詳細
        </h2>
      </div>

      {/* スクロール可能なコンテンツ */}
      <div ref={scrollRef} className="flex-1 overflow-y-auto px-5 py-5 space-y-6 touch-auto">
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
        <DueDateField
          dueDate={todo.due_date ?? null}
          onSave={(val) => fullUpdate(val ? { due_date: val } : { clear_due_date: true })}
        />
        <TagSelector todo={todo} onClose={onClose} />
      </div>

      {/* フッター */}
      <div className="flex items-center gap-2 px-5 py-4 border-t border-border shrink-0">
        <button
          onClick={() => fullUpdate({ status: isDone ? "pending" : "completed" })}
          className={cn(
            "flex-1 rounded-md px-4 py-2.5 text-sm font-medium transition-colors duration-150",
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
            "transition-colors duration-150 disabled:opacity-50"
          )}
          aria-label="タスクを削除"
        >
          <Trash2 className="h-4 w-4" />
        </button>
      </div>
    </div>
  )
}

// ─── 期限設定 UI ─────────────────────────────────────────────────────

function DueDateField({
  dueDate,
  onSave,
}: {
  dueDate: string | null
  onSave: (val: string | null) => Promise<void>
}) {
  const [value, setValue] = useState(
    dueDate ? new Date(dueDate).toISOString().split("T")[0] : ""
  )

  async function handleChange(e: React.ChangeEvent<HTMLInputElement>) {
    const v = e.target.value
    setValue(v)
    await onSave(v ? new Date(v).toISOString() : null)
  }

  return (
    <div className="w-full">
      <p className="mb-1.5 text-xs font-medium text-muted-foreground">期限</p>
      <div className="flex items-center gap-2">
        <input
          type="date"
          value={value}
          onChange={handleChange}
          className={cn(
            "rounded-md border border-border bg-background px-3 py-2 text-base",
            "focus:outline-none focus:ring-2 focus:ring-ring"
          )}
        />
        {value && (
          <button
            onClick={() => { setValue(""); onSave(null) }}
            className="text-xs text-muted-foreground hover:text-destructive transition-colors"
          >
            削除
          </button>
        )}
      </div>
    </div>
  )
}

// ─── タグ選択 UI ─────────────────────────────────────────────────────

function TagSelector({ todo, onClose: _onClose }: { todo: Todo; onClose: () => void }) {
  const queryClient = useQueryClient()
  const [showInput, setShowInput] = useState(false)
  const [newTagName, setNewTagName] = useState("")

  const { data: allTags = [] } = useQuery({
    queryKey: ["tags"],
    queryFn: tagApi.list,
  })

  const setTagsMutation = useMutation({
    mutationFn: (tagIds: string[]) => todoApi.setTags(todo.id, tagIds),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["todos"] }),
  })

  const createTagMutation = useMutation({
    mutationFn: tagApi.create,
    onSuccess: (newTag) => {
      queryClient.invalidateQueries({ queryKey: ["tags"] })
      // 作成後すぐにこの Todo にも付与
      const currentIds = (todo.tags ?? []).map((t) => t.id)
      setTagsMutation.mutate([...currentIds, newTag.id])
      setNewTagName("")
      setShowInput(false)
    },
  })

  const currentTagIds = new Set((todo.tags ?? []).map((t) => t.id))

  function toggleTag(tag: Tag) {
    const ids = new Set(currentTagIds)
    if (ids.has(tag.id)) {
      ids.delete(tag.id)
    } else {
      ids.add(tag.id)
    }
    setTagsMutation.mutate([...ids])
  }

  function handleCreateTag(e: React.FormEvent) {
    e.preventDefault()
    const name = newTagName.trim()
    if (!name) return
    createTagMutation.mutate(name)
  }

  return (
    <div className="w-full">
      <p className="mb-1.5 text-xs font-medium text-muted-foreground">タグ</p>

      {/* 付与済みタグ + 選択可能なタグ */}
      <div className="flex flex-wrap gap-1.5">
        {allTags.map((tag) => {
          const isSelected = currentTagIds.has(tag.id)
          return (
            <button
              key={tag.id}
              onClick={() => toggleTag(tag)}
              className={cn(
                "inline-flex items-center rounded-full px-2.5 py-1 text-xs transition-colors duration-150",
                isSelected
                  ? "bg-foreground/10 text-foreground font-medium"
                  : "bg-secondary text-muted-foreground hover:bg-secondary/80"
              )}
            >
              {isSelected && <span className="mr-1">✓</span>}
              {tag.name}
            </button>
          )
        })}

        {/* 新規タグ作成ボタン */}
        {!showInput && (
          <button
            onClick={() => setShowInput(true)}
            className="inline-flex items-center gap-1 rounded-full border border-dashed border-border px-2.5 py-1 text-xs text-muted-foreground hover:border-primary hover:text-foreground transition-colors duration-150"
          >
            <Plus className="h-3 w-3" />
            新規
          </button>
        )}
      </div>

      {/* 新規タグ入力フォーム */}
      {showInput && (
        <form onSubmit={handleCreateTag} className="mt-2 flex items-center gap-2">
          <input
            autoFocus
            type="text"
            value={newTagName}
            onChange={(e) => setNewTagName(e.target.value)}
            placeholder="タグ名..."
            className={cn(
              "flex-1 rounded-md border border-border bg-background px-2.5 py-1 text-sm",
              "focus:outline-none focus:ring-2 focus:ring-ring"
            )}
          />
          <button
            type="submit"
            disabled={!newTagName.trim()}
            className="rounded-md bg-primary px-2.5 py-1 text-xs text-primary-foreground disabled:opacity-40"
          >
            追加
          </button>
          <button
            type="button"
            onClick={() => { setShowInput(false); setNewTagName("") }}
            className="text-muted-foreground hover:text-foreground"
          >
            <X className="h-4 w-4" />
          </button>
        </form>
      )}
    </div>
  )
}

// ─── 編集可能フィールド ───────────────────────────────────────────────

function EditableField({
  label,
  value,
  multiline,
  placeholder,
  onSave,
  strikethrough,
}: {
  label: string
  value: string
  multiline: boolean
  placeholder?: string
  onSave: (value: string) => Promise<void>
  strikethrough: boolean
}) {
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
            onKeyDown={(e) => {
              // Ctrl/Cmd+Enter は改行として扱う（送信しない）
              if (e.key === "Enter" && (e.ctrlKey || e.metaKey)) {
                e.stopPropagation()
              }
            }}
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
          onClick={startEdit}
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
