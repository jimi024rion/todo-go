"use client"

import { useEffect, useRef, useState, useTransition } from "react"
import { Trash2, Calendar, GripVertical } from "lucide-react"
import { cn } from "@/lib/utils"
import type { Todo } from "@/types/todo"

function DueDateBadge({ dueDate }: { dueDate: string }) {
  const date = new Date(dueDate)
  const today = new Date(); today.setHours(0, 0, 0, 0)
  const tomorrow = new Date(today); tomorrow.setDate(today.getDate() + 1)
  const isOverdue = date < today
  const isToday = date >= today && date < tomorrow

  const label = isToday ? "今日" : date.toLocaleDateString("ja-JP", { month: "numeric", day: "numeric" })

  return (
    <span className={cn(
      "inline-flex items-center gap-0.5 rounded text-xs px-1.5 py-0.5",
      isOverdue ? "text-destructive bg-destructive/10" :
      isToday   ? "text-orange-600 bg-orange-50" :
                  "text-muted-foreground bg-secondary"
    )}>
      <Calendar className="h-3 w-3" />
      {label}
    </span>
  )
}

interface TodoCardProps {
  todo: Todo
  onToggle: (todo: Todo, done: boolean) => void
  onClick: (todo: Todo) => void
  onDelete: (id: string) => void
  dragHandleProps?: React.HTMLAttributes<HTMLDivElement>
  isDragging?: boolean
}

const SWIPE_THRESHOLD = 72   // この距離スワイプしたら削除ボタン表示
const DELETE_THRESHOLD = 140 // この距離まで引っ張ったら即削除

export function TodoCard({ todo, onToggle, onClick, onDelete, dragHandleProps, isDragging }: TodoCardProps) {
  const [isPending, startTransition] = useTransition()
  const isDone = todo.status === "completed"

  const cardRef = useRef<HTMLDivElement>(null)
  const touchStartX = useRef(0)
  const touchStartY = useRef(0)
  const [translateX, setTranslateX] = useState(0)
  const [isRevealed, setIsRevealed] = useState(false)
  const isHorizontalSwipe = useRef(false)

  function handleCheck(e: React.ChangeEvent<HTMLInputElement>) {
    e.stopPropagation()
    startTransition(() => { onToggle(todo, e.target.checked) })
  }

  // スワイプ処理（ネイティブイベントで passive:false）
  useEffect(() => {
    const card = cardRef.current
    if (!card) return

    function onTouchStart(e: TouchEvent) {
      touchStartX.current = e.touches[0].clientX
      touchStartY.current = e.touches[0].clientY
      isHorizontalSwipe.current = false
      card!.style.transition = "none"
    }

    function onTouchMove(e: TouchEvent) {
      const dx = e.touches[0].clientX - touchStartX.current
      const dy = e.touches[0].clientY - touchStartY.current

      // 最初の動きで方向を確定
      if (!isHorizontalSwipe.current && Math.abs(dx) < 5 && Math.abs(dy) < 5) return
      if (!isHorizontalSwipe.current) {
        isHorizontalSwipe.current = Math.abs(dx) > Math.abs(dy)
      }

      if (isHorizontalSwipe.current && dx < 0) {
        e.preventDefault()
        const clamped = Math.max(dx, -DELETE_THRESHOLD - 20)
        setTranslateX(clamped)
      }
    }

    function onTouchEnd() {
      card!.style.transition = "transform 0.25s ease"
      if (translateX < -DELETE_THRESHOLD) {
        // 即削除
        setTranslateX(-window.innerWidth)
        setTimeout(() => onDelete(todo.id), 200)
      } else if (translateX < -SWIPE_THRESHOLD) {
        // 削除ボタン表示
        setTranslateX(-SWIPE_THRESHOLD)
        setIsRevealed(true)
      } else {
        // 元に戻す
        setTranslateX(0)
        setIsRevealed(false)
      }
    }

    card.addEventListener("touchstart", onTouchStart, { passive: true })
    card.addEventListener("touchmove", onTouchMove, { passive: false })
    card.addEventListener("touchend", onTouchEnd, { passive: true })
    return () => {
      card.removeEventListener("touchstart", onTouchStart)
      card.removeEventListener("touchmove", onTouchMove)
      card.removeEventListener("touchend", onTouchEnd)
    }
  }, [todo.id, translateX, onDelete])

  function handleDeleteTap() {
    if (cardRef.current) {
      cardRef.current.style.transition = "transform 0.2s ease"
    }
    setTranslateX(-window.innerWidth)
    setTimeout(() => onDelete(todo.id), 200)
  }

  function resetSwipe() {
    setTranslateX(0)
    setIsRevealed(false)
  }

  return (
    <div className="relative overflow-hidden rounded-lg">
      {/* 削除ボタン（背景） */}
      <div
        className="absolute inset-y-0 right-0 flex items-center justify-center bg-destructive px-5 rounded-lg"
        style={{ width: SWIPE_THRESHOLD + 20 }}
        onClick={handleDeleteTap}
      >
        <Trash2 className="h-5 w-5 text-white" />
      </div>

      {/* カード本体 */}
      <div
        ref={cardRef}
        style={{ transform: `translateX(${translateX}px)` }}
        className={cn(
          "group flex cursor-pointer gap-3 rounded-lg border border-border bg-card p-4 shadow-sm transition-all duration-150",
          "hover:border-slate-300 hover:shadow-md",
          isPending && "opacity-60"
        )}
        onClick={() => {
          if (isRevealed) { resetSwipe(); return }
          onClick(todo)
        }}
        role="button"
        tabIndex={0}
        onKeyDown={(e) => e.key === "Enter" && !isRevealed && onClick(todo)}
      >
        {/* ドラッグハンドル（縦方向のみ移動。スワイプ削除と分離） */}
        {dragHandleProps && (
          <div
            {...dragHandleProps}
            className="flex shrink-0 touch-none items-center self-stretch px-0.5 text-muted-foreground/30 active:text-muted-foreground/60"
            onClick={(e) => e.stopPropagation()}
          >
            <GripVertical className="h-4 w-4" />
          </div>
        )}

        {/* チェックボックス */}
        <div className="mt-0.5 shrink-0" onClick={(e) => e.stopPropagation()}>
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
                  ? "border-foreground/40 bg-foreground/10"
                  : "border-border bg-background peer-focus-visible:ring-2 peer-focus-visible:ring-ring"
              )}
            >
              {isDone && (
                <svg className="h-3 w-3 text-foreground/50" viewBox="0 0 12 12" fill="none">
                  <path d="M2 6l3 3 5-5" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" />
                </svg>
              )}
            </span>
          </label>
        </div>

        {/* テキスト */}
        <div className="min-w-0 flex-1">
          <p className={cn(
            "text-sm leading-snug transition-colors duration-150",
            isDone ? "text-muted-foreground line-through" : "text-foreground"
          )}>
            {todo.title}
          </p>
          {todo.description && (
            <p className="mt-1 line-clamp-2 text-xs text-muted-foreground">{todo.description}</p>
          )}
          {todo.due_date && !isDone && (
            <div className="mt-1.5">
              <DueDateBadge dueDate={todo.due_date} />
            </div>
          )}
          {todo.tags && todo.tags.length > 0 && (
            <div className="mt-2 flex flex-wrap gap-1">
              {todo.tags.map((tag) => (
                <span key={tag.id} className="inline-flex items-center rounded-full bg-secondary px-2 py-0.5 text-xs text-muted-foreground">
                  {tag.name}
                </span>
              ))}
            </div>
          )}
        </div>
      </div>
    </div>
  )
}
